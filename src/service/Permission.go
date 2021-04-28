package service

import (
	"../../localization"
	"../../lib/core/database"
	"../../lib/core/errors"
	"../repository"
	"fmt"
)
var UserGroup = map[int]string{
	1 : "User",
	5 : "manager",
	6 : "admin",
	7 : "superadmin",
}

//RS - read Self, < AS - add self, < US - Update self, < DS - Delete self <
//R - read, < A - add, < U - Update, < D - Delete
var NormalUser = map[string]string {

	"UserDetail" : "2",
	"Me" : "2",

	"UserAddressList" : "2",
	"UserAddressDetail" : "2",
	"AddUserAddress": "2",
	"UpdateUserAddress": "2",
	"RemoveUserAddress": "2",

	"ProductDetailById" : "1",
	"ProductDetail" : "1",
	"SearchProduct" : "1",
	"SearchDoctor" : "1",
	"UploadImage" : "2",

}

var ManagerUser = map[string]string {

	"UserDetail" : "1",
	"SearchUser":"1",

	"UserAddressList" : "1",
	"UserAddressDetail" : "1",
	"AddUserAddress": "1",
	"UpdateUserAddress": "1",

	"ProductDetailById" : "1",
	"ProductDetail" : "1",
	"SearchProduct" : "1",
	"SearchDoctor" : "1",
	"UploadImage" : "1",

	"AddProduct": "3",
	"UpdateProduct": "3",

}

var AdminUser = map[string]string {

	"DeleteImage":"1",

	"UserDetail" : "1",
	"SearchUser":"1",

	"UserAddressList" : "1",
	"UserAddressDetail" : "1",
	"AddUserAddress": "1",
	"UpdateUserAddress": "1",
	"RemoveUserAddress":"1",

	"ProductDetailById" : "1",
	"ProductDetail" : "1",
	"SearchProduct" : "1",
	"SearchDoctor" : "1",
	"UploadImage" : "1",

	"AddProduct": "1",
	"UpdateProduct": "1",
	"RemoveProduct":"1",

	"AddCoupon":"1",
}

type PermissionUserGroup struct {

}

func IsAllowed(db *database.Database, local localization.Localization,  userId int64, action string) int {
	UserGroup, err := GetUserGroup(db, local, userId)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if _, ok := UserGroup["7"]; ok {
		return 1
	}
	if _, ok := UserGroup["6"]; ok {
		if _, ok := AdminUser[action]; ok {
			return 1
		}
	}
	if _, ok := UserGroup["5"]; ok {
		if _, ok := ManagerUser[action]; ok {
			if ManagerUser[action] == "1" {
				return 1
			} else if ManagerUser[action] == "2" {
				return 2
			} else {
				return 3
			}
		}
	}

	if NormalUser[action] == "1" {
		return 1
	} else {
		return 2
	}

	return 0
}

func GetUserGroup(db *database.Database, local localization.Localization,  userId int64) (groups map[string]int, err error) {

	groups, err = repository.GetUserMappingGroup(db, userId)
	if err != nil {
		//log error
		fmt.Println(err)
		return groups, errors.CustomErrors("InternalError")
	}

	return groups, err
}

func GetUserGroupType(db *database.Database, local localization.Localization,  userId int64)  string{
	UserGroup, err := GetUserGroup(db, local, userId)
	if err != nil {
		return "User"
	}
	if _, ok := UserGroup["7"]; ok {
		return "SuperAdmin"
	}
	if _, ok := UserGroup["6"]; ok {
		return "Admin"
	}
	if _, ok := UserGroup["5"]; ok {
		return "Manager"
	}

	return "User"
}