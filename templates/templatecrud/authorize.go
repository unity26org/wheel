package templatecrud

var AuthorizeContent = `	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/{{ .EntityName.SnakeCasePlural }}(\/){0,1}.*\z` + "`" + `),
			Methods:   []string{"GET", "POST", "DELETE", "PUT"},
			UserRoles: []string{"public", "signed_in", "admin"},
		})`
