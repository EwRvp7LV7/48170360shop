package postgres

import (
	"errors"
	"log"
)

func AddToBacketeDB(jsonbytes []byte) (jsonData []byte, err error) {

	q := `
SELECT json_agg(s) FROM (
WITH xid AS (SELECT * FROM  json_to_record($1) x(user_name text, goods_name text, goods_add integer))
SELECT (add_to_basket(xid.user_name, xid.goods_name, xid.goods_add)).* FROM xid
) AS s;`

	if err := db.QueryRowx(q, jsonbytes).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	return jsonData, nil

}

func GetGoodsListDB() (jsonData []byte, err error) {

	//выводит сразу все строки! Если Одна строка - массив с одним значением
	q := `
SELECT json_agg(s) FROM (
SELECT name, description, price FROM get_goods()
) AS s;`

	if err := db.QueryRowx(q).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	if len(jsonData) == 0 {
		err := errors.New("empty value entered")
		log.Printf("Error jsonData: %s", err)
		return nil, err
	}

	return jsonData, nil
}
