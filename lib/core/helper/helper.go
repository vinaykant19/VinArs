package helper

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strconv"
	"unsafe"
)
/*
 * helper functions using by base framework
 * @author: Vinaykant (vinaykantsahu@gmail.com)
 */
func CheckTypeByStringValue(val string, chkType string) bool {
	var pattern = regexp.MustCompile(``)
	switch chkType {
	case "byte":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "int":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "int8":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "int16":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "int32":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "int64":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "uint":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "uint8":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "uint16":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "uint32":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "uint64":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "uintptr":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "rune":
		pattern = regexp.MustCompile(`^[0-9]+$`)
		break
	case "float32":
		pattern = regexp.MustCompile(`^[+-]?([0-9]*[.])?[0-9]+`)
		break
	case "float64":
		pattern = regexp.MustCompile(`^[+-]?([0-9]*[.])?[0-9]+`)
		break
	case "bool":
		pattern = regexp.MustCompile(`^Yes|yes|YES|y|Y|No|no|NO|n|N|true|false|TRUE|FALSE$`)
		break
	case "string":
		return true // no need to check ;)
		break
	//case "complex128":
	//	pattern = regexp.MustCompile(`^([-]?[0-9]+\\.?[0-9]?)([-|+]+[0-9]+\\.?[0-9]?)[i$]+$`)
	//	break
	default:
		panic("Not supported type passed for API end point")
	}

	return pattern.MatchString(val)
}

func ToByte(val string) []byte {
	return *(*[]byte)(unsafe.Pointer(&val))
}
func ToIntiger(val string) int {
	newVal, _ :=  strconv.Atoi(val)
	return newVal
}
func ToIntiger8(val string) int8 {
	newVal, _ :=  strconv.ParseInt(val, 10, 8)
	return int8(newVal)
}
func ToIntiger16(val string) int16 {
	newVal, _ :=  strconv.ParseInt(val, 10, 16)
	return int16(newVal)
}
func ToIntiger32(val string) int32 {
	newVal, _ :=  strconv.ParseInt(val, 10, 32)
	return int32(newVal)
}
func ToIntiger64(val string) int64 {
	newVal, _ :=  strconv.ParseInt(val, 10, 64)
	return newVal
}
func ToUintiger(val string) uint {
	newVal, _ :=  strconv.Atoi(val)
	return uint(newVal)
}
func ToUintiger8(val string) uint8 {
	newVal, _ :=  strconv.ParseUint(val, 10, 8)
	return uint8(newVal)
}
func ToUintiger16(val string) uint16 {
	newVal, _ :=  strconv.ParseUint(val, 10, 16)
	return uint16(newVal)
}
func ToUintiger32(val string) uint32 {
	newVal, _ :=  strconv.ParseUint(val, 10, 32)
	return uint32(newVal)
}
func ToUintiger64(val string) uint64 {
	newVal, _ :=  strconv.ParseUint(val, 10, 64)
	return uint64(newVal)
}
func ToUintptr(val string) uintptr {
	return uintptr(ToIntiger(val))
}
func ToRune(val string) rune {
	return rune(ToIntiger(val))
}

func ToFloatVal32(val string) float32 {
	newVal, _ :=  strconv.ParseFloat(val, 32)
	return float32(newVal)
}
func ToFloatVal64(val string) float64 {
	newVal, _ :=  strconv.ParseFloat(val, 64)
	return newVal
}
func ToBoolean(val string) bool {
	newVal, _ := strconv.ParseBool(val)
	return newVal
}
//func ToComplex128(val string) complex128 {
//	newVal, _ :=  strconv.
//}

func BytesToString(data []byte) string {
	return string(data[:])
}

func MapDbType(typeExpr ast.Expr, length string) (string, string)  {
	dataType := ""
	newLength := ""

	//dataType, _ = typeExpr.(*ast.SelectorExpr)
	//newLength = "10"
	//return dataType, newLength
	//if IsCustomType(typeExpr) {
	//	dataType = "INTEGER"
	//	newLength = "10"
	//	fmt.Println("here ")
	//}

	switch typeExpr.(*ast.Ident).Name {
		case "byte":
			dataType = "BIT"
			newLength = ""
			break;
		case "int8":
			dataType = "TINYINT"
			newLength = ""
			break;
		case "uint8":
			dataType = "TINYINT"
			newLength = ""
			break;
		case "int16":
			dataType = "SMALLINT"
			newLength = ""
			break;
		case "uint16":
			dataType = "SMALLINT"
			newLength = ""
			break;
		case "int32":
			dataType = "MEDIUMINT"
			newLength = "11"
			break;
		case "uint32":
			dataType = "MEDIUMINT"
			newLength = "11"
			break;
		case "int64":
			dataType = "BIGINT"
			newLength = "100"
		case "uint64":
			dataType = "BIGINT"
			newLength = "100"
			break;
		case "int":
			dataType = "INTEGER"
			newLength = "10"
			break;
		case "uint":
			dataType = "INTEGER"
			newLength = "10"
			break;
		case "uintptr":
			dataType = "INTEGER"
			newLength = "10"
			break;
		case "rune":
			dataType = "INTEGER"
			newLength = "10"
			break;
		case "string":
			dataType = "VARCHAR"
			newLength = "255"
			break;
		case "bool":
			dataType = "bool"
			newLength = ""
			break;
		case "float32":
			dataType = "float"
			newLength = "10,2"
			break;
		case "float64":
			dataType = "float"
			newLength = "16,4"
			break;
	default:
		fmt.Println(typeExpr.(*ast.Ident).Name)
		fmt.Println(typeExpr.(*ast.Ident).IsExported())
		break
	}

	if length != "" {
		newLength = length
	}
	return dataType, newLength
}

func IsCustomType(t types.Type) bool {
	switch x := t.(type) {
	case *types.Basic:
		return false
	case *types.Slice:
		return false
	case *types.Map:
		return false
	case *types.Pointer:
		return IsCustomType(x.Elem())
	default:
		return true
	}
}

func GenerateHashedPassword(pwd string) string {

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func VerifyHashedPassword(hash string, orgPwd string) bool  {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(orgPwd))
	if err == nil {
		return true
	}

	return false
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func StringArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}