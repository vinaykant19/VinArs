package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Port              		string
	Api_URL					string
	Web_URL					string
	Static_Variable   		string
	Connection_String 		string
	Default_Response_type 	string
	Logpath			  		string
	Url_Param_With_Variable string
	DB_Driver 				string
	DB_User					string
	DB_Pwd			 		string
	DB_Host 				string
	DB_Port 				string
	DB_Name 				string
	API_KEY 				string
	Token_Life_In_Minutes 	string
	Non_Token_Calls_GET 	string
	Non_Token_Calls_POST	string
	FileUpload_MaxSize_MB 	string
	FileUpload_Path  		string
	SMTP_Server				string
	SMTP_Port				int
	SMTP_User				string
	SMTP_Pass				string
	ADMIN_EMAIL_GROUP		string
	FROM_SUPPORT_EMAIL		string
	FROM_ADMIN_EMAIL		string
	REPLY_TO_EMAIL			string
}

func SetConfig(configuration *Configuration, env string) (err error) {
	filename := "./configuration/config."+env+".json"
	file, err := os.Open(filename)
	if err != nil {  return err }
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {  return err }

	return nil
}