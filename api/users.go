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

//приходящая извне информация для операций с корзиной
type InputUserBasket struct {
	UserName  string `json:"user_name"`
	GoodsName string `json:"goods_name"`
	GoodsAdd  string `json:"goods_add"`
}

//прописывает роуты для запросов профиля пользователя.
func AddRouteInputUserBasket(r *chi.Mux) {

	r.Route("/api/user", func(r chi.Router) {

		r.Use(JWTSecurety) //здесь защита неавторизованного дальше не пустит

		r.Get("getgoodslist", GetGoodsList) //отдает список товаров
		//r.Get("buy", AddToBacket)   //склад уменьшается на размер корзины, корзина удаляется
		r.Post("add2basket", AddToBacket) //операции с корзиной
	})
}

func GetGoodsList(w http.ResponseWriter, r *http.Request) {

	jsonbytes, err := postgres.GetGoodsListDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbytes)
}

func AddToBacket(w http.ResponseWriter, r *http.Request) {

	var data InputUserBasket
	var buf bytes.Buffer

	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &data)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = data.Validate()
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

	//зачем то кодирует в Base64UTF
	//Array and slice values encode as JSON arrays, except that []byte encodes as a base64-encoded string, and a nil slice encodes as the null JSON value.
	// https://pkg.go.dev/encoding/json#Marshal
	//render.DefaultResponder(w, r, jsonbytes)

	// w.WriteHeader(http.StatusOK)
	// _, err = w.Write([]byte(jsonbytes))

	// if err != nil {
	// 	return
	// }

	//w.WriteHeader(http.StatusOK) //меняет application/json на text
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbytes)
}

//Validate валидация структуры.
func (up InputUserBasket) Validate() error {
	return validation.ValidateStruct(&up,
		validation.Field(&up.GoodsName, validation.Required, validation.Match(regexp.MustCompile("^[А-Яа-яЁё]{2,50}$"))),
		validation.Field(&up.GoodsAdd, validation.Required, validation.Match(regexp.MustCompile(`^[\-+]?\d+$`))))
}

//auxiliary:заглушка для интерфейса
func (*InputUserBasket) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

//auxiliary:заглушка для интерфейса
func (*InputUserBasket) Bind(r *http.Request) error {
	return nil
}
