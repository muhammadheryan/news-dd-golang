package repository

import "news-dd/domain"

// NewsRepository represent repository of  the tags
// Expect implementation by the infrastructure layer
type TagsRepository interface {
	Get(id int) (*domain.Tags, error)
	GetAll() ([]domain.Tags, error)
	Save(*domain.Tags) error
	Update(*domain.Tags) error
	Remove(id int) error
}
