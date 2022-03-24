package repository

import "news-dd/domain"

// NewsRepository represent repository of  the news
// Expect implementation by the infrastructure layer
type NewsRepository interface {
	Get(id int) (*domain.News, error)
	GetAll(topicId int, status string) ([]domain.News, error)
	Save(*domain.News) error
	Update(*domain.News) error
	Remove(id int) error
}
