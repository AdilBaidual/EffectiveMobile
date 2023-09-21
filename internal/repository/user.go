package repository

import (
	"EffectiveMobile/entity"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db    *sqlx.DB
	cache *RedisRepository
}

func NewUserPostgres(db *sqlx.DB, cache *RedisRepository) *UserPostgres {
	return &UserPostgres{db: db, cache: cache}
}

func (r *UserPostgres) Create(user entity.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, surname, patronymic, age, gender, nationality)"+
		" values ($1, $2, $3, $4, $5, $6) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nationality)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserPostgres) Delete(userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", usersTable)
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *UserPostgres) Update(userId int, user entity.User) error {
	query := fmt.Sprintf("UPDATE %s SET username=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE id=$7", usersTable)
	_, err := r.db.Exec(query, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nationality, userId)
	return err
}

func (r *UserPostgres) GetById(userId int) (entity.User, error) {
	var user entity.User

	cacheKey := fmt.Sprintf("user_%d", userId)
	cachedData, err := r.cache.Get(cacheKey)
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &user); err == nil {
			fmt.Println("THIS DATA FROM REDIS")
			return user, nil
		}
	}

	query := fmt.Sprintf("SELECT id, username, surname, patronymic, age, gender, nationality FROM %s WHERE id=$1", usersTable)
	err = r.db.Get(&user, query, userId)

	if err == nil {
		userData, _ := json.Marshal(user)
		r.cache.Set(cacheKey, string(userData))
	}

	return user, err
}

func (r *UserPostgres) GetAll(query string) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Select(&users, query)
	return users, err
}
