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

func (c webhookRepo) insert(wb *webhook.Webhook) error {
	return c.insert(wb)
}

func (c webhookRepo) Count(senderID string) (int, error) {
	count, err := c.collection.Find(bson.M{"senderId": senderID}).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c webhookRepo) Update(senderID string, update bson.M) error {
	return c.collection.Update(bson.M{"senderId": senderID}, update)
}
