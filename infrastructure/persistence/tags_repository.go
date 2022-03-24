package persistence

import (
	"github.com/jinzhu/gorm"
	"news-dd/domain"
	"news-dd/domain/repository"
	"time"
)

type TagsRepositoryImpl struct {
	Conn *gorm.DB
}

func NewTagsRepositoryWithRDB(conn *gorm.DB) repository.TagsRepository {
	return &TagsRepositoryImpl{Conn: conn}
}

// Get Tags by id return domain.tags
func (r *TagsRepositoryImpl) Get(id int) (*domain.Tags, error) {
	tags := &domain.Tags{}
	if err := r.Conn.Raw("SELECT * FROM tags WHERE id = ?", id).Scan(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetAll tags return all domain.tags
func (r *TagsRepositoryImpl) GetAll() ([]domain.Tags, error) {
	tags := []domain.Tags{}
	if err := r.Conn.Raw("SELECT * FROM tags").Scan(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

// Save to add tags
func (r *TagsRepositoryImpl) Save(tags *domain.Tags) error {
	if err := r.Conn.Save(&tags).Error; err != nil {
		return err
	}

	return nil
}

// Remove delete tags
func (r *TagsRepositoryImpl) Remove(id int) error {
	tags := &domain.Tags{}
	if err := r.Conn.First(&tags, id).Error; err != nil {
		return err
	}

	if err := r.Conn.Delete(&tags).Error; err != nil {
		return err
	}

	return nil
}

// Update data tags
func (r *TagsRepositoryImpl) Update(tags *domain.Tags) error {
	if err := r.Conn.Model(&tags).UpdateColumns(domain.Tags{Tag: tags.Tag, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}
