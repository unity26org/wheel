package templatecrud

var RoutesContent = `
	router.HandleFunc("/{{ .EntityName.SnakeCasePlural }}", handlers.{{ .EntityName.CamelCase }}List).Methods("GET")
	router.HandleFunc("/{{ .EntityName.SnakeCasePlural }}/{id}", handlers.{{ .EntityName.CamelCase }}Show).Methods("GET")
	router.HandleFunc("/{{ .EntityName.SnakeCasePlural }}", handlers.{{ .EntityName.CamelCase }}Create).Methods("POST")
	router.HandleFunc("/{{ .EntityName.SnakeCasePlural }}/{id}", handlers.{{ .EntityName.CamelCase }}Update).Methods("PUT")
	router.HandleFunc("/{{ .EntityName.SnakeCasePlural }}/{id}", handlers.{{ .EntityName.CamelCase }}Destroy).Methods("DELETE")`
