package mongo

import (
	"github.com/sch00lb0y/StockiumBot/webhook"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type webhookRepo struct {
	collection *mgo.Collection
	webhook.Repo
}

func NewRepo(c *mgo.Collection) webhook.Repo {
	return webhookRepo{collection: c}
}

func (c webhookRepo) Insert(wb webhook.Webhook) error {
	return c.collection.Insert(wb)
}

func (c webhookRepo) Count(senderID string) (int, error) {
	count, err := c.collection.Find(bson.M{"senderId": senderID}).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c webhookRepo) Update(senderID string, update bson.M) error {
	_, err := c.collection.Upsert(bson.M{"senderId": senderID}, update)
	return err
}
