package write

import (
	"encoding/json"
	"../../service"
	"fmt"
	response "../../../lib/core/response"
)
type loginSuccessWithData struct {
	AuthToken string
}

func (c *WriteController) Login () {
	var u service.User
	err := json.NewDecoder(c.Request.Body).Decode(&u)
	if err != nil {
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidInput", c.Language)
		return
	}
	if  c.UserId != 0 {
		c.Result, c.StatusCode = response.ErrorCodeResponse("AlreadyLoggedIn", c.Language)
		return
	}

	err, token := u.Login(c.DB, c.Language, c.Configuration, c.DeviceType, c.DeviceToken)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}
	if token != "" {
		c.Result = response.ResponseWithData(formatResponseLogin(token), c.Language.GetMessageTextValue("LoggedIn"))
		return
	}
}

func (c *WriteController) SocialLogin () {

}

func (c *WriteController) MobileLogin () {

}

func (c *WriteController) ResetPassword () {
	var p service.Password
	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidInput", c.Language)
		return
	}
	err = p.ResetPassword(c.DB, c.Language, c.Configuration)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}

	c.Result = response.ResponseWithoutData("SUCCESS", c.Language.GetMessageTextValue("ConfirmationLinkSentToYourEmail"))
}

func (c *WriteController) ChangePassword () {
	var p service.Password
	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidInput", c.Language)
		return
	}
	err = p.ChangePassword(c.DB, c.Language, c.Configuration, c.UserId)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}

	c.Result = response.ResponseWithoutData("SUCCESS", c.Language.GetMessageTextValue("PasswordUpdatedSuccessfully"))
}

func formatResponseLogin(result string) []interface{} {
	pRes := loginSuccessWithData{
		AuthToken: result,
	}
	var data []interface{}
	return append(data, pRes)
}