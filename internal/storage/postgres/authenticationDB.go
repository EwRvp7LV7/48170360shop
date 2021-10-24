package postgres

import (
	// "encoding/json"
	"errors"
	"log"
	"strings"
	//"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
)

func GetUIDPasswordHash(jsonbytes []byte) (string, string, int, error) {
	q := `
WITH xid AS (SELECT * FROM  json_to_record($1) as x(account text))
select * from b_users JOIN xid USING(account);`

	//приходящая извне информация для создания пользователя.
	type AuthenticationDB struct {
		UserId     string `db:"user_id"`
		UserName   string `db:"account"`
		Password   string `db:"password"`
		IsCustomer int    `db:"type_id"`
	}

	var adb AuthenticationDB

	if err := db.QueryRowx(q, jsonbytes).StructScan(&adb); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return "", "", 0, err
	}

	if len(strings.TrimSpace(adb.Password)) == 0 {
		err := errors.New("wrong user name")
		log.Printf("Error UserPass: %s", err)
		return "", "", 0, err
	}

	return adb.UserId, adb.Password, adb.IsCustomer, nil

}
