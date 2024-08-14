package mongo_repository

import (
	"context"
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewSearchResultMongoRepository(db *mongo.Database) SearchResultMongoRepository {
	return SearchResultMongoRepository{
		db: db,
	}
}

type SearchResultMongoRepository struct {
	db *mongo.Database
}

func (r *SearchResultMongoRepository) Store(sr process_entity.SearchResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.db.Collection(SEARCH_RESULT_BUNDLE_COLLECTION)

	_, err := collection.InsertOne(ctx, sr)
	if err != nil {
		return err
	}

	return nil
}

func (r *SearchResultMongoRepository) StoreManyResultItems(srs []process_entity.SearchResultItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.db.Collection(SEARCH_RESULT_COLLECTION)

	var data []interface{} = make([]interface{}, len(srs))
	for i, v := range srs {
		data[i] = v
	}

	_, err := collection.InsertMany(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
