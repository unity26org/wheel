package routes

var AuthorizePath = []string{"routes", "authorize.go"}

var AuthorizeContent = `package routes

import (
	"regexp"
)

type Permission struct {
	UrlRegexp *regexp.Regexp
	Methods   []string
	UserRoles []string
}

var Permissions []Permission

func GrantPermission(url string, method string, userRole string) bool {
	for _, permission := range Permissions {
		if permission.UrlRegexp.MatchString(url) && checkItem(method, permission.Methods) && checkItem(userRole, permission.UserRoles) {
			return true
		}
	}

	return false
}

func checkItem(currrent string, availables []string) bool {
	for _, item := range availables {
		if item == currrent {
			return true
		}
	}

	return false
}

func init() {
	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/{0,1}\z` + "`" + `),
			Methods:   []string{"GET"},
			UserRoles: []string{"public", "signed_in", "admin"},
		})

	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/sessions\/(sign_in|sign_up|password)(\/){0,1}.*\z` + "`" + `),
			Methods:   []string{"GET", "POST", "DELETE"},
			UserRoles: []string{"public"},
		})

	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/sessions\/sign_out(\/){0,1}.*\z` + "`" + `),
			Methods:   []string{"DELETE"},
			UserRoles: []string{"admin", "signed_in"},
		})

	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/sessions\/refresh(\/){0,1}.*\z` + "`" + `),
			Methods:   []string{"POST"},
			UserRoles: []string{"admin", "signed_in"},
		})

	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/users(\/){0,1}.*\z` + "`" + `),
			Methods:   []string{"GET", "POST", "DELETE", "PUT"},
			UserRoles: []string{"admin"},
		})

	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/myself(\/){0,1}.*\z` + "`" + `),
			Methods:   []string{"GET", "POST", "DELETE", "PUT"},
			UserRoles: []string{"admin", "signed_in"},
		})

}`
