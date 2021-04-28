package localization

import (
	"encoding/json"
	"os"
	"reflect"
)

type ErrorText struct {
	AlreadyLoggedIn string
	InternalError string
	InvalidAuthorization string
	InvalidInput   string
	EmailEmpty  string
	EmailInvalid  string
	EmailDuplicate  string
	MobileEmpty  string
	PasswordEmpty  string
	NameEmpty  string
	FirstNameEmpty  string
	LastNameEmpty  string
	GenderEmpty  string
	Address1Empty  string
	AddressEmpty  string
	CityEmpty  string
	CountryEmpty  string
	InvalidHash  string
	PermissionDenied   string
	InactiveUser  string
	NotFound  string

	InvalidFileType string
	InvalidFile string
	FileMaxSizeLimit string
	FileUploadLimit string
	ImageUrlEmpty  string
	ProductConnectionEmpty string
}

type EmailText struct {
	RegistrationConfirm string
	Welcome string
	SubjectConfirmYourRegistration string
	ConfirmYourRegistrationEmailText string
	SubjectResetPasswordRRequest  string
	ResetPasswordRRequestEmailText string
}

type MessageText struct {
	ThankYou string
	LoggedOut string
	LoggedIn  string
	RegisteredSuccessfully string
	ConfirmationLinkSentToYourEmail string
	UserConfirmedSuccessfully  string
	PasswordUpdatedSuccessfully string
	ProductImageAdded string
	ProductImageRemoved string
}


// // // //

type Localization struct {
	ErrorText ErrorText
	EmailText EmailText
	MessageText MessageText
}

func NewLocalization(lang string) (Localization,  error) {
	local := Localization{}
	filename := "./localization/" + lang + ".json"
	file, err := os.Open(filename)
	if err != nil {
		filename = "./localization/" + lang + ".json"
		file, err = os.Open(filename)
		if err != nil {
			return local, err
		}
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&local)
	if err != nil {  return local, err }

	return local,  nil
}

func (l *Localization) GetErrorValue(field string) string {
	v := l.ErrorText
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		return f.String()
	}

	return field
}

func (l *Localization) GetEmailTextValue(field string) string {
	v := l.EmailText
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		return f.String()
	}

	return field
}

func (l *Localization) GetMessageTextValue(field string) string {
	v := l.MessageText
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		return f.String()
	}

	return field
}