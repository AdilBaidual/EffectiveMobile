package entity

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required" db:"username"`
	Surname     string `json:"surname" binding:"required"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
