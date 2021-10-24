package api

import (
	"testing"
)

var defaultCorrectInputNewGoods = InputNewGoods{
	UserName:   "user1",
	GoodsName:  "морковь",
	GoodsDesc:  "свежая",
	GoodsAdd:   "-3",
	GoodsPrice: "120.47",
}

func TestInputUserInvalidFields(t *testing.T) {

	var err error

	input := InputNewGoods{}
	err = input.ValidateNewGoods()
	if err == nil {
		t.Error("Пустая структура прошла валидацию")
	}
	errstr := "Некорректное поле прошло валидацию"

	input = defaultCorrectInputNewGoods
	input.GoodsName = "Name Name"
	err = input.ValidateNewGoods()
	if err == nil {
		t.Error(errstr)
	}

	input = defaultCorrectInputNewGoods
	input.GoodsDesc = "+ggggg+ 57"
	err = input.ValidateNewGoods()
	if err == nil {
		t.Error(errstr)
	}

}
