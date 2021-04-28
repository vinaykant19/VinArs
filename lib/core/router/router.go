package router

import (
	authToken "../token"
	cnf "../../../configuration"
	customeResponse "../response"
	helper "../helper"
	readcontrol "../../../src/controller/read"
	strCore "../strings"
	writecontrol "../../../src/controller/write"
	"../database"
	"../../../localization"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)
/*
 * Router functions using by base framework, it is helping to select right function based on the URL
 * @author: Vinaykant (vinaykantsahu@gmail.com)
 */
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type Routes struct {
	Methods []string
}

//
//  returns a new instance of controller.
//
func Call(response http.ResponseWriter, request *http.Request, allCalls map[string]map[string][]string, configuration *cnf.Configuration) {
	isApiKey := false
	AuthorizationToken  := ""
	DeviceToken  := ""
	DeviceType  := ""
	Language := "english"
	// Loop over header names
	for name, values := range request.Header {
		// Loop over all values for the name.
		for _, value := range values {
			if name == "X-Api-Key" {
				isApiKey = helper.VerifyHashedPassword(value, configuration.API_KEY)
			}
			if name == "Devicetoken" {
				DeviceToken = value
			}
			if name == "Devicetype" {
				DeviceType = value
			}
			if name == "Authorization" {
				AuthorizationToken = value[7:len(value)]
			}
			if name == "Localization" {
				Language = value
			}
		}
	}
	local, _ := localization.NewLocalization(Language)

	returnResult, returnStatusCode := customeResponse.ErrorCodeResponse("NotFound", local)

	urlParts := strings.Split(request.URL.String(), "/")
	callName := urlParts[1]
	if strings.Trim(callName, " ") == "" {
		callName = "Home"
	}

	if isApiKey == false {
		returnResult, returnStatusCode = customeResponse.ErrorCodeResponse("InvalidAuthorization", local)
		response.WriteHeader(returnStatusCode)
		response.Header().Add("Content-Type", configuration.Default_Response_type)

		fmt.Fprintf(response, returnResult)
		return
	}
	method := request.Method
	functionName := strCore.ToCamelCase(callName);
	db := database.Database{}
	fmt.Println(method, callName, functionName, request.URL.String(), request.UserAgent(), DeviceType, DeviceToken)
	if method == "GET" {
		RController := readcontrol.ReadController{
			Response: &response,
			Request: request,
			Token:	AuthorizationToken,
			ResponseType: configuration.Default_Response_type,
			Configuration: configuration,
			DB: &db,
			Logger: nil,
			DeviceToken: DeviceToken,
			DeviceType: DeviceType,
			Language: local,
		}
		allowWithoutToken := strings.Split(configuration.Non_Token_Calls_GET, ",")

		if _, ok := allCalls["read"][functionName]; ok {
			db.Connect(configuration)
			userId, _ := authToken.ValidateTokenAndGetUserId(&db, local, configuration, AuthorizationToken)
			RController.UserId = userId
			if !helper.StringArrayContains(allowWithoutToken, functionName) && userId == 0 {
				returnResult, returnStatusCode = customeResponse.ErrorCodeResponse("InvalidAuthorization", local)
			} else {
				RController.ResponseType = configuration.Default_Response_type
				funcCall := reflect.ValueOf(&RController).MethodByName(functionName)
				arg := make([]reflect.Value, funcCall.Type().NumIn())
				status, newArg := matchInputParam(funcCall.Type(), urlParts, allCalls["read"][functionName], arg, configuration)
				if status {
					RController.StatusCode = 200;
					funcCall.Call(newArg)
					returnResult = RController.Result
					returnStatusCode = RController.StatusCode
				}
			}
		}

		response.Header().Add("Content-Type", RController.ResponseType)
	} else {
		WController := writecontrol.WriteController{
			Response:      &response,
			Request:       request,
			Token:         AuthorizationToken,
			ResponseType:  configuration.Default_Response_type,
			Configuration: configuration,
			DB:            &db,
			Logger:        nil,
			DeviceToken: DeviceToken,
			DeviceType: DeviceType,
			Language: local,
		}
		allowWithoutToken := strings.Split(configuration.Non_Token_Calls_POST, ",")
		if _, ok := allCalls["write"][functionName]; ok {
			db.Connect(configuration)
			userId, _ := authToken.ValidateTokenAndGetUserId(&db, local, configuration, AuthorizationToken)
			WController.UserId = userId
			if !helper.StringArrayContains(allowWithoutToken, functionName) && userId == 0 {
				returnResult, returnStatusCode = customeResponse.ErrorCodeResponse("InvalidAuthorization", local)
			} else {
				WController.ResponseType = configuration.Default_Response_type
				funcCall := reflect.ValueOf(&WController).MethodByName(functionName)
				arg := make([]reflect.Value, funcCall.Type().NumIn())
				status, newArg := matchInputParam(funcCall.Type(), urlParts, allCalls["write"][functionName], arg, configuration)
				if status {
					WController.StatusCode = 200;
					funcCall.Call(newArg)

					returnResult = WController.Result
					returnStatusCode = WController.StatusCode
				}
			}
		}

		response.Header().Add("Content-Type", WController.ResponseType)
	}
	if db.IsLive() {
		db.Transaction("commit")
		db.Disconnect()
	}
	//
	response.WriteHeader(returnStatusCode)

	fmt.Fprintf(response, returnResult)
}

