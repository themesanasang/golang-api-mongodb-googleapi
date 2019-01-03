package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/themesanasang/testdb/model"
	"gopkg.in/mgo.v2/bson"
)

func (h *Handler) PokemonsAll(c echo.Context) (err error) {
	pokemons := []*model.Pokemons{}
	db := h.DB.Clone()

	//ดึงข้อมูลทั้งหมด
	if err = db.DB("pokemondb").C("pokemons").
		Find(nil).
		All(&pokemons); err != nil {
		return
	}
	defer db.Close()

	return c.JSON(http.StatusOK, pokemons)
}

func (h *Handler) PokemonsFindId(c echo.Context) (err error) {
	id := bson.ObjectIdHex(c.Param("id"))
	pokemon := model.Pokemons{}
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("pokemondb").C("pokemons").
		Find(bson.M{"_id": id}).
		One(&pokemon); err != nil {
		return
	}

	return c.JSON(http.StatusOK, pokemon)
}

func (h *Handler) PokemonsFindName(c echo.Context) (err error) {
	name := c.Param("name")
	pokemon := model.Pokemons{}
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("pokemondb").C("pokemons").
		Find(bson.M{"name": name}).
		One(&pokemon); err != nil {
		return
	}

	return c.JSON(http.StatusOK, pokemon)
}

func (h *Handler) CreatePokemon(c echo.Context) (err error) {
	p := new(model.Pokemons)
	if err = c.Bind(p); err != nil {
		return
	}

	db := h.DB.Clone()
	defer db.Close()

	// Save post in database
	if err = db.DB("pokemondb").C("pokemons").Insert(p); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, p)
}

func (h *Handler) PokemonsUpdate(c echo.Context) (err error) {
	id := bson.ObjectIdHex(c.Param("id"))

	p := new(model.Pokemons)
	if err := c.Bind(p); err != nil {
		return err
	}

	db := h.DB.Clone()
	defer db.Close()

	if err = db.DB("pokemondb").C("pokemons").
		UpdateId(id, bson.M{"$set": bson.M{"element": p.Element}}); err != nil {
		return
	}

	return c.JSON(http.StatusOK, p)
}

func (h *Handler) PokemonsDelete(c echo.Context) (err error) {
	id := bson.ObjectIdHex(c.Param("id"))

	db := h.DB.Clone()
	defer db.Close()

	if err = db.DB("pokemondb").C("pokemons").
		Remove(bson.M{"_id": id}); err != nil {
		return
	}

	return c.JSON(http.StatusOK, id)
}
