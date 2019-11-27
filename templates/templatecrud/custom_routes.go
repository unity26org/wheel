package templatecrud

var CustomRoutesContent = `
  {{- range $element := .EntityColumns }}
	router.HandleFunc("/{{ $.EntityName.SnakeCasePlural }}/{{ $element.NameSnakeCase }}", handlers.{{ $.EntityName.CamelCase }}{{ $element.Name }}).Methods("GET")
  {{- end }}
`
