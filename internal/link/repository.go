package link

import (
	"url-shortener/pkg/db"

	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	Db db.Db
}

func NewLinkRepository(db *db.Db) *LinkRepository {
	return &LinkRepository{
		Db: *db,
	}
}

func (repo *LinkRepository) GetCount() int64 {
	var count int64

	repo.Db.DB.
		Table("links").
		Count(&count)

	return count
}

func (repo *LinkRepository) GetLinks(limit int, offset int) []Link {
	var links []Link

	repo.Db.DB.
		Table("links").
		Where("deleted_at is null").
		Order("id asc").
		Limit(limit).
		Offset(offset).
		Scan(&links)

	return links
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Db.DB.Create(link)

	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil
}

func (repo *LinkRepository) FindByHash(hash string) (*Link, error) {
	var link Link

	result := repo.Db.First(&link, "hash = ?", hash)

	if result.Error != nil {
		return nil, result.Error
	}

	return &link, nil
}

func (repo *LinkRepository) FindById(id uint) (*Link, error) {
	var link Link

	result := repo.Db.First(&link, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.Db.Clauses(clause.Returning{}).Updates(link)

	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.Db.Delete(&Link{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
