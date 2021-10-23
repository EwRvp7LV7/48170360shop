package postgres

import (
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	//"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
)

func GetUIDPasswordHash(jsonbytes []byte) (string, string, bool, error) {
	q := `
WITH xid AS (SELECT * FROM  json_to_record($1) as x(user_name text))
select user_registered.id_user, user_registered.password_user, user_registered.is_staff from user_registered, xid 
where user_registered.user_name =  xid.user_name`

	//InputUserProfile приходящая извне информация для создания пользователя.
	type AuthenticationDB struct {
		UserName string `db:"id_user"`
		// Email         string `json:"e_mail"`
		Password string `db:"password_user"`
		IsAdmin  bool   `db:"is_staff"`
	}

	var adb AuthenticationDB

	if err := db.QueryRowx(q, jsonbytes).StructScan(&adb); err != nil {
		log.Printf("Error ToFUserProfileDB: %s", err)
		return "", "", false, err
	}

	if len(strings.TrimSpace(adb.Password)) == 0 {
		err := errors.New("wrong user name")
		log.Printf("Error UserPass: %s", err)
		return "", "", false, err
	}
	fmt.Printf("iseq '%s','%s'", adb.UserName, adb.Password) //iseq ZcexR6Ebv3KYKl ZcexR6Ebv3KYKl

	return adb.UserName, adb.Password, adb.IsAdmin, nil

}
