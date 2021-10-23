package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"


	"github.com/EwRvp7LV7/48170360shop/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//приходящая извне информация для валидации и добавления товара на склад
type InputNewGoods struct {
	UserName  string `json:"user_name"`
	GoodsName string `json:"goods_name"`
	GoodsDesc  string `json:"goods_descrription"`
	GoodsAdd  int `json:"goods_add"`
	GoodsPrice  float32 `json:"goods_price"`
}


//прописывает роуты для запросов профиля пользователя.
func AddRouteInputManager(r *chi.Mux) {

	r.Route("/api/manager", func(r chi.Router) {

		r.Use(JWTSecurety) //здесь защита неавторизованного дальше не пустит

		r.Get("getbaskets", GetBasketsList) //отдает список товаров
		r.Post("newgoods", NewGoods)   //Добавить новые товары
		r.Post("addgoods2store", AddGoodsToStore)   //Изменить количество товара на складе
	})
}

func GetBasketsList(w http.ResponseWriter, r *http.Request) {

	_, claims, _ := jwtauth.FromContext(r.Context()) //можно не проверять err тк JWTSecurety

	jsonbytes, err := postgres.GetBasketsListtDB(claims["user_name"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbytes)
}

func NewGoods(w http.ResponseWriter, r *http.Request) {

	var data InputNewGoods
	var buf bytes.Buffer

	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &data)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = data.ValidateNewGoods()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context()) //можно не проверять err тк JWTSecurety
	data.UserName = claims["user_name"].(string)

	inputjsonbytes, _ := json.Marshal(data)
	jsonbytes, err := postgres.AddToBacketeDB(inputjsonbytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}


	//w.WriteHeader(http.StatusOK) //меняет application/json на text
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbytes)
}

func AddGoodsToStore(w http.ResponseWriter, r *http.Request) {

	var data InputNewGoods
	var buf bytes.Buffer

	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &data)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = data.ValidateGoodsAdd()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context()) //можно не проверять err тк JWTSecurety
	data.UserName = claims["user_name"].(string)

	inputjsonbytes, _ := json.Marshal(data)
	jsonbytes, err := postgres.AddGoodsToStoreDB(inputjsonbytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}


	//w.WriteHeader(http.StatusOK) //меняет application/json на text
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbytes)
}


//Validate валидация структуры.
func (up InputNewGoods) ValidateNewGoods() error {
	return validation.ValidateStruct(&up,
		validation.Field(&up.GoodsName, validation.Required, validation.Match(regexp.MustCompile("^[А-Яа-яЁё]{2,50}$"))),
		validation.Field(&up.GoodsDesc, validation.Required, validation.Match(regexp.MustCompile("^[А-Яа-яЁё]{2,150}$"))),
		validation.Field(&up.GoodsAdd, validation.Required, validation.Match(regexp.MustCompile(`^\d+$`))),
		validation.Field(&up.GoodsPrice, validation.Required, validation.Match(regexp.MustCompile(`^\d+\.\d\d$`))))
}

//Validate валидация структуры.
func (up InputNewGoods) ValidateGoodsAdd() error {
	return validation.ValidateStruct(&up,
		validation.Field(&up.GoodsName, validation.Required, validation.Match(regexp.MustCompile("^[А-Яа-яЁё]{2,50}$"))),
		validation.Field(&up.GoodsAdd, validation.Required, validation.Match(regexp.MustCompile(`^\d+$`))))
}

//auxiliary:заглушка для интерфейса
func (*InputNewGoods) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

//auxiliary:заглушка для интерфейса
func (*InputNewGoods) Bind(r *http.Request) error {
	return nil
}
