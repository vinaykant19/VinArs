package repository

import (
	"../../lib/core/database"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func InsertUser(
	db *database.Database,
	email string,
	name string,
	password string,
	status int,
	ipAddress string,
	hasCode string,
) (userId int64, err error) {

	var data = map[string]string{
		"email":             email,
		"password":          password,
		"name":              name,
		"status":            strconv.Itoa(status),
		"register_ip":       ipAddress,
		"register_datetime": time.Now().Format("2006-01-02 15:04:05"),
		"update_datetime":   time.Now().Format("2006-01-02 15:04:05"),
		"hashCode":          hasCode,
	}

	userId, err = db.Insert("user", data)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func UpdateUser(
	db *database.Database,
	userId int64,
	data map[string]string,
) (err error) {
	var condition = map[string]string{
		"id": strconv.FormatInt(userId, 10),
	}
	_, err = db.Update("user", data, condition)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(db *database.Database,
	email string,
) (userId int64, name string, password string, status int8, err error) {

	var condition = map[string]string{
		"email": email,
	}
	var fields = []string{
		"id",
		"email",
		"name",
		"password",
		"status",
	}
	var orderBy = []string{"id"}
	rows, err := db.Select("user", fields, condition, orderBy, 1)
	if err != nil {
		return 0, "", "", 0, err
	}

	if rows.Next() {
		err := rows.Scan(&userId, &email, &name, &password, &status)
		if err == nil {
			return userId, name, password, status, nil
		}
	}
	defer rows.Close()
	return 0, "", "", 0, err
}

func IsUserExistByEmail(
	db *database.Database,
	email string,
) (bool, error) {
	var condition = map[string]string{
		"email": email,
	}
	var fields = []string{
		"id",
	}
	var orderBy = []string{"id"}
	rows, err := db.Select("user", fields, condition, orderBy, 1)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		defer rows.Close()
		return true, nil
	}
	defer rows.Close()
	return false, nil
}

func GetUserById(db *database.Database,
	userId int64,
) (email string, name string, status string, password string, err error) {

	var condition = map[string]string{
		"id": strconv.FormatInt(userId, 10),
	}
	var fields = []string{
		"email",
		"name",
		"status",
		"password",
	}
	var orderBy = []string{"id"}
	rows, err := db.Select("user", fields, condition, orderBy, 1)
	if err != nil {
		return email, name, status, password, err
	}

	if rows.Next() {
		err = rows.Scan(&email, &name, &status, &password)
	}
	defer rows.Close()
	return email, name, status, password, err
}

func GetUserByHash(db *database.Database,
	hash string,
) (userId int64, err error) {

	var condition = map[string]string{
		"hashCode": hash,
	}
	var fields = []string{
		"id",
	}
	var orderBy = []string{"id"}
	rows, err := db.Select("user", fields, condition, orderBy, 1)
	if err != nil {
		return userId, err
	}

	if rows.Next() {
		err = rows.Scan(&userId)
	}
	defer rows.Close()
	return userId, err
}

func SearchUser(
	db *database.Database,
	name string,
	email string,
) (userList []map[string]string, err error) {
	var id, register_datetime, update_datetime, register_ip, status string
	var isAnyCondition bool

	isAnyCondition = false
	var args []interface{}

	var query = "SELECT id, name, email, register_datetime, update_datetime," +
		" register_ip, status FROM user WHERE 1 "

	if len(strings.Trim(name, "%20")) > 0 && name != "_" {
		query += " AND name like '%" + strings.Trim(name, "%20") + "%' "
		isAnyCondition = true
	}

	if len(strings.Trim(email, "%20")) > 0 {
		query += " AND email = ? "
		args = append(args, strings.Trim(email, "%20"))
		isAnyCondition = true
	}
	if !isAnyCondition {
		return userList, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return userList, err
	}

	var user map[string]string
	for rows.Next() {
		err := rows.Scan(&id, &name, &email, &register_datetime, &update_datetime, &register_ip, &status)

		if err == nil {
			user = map[string]string{
				"id":                id,
				"name":              name,
				"email":             email,
				"register_datetime": register_datetime,
				"update_datetime":   update_datetime,
				"register_ip":       register_ip,
				"status":            status,
			}
			userList = append(userList, user)
		}
	}
	defer rows.Close()

	return userList, err
}

func AllUser(db *database.Database) (userList []map[string]string, err error) {
	var id, name, email, register_datetime, update_datetime, register_ip, status string

	var args []interface{}

	var query = "SELECT id, name, email, register_datetime, update_datetime," +
		" register_ip, status FROM user WHERE 1 "

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return userList, err
	}

	var user map[string]string
	for rows.Next() {
		err := rows.Scan(&id, &name, &email, &register_datetime, &update_datetime, &register_ip, &status)

		if err == nil {
			user = map[string]string{
				"id":                id,
				"name":              name,
				"email":             email,
				"register_datetime": register_datetime,
				"update_datetime":   update_datetime,
				"register_ip":       register_ip,
				"status":            status,
			}
			userList = append(userList, user)
		}
	}
	defer rows.Close()

	return userList, err
}
