package service

import (
	"../../localization"
	"../../lib/core/database"
	"../../lib/core/errors"
	"../repository"
	"fmt"
	"strconv"
	"time"
)

type Image struct {
	Id int64
	ConnectId int64
	ImgUrl string
	Status int
	AddedDateTime string
}

func (i *Image) DeleteProductImage(db *database.Database) (err error) {

	err = repository.DeleteProductImage(db, i.ConnectId, i.Id)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError")
	}

	return nil
}

func (i *Image) AddProductImage (db *database.Database) (err error, imageId int64) {
	//validate input
	if len(i.ImgUrl) < 1 {
		return errors.CustomErrors("ImageUrlEmpty"),0
	}
	if i.ConnectId < 1 {
		return errors.CustomErrors("ProductConnectionEmpty"),0
	}

	i.AddedDateTime  = time.Now().Format("2006-01-02 15:04:05")

	//call db query
	var data = map[string]string {
		"img_url": i.ImgUrl,
		"product_id": strconv.FormatInt(i.ConnectId, 10),
		"added_datetime": i.AddedDateTime,
	}

	imageId, err = repository.InsertProductImage(db, data)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError"),0
	}

	return nil, imageId
}

func (i *Image) GetAllImagesByProductId(db *database.Database, local localization.Localization,  productId int64) (result []Image, err error) {
	data, err := repository.GetImageByProductId(db, productId)
	if err != nil {
		//log error
		fmt.Println(err)
		return result, errors.CustomErrors("InternalError")
	}

	for _, data := range data {
		var image = Image{}
		image.Id, _ = strconv.ParseInt(data["id"], 10, 64)
		image.ImgUrl = data["img_url"]
		image.AddedDateTime = data["added_datetime"]

		result = append(result, image)
	}

	return result, err
}
