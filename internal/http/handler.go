package http

import (
	"EffectiveMobile/entity"
	"EffectiveMobile/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoute() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.DELETE("/deleteUser/:id", h.deleteUserHandler)
		api.POST("/updateUser/:id", h.updateUserHandler)
		api.POST("/createUser", h.createUserHandler)
		api.GET("/getUsers", h.getUsersHandler)
		api.GET("/getUser/:id", h.getUserByIdHandler)
	}

	return router
}

func (h *Handler) deleteUserHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		logrus.Error(err.Error())
		return
	}

	err = h.services.User.DeleteUser(userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		logrus.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

type UpdateUser struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

func (h *Handler) updateUserHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		logrus.Error(err.Error())
		return
	}
	var tmpUser UpdateUser
	if err = c.BindJSON(&tmpUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		logrus.Error(err.Error())
		return
	}

	user := entity.User{
		Name:        tmpUser.Name,
		Surname:     tmpUser.Surname,
		Patronymic:  tmpUser.Patronymic,
		Age:         tmpUser.Age,
		Gender:      tmpUser.Gender,
		Nationality: tmpUser.Nationality,
	}

	err = h.services.User.UpdateUser(userId, user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		logrus.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) createUserHandler(c *gin.Context) {
	var tmpUser entity.User
	if err := c.BindJSON(&tmpUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		logrus.Error(err.Error())
		return
	}

	id, err := h.services.User.CreateUser(tmpUser)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		logrus.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, id)
}

func (h *Handler) getUsersHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	minAgeStr := c.DefaultQuery("min_age", "0")
	maxAgeStr := c.DefaultQuery("max_age", "100")
	ageStr := c.DefaultQuery("age", "0")
	gender := c.DefaultQuery("gender", "")

	minAge, err := strconv.Atoi(minAgeStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if age != 0 {
		minAge = age
		maxAge = age
	}

	options := entity.Options{
		Page:    page,
		PerPage: perPage,
		MaxAge:  maxAge,
		MinAge:  minAge,
		Age:     age,
		Gender:  gender,
	}

	users, err := h.services.GetUsers(options)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		logrus.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUserByIdHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		logrus.Error(err.Error())
		return
	}

	user, err := h.services.User.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		logrus.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
