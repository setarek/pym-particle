package mongo

import (
	"context"
	"fmt"
	"github.com/setarek/pym-particle-microservice/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

var linkCollectionName = "link"

type MongoDB struct {
	host string
	port string
	db   string
	user string
	pass string

	conn *mongo.Client
}

func newMongoDB(host, port, db, user, pass string) *MongoDB {
	return &MongoDB{
		host: host,
		port: port,
		db:   db,
		user: user,
		pass: pass,
	}
}

func (m *MongoDB) ping(dbName string) {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				pingContext, can := context.WithTimeout(context.Background(), 5*time.Second)
				defer can()
				err := m.conn.Ping(pingContext, readpref.Primary())
				if err != nil {

					time.Sleep(5 * time.Second)
					go m.connect(dbName)
					<-quit
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *MongoDB) connect(dbName string) *mongo.Database{
	var err error

MONGOTRY:

	ctx, can := context.WithTimeout(context.Background(), 30*time.Second)
	defer can()
	m.conn, err = mongo.NewClient(
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", m.host, m.port) + "/?direct=true"),
		options.Client().SetReadPreference(readpref.Primary()),
		options.Client().SetDirect(true),
	)

	if err != nil {
		goto MONGOTRY
	}


	err = m.conn.Connect(ctx)
	if err != nil {
		goto MONGOTRY
	}
	m.ping(dbName)


	return m.conn.Database(dbName)
}

func InitDB(config *config.Config) (*mongo.Database, error){

	Client := newMongoDB(
		config.GetString("db_host"),
		config.GetString("db_port"),
		config.GetString("db_name"),
		config.GetString("db_user"),
		config.GetString("db_password"),
	)

	ctx, can := context.WithTimeout(context.Background(), 10*time.Second)
	defer can()

	Client.ping(config.GetString("db_name"))

	mdb := Client.connect(config.GetString("db_name"))

	if err := mdb.RunCommand(
		ctx,
		bsonx.Doc{
			{Key: "create", Value: bsonx.String(linkCollectionName)},
		},
	).Err(); err == nil {
		mdb.Collection(linkCollectionName).Indexes().CreateMany(
			ctx,
			[]mongo.IndexModel{
				mongo.IndexModel{
					Keys: bsonx.Doc{
						{Key: "shorten", Value: bsonx.String("text")},
					},
					Options: options.Index().SetUnique(true),
				},
			},
		)

	}

	return Client.connect(config.GetString("db_name")), nil
}