package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"first-project-go/config"
	"first-project-go/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // Menyiapkan upgrader WebSocket dengan CheckOrigin yang mengizinkan semua origin

func GetHistory(c *gin.Context) { // Handler untuk mengambil riwayat chat
	myID := c.Param("myID")                                                                           // Mendapatkan myID dari parameter URL
	targetID := c.Param("targetID")                                                                   // Mendapatkan targetID dari parameter URL
	filter := bson.M{"$or": []bson.M{{"from": myID, "to": targetID}, {"from": targetID, "to": myID}}} // Membuat filter query MongoDB untuk mencari pesan antara dua user

	cursor, _ := config.MsgCollection.Find(context.TODO(), filter) // Menjalankan query Find ke koleksi messages
	var history []models.Message                                   // Menyiapkan slice untuk menampung hasil query
	cursor.All(context.TODO(), &history)                           // Mengambil semua dokumen hasil query dan memasukkannya ke history
	c.JSON(http.StatusOK, history)                                 // Mengirimkan response JSON berisi riwayat chat
}

func GetContacts(c *gin.Context) { // Handler untuk mengambil daftar kontak yang pernah chat
	myID := c.Param("myID")                                         // Mendapatkan myID dari parameter URL
	filter := bson.M{"$or": []bson.M{{"from": myID}, {"to": myID}}} // Membuat filter untuk mencari semua pesan yang melibatkan user ini
	cursor, _ := config.MsgCollection.Find(context.TODO(), filter)  // Menjalankan query

	contactMap := make(map[string]bool)   // Membuat map untuk menyimpan ID unik kontak
	var messages []models.Message         // Slice untuk menampung pesan
	cursor.All(context.TODO(), &messages) // Mengambil semua pesan

	for _, m := range messages { // Iterasi semua pesan
		if m.From != myID {
			contactMap[m.From] = true
		} // Jika pengirim bukan saya, tambahkan pengirim ke kontak
		if m.To != myID {
			contactMap[m.To] = true
		} // Jika penerima bukan saya, tambahkan penerima ke kontak
	}

	var contacts []string // Slice untuk hasil akhir daftar kontak
	for k := range contactMap {
		contacts = append(contacts, k)
	} // Memindahkan key map ke slice contacts
	c.JSON(http.StatusOK, contacts) // Mengirimkan response JSON
}

func HandleChat(c *gin.Context) { // Handler untuk koneksi WebSocket chat
	myID := c.Param("myID")                             // Mendapatkan myID dari parameter URL
	ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil) // Meng-upgrade koneksi HTTP ke WebSocket
	defer ws.Close()                                    // Menutup koneksi WebSocket saat fungsi berakhir

	sub, _ := config.NC.Subscribe("user."+myID, func(m *nats.Msg) { // Subscribe ke subject NATS spesifik untuk user ini
		ws.WriteMessage(websocket.TextMessage, m.Data) // Mengirim pesan yang diterima dari NATS ke klien WebSocket
	})
	defer sub.Unsubscribe() // Unsubscribe dari NATS saat koneksi ditutup

	for { // Loop tak terbatas untuk membaca pesan dari WebSocket
		_, payload, err := ws.ReadMessage() // Membaca pesan dari WebSocket
		if err != nil {
			break
		} // Jika ada error (misal koneksi putus), keluar dari loop

		var msg models.Message        // Menyiapkan variabel untuk menampung pesan
		json.Unmarshal(payload, &msg) // Decode JSON payload ke struct Message
		msg.Timestamp = time.Now()    // Set timestamp pesan ke waktu sekarang

		config.MsgCollection.InsertOne(context.TODO(), msg) // Simpan pesan ke MongoDB
		data, _ := json.Marshal(msg)                        // Encode struct Message kembali ke JSON
		config.NC.Publish("user."+msg.To, data)             // Publish pesan ke NATS dengan subject penerima (agar diterima oleh subscriber penerima)
		ws.WriteMessage(websocket.TextMessage, data)        // Mengirim balik pesan ke pengirim (untuk konfirmasi/tampilan UI)
	}
}
