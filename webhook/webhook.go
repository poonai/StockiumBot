package webhook

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Webhook is the struct for facebook users
type Webhook struct {
	SenderID  string        `bson:"senderId" json:"senderId"`
	Portfolio []string      `bson:"portfolio" json:"portolio"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	ID        bson.ObjectId `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
}