func matchInputParam(f reflect.Type, urlParts []string, paramNames []string, arg []reflect.Value, configuration *cnf.Configuration) (bool, []reflect.Value) {
	if configuration.Url_Param_With_Variable == "Yes" {
		return matchInputParamWithVariableName(f, urlParts, paramNames, arg)
	}

	totalFuncParam := f.NumIn()
	totalUrlParam := len(urlParts) - 2
	if (totalFuncParam != totalUrlParam) {
		return false, arg
	}

	for i := 0; i < totalFuncParam; i++ {
		j := i + 2
		fmt.Printf(" %d. %v  : %v (%v) \n", i,  f.In(i), urlParts[j], reflect.TypeOf(urlParts[j]).Kind())
		if ! helper.CheckTypeByStringValue(urlParts[j], f.In(i).String()) {
			return false, arg
		}
		//change the type of value
		switch f.In(i).String() {
		case "byte":
			arg[i] = reflect.ValueOf(helper.ToByte(urlParts[j]))
			break
		case "int":
			arg[i] = reflect.ValueOf(helper.ToIntiger(urlParts[j]))
			break
		case "int8":
			arg[i] = reflect.ValueOf(helper.ToIntiger8(urlParts[j]))
			break
		case "int16":
			arg[i] = reflect.ValueOf(helper.ToIntiger16(urlParts[j]))
			break
		case "int32":
			arg[i] = reflect.ValueOf(helper.ToIntiger32(urlParts[j]))
			break
		case "int64":
			arg[i] = reflect.ValueOf(helper.ToIntiger64(urlParts[j]))
			break
		case "uint":
			arg[i] = reflect.ValueOf(helper.ToUintiger(urlParts[j]))
			break
		case "uint8":
			arg[i] = reflect.ValueOf(helper.ToUintiger8(urlParts[j]))
			break
		case "uint16":
			arg[i] = reflect.ValueOf(helper.ToUintiger16(urlParts[j]))
			break
		case "uint32":
			arg[i] = reflect.ValueOf(helper.ToUintiger32(urlParts[j]))
			break
		case "uint64":
			arg[i] = reflect.ValueOf(helper.ToUintiger64(urlParts[j]))
			break
		case "uintptr":
			arg[i] = reflect.ValueOf(helper.ToUintptr(urlParts[j]))
			break
		case "rune":
			arg[i] = reflect.ValueOf(helper.ToRune(urlParts[j]))
			break
		case "float32":
			arg[i] = reflect.ValueOf(helper.ToFloatVal32(urlParts[j]))
			break
		case "float64":
			arg[i] = reflect.ValueOf(helper.ToFloatVal64(urlParts[j]))
			break
		case "bool":
			arg[i] = reflect.ValueOf(helper.ToBoolean(urlParts[j]))
			break
		case "string":
			arg[i] = reflect.ValueOf(urlParts[j]) // no need to check ;)
			break
		//case "complex128":
		//	pattern = regexp.MustCompile(`^([-]?[0-9]+\\.?[0-9]?)([-|+]+[0-9]+\\.?[0-9]?)[i$]+$`)
		//	break
		default:
			panic("Not supported type passed for API end point")
		}
	}

	return true, arg
}

