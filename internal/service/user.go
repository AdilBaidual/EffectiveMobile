package service

import (
	"EffectiveMobile/entity"
	"EffectiveMobile/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type UserService struct {
	repo repository.User
}

const (
	ageAPI         = "https://api.agify.io/?name="
	genderAPI      = "https://api.genderize.io/?name="
	nationalityAPI = "https://api.nationalize.io/?name="
)

type Age struct {
	Count *int    `json:"count"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
}

type Gender struct {
	Count       *int     `json:"count"`
	Name        *string  `json:"name"`
	Gender      *string  `json:"gender"`
	Probability *float32 `json:"probability"`
}

type Nationality struct {
	Count   *int      `json:"count"`
	Name    *string   `json:"name"`
	Country []Country `json:"country"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user entity.User) (int, error) {
	id, err := s.repo.Create(user)
	return id, err
}

func (s *UserService) EnrichUser(fio entity.Fio) (int, error) {
	var user entity.User
	user.Name = fio.Name
	user.Surname = fio.Surname
	user.Patronymic = fio.Patronymic

	age, err := getAge(fio.Name)
	if err != nil {
		logrus.Warn(err.Error())
	}
	user.Age = age

	gender, err := getGender(fio.Name)
	if err != nil {
		logrus.Warn(err.Error())
	}
	user.Gender = gender

	nationality, err := getNationality(fio.Name)
	if err != nil {
		logrus.Warn(err.Error())
	}
	user.Nationality = nationality

	fmt.Println("User after enrich: ", user)

	id, err := s.CreateUser(user)

	return id, err
}

func getAge(name string) (int, error) {
	var result int
	response, err := http.Get(ageAPI + name)
	defer response.Body.Close()

	if err != nil {
		return result, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	var age Age
	err = json.Unmarshal(body, &age)
	if err != nil {
		return result, err
	}

	if age.Age == nil {
		err = errors.New("User age not found")
		return result, err
	}
	result = *age.Age
	return result, nil
}

func getGender(name string) (string, error) {
	var result string
	response, err := http.Get(genderAPI + name)
	defer response.Body.Close()

	if err != nil {
		return result, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	var gender Gender
	err = json.Unmarshal(body, &gender)
	if err != nil {
		return result, err
	}

	if gender.Gender == nil {
		err = errors.New("User gender not found")
		return result, err
	}

	result = *gender.Gender
	return result, nil
}

func getNationality(name string) (string, error) {
	var result string
	response, err := http.Get(nationalityAPI + name)
	defer response.Body.Close()

	if err != nil {
		return result, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	var nationality Nationality
	err = json.Unmarshal(body, &nationality)
	if err != nil {
		return result, err
	}

	if len(nationality.Country) == 0 {
		err = errors.New("User nationality not found")
		return result, err
	}

	maxProbability := -1.0
	var maxCountryID string
	for _, country := range nationality.Country {
		if country.Probability > maxProbability {
			maxProbability = country.Probability
			maxCountryID = country.CountryId
		}
	}

	result = maxCountryID
	return result, nil
}

func (s *UserService) DeleteUser(userId int) error {
	err := s.repo.Delete(userId)
	return err
}

func (s *UserService) UpdateUser(userId int, updatedUser entity.User) error {
	user, err := s.repo.GetById(userId)
	if err != nil {
		return err
	}

	if updatedUser.Name == "" {
		updatedUser.Name = user.Name
	}

	if updatedUser.Surname == "" {
		updatedUser.Surname = user.Surname
	}

	if updatedUser.Patronymic == "" {
		updatedUser.Patronymic = user.Patronymic
	}

	if updatedUser.Age == 0 {
		updatedUser.Age = user.Age
	}

	if updatedUser.Gender == "" {
		updatedUser.Gender = user.Gender
	}

	if updatedUser.Nationality == "" {
		updatedUser.Nationality = user.Nationality
	}

	err = s.repo.Update(userId, updatedUser)
	return err
}

func (s *UserService) GetUserById(userId int) (entity.User, error) {
	user, err := s.repo.GetById(userId)
	return user, err
}

func (s *UserService) GetUsers(options entity.Options) ([]entity.User, error) {
	offset := (options.Page - 1) * options.PerPage
	query := fmt.Sprintf("SELECT * FROM users WHERE age >= %d AND age <= %d", options.MinAge, options.MaxAge)
	if options.Gender != "" {
		query += fmt.Sprintf("AND gender = %s", options.Gender)
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", options.PerPage, offset)

	users, err := s.repo.GetAll(query)

	return users, err
}
