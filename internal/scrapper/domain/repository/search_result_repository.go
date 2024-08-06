package repository

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"

type SearchResultRepository interface {
	Store(sr process_entity.SearchResult) error
}
