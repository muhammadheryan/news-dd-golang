package persistence

import (
	"github.com/jinzhu/gorm"
	"news-dd/domain"
	"news-dd/domain/repository"
	"time"
)

type TopicRepositoryImpl struct {
	Conn *gorm.DB
}

func NewTopicRepositoryWithRDB(conn *gorm.DB) repository.TopicRepository {
	return &TopicRepositoryImpl{Conn: conn}
}

// Get topic by id return domain.topic
func (r *TopicRepositoryImpl) Get(id int) (*domain.Topic, error) {
	topic := &domain.Topic{}
	if err := r.Conn.Raw("SELECT * FROM topic WHERE id = ?", id).Scan(&topic).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

// GetAll topic return all domain.topic
func (r *TopicRepositoryImpl) GetAll() ([]domain.Topic, error) {
	topics := []domain.Topic{}
	if err := r.Conn.Raw("SELECT * FROM topic").Scan(&topics).Error; err != nil {
		return nil, err
	}

	return topics, nil
}

// Save to add topic
func (r *TopicRepositoryImpl) Save(topic *domain.Topic) error {
	if err := r.Conn.Save(&topic).Error; err != nil {
		return err
	}

	return nil
}

// Remove delete topic
func (r *TopicRepositoryImpl) Remove(id int) error {
	topic := &domain.Topic{}
	if err := r.Conn.First(&topic, id).Error; err != nil {
		return err
	}

	if err := r.Conn.Delete(&topic).Error; err != nil {
		return err
	}

	return nil
}

// Update data topic
func (r *TopicRepositoryImpl) Update(topic *domain.Topic) error {
	if err := r.Conn.Model(&topic).UpdateColumns(domain.Topic{Topic: topic.Topic, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}

	return nil
}
