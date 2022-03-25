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

// GetTags returns a tags by id
func GetTags(id int) (*domain.Tags, error) {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("tags:%d", id)
	defer client.Close()

	resultRedis, err := client.Get(redisKey).Result()
	if err == nil {
		var result *domain.Tags
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

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	result, err := repo.Get(id)

	redisValue := helper.StructToJson(result)
	client.Set(redisKey, redisValue, time.Hour*12)

	return result, err
}

// GetAllTags return all tags
func GetAllTags() ([]domain.Tags, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	return repo.GetAll()
}

// AddTags saves new tags
func AddTags(tag string) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	u := &domain.Tags{
		Tag:        tag,
		CreateTime: time.Now(),
	}
	return repo.Save(u)
}

// RemoveTags do remove tags by id
func RemoveTags(id int) error {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("tags:%d", id)
	defer client.Close()

	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	err = repo.Remove(id)
	if err == nil {
		client.Del(redisKey)
	}
	return err
}

// UpdateTags do update tags by id
func UpdateTags(p domain.Tags, id int) error {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("tags:%d", id)
	defer client.Close()

	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	err = repo.Update(&p)
	if err == nil {
		client.Del(redisKey)
	}

	return err
}
