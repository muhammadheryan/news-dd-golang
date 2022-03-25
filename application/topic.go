package application

import (
	"encoding/json"
	"fmt"
	"news-dd/config"
	"news-dd/domain"
	"news-dd/infrastructure/helper"
	"news-dd/infrastructure/persistence"
	"time"
)

// GetTopic returns a topic by id
func GetTopic(id int) (*domain.Topic, error) {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("topic:%d", id)
	defer client.Close()

	resultRedis, err := client.Get(redisKey).Result()
	if err == nil {
		var result *domain.Topic
		json.Unmarshal([]byte(resultRedis), &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	result, err := repo.Get(id)

	redisValue := helper.StructToJson(result)
	client.Set(redisKey, redisValue, time.Hour*12)

	return result, err
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
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("topic:%d", id)
	defer client.Close()

	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	err = repo.Remove(id)
	if err == nil {
		client.Del(redisKey)
	}
	return err
}

// UpdateTopic do update topic by id
func UpdateTopic(p domain.Topic, id int) error {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("topic:%d", id)
	defer client.Close()

	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTopicRepositoryWithRDB(conn)
	p.Id = id
	err = repo.Update(&p)
	if err == nil {
		client.Del(redisKey)
	}

	return err
}
