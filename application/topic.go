package application

import (
	"news-dd/config"
	"news-dd/domain"
	"news-dd/infrastructure/persistence"
	"time"
)

// GetTopic returns a topic by id
func GetTopic(id int) (*domain.Topic, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	return repo.Get(id)
}

// GetAllTopic return all topics
func GetAllTopic() ([]domain.Topic, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	return repo.GetAll()
}

// AddTopic saves new topic
func AddTopic(topic string) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	u := &domain.Topic{
		Topic:      topic,
		CreateTime: time.Now(),
	}
	return repo.Save(u)
}

// RemoveTopic do remove topic by id
func RemoveTopic(id int) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	return repo.Remove(id)
}

// UpdateTopic do update topic by id
func UpdateTopic(p domain.Topic, id int) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	p.Id = id

	return repo.Update(&p)
}
