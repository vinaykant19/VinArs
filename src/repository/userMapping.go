package repository

import (
	"../../lib/core/database"
	"strconv"
)
func InsertUserMappingGroup (
	db *database.Database,
	data map[string]string,
) (id int64, err error) {
	id, err = db.Insert("user_type",  data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//1= user; 2= agent; 3= lab; 4= doctor; 5=manager; 6=admin; 7=superAdmin
func GetUserMappingGroup(
	db *database.Database,
	userId int64) (userMappingGroup map[string]int, err error) {
	var usertype string
	var condition = map[string]string {
		"user_id" : strconv.FormatInt(userId, 10),
	}
	userMappingGroup = map[string]int{}

	var fields = []string{
		"user_type",
	}
	var orderBy = []string{"id"}
	rows, err := db.Select("user_type", fields , condition, orderBy, 1)
	if err != nil {
		return userMappingGroup, err
	}

 	for rows.Next() {
		err = rows.Scan(&usertype)
		if err == nil {
			userMappingGroup[usertype] = 1
		}
	}
	defer rows.Close()

	return userMappingGroup, err
}
