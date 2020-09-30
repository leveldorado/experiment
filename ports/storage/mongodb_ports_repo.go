package storage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/leveldorado/experiment/grpc/portspb"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ExperimentDatabaseName = "experiment"
	PortCollectionName     = "port"
)

/*
port is duplicate of portpg.Port customized for mongodb use
at the moment that's bson tag field, but there may be done more customization non-related to business logic,
e.g. coordinates maybe  required to be  use in geo search and for that purpose they had to be placed in mongodb doc in specific format
*/
type port struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	City        string    `bson:"city"`
	Country     string    `bson:"country"`
	Alias       []string  `bson:"alias"`
	Regions     []string  `bson:"regions"`
	Coordinates []float32 `bson:"coordinates"`
	Province    string    `bson:"province"`
	Timezone    string    `bson:"timezone"`
	Unlocks     []string  `bson:"unlocks"`
}

func fromPortPb(p *portspb.Port) port {
	return port{
		ID:          p.Id,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocks:     p.Unlocks,
	}
}

func (p port) toPortPb() *portspb.Port {
	return &portspb.Port{
		Id:          p.ID,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocks:     p.Unlocks,
	}
}

type MongodbPortsRepo struct {
	c *mongo.Collection
}

func NewMongodbPortsRepo(cl *mongo.Client, databaseName, collectionName string) *MongodbPortsRepo {
	return &MongodbPortsRepo{c: cl.Database(databaseName).Collection(collectionName)}
}

func (r *MongodbPortsRepo) Get(ctx context.Context, id string) (*portspb.Port, error) {
	p := port{}
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf(`failed to get port by id %s: [err: %w]`, id, err)
	}
	return p.toPortPb(), nil
}

func (r *MongodbPortsRepo) List(ctx context.Context) (<-chan *portspb.Port, error) {
	portsChan := make(chan *portspb.Port)
	cursor, err := r.c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	go func() {
		for cursor.Next(ctx) {
			p := port{}
			if err := cursor.Decode(&p); err != nil {
				/*
						logging is dependence.
						and should log level (info, warn, error e.t.c) and additional information like correlation id
					    proper logging skipped for now
				*/
				log.Println(`failed to decode port document`, err)
				close(portsChan)
				return
			}
			select {
			case <-ctx.Done():
				close(portsChan)
				return
			case portsChan <- p.toPortPb():
			}
		}
		close(portsChan)
	}()
	return portsChan, nil
}

func (r *MongodbPortsRepo) Save(ctx context.Context, p *portspb.Port) error {
	opts := options.Update().SetUpsert(true)
	f := bson.M{"_id": p.Id}
	u := bson.M{"$set": fromPortPb(p)}
	if _, err := r.c.UpdateOne(ctx, f, u, opts); err != nil {
		return fmt.Errorf(`failed to upsert port record: [opts: %+v, filter: %+v, update: %+v, error: %w]`, opts, f, u, err)
	}
	return nil
}
