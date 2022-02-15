package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/setarek/pym-particle-microservice/internal/shotener/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ShortenerRepository struct {
	db    *mongo.Database
}

var LinkCollectionName = "link"

func NewShortenerRepository(db *mongo.Database) *ShortenerRepository {
	return &ShortenerRepository{db: db}
}

func (r *ShortenerRepository) CreateLinkInfo(ctx context.Context, linkInfo model.Link)  error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "ShortenerRepository.CreateLinkInfo")
	defer span.Finish()

	doc := linkInfo.ToBson()
	if _, err := r.db.Collection(LinkCollectionName).InsertOne(ctx, doc); err != nil {
		return err
	}
	return nil
}

func (r *ShortenerRepository) UpdateVisitInfo(ctx context.Context, shorten string, visitCount int64)  error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "ShortenerRepository.CreateLinkInfo")
	defer span.Finish()

	filter := bson.M{
		"shorten": bson.M{
			"$eq":shorten,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"visited_at": time.Now(),
		},
		"$inc": bson.M{
			"visitor_count": visitCount,
		},
	}
	if _, err := r.db.Collection(LinkCollectionName).UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (r *ShortenerRepository) DeleteLinks(ctx context.Context, expireLimit time.Time) error {

	_, err := r.db.Collection(LinkCollectionName).DeleteMany(ctx,
		bson.M{"created_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(expireLimit),}},
			)

	if err != nil {
		return err
	}
	return nil

}
