package mongo_repository

import (
	"context"
	"time"

	"github.com/charmingruby/serpright/internal/common/infra/database/mongo_collection"
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

	collection := r.db.Collection(mongo_collection.SEARCH_RESULT_COLLECTION)

	_, err := collection.InsertOne(ctx, sr)
	if err != nil {
		return err
	}

	return nil
}
