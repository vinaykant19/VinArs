package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	cnf "../../configuration"
	"../../lib/core/database"
	"../../lib/core/errors"
	helper "../../lib/core/helper"
	authToken "../../lib/core/token"
	"../../lib/core/validator"
	"../../localization"
	"../repository"
)
type Password struct {
	Email string
	HashCode string
	OldPassword string
	Password string
}

type User struct {
	Id int64
	Email string
	Name string
	Password string
	Status int
	RegisterDateTime string
	UpdateDateTime string
	RegisterIpAddress string
	HashCode string
	AgencyId int64
	AgencyCity string
	UserType string
}

func (u *User) AddUser(db *database.Database, local localization.Localization,  conf *cnf.Configuration) (err error, userId int64){
	//validate input
	if len(u.Email) < 1 {
		return errors.CustomErrors("EmailEmpty"),0
	}
	if len(u.Password) < 1 {
		return errors.CustomErrors("PasswordEmpty"),0
	}
	if len(u.Name) < 1 {
		return errors.CustomErrors("NameEmpty"),0
	}
	if !validator.IsValidEmail(u.Email) {
		return errors.CustomErrors("EmailInvalid"),0
	}

	//check if email already exist
	ok, err := repository.IsUserExistByEmail(db, u.Email)
	if ok == true {
		return errors.CustomErrors("EmailDuplicate"),0
	}
	//prepare data set
	u.Status = 1
	u.Password = helper.GenerateHashedPassword(u.Password)
	u.RegisterDateTime  = time.Now().Format("2006-01-02 15:04:05")
	//call db query

	uid := uuid.New()
	u.HashCode = uid.String();
	userId, err = repository.InsertUser(db, u.Email, u.Name, u.Password, u.Status, u.RegisterIpAddress, u.HashCode)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError"),0
	}

	//send confirmation link
	var link = conf.Web_URL + "/confirm-user/" + u.HashCode
	subject := local.GetEmailTextValue("SubjectConfirmYourRegistration")
	msgBody :=  local.GetEmailTextValue("ConfirmYourRegistrationEmailText")
	msgBody = strings.ReplaceAll(msgBody, "#{NAME}#", u.Name)
	msgBody = strings.ReplaceAll(msgBody, "#{LINK}#", link)
	toEmail :=make(map[string]string)
	ccEmail :=make(map[string]string)
	bccEmail :=make(map[string]string)
	files :=make(map[int]string)

	toEmail[u.Name] =  u.Email
	SendEmail(conf, conf.FROM_SUPPORT_EMAIL, toEmail, ccEmail, bccEmail, subject, msgBody, true, files)

	return nil, userId
}

func (u *User) Login(db *database.Database, local localization.Localization,  configuration *cnf.Configuration, deviceName string, deviceToken  string) (err error, token string){
	//validate input
	if len(u.Email) < 1 {
		return errors.CustomErrors("EmailEmpty"),""
	}
	if len(u.Password) < 1 {
		return errors.CustomErrors("PasswordEmpty"),""
	}
	if !validator.IsValidEmail(u.Email) {
		return errors.CustomErrors("EmailInvalid"),""
	}

	//fetch user if email exist
	userId, _, password, status, err := repository.GetUserByEmail(db, u.Email)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError"),""
	}
	if ! helper.VerifyHashedPassword(password, u.Password) {
		return errors.CustomErrors("InvalidAuthorization"),""
	}
	if status != 1 {
		fmt.Println(status)
		return errors.CustomErrors("InactiveUser"),""
	}
	token = authToken.GetNewTokenForUser(db, local, configuration, userId, u.Email)

	if deviceToken != "" && deviceName != ""{
		err = repository.ReplaceDeviceToken(db, userId, deviceName, deviceToken)
		println(err)
	}
	return nil, token
}

func (u *User) ConfirmUser(hash string, db *database.Database, local localization.Localization) (err error, userId int64) {
	userId, err = repository.GetUserByHash(db, hash)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError"),userId
	}

	if userId == 0 {
		return errors.CustomErrors("InvalidHash"),userId
	}

	u.UpdateDateTime  = time.Now().Format("2006-01-02 15:04:05")
	//call db query
	var data = map[string]string {
		"status": "1",
		"hashCode": "",
		"update_dateTime": u.UpdateDateTime,
	}

	err = repository.UpdateUser(db, userId, data)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError"), 0
	}

	return err, userId
}

func (u *User) AllUsers(
	db *database.Database,
) (result []User, err error) {

	userData, err := repository.AllUser(db)
	if err != nil {
		//log error
		fmt.Println(err)
		return result, errors.CustomErrors("InternalError")
	}

	for _, data := range userData {
 		var user = User{}
		user.Id, _ = strconv.ParseInt(data["id"], 10, 64)
		user.Name =  data["name"]
		user.Email =  data["email"]
		user.Status, _ = strconv.Atoi(data["status"])
		user.RegisterDateTime =  data["register_datetime"]
		user.UpdateDateTime =  data["update_datetime"]
		user.RegisterIpAddress =  data["register_ip"]

		result = append(result, user)
	}

	return result, err
}


