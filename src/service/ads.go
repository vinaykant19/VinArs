package service

import (
	"../../localization"
	"../../lib/core/database"
	"../../lib/core/errors"
	"../repository"
	"fmt"
	"strconv"
)
type Ads struct {
	Id      		int64
	Heading			string
	Description		string
	Type			string //product/weblink
	WebLink			string
	ProductId 		int64
	StartFrom 		string
	ActiveUpTo		string
	City 			string
	Priority 		int
	Status          int
	AddedDateTime   string
	UpdateDateTime  string
	AddedBy         int64
}

func (ad Ads) GetAds(db *database.Database, local localization.Localization,  userId int64, city string, searchText string, productId int64) (adsList []Ads, err error) {

	//call db query
	adsListItems, err := repository.GetAdsList(db, userId, city, searchText, productId)
	if err != nil {
		//log error
		fmt.Println(err)
		return adsList, errors.CustomErrors("InternalError")
	}
	for _, data := range adsListItems {
		var ads = RepositoryDataToAds(data)

		adsList = append(adsList, ads)
	}

	return adsList, err
}

func RepositoryDataToAds( data map[string]string) Ads {
	var ads = Ads{}

	ads.Id, _ = strconv.ParseInt(data["id"], 10, 64)
	ads.Heading = data["heading"]
	ads.Description = data["description"]
	ads.Type = data["type"]
	ads.WebLink = data["webLink"]
	ads.City = data["city"]
	ads.StartFrom = data["start_from"]
	ads.ActiveUpTo = data["active_upto"]
	ads.ProductId, _ = strconv.ParseInt(data["product_id"], 10, 64)
	ads.Priority, _ = strconv.Atoi(data["priority"])
	ads.Status, _ = strconv.Atoi(data["status"])
	ads.AddedDateTime = data["added_datetime"]
	ads.UpdateDateTime = data["last_updated"]
	ads.AddedBy, _ = strconv.ParseInt(data["last_updated_by"], 10, 64)

	return ads
}