package models

import "time"

type Message struct {
	From      string    `json:"from" bson:"from"`
	To        string    `json:"to" bson:"to"`
	Content   string    `json:"content" bson:"content"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}