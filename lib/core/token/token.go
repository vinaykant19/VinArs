package token

import (
	"../database"
	"../helper"
	"../errors"
	"../../../localization"
	cnf "../../../configuration"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"github.com/google/uuid"
)

//encryption
func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func createHashMd5(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher(createHash(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := createHash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}

// token system
func GetNewTokenForUser(db *database.Database, local localization.Localization,  conf *cnf.Configuration, userId int64, loginId string)  string {
	var token = ""
	//create new uid
	uid := uuid.New()
	//encrypt uid with hash as salt STR0
	str0 := uid.String() //helper.BytesToString(encrypt([]byte(uid.String()), "MPathIshTanVinars2009!"))
	//bcrypt userId STR1
	userIdStr := strconv.FormatInt(userId, 10)
	str1 := helper.GenerateHashedPassword(userIdStr)
	//bcrypt loginId STR2
	str2 := helper.GenerateHashedPassword(loginId)
	//bcrypt userId_loginId STR3
	str3 := helper.GenerateHashedPassword(userIdStr + "_" + loginId)
	//bcrypt userId again STR4
	str4 := helper.GenerateHashedPassword(userIdStr)
	//Token = STR0+STR1+STR2+STR3+STR4
	token = str0+str1+str2+str3+str4
	//save in db sessionId, userId, LoginId, valid_until = now()+5Min

	err := addToken(db, local, str0, userId, loginId, helper.ToIntiger(conf.Token_Life_In_Minutes))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return token
}

func ValidateTokenAndGetUserId(db *database.Database, local localization.Localization,  conf *cnf.Configuration, token string) (userId int64, sessionId  string) {
	if len(token)< 276 {
		return 0, ""
	}
	//hash secret key
	//split token to get encrypted session
	str0 := token[0:36]
	//decrypt session Id with hash as salt
	//check in db for sessionId if not expire get userId and loginId
	userId, loginId, err := getToken(db, local, str0)
	if err != nil {
		fmt.Println(err)
		return 0, str0
	}
	if userId == 0 {
		return userId, str0
	}
	//split hashed  userId1  and verify with  saved userId
	str1 := token[36:96]
	if ! helper.VerifyHashedPassword(str1, strconv.FormatInt(userId, 10)) {
		return 0, str0
	}
	//split hashed loginId verify with  saved loginId
	str2 := token[96:156]
	if ! helper.VerifyHashedPassword(str2, loginId) {
		return 0, str0
	}
	//split hashed  userId_loginId verify with  userId_loginId from db result
	str3 := token[156:216]
	if ! helper.VerifyHashedPassword(str3, strconv.FormatInt(userId, 10)+"_"+loginId) {
		return 0, str0
	}
	//bcrypt hashed userId2  and verify with  saved userId
	str4 := token[216:276]
	if ! helper.VerifyHashedPassword(str4, strconv.FormatInt(userId, 10)) {
		return 0, str0
	}
	//update in db for sessionId- valid_until = now()+5Min
	err = refreshToken(db, local, str0, helper.ToIntiger(conf.Token_Life_In_Minutes))
	if err != nil {
		fmt.Println(err)
		return 0, str0
	}

	return userId, str0
}

func DeleteToken(db *database.Database, local localization.Localization,  conf *cnf.Configuration, token string) error {
	userId, id := ValidateTokenAndGetUserId(db, local, conf, token)
	if userId == 0 {
		return errors.CustomErrors("InvalidAuthorization")
	}

	err := deleteToken(db, local, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/// db  calls
func addToken(db *database.Database, local localization.Localization,  id string, userId int64, loginId string, life int) error {
	var data = map[string]string {
		"id" : id,
		"user_id" : strconv.FormatInt(userId, 10),
		"login_id" : loginId,
		"valid_until" : time.Now().Add(250 * time.Minute).Format("2006-01-02 15:04:05"),
	}

	_, err := db.Insert("keyclock",  data)
	if err != nil {
		return err
	}

	return nil
}

func refreshToken(db *database.Database, local localization.Localization,  id string, life int) error {
	var condition = map[string]string {
		"id" : id,
	}

	var updateData = map[string]string {
		"valid_until" : time.Now().Add(2500 * time.Minute).Format("2006-01-02 15:04:05"),
	}

	_, err := db.Update("keyclock",  updateData, condition)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func getToken(db *database.Database, local localization.Localization,  id string) (userId int64, loginId string, err error) {
	var condition []interface{}
	condition = append(condition, id)
	condition = append(condition, time.Now().Format("2006-01-02 15:04:05"))

	sql := "SELECT user_id, login_id from keyclock where id=? AND valid_until >=?"
	rows, err := db.Query(sql, condition...)
	if err != nil {
		return 0, "", err
	}
	for rows.Next() {
		err := rows.Scan(&userId, &loginId)
		if err != nil {
			fmt.Println(err)
			return 0, "", err
		}
	}
	rows.Close()

	return userId, loginId, nil
}

func deleteToken(db *database.Database, local localization.Localization,  id string) (err error) {
	var condition = map[string]string {
		"id" : id,
	}
	_, err = db.Delete("keyclock", condition)
	if err != nil {
		return  err
	}

	return nil
}
//9d26c00e-121a-4663-9c5f-8a231cdf33df
//$2a$04$D5bFzfJW/ygy7lCvS/a3Dee.m3yWeYBKfjeuTP9byy4GsNPHJh086
//$2a$04$Zw5qoJfpm4KDzG9lDNn0KO.O/rVox467zycXvvutbYVI/btk6jmCK
//$2a$04$khd3m5tmSabyE05VwsmXieVAvIQO/8e3fh2d9De7yaIc22ILNQpYe
//$2a$04$U2DK1DgQu6fFEOHIIY306uhBqPtzi0jn0rwfGWwRfTkKGYtGQoeeS