package postgres

import (
	"log"
)

func GetBasketsListtDB(user_name string) (jsonData []byte, err error) {

	//выводит сразу все строки! Если Одна строка - массив с одним значением
	q := `
SELECT json_agg(s) FROM (
SELECT name, description, price FROM get_goods($1)
) AS s;`

	if err := db.QueryRowx(q, user_name).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	return jsonData, nil
}

func NewGoods(jsonbytes []byte) (jsonData []byte, err error) {

	q := `
SELECT json_agg(s) FROM (
WITH xid AS (SELECT * FROM  json_to_record($1) x(user_name text, goods_name text, good_descrription text, goods_add integer, good_price money))
SELECT (new_goods(xid.user_name, xid.goods_name, xid.good_descrription, xid.goods_add, xid.good_price)).* FROM xid
) AS s;`

	if err := db.QueryRowx(q, jsonbytes).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	return jsonData, nil

}

func AddGoodsToStoreDB(jsonbytes []byte) (jsonData []byte, err error) {

	q := `
SELECT json_agg(s) FROM (
WITH xid AS (SELECT * FROM  json_to_record($1) x(user_name text, goods_name text, good_descrription text, goods_add integer, good_price money))
SELECT (add_goods2store(xid.user_name, xid.goods_name, xid.good_descrription, xid.goods_add).* FROM xid
) AS s;`

	if err := db.QueryRowx(q, jsonbytes).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	return jsonData, nil

}
