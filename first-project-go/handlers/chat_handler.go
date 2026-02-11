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

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func GetHistory(c *gin.Context) {
	myID := c.Param("myID")
	targetID := c.Param("targetID")
	filter := bson.M{"$or": []bson.M{{"from": myID, "to": targetID}, {"from": targetID, "to": myID}}}
	
	cursor, _ := config.MsgCollection.Find(context.TODO(), filter)
	var history []models.Message
	cursor.All(context.TODO(), &history)
	c.JSON(http.StatusOK, history)
}

func GetContacts(c *gin.Context) {
	myID := c.Param("myID")
	filter := bson.M{"$or": []bson.M{{"from": myID}, {"to": myID}}}
	cursor, _ := config.MsgCollection.Find(context.TODO(), filter)
	
	contactMap := make(map[string]bool)
	var messages []models.Message
	cursor.All(context.TODO(), &messages)

	for _, m := range messages {
		if m.From != myID { contactMap[m.From] = true }
		if m.To != myID { contactMap[m.To] = true }
	}

	var contacts []string
	for k := range contactMap { contacts = append(contacts, k) }
	c.JSON(http.StatusOK, contacts)
}

func HandleChat(c *gin.Context) {
	myID := c.Param("myID")
	ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	defer ws.Close()

	sub, _ := config.NC.Subscribe("user."+myID, func(m *nats.Msg) {
		ws.WriteMessage(websocket.TextMessage, m.Data)
	})
	defer sub.Unsubscribe()

	for {
		_, payload, err := ws.ReadMessage()
		if err != nil { break }

		var msg models.Message
		json.Unmarshal(payload, &msg)
		msg.Timestamp = time.Now()

		config.MsgCollection.InsertOne(context.TODO(), msg)
		data, _ := json.Marshal(msg)
		config.NC.Publish("user."+msg.To, data)
		ws.WriteMessage(websocket.TextMessage, data)
	}
}