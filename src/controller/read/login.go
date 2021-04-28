package read
import (
	response "../../../lib/core/response"
	authToken "../../../lib/core/token"
	"fmt"
)

func (c *ReadController) Logout () {
	if  c.UserId == 0 {
		c.Result, _ = response.ErrorCodeResponse("InvalidAuthorization", c.Language)
		return
	}
	err := authToken.DeleteToken(c.DB, c.Language, c.Configuration, c.Token)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, _ = response.ErrorResponse(err, c.Language)
		return
	}
	c.Result = response.ResponseWithoutData("SUCCESS", c.Language.GetMessageTextValue("LoggedOut"))
	return
}
