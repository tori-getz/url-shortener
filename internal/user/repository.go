package user

import "url-shortener/pkg/db"

type UserRepository struct {
	Db db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Db: *db,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Db.DB.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User

	result := repo.Db.First(&user, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
