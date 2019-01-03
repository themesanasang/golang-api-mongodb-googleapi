package handler

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/themesanasang/testdb/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (h *Handler) UserAll(c echo.Context) (err error) {
	user := []*model.User{}
	db := h.DB.Clone()

	//ดึงข้อมูลทั้งหมด
	if err = db.DB("pokemondb").C("user").
		Find(nil).
		All(&user); err != nil {
		return
	}
	defer db.Close()

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c echo.Context) (err error) {

	u := &model.User{ID: bson.NewObjectId()}
	if err = c.Bind(u); err != nil {
		return
	}

	db := h.DB.Clone()
	defer db.Close()

	//t := time.Now()
	//format => "created_at": "2018-12-31T15:13:40.692895876+07:00"

	// Save post in database
	if err = db.DB("pokemondb").C("user").Insert(u); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) Login(c echo.Context) (err error) {
	// Bind
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	// Find user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("pokemondb").C("user").
		Find(bson.M{"username": u.Username, "password": u.Password}).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid username or password"}
		}
		return
	}

	//-----
	// JWT
	//-----

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}

func userIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}
