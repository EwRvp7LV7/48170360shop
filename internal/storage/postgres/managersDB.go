package postgres

import (
	"log"
)

func GetBasketsListDB(user_name string) (jsonData []byte, err error) {

	//выводит сразу все строки! Если Одна строка - массив с одним значением
	q := `
SELECT json_agg(s) FROM (
SELECT * FROM get_baskets($1)
) AS s;`

	if err := db.QueryRowx(q, user_name).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	return jsonData, nil
}

func NewGoodsDB(jsonbytes []byte) (jsonData []byte, err error) {

	q := `
WITH xid AS (SELECT * FROM  json_to_record($1) x(user_name text, goods_name text, goods_descrription text, goods_add integer, goods_price money))
SELECT new_goods(xid.user_name, xid.goods_name, xid.goods_descrription, xid.goods_add, xid.goods_price) FROM xid;`

	if err := db.QueryRowx(q, jsonbytes).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	return jsonData, nil

}

func AddGoodsToStoreDB(jsonbytes []byte) (jsonData []byte, err error) {

	q := `
WITH xid AS (SELECT * FROM  json_to_record($1) x(user_name text, goods_name text, goods_descrription text, goods_add integer, goods_price money))
SELECT add_goods2store(xid.user_name, xid.goods_name, xid.goods_add) FROM xid;`

	if err := db.QueryRowx(q, jsonbytes).Scan(&jsonData); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return nil, err
	}

	return jsonData, nil

}
