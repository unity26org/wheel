package config

var DatabasePath = []string{"config", "database.yml"}

var DatabaseContent = `adapter: {{ .Database }}
host: localhost
{{- if eq .Database "postgres" }}
port: 5432
{{- else if eq .Database "mysql" }}
port: 3306
{{- end }}
user: root
dbname: {{ .AppName }}
password: 
sslmode: disable
pool: 10`
