package write

import (
	"encoding/json"
	"../../service"
	"fmt"
	"strconv"
	helper "../../../lib/core/helper"
	response "../../../lib/core/response"
)
type registerSuccess struct {
	UserId string
}
func (c *WriteController) UserRegister () {
	var u service.User
	err := json.NewDecoder(c.Request.Body).Decode(&u)
	if err != nil {
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidInput", c.Language)
		return
	}
	u.RegisterIpAddress = helper.ReadUserIP(c.Request)

	err, userId := u.AddUser(c.DB, c.Language, c.Configuration)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}
	if userId != 0 {
		c.Result = response.ResponseWithData(formatResponseRegister(strconv.FormatInt(userId, 10)), c.Language.GetMessageTextValue("RegisteredSuccessfully"))
		return
	}
}

func (c *WriteController) ConfirmUser(hash string)  {
	var u service.User
	err, userId := u.ConfirmUser(hash, c.DB, c.Language)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}
	if userId != 0 {
		c.Result = response.ResponseWithoutData("SUCCESS", c.Language.GetMessageTextValue("UserConfirmedSuccessfully"))
		return
	}
}

func formatResponseRegister(result string) []interface{} {
	pRes := registerSuccess{
		UserId: result,
	}
	var data []interface{}
	return append(data, pRes)
}
