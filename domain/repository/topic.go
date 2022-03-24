package repository

import "news-dd/domain"

// NewsRepository represent repository of  the topic
// Expect implementation by the infrastructure layer
type TopicRepository interface {
	Get(id int) (*domain.Topic, error)
	GetAll() ([]domain.Topic, error)
	Save(*domain.Topic) error
	Update(*domain.Topic) error
	Remove(id int) error
}
