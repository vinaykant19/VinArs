package response

import (
	"encoding/json"

	"../../../localization"
	"../errors"
)

type ResArray struct {
	Data []interface{}
	Status string
	StatusCode string
	Msg  string
}

func ResponseWithData(data []interface{}, msg string) string{
	resArray := ResArray{
		Data: data,
		Status:"success",
		StatusCode:"SUCCESS",
		Msg:msg,
	}

	fmtRes,_ := json.Marshal(resArray)
	return string(fmtRes)
}

func ResponseWithoutData(statusCode string, msg string) string{
	var data []interface{}
	status := "Failed"
	if statusCode == "SUCCESS" {
		status = "success"
	}
	resArray := ResArray{
		Data: data,
		Status:status,
		StatusCode:statusCode,
		Msg:msg,
	}
	fmtRes,_ := json.Marshal(resArray)
	return string(fmtRes)
}

func ErrorResponse(err error, lang localization.Localization) (string, int){
	var errors = errors.CustomError{}
	data := []byte(err.Error())
	err = json.Unmarshal(data, &errors)
	return ResponseWithoutData(errors.Id, lang.GetErrorValue(errors.Detail)), int(errors.Code)
}

func ErrorCodeResponse(errCode string, lang localization.Localization) (string, int){
	var customError = errors.CustomError{}

	data := []byte(errors.CustomErrors(errCode).Error())
	_ = json.Unmarshal(data, &customError)
	return ResponseWithoutData(customError.Id , lang.GetErrorValue(customError.Detail)), int(customError.Code)
}
