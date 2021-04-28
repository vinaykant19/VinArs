package errors
import (
	"github.com/micro/go-micro/errors"
)
type CustomError struct{
	Id string `json:"id"`
	Code int32 `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

var errorCodes = map[string]int32{
	"AlreadyLoggedIn"  : 200,

	"InvalidInput" : 400,
	"EmailEmpty": 400,
	"EmailInvalid": 400,
	"EmailDuplicate": 400,
	"MobileEmpty": 400,
	"PasswordEmpty": 400,
	"NameEmpty": 400,
	"FirstNameEmpty": 400,
	"LastNameEmpty": 400,
	"GenderEmpty": 400,
	"Address1Empty": 400,
	"AddressEmpty": 400,
	"CityEmpty": 400,
	"CountryEmpty": 400,

	"InvalidFileType": 400,
	"InvalidFile": 400,
	"FileMaxSizeLimit": 400,
	"FileUploadLimit": 400,

	"InvalidAuthorization" : 401,
	"InvalidHash": 401,
	"PermissionDenied" : 401,

	"InactiveUser" : 402,

	"NotFound": 404,

	"InternalError" : 500,
}

func CustomErrors(id string) error  {
	if _, ok := errorCodes[id]; ok {
		return errors.New(id, id, errorCodes[id])
	}

	return errors.New(id, "InternalError", errorCodes["InternalError"])
}
