package config

import (
	"context"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	NC            *nats.Conn
	MsgCollection *mongo.Collection
)

func Init() {
	// Koneksi NATS
	NC, _ = nats.Connect(nats.DefaultURL)

	// Koneksi MongoDB
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	MsgCollection = client.Database("whatsapp_db").Collection("messages")
}