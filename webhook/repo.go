package webhook

import (
	"gopkg.in/mgo.v2/bson"
)

// Repo is the interface of webhook mongo db
type Repo interface {
	Insert(wb Webhook) error
	Count(SenderID string) (int, error)
	Update(SenderID string, update bson.M) error
	Select(SenderID string, selectQ bson.M, wb *Webhook) error
}
