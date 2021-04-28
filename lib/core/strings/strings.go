package strings

import (
	"regexp"
	"strings"
)
/*
 * string functions using by base framework
 * @author: Vinaykant (vinaykantsahu@gmail.com)
 */
func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake  = matchAllCap.ReplaceAllString(snake, "${1}-${2}")

	return strings.ToLower(snake)
}

func ToCamelCase(str string) string {
	var pattern = regexp.MustCompile("(^[A-Za-z])|-([A-Za-z])")

	return pattern.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s,"-","",-1))
	})
}
