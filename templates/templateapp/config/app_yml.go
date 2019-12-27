package config

var AppPath = []string{"config", "app.yml"}

var AppContent = `app_name: "{{ .AppName }}"
app_repository: "{{ .AppRepository }}"
frontend_base_url: "https://example.com"
secret_key: "{{ .SecretKey }}"
reset_password_expiration_seconds: 172800
token_expiration_seconds: 7200
pagination:
  default: 20
  maximum: 50
locales:
  - "en"
  - "pt-BR"`
