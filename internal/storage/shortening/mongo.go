package shortening

import (
	"context"
	"fmt"
	"github.com/winterochek/go-linker/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type mgo struct {
	db *mongo.Database
}

func NewMongoDB(client *mongo.Client) *mgo {
	return &mgo{db: client.Database("linker")}
}

func (m *mgo) col() *mongo.Collection {
	return m.db.Collection("shortenings")
}

func (m *mgo) Put(ctx context.Context, shortening model.Shortening) (*model.Shortening, error) {
	const op = "shortening.mgo.Put"
	shortening.CreatedAt = time.Now().UTC()

	// check is there already a document with this identifier
	count, err := m.col().CountDocuments(ctx, bson.M{"_id": shortening.Identifier})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return nil, fmt.Errorf("%s: %w", op, model.ErrIdentifierExists)
	}

	// add in collection
	_, err = m.col().InsertOne(ctx, mgoShorteningFromModel(shortening))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shortening, nil
}

func (m *mgo) Get(ctx context.Context, shorteningID string) (*model.Shortening, error) {
	const op = "shortening.mgo.Get"

	var shortening mgoShortening
	if err := m.col().FindOne(ctx, bson.M{"_id": shorteningID}).Decode(&shortening); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", op, model.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return modelShorteningFromMgo(shortening), nil
}

func (m *mgo) IncrementVisits(ctx context.Context, shorteningID string) error {
	const op = "shortening.mgo.IncrementVisits"

	var (
		filter = bson.M{"_id": shorteningID}
		update = bson.M{
			"$inc": bson.M{"visits": 1},
			"$set": bson.M{"updated_at": time.Now().UTC()},
		}
	)

	_, err := m.col().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type mgoShortening struct {
	Identifier  string    `bson:"_id"`
	CreatedBy   string    `bson:"created_by"`
	OriginalURL string    `bson:"original_url"`
	Visits      int64     `bson:"visits"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func mgoShorteningFromModel(shortening model.Shortening) mgoShortening {
	return mgoShortening{
		Identifier:  shortening.Identifier,
		OriginalURL: shortening.OriginalURL,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
		CreatedBy:   shortening.CreatedBy,
	}
}

func modelShorteningFromMgo(shortening mgoShortening) *model.Shortening {
	return &model.Shortening{
		Identifier: shortening.Identifier,
		CreatedBy:   shortening.CreatedBy,
		OriginalURL: shortening.OriginalURL,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}