func (u *User) SearchUser(
	db *database.Database,
	name string,
	email string,
) (result []User, err error) {

	userData, err := repository.SearchUser(db, name, email)
	if err != nil {
		//log error
		fmt.Println(err)
		return result, errors.CustomErrors("InternalError")
	}

	for _, data := range userData {
		var user = User{}
		user.Id, _ = strconv.ParseInt(data["id"], 10, 64)
		user.Name =  data["name"]
		user.Email =  data["email"]
		user.Status, _ = strconv.Atoi(data["status"])
		user.RegisterDateTime =  data["register_datetime"]
		user.UpdateDateTime =  data["update_datetime"]
		user.RegisterIpAddress =  data["register_ip"]

		result = append(result, user)
	}

	return result, err
}

func (p *Password) ResetPassword(
	db *database.Database,
	local localization.Localization,
	conf *cnf.Configuration,
) (err error) {
	//validate input
	if len(p.Email) < 1 {
		return errors.CustomErrors("EmailEmpty")
	}

	//fetch user if email exist
	userId, name, _, _, err := repository.GetUserByEmail(db, p.Email)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError")
	}
	if (userId == 0) {
		return errors.CustomErrors("InvalidAuthorization")
	}
	uid := uuid.New()
	var hashCode = uid.String();
	var updateDateTime  = time.Now().Format("2006-01-02 15:04:05")
	var data = map[string]string {
		"hashCode": hashCode,
		"update_dateTime": updateDateTime,
	}

	err = repository.UpdateUser(db, userId, data)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError")
	}

	//send confirmation link
	var link = conf.Web_URL + "/resetPassword/" + hashCode
	subject := local.GetEmailTextValue("SubjectResetPasswordRRequest")
	msgBody :=  local.GetEmailTextValue("ResetPasswordRRequestEmailText")
	msgBody = strings.ReplaceAll(msgBody, "#{NAME}#", name)
	msgBody = strings.ReplaceAll(msgBody, "#{LINK}#", link)
	toEmail :=make(map[string]string)
	ccEmail :=make(map[string]string)
	bccEmail :=make(map[string]string)
	files :=make(map[int]string)

	toEmail[name] =  p.Email
	SendEmail(conf, conf.FROM_SUPPORT_EMAIL, toEmail, ccEmail, bccEmail, subject, msgBody, true, files)

	return nil
}

func (p *Password) ChangePassword(
	db *database.Database,
    local localization.Localization,
	conf *cnf.Configuration,
	userId int64,
) (err error) {
	//validate input
	if len(p.OldPassword) < 1 && len(p.HashCode) < 1{
		return errors.CustomErrors("InvalidAuthorization")
	}

	if len(p.Password) < 1{
		return errors.CustomErrors("PasswordEmpty")
	}
	var updateDateTime  = time.Now().Format("2006-01-02 15:04:05")
	var hashedPassword = helper.GenerateHashedPassword(p.Password)
	if len(p.HashCode) > 1{
		userId, err = repository.GetUserByHash(db, p.HashCode)
		if err != nil {
			//log error
			fmt.Println(err)
			return errors.CustomErrors("InternalError")
		}

		if userId == 0 {
			return errors.CustomErrors("InvalidHash")
		}

	} else if len(p.OldPassword) > 1 && userId != 0 {
		//fetch user if  exist
		_, _, _, password, err := repository.GetUserById(db, userId)
		if err != nil {
			//log error
			fmt.Println(err)
			return errors.CustomErrors("InternalError")
		}

		if !helper.VerifyHashedPassword(password, p.OldPassword) {
			return errors.CustomErrors("InvalidAuthorization")
		}
	} else {
		return errors.CustomErrors("InvalidAuthorization")
	}

	var data = map[string]string {
		"password": hashedPassword,
		"update_dateTime": updateDateTime,
		"hashCode": "",
	}
	err = repository.UpdateUser(db, userId, data)
	if err != nil {
		//log error
		fmt.Println(err)
		return errors.CustomErrors("InternalError")
	}

	return nil
}

func GetUserById(db *database.Database, local localization.Localization,  userId int64) (user User) {
	//fetch user if  exist
	email, name, status, _, err := repository.GetUserById(db, userId)
	if err != nil {
		//log error
		fmt.Println(err)
		return user
	}

	user.Id = userId
	user.Name = name
	user.Email = email
	user.Status, _ = strconv.Atoi(status)

	user.UserType = GetUserGroupType(db, local, userId)

	return user
}


