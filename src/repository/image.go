package repository

import (
	"../../lib/core/database"
	"strconv"
)

func DeleteLabImage(db *database.Database, labId int64,imageId int64) ( err error) {
	var condition = map[string]string {
		"id" : strconv.FormatInt(imageId, 10),
		"lab_id" : strconv.FormatInt(labId, 10),
	}
	_, err = db.Delete("pathlab_img", condition)
	if err != nil  {
		return err
	}

	return nil
}

func InsertLabImage(
	db *database.Database,
	data map[string]string,
) (imageId int64, err error) {
	imageId, err = db.Insert("pathlab_img",  data)
	if err != nil {
		return 0, err
	}

	return imageId, nil
}


func DeleteProductImage(db *database.Database, productId int64,imageId int64) ( err error) {
	var condition = map[string]string {
		"id" : strconv.FormatInt(imageId, 10),
		"product_id" : strconv.FormatInt(productId, 10),
	}
	_, err = db.Delete("products_img", condition)
	if err != nil  {
		return err
	}

	return nil
}

func InsertProductImage(
	db *database.Database,
	data map[string]string,
) (imageId int64, err error) {
	imageId, err = db.Insert("products_img",  data)
	if err != nil {
		return 0, err
	}

	return imageId, nil
}

func GetImageByProductId(db *database.Database, productId int64) (images []map[string]string, err error) {
	var id, imgUrl, addedDateTime string
	var condition = map[string]string {
		"product_id" : strconv.FormatInt(productId, 10),
	}

	var fields = []string{
		"id",
		"img_url",
		"added_datetime",
	}
	var orderBy = []string{"id"}
	rows, err := db.Select("products_img", fields , condition, orderBy, 0)
	if err != nil {
		return images, err
	}

	for rows.Next() {
		err := rows.Scan(&id, &imgUrl, &addedDateTime)
		if err == nil {
			image := map[string]string{
				"id":               id,
				"img_url":          imgUrl,
				"added_datetime":   addedDateTime,
			}
			images = append(images, image)
		}
	}
	defer rows.Close()
	return images, err
}

func GetImageByLabId(db *database.Database, labId int64) (images []map[string]string, err error) {
	var id, imgUrl, addedDateTime string
	var condition = map[string]string {
		"lab_id" : strconv.FormatInt(labId, 10),
	}

	var fields = []string{
		"id",
		"img_url",
		"added_datetime",
	}
	var orderBy = []string{"id"}
	rows, err := db.Select("pathlab_img", fields , condition, orderBy, 0)
	if err != nil {
		return images, err
	}

	for rows.Next() {
		err := rows.Scan(&id, &imgUrl, &addedDateTime)
		if err == nil {
			image := map[string]string{
				"id":               id,
				"img_url":          imgUrl,
				"added_datetime":   addedDateTime,
			}
			images = append(images, image)
		}
	}
	defer rows.Close()
	return images, err
}
