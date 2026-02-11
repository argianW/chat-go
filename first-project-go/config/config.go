package config // Mendefinisikan package config yang berisi konfigurasi aplikasi

import ( // Melakukan import library yang dibutuhkan
	"context" // Mengimport package context untuk mengatur deadline, pembatalan sinyal, dan nilai-nilai request-scoped
	"github.com/nats-io/nats.go" // Mengimport driver client NATS untuk Go
	"go.mongodb.org/mongo-driver/mongo" // Mengimport driver MongoDB untuk Go
	"go.mongodb.org/mongo-driver/mongo/options" // Mengimport options untuk konfigurasi client MongoDB
)

var ( // Mendeklarasikan variabel global untuk koneksi database
	NC            *nats.Conn // Variabel global untuk menyimpan koneksi ke NATS server
	MsgCollection *mongo.Collection // Variabel global untuk menyimpan referensi ke koleksi MongoDB 'messages'
)

func Init() { // Fungsi Init untuk menginisialisasi koneksi ke database dan message broker
	// Koneksi NATS
	NC, _ = nats.Connect(nats.DefaultURL) // Membuat koneksi ke server NATS menggunakan URL default

	// Koneksi MongoDB
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017")) // Membuat client MongoDB dan melakukan koneksi ke server MongoDB lokal
	MsgCollection = client.Database("whatsapp_db").Collection("messages") // Mengarahkan ke database 'whatsapp_db' dan koleksi 'messages'
}