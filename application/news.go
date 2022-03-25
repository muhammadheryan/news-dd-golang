package application

import (
	"news-dd/config"
	"news-dd/domain"
	"news-dd/infrastructure/persistence"
)

// GetNews returns a tags by id
func GetNews(id int) (*domain.News, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	return repo.Get(id)
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
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	return repo.Remove(id)
}

// UpdateNews do remove news by id
func UpdateNews(p domain.News, id int) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := persistence.NewNewsRepositoryWithRDB(conn)
	p.Id = id

	return repo.Update(&p)
}
