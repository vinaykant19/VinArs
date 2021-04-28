package repository

import (
	"../../lib/core/database"
)
func  ReplaceDeviceToken(
	db *database.Database,
	userId int64,
	deviceName string,
	deviceToken string,
) (err error) {
	var args []interface{}
	args = append(args, userId)
	args = append(args, deviceName)
	args = append(args, deviceToken)

	var sql =  "REPLACE into user_device_token (user_id, device_name, device_token) values(?, ?, ?)"
	_, err = db.Query(sql, args...)

	return err
}
