package config

var DatabasePath = []string{"config", "database.yml"}

var DatabaseContent = `host: localhost
port: 5432
user: root
dbname: {{ .AppName }}
password: Secret123!
sslmode: disable
pool: 10`
