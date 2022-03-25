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

// GetNews returns a tags by id
func GetNews(id int) (*domain.News, error) {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("news:%d", id)
	defer client.Close()

	resultRedis, err := client.Get(redisKey).Result()
	if err == nil {
		var result *domain.News
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

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	result, err := repo.Get(id)

	redisValue := helper.StructToJson(result)
	client.Set(redisKey, redisValue, time.Hour*12)

	return result, err
}

// GetAllTags return all news
func GetAllNews(topicId int, status string) ([]domain.News, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	return repo.GetAll(topicId, status)
}

// AddNews saves new news
func AddNews(p domain.News) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	return repo.Save(&p)
}

// RemoveNews do remove news by id
func RemoveNews(id int) error {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("news:%d", id)
	defer client.Close()

	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	err = repo.Remove(id)
	if err == nil {
		client.Del(redisKey)
	}
	return err
}

// UpdateNews do remove news by id
func UpdateNews(p domain.News, id int) error {
	client := config.ConnectRedis()
	redisKey := fmt.Sprintf("news:%d", id)
	defer client.Close()

	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	p.Id = id
	err = repo.Update(&p)
	if err == nil {
		client.Del(redisKey)
	}

	return err
}
