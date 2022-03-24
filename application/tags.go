package application

import (
	"news-dd/config"
	"news-dd/domain"
	"news-dd/infrastructure/persistence"
	"time"
)

// GetTags returns a tags by id
func GetTags(id int) (*domain.Tags, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	return repo.Get(id)
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
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	return repo.Remove(id)
}

// UpdateTags do update tags by id
func UpdateTags(p domain.Tags, id int) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewTagsRepositoryWithRDB(conn)
	p.Id = id

	return repo.Update(&p)
}
