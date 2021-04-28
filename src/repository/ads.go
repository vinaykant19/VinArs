package repository

import (
	"../../lib/core/database"
)

func GetAdsList(db *database.Database, userId int64, cityName string, searchText string, productId int64) (adsList []map[string]string, err error) {
	var args []interface{}
	var id, heading, description, adType, webLink, city, start_from, active_upto, product_id, priority, status, added_datetime, last_updated, last_updated_by string
	args = append(args, userId)
	args = append(args, cityName)
	args = append(args, searchText)
	args = append(args, productId)
	var sql =  "CALL Search_ads(?, ?, ?, ?)"
	rows, err := db.Query(sql, args...)
	if err != nil {
		return adsList, err
	}
	for rows.Next() {
		err = rows.Scan(&id, &heading, &description, &adType, &webLink, &city, &start_from,
			&active_upto, &product_id, &priority, &status, &added_datetime, &last_updated, &last_updated_by)
		if err == nil {
			var ad = map[string]string{
				"id": id,
				"heading": heading,
				"description": description,
				"type": adType,
				"webLink": webLink,
				"city": city,
				"start_from": start_from,
				"active_upto": active_upto,
				"product_id": product_id,
				"priority": priority,
				"status": status,
				"added_datetime": added_datetime,
				"last_updated": last_updated,
				"last_updated_by": last_updated_by,
			}
			adsList = append(adsList, ad)
		}
	}

	defer rows.Close()
	return adsList, err
}