func matchInputParamWithVariableName(f reflect.Type, urlParts []string, paramNames []string, arg []reflect.Value) (bool, []reflect.Value) {
	totalFuncParam := f.NumIn()
	totalUrlParam := len(urlParts) - 2
	if (totalFuncParam * 2 != totalUrlParam) {
		return false, arg
	}
	k := 0
	for i := 2; i < len(urlParts); i = i + 2 {
		j := i + 1
		if paramNames[k] != urlParts[i] { //strict check for url
			return false, arg
		}
		if ! helper.CheckTypeByStringValue(urlParts[j], f.In(k).String()) {
			return false, arg
		}

		//change the type of value
		switch f.In(k).String() {
		case "byte":
			arg[k] = reflect.ValueOf(helper.ToByte(urlParts[j]))
			break
		case "int":
			arg[k] = reflect.ValueOf(helper.ToIntiger(urlParts[j]))
			break
		case "int8":
			arg[k] = reflect.ValueOf(helper.ToIntiger8(urlParts[j]))
			break
		case "int16":
			arg[k] = reflect.ValueOf(helper.ToIntiger16(urlParts[j]))
			break
		case "int32":
			arg[k] = reflect.ValueOf(helper.ToIntiger32(urlParts[j]))
			break
		case "int64":
			arg[k] = reflect.ValueOf(helper.ToIntiger64(urlParts[j]))
			break
		case "uint":
			arg[k] = reflect.ValueOf(helper.ToUintiger(urlParts[j]))
			break
		case "uint8":
			arg[k] = reflect.ValueOf(helper.ToUintiger8(urlParts[j]))
			break
		case "uint16":
			arg[k] = reflect.ValueOf(helper.ToUintiger16(urlParts[j]))
			break
		case "uint32":
			arg[k] = reflect.ValueOf(helper.ToUintiger32(urlParts[j]))
			break
		case "uint64":
			arg[k] = reflect.ValueOf(helper.ToUintiger64(urlParts[j]))
			break
		case "uintptr":
			arg[k] = reflect.ValueOf(helper.ToUintptr(urlParts[j]))
			break
		case "rune":
			arg[k] = reflect.ValueOf(helper.ToRune(urlParts[j]))
			break
		case "float32":
			arg[k] = reflect.ValueOf(helper.ToFloatVal32(urlParts[j]))
			break
		case "float64":
			arg[k] = reflect.ValueOf(helper.ToFloatVal64(urlParts[j]))
			break
		case "bool":
			arg[k] = reflect.ValueOf(helper.ToBoolean(urlParts[j]))
			break
		case "string":
			arg[k] = reflect.ValueOf(urlParts[j]) // no need to check ;)
			break
		//case "complex128":
		//	pattern = regexp.MustCompile(`^([-]?[0-9]+\\.?[0-9]?)([-|+]+[0-9]+\\.?[0-9]?)[i$]+$`)
		//	break
		default:
			panic("Not supported type passed for API end point")
		}
		//fmt.Printf("%d. func Arg: %v func arg Type: %v - param Name: %v Param value: %v param value type: %v \n",
		//	k,  paramNames[k], f.In(k), urlParts[i], urlParts[j], reflect.TypeOf(urlParts[j]).Kind())
		k++
	}

	return true, arg
}

//
// callPerser
//
func CallParser(callType string, allCalls map[string][]string )  map[string][]string {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, basepath + "/../../../src/controller/" + callType, nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
	}
	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					paramArray := []string{}
					for arg := range fn.Type.Params.List {
						paramArray = append(paramArray, fn.Type.Params.List[arg].Names[0].String())
					}
					allCalls[fn.Name.String()] = paramArray
				}
			}
		}
	}

	return allCalls
}

