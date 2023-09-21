package service

import (
	"EffectiveMobile/entity"
	"EffectiveMobile/internal/repository"
	mock_repository "EffectiveMobile/internal/repository/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockUser, user entity.User)

	testTable := []struct {
		name         string
		input        entity.User
		mockBehavior mockBehavior
		id           int
		err          error
	}{
		{
			name: "OK",
			input: entity.User{
				Name:    "Test1",
				Surname: "SurTest1",
			},
			mockBehavior: func(s *mock_repository.MockUser, user entity.User) {
				s.EXPECT().Create(user).Return(1, nil)
			},
			id:  1,
			err: nil,
		},
		{
			name: "ERROR",
			input: entity.User{
				Name: "Test1",
			},
			mockBehavior: func(s *mock_repository.MockUser, user entity.User) {
				s.EXPECT().Create(user).Return(0, errors.New("Field not defined"))
			},
			id:  0,
			err: errors.New("Some error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_repository.NewMockUser(c)
			testCase.mockBehavior(user, testCase.input)

			repo := &repository.Repository{User: user}
			service := NewService(repo)

			id, err := service.CreateUser(testCase.input)

			assert.Equal(t, testCase.id, id)
			assert.Equal(t, testCase.err, err)
		})
	}
}

func TestEnrichUser(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockUser, fio entity.Fio)

	testTable := []struct {
		name         string
		input        entity.Fio
		mockBehavior mockBehavior
		id           int
		err          error
	}{
		{
			name: "OK",
			input: entity.Fio{
				Name:       "John",
				Surname:    "Doe",
				Patronymic: "Smith",
			},
			mockBehavior: func(s *mock_repository.MockUser, fio entity.Fio) {
				user := entity.User{
					Name:        fio.Name,
					Surname:     fio.Surname,
					Patronymic:  fio.Patronymic,
					Age:         72,
					Gender:      "male",
					Nationality: "IE",
				}
				s.EXPECT().Create(user).Return(1, nil)
			},
			id:  1,
			err: nil,
		},
		{
			name: "ERROR",
			input: entity.Fio{
				Name: "John",
			},
			mockBehavior: func(s *mock_repository.MockUser, fio entity.Fio) {
				user := entity.User{
					Name:        fio.Name,
					Age:         72,
					Gender:      "male",
					Nationality: "IE",
				}
				s.EXPECT().Create(user).Return(0, errors.New("Field not defined"))
			},
			id:  0,
			err: errors.New("Field not defined"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_repository.NewMockUser(c)
			testCase.mockBehavior(user, testCase.input)

			repo := &repository.Repository{User: user}
			service := NewService(repo)

			id, err := service.EnrichUser(testCase.input)

			assert.Equal(t, testCase.id, id)
			assert.Equal(t, testCase.err, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockUser, id int)

	testTable := []struct {
		name         string
		input        int
		mockBehavior mockBehavior
		err          error
	}{
		{
			name:  "OK",
			input: 1,
			mockBehavior: func(s *mock_repository.MockUser, id int) {
				s.EXPECT().Delete(id).Return(nil)
			},
			err: nil,
		},
		{
			name:  "ERROR",
			input: 2,
			mockBehavior: func(s *mock_repository.MockUser, id int) {
				s.EXPECT().Delete(id).Return(errors.New("User not found"))
			},
			err: errors.New("User not found"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_repository.NewMockUser(c)
			testCase.mockBehavior(user, testCase.input)

			repo := &repository.Repository{User: user}
			service := NewService(repo)

			err := service.DeleteUser(testCase.input)

			assert.Equal(t, testCase.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockUser, userId int, updatedUser entity.User)

	testTable := []struct {
		name         string
		userId       int
		user         entity.User
		mockBehavior mockBehavior
		expectedErr  error
	}{
		{
			name:   "OK",
			userId: 1,
			user:   entity.User{Name: "UpdatedName"},
			mockBehavior: func(s *mock_repository.MockUser, userId int, updatedUser entity.User) {
				user := entity.User{
					Name:        "OriginalName",
					Surname:     "OriginalSurname",
					Age:         30,
					Gender:      "Male",
					Nationality: "American",
				}
				s.EXPECT().GetById(userId).Return(user, nil)
				s.EXPECT().Update(userId, gomock.Eq(entity.User{
					Name:        "UpdatedName",
					Surname:     "OriginalSurname",
					Age:         30,
					Gender:      "Male",
					Nationality: "American",
				})).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Error",
			userId: 2,
			user:   entity.User{Name: "UpdatedName"},
			mockBehavior: func(s *mock_repository.MockUser, userId int, updatedUser entity.User) {
				s.EXPECT().GetById(userId).Return(entity.User{}, errors.New("User not found"))
			},
			expectedErr: errors.New("User not found"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_repository.NewMockUser(c)
			testCase.mockBehavior(user, testCase.userId, testCase.user)

			repo := &repository.Repository{User: user}
			service := NewService(repo)

			err := service.UpdateUser(testCase.userId, testCase.user)

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestGetUserById(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockUser, userId int)

	testTable := []struct {
		name         string
		userId       int
		mockBehavior mockBehavior
		expectedUser entity.User
		expectedErr  error
	}{
		{
			name:   "OK",
			userId: 1,
			mockBehavior: func(s *mock_repository.MockUser, userId int) {
				// Мокируем вызов GetById и ожидаем, что будет возвращен пользователь
				expectedUser := entity.User{
					Name:        "John",
					Surname:     "Doe",
					Age:         30,
					Gender:      "Male",
					Nationality: "American",
				}
				s.EXPECT().GetById(userId).Return(expectedUser, nil)
			},
			expectedUser: entity.User{
				Name:        "John",
				Surname:     "Doe",
				Age:         30,
				Gender:      "Male",
				Nationality: "American",
			},
			expectedErr: nil,
		},
		{
			name:   "Error - User Not Found",
			userId: 2,
			mockBehavior: func(s *mock_repository.MockUser, userId int) {
				// Мокируем вызов GetById и ожидаем, что будет возвращена ошибка
				s.EXPECT().GetById(userId).Return(entity.User{}, errors.New("User not found"))
			},
			expectedUser: entity.User{}, // Ожидается пустой пользователь
			expectedErr:  errors.New("User not found"),
		},
		// Добавьте другие тестовые случаи по необходимости
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			User := mock_repository.NewMockUser(c)
			testCase.mockBehavior(User, testCase.userId)

			repo := &repository.Repository{User: User}
			service := NewService(repo)

			user, err := service.GetUserById(testCase.userId)

			assert.Equal(t, testCase.expectedUser, user)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestGetUsers(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockUser, options entity.Options)

	testTable := []struct {
		name          string
		options       entity.Options
		mockBehavior  mockBehavior
		expectedUsers []entity.User
		expectedErr   error
	}{
		{
			name: "OK",
			options: entity.Options{
				Page:    1,
				PerPage: 10,
				MinAge:  20,
				MaxAge:  30,
				Gender:  "Male",
			},

			mockBehavior: func(s *mock_repository.MockUser, options entity.Options) {
				// Мокируем вызов GetAll и ожидаем, что будет возвращен список пользователей
				expectedUsers := []entity.User{
					{Id: 1, Name: "John", Surname: "Doe", Age: 25, Gender: "Male"},
					{Id: 2, Name: "Jane", Surname: "Smith", Age: 28, Gender: "Male"},
				}
				query := "SELECT * FROM users WHERE age >= 20 AND age <= 30 AND gender = Male LIMIT 10 OFFSET 0"
				s.EXPECT().GetAll(query).Return(expectedUsers, nil)
			},
			expectedUsers: []entity.User{
				{Id: 1, Name: "John", Surname: "Doe", Age: 25, Gender: "Male"},
				{Id: 2, Name: "Jane", Surname: "Smith", Age: 28, Gender: "Male"},
			},
			expectedErr: nil,
		},
		{
			name: "Error - Database Error",
			options: entity.Options{
				Page:    1,
				PerPage: 10,
				MinAge:  20,
				MaxAge:  30,
				Gender:  "Female",
			},
			mockBehavior: func(s *mock_repository.MockUser, options entity.Options) {
				query := "SELECT * FROM users WHERE age >= 20 AND age <= 30 AND gender = Female LIMIT 10 OFFSET 0"
				s.EXPECT().GetAll(query).Return(nil, errors.New("Database error"))
			},
			expectedUsers: nil,
			expectedErr:   errors.New("Database error"),
		},
		// Добавьте другие тестовые случаи по необходимости
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			User := mock_repository.NewMockUser(c)
			testCase.mockBehavior(User, testCase.options)

			repo := &repository.Repository{User: User}
			service := NewService(repo)

			users, err := service.GetUsers(testCase.options)

			assert.Equal(t, testCase.expectedUsers, users)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
