package read

import (
	response "../../../lib/core/response"
	"../../service"
	"fmt"
)
type userDetailSuccessWithData struct {
	User service.User
}
type userListSuccessWithData struct {
	User []service.User
}

func (c *ReadController) UserDetail (userId int64) {
	var u service.User
	permission := service.IsAllowed(c.DB, c.Language, c.UserId, "UserDetail")
	if permission == 0 {
		c.Result, c.StatusCode = response.ErrorCodeResponse("PermissionDenied", c.Language)
		return
	}
	if permission == 2 {
		if userId !=c.UserId {
			c.Result, c.StatusCode = response.ErrorCodeResponse("PermissionDenied", c.Language)
			return
		}
	}
	u.Id = userId

	result := service.GetUserById(c.DB, c.Language, userId)

	data := formatUserResponse(result)
	c.Result = response.ResponseWithData(data, "")
	return
}

func (c *ReadController) Me () {
	var u service.User
	u.Id = c.UserId

	result := service.GetUserById(c.DB, c.Language, c.UserId)

	data := formatUserResponse(result)
	c.Result = response.ResponseWithData(data, "")
	return
}

func (c *ReadController) SearchUser (name string, email string) {
	var u service.User
	permission := service.IsAllowed(c.DB, c.Language, c.UserId, "SearchUser")
	if permission == 0 {
		c.Result, c.StatusCode = response.ErrorCodeResponse("PermissionDenied", c.Language)
		return
	}
	result, err := u.SearchUser(c.DB, name, email)

	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}

	data := formatUserListResponse(result)
	c.Result = response.ResponseWithData(data, "")
	return
}

func (c *ReadController) AllUser () {
	var u service.User
	permission := service.IsAllowed(c.DB, c.Language, c.UserId, "SearchUser")
	if permission == 0 {
		c.Result, c.StatusCode = response.ErrorCodeResponse("PermissionDenied", c.Language)
		return
	}
	result, err := u.AllUsers(c.DB)

	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}

	data := formatUserListResponse(result)
	c.Result = response.ResponseWithData(data, "")
	return
}

func formatUserListResponse(result []service.User) []interface{} {
	pRes := userListSuccessWithData{
		User: result,
	}
	var data []interface{}
	return append(data, pRes)
}

func formatUserResponse(result service.User) []interface{} {
	pRes := userDetailSuccessWithData{
		User: result,
	}
	var data []interface{}
	return append(data, pRes)
}
