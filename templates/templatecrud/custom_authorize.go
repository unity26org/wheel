package templatecrud

var CustomAuthorizeContent = `
{{- range $element := .EntityColumns }}  
	Permissions = append(Permissions,
		Permission{
			UrlRegexp: regexp.MustCompile(` + "`" + `\A\/{{ $.EntityName.SnakeCasePlural }}\/{{ $element.NameSnakeCase }}(\/){0,1}.*\z` + "`" + `),
			Methods:   []string{"GET"},
			UserRoles: []string{"public", "signed_in", "admin"},
	})
  {{- end }}`
