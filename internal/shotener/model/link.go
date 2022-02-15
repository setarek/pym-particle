package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Link struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id"`
	Shorten        string    `json:"shorten" bson:"shorten"`
	OriginalLink   string    `json:"original_link" bson:"original_link"`
	VisitorCount   int64     `json:"visitor_count" bson:"visitor_count"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastVisitedAt  time.Time `json:"last_visited_at" bson:"last_visited_at"`

}

func (l Link) ToBson() interface{} {
	link := bson.M{}
	link["original_link"] = l.OriginalLink
	link["shorten"] = l.Shorten
	link["visitor_count"] = 0
	link["created_at"] = time.Now()
	link["last_visited_at"] = time.Time{}
	return link
}