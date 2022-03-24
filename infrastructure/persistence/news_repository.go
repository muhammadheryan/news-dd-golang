package persistence

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"news-dd/domain"
	"news-dd/domain/repository"
	"time"
)

type NewsRepositoryImpl struct {
	Conn *gorm.DB
}

func NewNewsRepositoryWithRDB(conn *gorm.DB) repository.NewsRepository {
	return &NewsRepositoryImpl{Conn: conn}
}

// Get News by id return domain.News
func (r *NewsRepositoryImpl) Get(id int) (*domain.News, error) {
	news := &domain.News{}
	if err := r.Conn.Raw("SELECT * FROM news WHERE id = ?", id).Scan(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

// GetAll News return all domain.News
func (r *NewsRepositoryImpl) GetAll(topicId int, status string) ([]domain.News, error) {
	whereCond := "WHERE 1 = 1"

	if topicId > 0 {
		whereCond = fmt.Sprintf("%s AND topic = %d", whereCond, topicId)
	}

	if status != "" {
		whereCond = fmt.Sprintf(`%s AND status = "%s"`, whereCond, status)
	}

	sql := fmt.Sprintf("SELECT * FROM news %s", whereCond)

	news := []domain.News{}
	if err := r.Conn.Raw(sql).Scan(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

// Save to add News
func (r *NewsRepositoryImpl) Save(news *domain.News) error {
	if err := r.Conn.Save(&news).Error; err != nil {
		return err
	}

	return nil
}

// Remove delete News
func (r *NewsRepositoryImpl) Remove(id int) error {
	news := &domain.News{}
	if err := r.Conn.First(&news, id).Error; err != nil {
		return err
	}

	if err := r.Conn.Delete(&news).Error; err != nil {
		return err
	}

	return nil
}

// Update data News
func (r *NewsRepositoryImpl) Update(news *domain.News) error {
	if err := r.Conn.Model(&news).UpdateColumns(domain.News{Title: news.Title, Content: news.Content, Status: news.Status, Topic: news.Topic, Tags: news.Tags, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}

	return nil
}
