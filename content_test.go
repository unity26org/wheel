package main

import ()

var helpContent = `Usage:
  wheel new APP_PATH [options]             # Creates new app

  wheel generate SUBJECT NAME ATTRIBUTES   # Adds new CRUD to an existing app. 
                                           # SUBJECT: scaffold/model/entity/handler. 
                                           # NAME: name of the model, entity or handler
                                           # ATTRIBUTES: when not a handler, is a pair of column name
                                           # and column type separated by ":" i.e. description:string
                                           # Available types are: 
                                           # string/text/integer/decimal/datetime/bool/reference.
                                           # When a handler "attributes" are functions inside handler.
                                           
Options:
  -G, [--skip-git]                         # Skip .gitignore file

More:
  -h, [--help]                             # Show this help message and quit
  -v, [--version]                          # Show Wheel version number and quit`

var versionContent = `Wheel 1.0`

var gitIgnoreContent = `/log/*.log`

var mainContent = `package main

import (
	"flag"
	"net/http"
	"test_repository_hub.com/test_account/test_project/commons/app/model"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"test_repository_hub.com/test_account/test_project/config"
	"test_repository_hub.com/test_account/test_project/db/schema"
	"test_repository_hub.com/test_account/test_project/routes"
)

func main() {
	var mode string
	var port string
	var host string

	flag.StringVar(&mode, "mode", "server", "run mode (options: server/migrate)")
	flag.StringVar(&host, "host", "localhost", "http server host")
	flag.StringVar(&port, "port", "8081", "http server port")
	flag.Parse()

	log.Info.Println("starting app", config.App.AppName)

	model.Connect()

	if mode == "migrate" {
		schema.Migrate()
	} else if mode == "s" || mode == "server" {
		log.Fatal.Println(http.ListenAndServe(host+":"+port, routes.Routes(host, port)))
	} else {
		log.Fatal.Println("invalid run mode, please, use \"--help\" for more details")
	}
}`

var routesContentV1 = `package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"test_repository_hub.com/test_account/test_project/app/handlers"
	"test_repository_hub.com/test_account/test_project/commons/app/handler"
	"test_repository_hub.com/test_account/test_project/commons/log"
)

func Routes(host string, port string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// middlewares
	log.Info.Println("setting up middlewares")
	router.Use(loggingMiddleware)
	router.Use(authorizeMiddleware)

	log.Info.Println("setting up routes")
	log.Info.Println("listening on " + host + ":" + port + ", CTRL+C to stop")

	router.NotFoundHandler = http.HandlerFunc(handler.Error404)
	router.HandleFunc("/", handler.ApiRoot).Methods("GET")

	// sessions
	router.HandleFunc("/sessions/sign_in", handlers.SessionSignIn).Methods("POST")
	router.HandleFunc("/sessions/sign_out", handlers.SessionSignOut).Methods("DELETE")
	router.HandleFunc("/sessions/sign_up", handlers.SessionSignUp).Methods("POST")
	router.HandleFunc("/sessions/password", handlers.SessionPassword).Methods("POST")
	router.HandleFunc("/sessions/password", handlers.SessionRecovery).Methods("PUT")
	router.HandleFunc("/sessions/refresh", handlers.SessionRefresh).Methods("POST")

	// user
	router.HandleFunc("/myself", handlers.MyselfShow).Methods("GET")
	router.HandleFunc("/myself", handlers.MyselfUpdate).Methods("PUT")
	router.HandleFunc("/myself/password", handlers.MyselfUpdatePassword).Methods("PUT")
	router.HandleFunc("/myself", handlers.MyselfDestroy).Methods("DELETE")

	// admin
	router.HandleFunc("/users", handlers.UserList).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.UserShow).Methods("GET")
	router.HandleFunc("/users", handlers.UserCreate).Methods("POST")
	router.HandleFunc("/users/{id}", handlers.UserUpdate).Methods("PUT")
	router.HandleFunc("/users/{id}/password", handlers.UserUpdatePassword).Methods("PUT")
	router.HandleFunc("/users/{id}", handlers.UserDestroy).Methods("DELETE")

	return router
}`

var middlewareContent = `package routes

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"test_repository_hub.com/test_account/test_project/app/handlers"
	"test_repository_hub.com/test_account/test_project/app/user"
	"test_repository_hub.com/test_account/test_project/commons/app/handler"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"test_repository_hub.com/test_account/test_project/db/entities"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info.Println(r.Method + ": " + filterUrlValues(r.URL.Path, r.URL.Query()) + " for " + r.RemoteAddr)

		if strings.Trim(r.Header.Get("Content-Type"), " \n") == "application/json" {
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				log.Info.Println("Body JSON:", filterJsonValues(string(body)))
				// put the body content back
				r.Body = ioutil.NopCloser(strings.NewReader(string(body)))
			} else {
				log.Error.Println("loggingMiddlware: ", err)
			}
		} else {
			r.ParseMultipartForm(100 * 1024)
			log.Info.Println("Form-data: " + filterFormValues(r.Form))
		}

		next.ServeHTTP(w, r)
	})
}

func authorizeMiddleware(next http.Handler) http.Handler {
	var userId uint
	var err error
	var userRole string
	var signedInUser entities.User

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId = 0
		err = nil
		userRole = "public"

		userId, err = checkToken(r.Header.Get("Authorization"))
		if err == nil {
			if signedInUser, err = checkSignedInUser(userId); err == nil {
				userRole = "signed_in"
			}

			if userRole == "signed_in" && checkAdminUser(user.Current) {
				userRole = "admin"
			}
		}

		if GrantPermission(r.RequestURI, r.Method, userRole) {
			next.ServeHTTP(w, r)
		} else if userRole == "public" {
			handler.Error401(w, r)
		} else {
			handler.Error403(w, r)
		}
	})
}

func checkAdminUser(signedInUser entities.User) bool {
	if signedInUser.Admin {
		log.Info.Println("admin access granted to user")
		return true
	} else {
		log.Info.Println("admin access denied to user")
		return false
	}
}

func checkSignedInUser(userId uint) (entities.User, error) {
	log.Info.Printf("checking user id: %d...\n", userId)

	if userId == 0 {
		log.Info.Println("user is not available")
		return entities.User{}, errors.New("user is not available")
	} else if signedInUser, err := user.Find(userId); err == nil {
		user.SetCurrent(userId)
		log.Info.Println("user was found")
		return signedInUser, nil
	} else {
		log.Info.Println("user was not found")
		return signedInUser, errors.New("user was not found")
	}
}

func checkToken(token string) (uint, error) {
	log.Info.Println("checking token...")

	if token == "" {
		log.Info.Println("token was not sent")
		return 0, nil
	} else {
		return validateToken(token)
	}
}

func validateToken(token string) (uint, error) {
	log.Info.Println("validating token...")

	userId, err := handlers.SessionCheck(token)
	if err == nil {
		log.Info.Println("token is valid")
		return userId, nil
	} else {
		log.Info.Println("invalid token")
		return 0, errors.New("invalid token")
	}
}

func filterParamsValues(queries map[string][]string) map[string][]string {
	var filter = regexp.MustCompile(` + "`" + `(?i)(password)|(token)` + "`" + `)
	queries_filtered := make(map[string][]string)

	for key := range queries {
		if filter.MatchString(key) {
			queries_filtered[key] = []string{"[FILTERED]"}
		} else {
			queries_filtered[key] = []string{}
			for _, element := range queries[key] {
				queries_filtered[key] = append(queries_filtered[key], element)
			}
		}

	}

	return queries_filtered
}

func filterUrlValues(path string, queries map[string][]string) string {
	var firstParam = true
	queries_filtered := filterParamsValues(queries)

	for key := range queries_filtered {
		if firstParam {
			path = path + "?"
			firstParam = false
		} else {
			path = path + "&"
		}

		path = path + key + "=" + strings.Join(queries_filtered[key], " ")
	}

	return path
}

func filterFormValues(queries map[string][]string) string {
	var buffer bytes.Buffer
	var index int
	queries_filtered := filterParamsValues(queries)

	index = 0
	buffer.WriteString("{ ")

	for key := range queries_filtered {
		buffer.WriteString("\"")
		buffer.WriteString(key)
		buffer.WriteString("\": \"")

		buffer.WriteString(strings.Join(queries_filtered[key], " "))
		buffer.WriteString("\"")

		if (index + 1) != len(queries_filtered) {
			buffer.WriteString(", ")
		}

		index++
	}

	buffer.WriteString(" }")

	return buffer.String()
}

func filterJsonValues(inputJson string) string {
	type Point struct {
		StartAt int
		EndAt   int
	}

	var filter = regexp.MustCompile(` + "`" + `(?i)(password)|(token)` + "`" + `)
	var regexpWhiteSpace = regexp.MustCompile(` + "`" + `[\s\t\n]{1}` + "`" + `)
	var stack []string
	var key, currentChar string
	var valueStartAt, valueEndAt, keyStartAt, keyEndAt int
	var points []Point

	substring := []rune(inputJson)
	isCharBeforeEscape := false
	isKey := false
	isValue := false

	for i := 0; i < len(inputJson); i++ {
		currentChar = string(inputJson[i])

		if regexpWhiteSpace.MatchString(currentChar) && (len(stack) == 0 || stack[len(stack)-1] != ` + "`" + `"` + "`" + ` || stack[len(stack)-1] == ` + "`" + `:` + "`" + `) {
			currentChar = currentChar
		} else if currentChar == "{" && (len(stack) == 0 || stack[len(stack)-1] == "{") {
			stack = append(stack, "{")
		} else if currentChar == ":" && (len(stack) == 0 || stack[len(stack)-1] == "{") {
			stack = append(stack, ":")
		} else if currentChar == "}" && (len(stack) > 0 && stack[len(stack)-1] == "{") {
			stack = stack[:len(stack)-1]
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ":") && !isCharBeforeEscape {
			stack = stack[:len(stack)-1]
			stack = append(stack, ` + "`" + `"` + "`" + `)
			isValue = true
			isKey = false
			valueStartAt = i + 1
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ` + "`" + `"` + "`" + `) && isValue && !isCharBeforeEscape {
			stack = stack[:len(stack)-1]
			isValue = false
			isKey = false
			valueEndAt = i

			if filter.MatchString(key) {
				points = append(points, Point{StartAt: valueStartAt, EndAt: valueEndAt})
			}
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == "{") && !isCharBeforeEscape {
			isKey = true
			isValue = false
			stack = append(stack, ` + "`" + `"` + "`" + `)
			keyStartAt = i + 1
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ` + "`" + `"` + "`" + `) && isKey && !isCharBeforeEscape {
			isKey = false
			isValue = false
			stack = stack[:len(stack)-1]
			keyEndAt = i
			key = string(substring[keyStartAt:keyEndAt])
		}

		if currentChar == ` + "`" + `\` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ` + "`" + `"` + "`" + `) {
			isCharBeforeEscape = true
		} else if isCharBeforeEscape {
			isCharBeforeEscape = false
		}
	}

	if len(points) > 0 {
		for i := len(points) - 1; i >= 0; i-- {
			inputJson = string(substring[0:points[i].StartAt]) + "[FILTERED]" + string(substring[points[i].EndAt:len(inputJson)])
			substring = []rune(inputJson)
		}
	}

	return inputJson
}`

var authorizeContentV1 = `package routes

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

var localePtBrContent = `# encoding: UTF-8
welcome: 'Bem vindo(a)'
password_recovery_instructions: 'Instruções de Recuperação de Senha'`

var localeEnContent = `# encoding: UTF-8
welcome: 'Welcome'
password_recovery_instructions: 'Password Recovery Instructions'`

var configAppContent = `app_name: "test_project"
app_repository: "test_repository_hub.com/test_account/test_project"
frontend_base_url: "https://example.com"
secret_key: "0B7f3892773a0be1E4C6c992f9D4BcdB96Bf2277665133d8CdE18FfED0ECcF7d4a5e7114d4c5fDD3430ad9d87A1C705aACfBCBBD928B38248aBd48Ff8bDdA5E0"
reset_password_expiration_seconds: 172800
token_expiration_seconds: 7200
pagination:
  default: 20
  maximum: 50
locales:
  - "en"
  - "pt-BR"`

var configConfigContent = `package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"test_repository_hub.com/test_account/test_project/commons/log"
)

type Pagination struct {
	Default int ` + "`" + `yaml:"default"` + "`" + `
	Maximum int ` + "`" + `yaml:"maximum"` + "`" + `
}

type AppConfig struct {
	AppName                        string     ` + "`" + `yaml:"app_name"` + "`" + `
	AppRepository                  string     ` + "`" + `yaml:"app_repository"` + "`" + `
	SecretKey                      string     ` + "`" + `yaml:"secret_key"` + "`" + `
	ResetPasswordExpirationSeconds int        ` + "`" + `yaml:"reset_password_expiration_seconds"` + "`" + `
	ResetPasswordUrl               string     ` + "`" + `yaml:"reset_password_url"` + "`" + `
	TokenExpirationSeconds         int        ` + "`" + `yaml:"token_expiration_seconds"` + "`" + `
	Pagination                     Pagination ` + "`" + `yaml:"pagination"` + "`" + `
	Locales                        []string   ` + "`" + `yaml:"locales"` + "`" + `
}

var App AppConfig

func readAppConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/app.yml")
	if err != nil {
		log.Error.Fatal(err)
	}

	return data
}

func init() {
	err := yaml.Unmarshal(readAppConfigFile(), &App)
	if err != nil {
		log.Error.Fatalf("error: %v", err)
	}
}`

var configDatabaseContent = `host: localhost
port: 5432
user: root
dbname: test_project
password: Secret123!
sslmode: disable
pool: 10`

var configEmailContent = `user: 'no-reply0@example.com'
name: 'No Reply'
password: 'Secret123!'
address: 'smtp.example.com'
port: 25`

var mailerContent = `package mailer

import (
	"crypto/tls"
	"fmt"
	"github.com/adilsonchacon/blog/commons/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	User     string
	Name     string
	Password string
	Address  string
	Port     int
}

var (
	from       mail.Address
	to         []mail.Address
	cc         []mail.Address
	bcc        []mail.Address
	validEmail = regexp.MustCompile(` + "`" + `\A[^@]+@([^@\.]+\.)+[^@\.]+\z` + "`" + `)
)

func SetFrom(name string, tEmail string) {
	if validEmail.MatchString(tEmail) {
		from = mail.Address{name, tEmail}
	} else {
		log.Error.Println("email is invalid")
	}
}

func AddTo(name string, tEmail string) {
	if validEmail.MatchString(tEmail) {
		to = append(to, mail.Address{name, tEmail})
	} else {
		log.Error.Println("email is invalid")
	}
}

func AddCc(name string, tEmail string) {
	if validEmail.MatchString(tEmail) {
		cc = append(cc, mail.Address{name, tEmail})
	} else {
		log.Error.Println("email is invalid")
	}
}

func AddBcc(name string, tEmail string) {
	if validEmail.MatchString(tEmail) {
		bcc = append(bcc, mail.Address{name, tEmail})
	} else {
		log.Error.Println("email is invalid")
	}
}

func ResetReceipts() {
	SetFrom("", "")
	to = to[:0]
	cc = cc[:0]
	bcc = bcc[:0]
}

func Send(subject string, body string, html bool) {
	receipts := Receipts()
	headers := make(map[string]string)
	config, err := loadConfigFile()
	if err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	if from.String() == "<@>" {
		from = mail.Address{config.Name, config.User}
	}

	// Setup header
	headers["From"] = from.String()

	if len(to) > 0 {
		headers["To"] = stringfy("to")
	}

	if len(cc) > 0 {
		headers["Cc"] = stringfy("cc")
	}

	if len(bcc) > 0 {
		headers["Bcc"] = stringfy("bcc")
	}

	if html {
		headers["MIME-version"] = "1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	}

	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := config.Address + ":" + strconv.Itoa(config.Port)
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", config.User, config.Password, host)

	err = checkHost(host)
	if err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	// Auth
	if err = client.Auth(auth); err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	// To
	if err = client.Mail(from.Address); err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	// To, Cc and Bcc
	for i := 0; i < len(receipts); i++ {
		if err = client.Rcpt(receipts[0]); err != nil {
			log.Error.Println("mailer.Send", err)
			return
		}
	}

	// Data
	socket, err := client.Data()
	if err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	_, err = socket.Write([]byte(message))
	if err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	err = socket.Close()
	if err != nil {
		log.Error.Println("mailer.Send", err)
		return
	}

	client.Quit()
}

func stringfy(tType string) string {
	var isTo = regexp.MustCompile(` + "`" + `(?i)\Ato\z` + "`" + `)
	var isCc = regexp.MustCompile(` + "`" + `(?i)\Acc\z` + "`" + `)
	var isBcc = regexp.MustCompile(` + "`" + `(?i)\Abcc\z` + "`" + `)
	var isFrom = regexp.MustCompile(` + "`" + `(?i)\Afrom\z` + "`" + `)
	var auxEmail []mail.Address
	var auxString []string

	if isTo.MatchString(tType) {
		auxEmail = to
	} else if isCc.MatchString(tType) {
		auxEmail = cc
	} else if isBcc.MatchString(tType) {
		auxEmail = bcc
	} else if isFrom.MatchString(tType) {
		auxEmail = append(auxEmail, from)
	} else {
		log.Error.Println("Invalid param for stringfy, available params are to, cc, bcc and from")
	}

	for i := 0; i < len(auxEmail); i++ {
		auxString = append(auxString, auxEmail[i].Address)
	}

	return strings.Join(auxString, "; ")
}

func Receipts() []string {
	var receipts []string

	for i := 0; i < len(to); i++ {
		receipts = append(receipts, to[0].Address)
	}

	for i := 0; i < len(cc); i++ {
		receipts = append(receipts, cc[0].Address)
	}

	for i := 0; i < len(bcc); i++ {
		receipts = append(receipts, bcc[0].Address)
	}

	return receipts
}

func loadConfigFile() (Config, error) {
	config := Config{}

	content, err := readConfigFile()
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func readConfigFile() ([]byte, error) {
	data, err := ioutil.ReadFile("./config/email.yml")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func checkHost(host string) error {
	var err error

	ip := net.ParseIP(host)
	if ip == nil {
		_, err = net.LookupHost(host)
		if err != nil {
			return err
		}
	} else {
		_, err = net.LookupAddr(host)
		if err != nil {
			return err
		}
	}

	return nil
}`

var commonsLogContent = `package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

var writer = setWriter()

var (
	Debug = log.New(writer, fmt.Sprintf("\033[94m[DEBUG]\033[39m "), log.LstdFlags)
	Info  = log.New(writer, fmt.Sprintf("\033[36m[INFO]\033[39m "), log.LstdFlags)
	Warn  = log.New(writer, fmt.Sprintf("\033[93m[WARN]\033[39m "), log.LstdFlags)
	Error = log.New(writer, fmt.Sprintf("\033[91m[ERROR]\033[39m "), log.LstdFlags)
	Fatal = log.New(writer, fmt.Sprintf("\033[90m[FATAL]\033[39m "), log.LstdFlags)
)

func setWriter() io.Writer {
	os.MkdirAll("./log", 0755)

	fileHandler, err := os.OpenFile("./log/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return io.MultiWriter(os.Stdout, fileHandler)
}`

var commonsLocaleContent = `package locale

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"test_repository_hub.com/test_account/test_project/commons/log"
)

type Keys struct {
	Welcome                      string ` + "`" + `yaml:"welcome"` + "`" + `
	PasswordRecoveryInstructions string ` + "`" + `yaml:"password_recovery_instructions"` + "`" + `
}

var I18n Keys

func Load(locale string) {
	err := yaml.Unmarshal(readLocaleFile(locale), &I18n)
	if err != nil {
		log.Error.Fatal(err)
	}
}

func readLocaleFile(locale string) []byte {
	data, err := ioutil.ReadFile("./config/locales/" + locale + ".yml")
	if err != nil {
		log.Error.Fatal(err)
	}

	return data
}`

var commonsCryptoContent = `package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"time"
)

func SetPassword(password string) string {
	var err error

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error.Println(err)
	}

	return string(hash)
}

func CheckPassword(password string, hashedPassword string) bool {
	var err error

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func RandString(size int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func EncryptText(clearText string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(clearText))

	return hex.EncodeToString(mac.Sum(nil))
}`

var commonsConversorContent = `package conversor

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StringToInterface(contentType string, contentValue string) (interface{}, error) {
	var regexpBooleanType = regexp.MustCompile(` + "`" + `bool` + "`" + `)
	var regexpIntType = regexp.MustCompile(` + "`" + `int` + "`" + `)
	var regexpFloatType = regexp.MustCompile(` + "`" + `float|double|decimal` + "`" + `)
	var regexpDateTimeType = regexp.MustCompile(` + "`" + `time|date` + "`" + `)
	var returnInterface interface{}
	var err error

	if regexpIntType.MatchString(contentType) {
		returnInterface, err = StringToInt(contentValue)
	} else if regexpFloatType.MatchString(contentType) {
		returnInterface, err = StringToFloat(contentValue)
	} else if regexpDateTimeType.MatchString(contentType) {
		returnInterface, err = time.Parse(time.RFC3339, contentValue)
	} else if regexpBooleanType.MatchString(contentType) {
		returnInterface, err = StringToBoolean(contentValue)
	} else {
		returnInterface = contentValue
		err = nil
	}

	return returnInterface, err
}

func StringToInt(valueContent string) (interface{}, error) {
	return strconv.ParseInt(valueContent, 10, 64)
}

func StringToFloat(valueContent string) (interface{}, error) {
	return strconv.ParseFloat(valueContent, 64)
}

func StringToBoolean(valueContent string) (interface{}, error) {
	var checkFalse = regexp.MustCompile(` + "`" + `\A(0|f|false|no)\z` + "`" + `)

	valueContent = strings.TrimSpace(valueContent)

	if valueContent != "" {
		if checkFalse.MatchString(valueContent) {
			return false, nil
		} else {
			return true, nil
		}
	} else {
		return false, nil
	}
}`

var commonsAppViewContent = `package view

import ()

type DefaultMessage struct {
	Message SystemMessage ` + "`" + `json:"system_message"` + "`" + `
}

type SystemMessage struct {
	Type    string ` + "`" + `json:"type"` + "`" + `
	Content string ` + "`" + `json:"content"` + "`" + `
}

type ErrorMessage struct {
	Message SystemMessage ` + "`" + `json:"system_message"` + "`" + `
	Errors  []string      ` + "`" + `json:"errors"` + "`" + `
}

type MainPagination struct {
	CurrentPage  int ` + "`" + `json:"current_page"` + "`" + `
	TotalPages   int ` + "`" + `json:"total_pages"` + "`" + `
	TotalEntries int ` + "`" + `json:"total_entries"` + "`" + `
}

func SetSystemMessage(mType string, content string) SystemMessage {
	return SystemMessage{Type: mType, Content: content}
}

func SetDefaultMessage(mType string, content string) DefaultMessage {
	return DefaultMessage{Message: SetSystemMessage(mType, content)}
}

func SetErrorMessage(mType string, content string, errs []error) ErrorMessage {
	stringErrors := []string{}

	for _, value := range errs {
		stringErrors = append(stringErrors, value.Error())
	}

	return ErrorMessage{Message: SetSystemMessage(mType, content), Errors: stringErrors}
}

func SetUnauthorizedErrorMessage() DefaultMessage {
	return SetDefaultMessage("alert", "401 Unauthorized")
}

func SetForbiddenErrorMessage() DefaultMessage {
	return SetDefaultMessage("alert", "403 Forbidden")
}

func SetNotFoundErrorMessage() DefaultMessage {
	return SetDefaultMessage("alert", "404 Not found")
}

func SetBadRequestErrorMessage() DefaultMessage {
	return SetDefaultMessage("alert", "400 Bad Request")
}

func SetBadRequestInvalidJsonErrorMessage() DefaultMessage {
	return SetDefaultMessage("alert", "400 Bad Request - could not parse JSON")
}`

var commonsAppSearchEngineContent = `package model

import (
	"regexp"
	"strconv"
	"strings"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"time"
)

func (q *Query) SearchEngine(criteria map[string]string) {
	query, values := Criteria(q.Table, criteria, "AND")

	q.Db = q.Db.Where(query, values...)
}

func Criteria(table interface{}, criteria map[string]string, logic string) (string, []interface{}) {
	var queries []string
	var values []interface{}
	var query string
	var value interface{}
	var err error
	var hasInputValue = regexp.MustCompile(` + "`" + `\?` + "`" + `)

	for criterionKey, criterionValue := range criteria {
		query, value, err = handleCriterion(table, criterionKey, criterionValue)
		if err == nil {
			queries = append(queries, query)

			if hasInputValue.MatchString(query) {
				values = append(values, value)
			}

		} else {
			log.Error.Println("SearchEngine Query", err)
		}
	}

	query = strings.Join(queries, " "+logic+" ")

	return query, values
}

func handleCriterion(table interface{}, key string, value string) (string, interface{}, error) {
	var column, columnType, query string
	var interfaceValue interface{}
	var regexpInclusionQuery = regexp.MustCompile(` + "`" + `IN\s\(\?\)` + "`" + `)
	var regexpIsNullQuery = regexp.MustCompile(` + "`" + `IS\s(NOT\s){0,1}NULL` + "`" + `)
	var err error

	names := strings.Split(key, "_")

	column, query = strings.Join(names[:len(names)-1], "_"), names[len(names)-1]

	columnType, err = GetColumnType(table, column)
	if err != nil {
		log.Error.Println("SearchEngine handleCriterion", err)
		return "", "", err
	}

	query, value = translateQuery(column, query, value)

	if regexpIsNullQuery.MatchString(query) {
		interfaceValue = 1
	} else {
		interfaceValue, err = valueToInterface(columnType, value, regexpInclusionQuery.MatchString(query))
		if err != nil {
			log.Error.Println("SearchEngine handleCriterion", err)
			return "", "", err
		}
	}

	return query, interfaceValue, nil
}

func valueToInterface(columnType string, valueContent string, isQueryInclusion bool) (interface{}, error) {
	var regexpBooleanType = regexp.MustCompile(` + "`" + `bool` + "`" + `)
	var regexpIntType = regexp.MustCompile(` + "`" + `int` + "`" + `)
	var regexpFloatType = regexp.MustCompile(` + "`" + `float|double` + "`" + `)
	var regexpDateTimeType = regexp.MustCompile(` + "`" + `time|date` + "`" + `)
	var returnValue interface{}
	var err error

	if regexpIntType.MatchString(columnType) {
		returnValue, err = convertoToInt(valueContent, isQueryInclusion)
		if err != nil {
			return "", err
		}
	} else if regexpFloatType.MatchString(columnType) {
		returnValue, err = convertToFloat(valueContent, isQueryInclusion)
		if err != nil {
			return "", err
		}
	} else if regexpDateTimeType.MatchString(columnType) {
		returnValue, err = time.Parse(time.RFC3339, valueContent)
		if err != nil {
			return "", err
		}
	} else if regexpBooleanType.MatchString(columnType) {
		returnValue, err = convertToBoolean(valueContent)
		if err != nil {
			return "", err
		}
	} else if isQueryInclusion {
		returnValue = regexpSplit(valueContent, ` + "`" + `\s*,\s*` + "`" + `)
		err = nil
	} else {
		returnValue = valueContent
		err = nil
	}

	return returnValue, err
}

func translateQuery(column string, query string, value string) (string, string) {
	switch query {
	case "cont":
		query = "ILIKE ?"
		value = "%" + value + "%"
	case "notcont":
		query = "NOT ILIKE ?"
		value = "%" + value + "%"
	case "start":
		query = "ILIKE ?"
		value = value + "%"
	case "notstart":
		query = "NOT ILIKE ?"
		value = value + "%"
	case "end":
		query = "ILIKE ?"
		value = "%" + value
	case "notend":
		query = "NOT ILIKE ?"
		value = "%" + value
	case "gt":
		query = "> ?"
	case "lt":
		query = "< ?"
	case "gteq":
		query = ">= ?"
	case "lteq":
		query = "<= ?"
	case "null":
		query = "IS NULL"
	case "isnull":
		query = "IS NULL"
	case "notnull":
		query = "IS NOT NULL"
	case "in":
		query = "IN (?)"
	case "notin":
		query = "NOT IN (?)"
	case "true":
		query = "= 't'"
	case "false":
		query = "= 'f'"
	case "noteq":
		query = "<> ?"
	default:
		query = "= ?"
	}

	return column + " " + query, value
}

func convertoToInt(valueContent string, isQueryInclusion bool) (interface{}, error) {
	var newValues []int64

	newValue, err := strconv.ParseInt(valueContent, 10, 64)

	if err != nil {
		if !isQueryInclusion {
			return 0, err
		}

		newValues, err = splitStringToIntArray(valueContent)
		if err != nil {
			return 0, err
		} else {
			return newValues, nil
		}
	} else {
		return newValue, nil
	}
}

func convertToFloat(valueContent string, isQueryInclusion bool) (interface{}, error) {
	var newValues []float64

	newValue, err := strconv.ParseFloat(valueContent, 64)

	if err != nil {
		if !isQueryInclusion {
			return 0, err
		}

		newValues, err = splitStringToFloatArray(valueContent)
		if err != nil {
			return 0, err
		} else {
			return newValues, nil
		}
	} else {
		return newValue, nil
	}
}

func convertToStringIn(valueContent string) (interface{}, error) {
	stringValues := regexpSplit(valueContent, ` + "`" + `\s*,\s*` + "`" + `)
	return stringValues, nil
}

func convertToBoolean(valueContent string) (interface{}, error) {
	var checkFalse = regexp.MustCompile(` + "`" + `\A(0|f|false|no)\z` + "`" + `)

	valueContent = strings.TrimSpace(valueContent)

	if valueContent != "" {
		if checkFalse.MatchString(valueContent) {
			return false, nil
		} else {
			return true, nil
		}
	} else {
		return false, nil
	}
}

func splitStringToIntArray(value string) ([]int64, error) {
	var intValues []int64
	var intValue int64
	var i int
	var err error

	stringValues := regexpSplit(value, ` + "`" + `\s*,\s*` + "`" + `)

	for i = 0; i < len(stringValues); i++ {
		intValue, err = strconv.ParseInt(stringValues[i], 10, 64)
		if err != nil {
			intValues = nil
			return intValues, err
		}

		intValues = append(intValues, intValue)
	}

	return intValues, nil
}

func splitStringToFloatArray(value string) ([]float64, error) {
	var floatValues []float64
	var floatValue float64
	var i int
	var err error

	stringValues := regexpSplit(value, ` + "`" + `\s*,\s*` + "`" + `)

	for i = 0; i < len(stringValues); i++ {
		floatValue, err = strconv.ParseFloat(stringValues[i], 64)
		if err != nil {
			floatValues = nil
			return floatValues, err
		}

		floatValues = append(floatValues, floatValue)
	}

	return floatValues, nil
}

func regexpSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	lastStart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[lastStart:element[0]]
		lastStart = element[1]
	}
	result[len(indexes)] = text[lastStart:len(text)]
	return result
}`

var commonsAppPaginationContent = `package model

import (
	"strconv"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"test_repository_hub.com/test_account/test_project/config"
)

func (q *Query) Pagination(page, perPage string) (int, int, int) {
	type counter struct {
		Entries int
	}

	var currentPage, totalPages, entriesPerPage int
	var result counter

	currentPage = handleCurrentPage(page)
	entriesPerPage = handleEntriesPerPage(perPage)

	q.Db.Table(TableName(q.Table)).Order("", true).Select("COUNT(*) AS entries").Scan(&result)

	totalPages = result.Entries / entriesPerPage
	if (result.Entries % entriesPerPage) > 0 {
		totalPages++
	}

	offset := (currentPage - 1) * entriesPerPage
	q.Db = q.Db.Offset(offset).Limit(entriesPerPage)

	return currentPage, totalPages, result.Entries
}

func handleCurrentPage(page string) int {
	var currentPage int
	var err error

	currentPage, err = strconv.Atoi(page)
	if err != nil {
		currentPage = 1
	}

	return currentPage
}

func handleEntriesPerPage(perPage string) int {
	var entriesPerPage int
	var err error

	entriesPerPage, err = strconv.Atoi(perPage)
	if err != nil {
		entriesPerPage = config.App.Pagination.Default
	}

	if entriesPerPage > config.App.Pagination.Maximum {
		entriesPerPage = config.App.Pagination.Maximum
		log.Warn.Printf("Maximum value for entries per page can't be greater than %d", config.App.Pagination.Maximum)
	}

	return entriesPerPage
}`

var commonsAppOrderingContent = `package model

import ()

func (q *Query) Ordering(order string) {
	if order != "" {
		q.Db = q.Db.Order(order)
	}
}`

var commonsAppModelContent = `package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
	"test_repository_hub.com/test_account/test_project/commons/conversor"
	"test_repository_hub.com/test_account/test_project/commons/log"
)

var Db *gorm.DB

var Errors []string

type Query struct {
	Db    *gorm.DB
	Table interface{}
}

func Connect() {
	var err error

	dbConfig := loadDatabaseConfigFile()
	Db, err = gorm.Open("postgres", stringfyDatabaseConfigFile(dbConfig))

	if err != nil {
		log.Fatal.Println(err)
		panic("failed connect to database")
	} else {
		log.Info.Println("connected to the database successfully")
	}

	pool, err := strconv.Atoi(dbConfig["pool"])
	if err != nil {
		log.Fatal.Println(err)
	} else {
		log.Info.Printf("database pool of connections: %d", pool)
	}

	Db.DB().SetMaxIdleConns(pool)
}

func Disconnect() {
	defer Db.Close()
}

func TableName(table interface{}) string {
	return Db.NewScope(table).GetModelStruct().TableName(Db)
}

func GetColumnType(table interface{}, columnName string) (string, error) {
	field, ok := Db.NewScope(table).FieldByName(columnName)

	if ok {
		return field.Field.Type().String(), nil
	} else {
		return "", errors.New("column was not found")
	}
}

func GetColumnValue(table interface{}, columnName string) (interface{}, error) {
	field, ok := Db.NewScope(table).FieldByName(columnName)

	if ok {
		return field.Field.Interface(), nil
	} else {
		return "", errors.New("column was not found")
	}
}

func SetColumnValue(table interface{}, columnName string, value string) error {
	field, ok := Db.NewScope(table).FieldByName(columnName)

	if ok {
		columnType, _ := GetColumnType(table, columnName)
		valueInterface, _ := conversor.StringToInterface(columnType, value)
		return field.Set(valueInterface)
	} else {
		return errors.New("column was not found")
	}
}

func ColumnsFromTable(table interface{}, all bool) []string {
	var columns []string
	fields := Db.NewScope(table).Fields()

	for _, field := range fields {
		if !all && ((field.Names[0] == "Model") || (field.Relationship != nil)) {
			continue
		}
		columns = append(columns, field.DBName)
	}

	return columns
}

// PACKAGE METHODS

func loadDatabaseConfigFile() map[string]string {
	config := make(map[string]string)

	err := yaml.Unmarshal(readDatabaseConfigFile(), &config)
	if err != nil {
		log.Fatal.Printf("error: %v\n", err)
	}

	if config["pool"] == "" {
		config["pool"] = "5"
	}

	return config
}

func readDatabaseConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/database.yml")
	if err != nil {
		log.Fatal.Println(err)
	}

	return data
}

func stringfyDatabaseConfigFile(mapped map[string]string) string {
	var arr []string

	for key, value := range mapped {
		if key != "pool" {
			arr = append(arr, key+"='"+value+"'")
		}
	}

	return strings.Join(arr, " ")
}`

var commonsAppHandleContent = `package handler

import (
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"test_repository_hub.com/test_account/test_project/commons/app/view"
	"test_repository_hub.com/test_account/test_project/commons/log"
)

func ApiRoot(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("handler: ApiRoot")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "Yeah! Your API is working!"))
}

func Error401(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: Error401")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	json.NewEncoder(w).Encode(view.SetUnauthorizedErrorMessage())
}

func Error403(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: Error403")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	json.NewEncoder(w).Encode(view.SetForbiddenErrorMessage())
}

func Error404(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: Error404")
	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(404)

	json.NewEncoder(w).Encode(view.SetNotFoundErrorMessage())
}

func Error400(w http.ResponseWriter, r *http.Request, jsonParseError bool) {
	log.Info.Println("Handler: Error400")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	if jsonParseError {
		json.NewEncoder(w).Encode(view.SetBadRequestInvalidJsonErrorMessage())
	} else {
		json.NewEncoder(w).Encode(view.SetBadRequestErrorMessage())
	}
}

func QueryParamsToMapCriteria(param string, mapParams map[string][]string) map[string]string {
	var criteria, value string
	var checkParam, removePrefix, removeSufix *regexp.Regexp
	var err error

	query := make(map[string]string)

	checkParam, err = regexp.Compile(param + ` + "`" + `\[[a-zA-Z0-9\-\_]+\](\[\]){0,1}` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	removePrefix, err = regexp.Compile(param + ` + "`" + `\[` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	removeSufix, err = regexp.Compile(` + "`" + `\](\[\]){0,1}` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	for key := range mapParams {
		if checkParam.MatchString(key) {
			criteria = key
			criteria = removeSufix.ReplaceAllString(criteria, "")
			criteria = removePrefix.ReplaceAllString(criteria, "")
			value = strings.Join(mapParams[key], ",")

			query[criteria] = value
		}
	}

	return query
}

func SetPermittedParamsToEntity(params interface{}, entity interface{}) {
	setPermittedParams(params, entity, []string{}, []string{})
}

func SetPermittedParamsToEntityWithExceptions(params interface{}, entity interface{}, excepts []string) {
	setPermittedParams(params, entity, excepts, []string{})
}

func SetPermittedParamsToEntityButOnly(params interface{}, entity interface{}, only []string) {
	setPermittedParams(params, entity, []string{}, only)
}

func setPermittedParams(params interface{}, entity interface{}, excepts []string, only []string) {
	valParams := reflect.ValueOf(params).Elem()
	valEntity := reflect.ValueOf(entity).Elem()
	returnEntity := reflect.ValueOf(entity)

	for i := 0; i < valParams.NumField(); i++ {
		paramValueField := valParams.Field(i)
		paramTypeField := valParams.Type().Field(i)

		if !inExcept(excepts, paramTypeField.Name) && inOnly(only, paramTypeField.Name) {
			for j := 0; j < valEntity.NumField(); j++ {
				entityTypeField := valEntity.Type().Field(j)
				if paramTypeField.Name == entityTypeField.Name {
					returnEntity.Elem().Field(j).Set(paramValueField)
				}
			}
		}
	}
}

func inExcept(excepts []string, value string) bool {
	if len(excepts) == 0 {
		return false
	} else {
		for _, element := range excepts {
			if value == element {
				return true
			}
		}

		return false
	}
}

func inOnly(only []string, value string) bool {
	if len(only) == 0 {
		return true
	} else {
		for _, element := range only {
			if value == element {
				return true
			}
		}

		return false
	}
}`

var appUserModelContent = `package user

import (
	"errors"
	"regexp"
	"strings"
	"test_repository_hub.com/test_account/test_project/commons/app/model"
	"test_repository_hub.com/test_account/test_project/commons/crypto"
	"test_repository_hub.com/test_account/test_project/config"
	"test_repository_hub.com/test_account/test_project/db/entities"
	"time"
)

const NotFound = "user was not found"

var Current entities.User

func Find(id interface{}) (entities.User, error) {
	var user entities.User
	var err error

	model.Db.First(&user, id, "deleted_at IS NULL")
	if model.Db.NewRecord(user) {
		err = errors.New(NotFound)
	}

	return user, err
}

func FindAll() []entities.User {
	var users []entities.User

	model.Db.Order("name").Find(&users, "deleted_at IS NULL")

	return users
}

func IsValid(user *entities.User) (bool, []error) {
	var count int
	var errs []error
	var validEmail = regexp.MustCompile(` + "`" + `\A[^@]+@([^@\.]+\.)+[^@\.]+\z` + "`" + `)

	if len(user.Name) == 0 {
		errs = append(errs, errors.New("name can't be blank"))
	} else if len(user.Name) > 255 {
		errs = append(errs, errors.New("name is too long"))
	}

	if len(user.Email) == 0 {
		errs = append(errs, errors.New("email can't be blank"))
	} else if len(user.Email) > 255 {
		errs = append(errs, errors.New("email is too long"))
	} else if !validEmail.MatchString(user.Email) {
		errs = append(errs, errors.New("email is invalid"))
	} else if model.Db.Model(&entities.User{}).Where("id <> ? AND email = ? AND deleted_at IS NULL", user.ID, user.Email).Count(&count); count > 0 {
		errs = append(errs, errors.New("email has already been taken"))
	}

	if len(user.Password) < 8 {
		errs = append(errs, errors.New("password is too short minimum is 8 characters"))
	} else if len(user.Password) > 255 {
		errs = append(errs, errors.New("password is too long"))
	}

	if !isLocaleValid(user.Locale) {
		errs = append(errs, errors.New("locale is invalid"))
	}

	return (len(errs) == 0), errs
}

func Update(user *entities.User) (bool, []error) {
	var newValue, currentValue interface{}
	var valid bool
	var errs []error

	mapUpdate := make(map[string]interface{})

	currentUser, findErr := Find(user.ID)
	if findErr != nil {
		return false, []error{findErr}
	}

	valid, errs = IsValid(user)

	if valid {
		columns := model.ColumnsFromTable(user, false)
		for _, column := range columns {
			newValue, _ = model.GetColumnValue(user, column)
			currentValue, _ = model.GetColumnValue(currentUser, column)

			if newValue != currentValue {
				mapUpdate[column] = newValue

				if column == "password" {
					mapUpdate[column] = crypto.SetPassword(mapUpdate[column].(string))
				}
			}
		}

		if len(mapUpdate) > 0 {
			model.Db.Model(&user).Updates(mapUpdate)
		}

	}

	return valid, errs
}

func Create(user *entities.User) (bool, []error) {
	valid, errs := IsValid(user)
	if valid && model.Db.NewRecord(user) {
		user.Password = crypto.SetPassword(user.Password)

		model.Db.Create(&user)

		if model.Db.NewRecord(user) {
			errs = append(errs, errors.New("database error"))
			return false, errs
		}
	}

	return valid, errs
}

func Save(user *entities.User) (bool, []error) {
	if model.Db.NewRecord(user) {
		return Create(user)
	} else {
		return Update(user)
	}
}

func Destroy(user *entities.User) bool {
	if model.Db.NewRecord(user) {
		return false
	} else {
		model.Db.Delete(&user)
		return true
	}
}

func FindByEmail(email string) (entities.User, error) {
	var user entities.User
	var err error

	model.Db.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if model.Db.NewRecord(user) {
		user = entities.User{}
		err = errors.New(NotFound)
	}

	return user, err
}

func FindByResetPasswordToken(token string) (entities.User, error) {
	var user entities.User
	var err error

	enconded_token := crypto.EncryptText(token, config.App.SecretKey)
	two_days_ago := time.Now().Add(time.Second * time.Duration(config.App.ResetPasswordExpirationSeconds) * (-1))

	model.Db.Where("reset_password_token = ? AND reset_password_sent_at >= ? AND deleted_at IS NULL", enconded_token, two_days_ago).First(&user)
	if model.Db.NewRecord(user) {
		user = entities.User{}
		err = errors.New(NotFound)
	}

	return user, err
}

func Paginate(criteria map[string]string, order, page, perPage string) ([]entities.User, int, int, int) {
	var users []entities.User
	var user entities.User

	q := model.Query{Db: model.Db, Table: &user}
	q.SearchEngine(criteria)
	q.Ordering(order)
	currentPage, totalPages, totalEntries := q.Pagination(page, perPage)

	q.Db.Find(&users, "deleted_at IS NULL")

	return users, currentPage, totalPages, totalEntries
}

func Authenticate(email string, password string) (entities.User, error) {
	user, err := FindByEmail(email)

	if model.Db.NewRecord(user) || !crypto.CheckPassword(password, user.Password) {
		user = entities.User{}
		err = errors.New("invalid credentials")
	}

	return user, err
}

func IsNil(user *entities.User) bool {
	return model.Db.NewRecord(user)
}

func Exists(user *entities.User) bool {
	return !IsNil(user)
}

func SetCurrent(id interface{}) error {
	var err error
	Current, err = Find(id)

	return err
}

func IdExists(id interface{}) bool {
	_, err := Find(id)

	return (err == nil)
}

func SetRecovery(user *entities.User) (string, []error) {
	token := crypto.RandString(20)

	if model.Db.NewRecord(user) {
		return "", []error{errors.New(NotFound)}
	} else {
		t := time.Now()
		user.ResetPasswordSentAt = &t
		user.ResetPasswordToken = crypto.EncryptText(token, config.App.SecretKey)

		valid, errs := Save(user)

		if valid {
			return token, errs
		} else {
			return "", errs
		}
	}
}

func ClearRecovery(user *entities.User) (bool, []error) {
	if model.Db.NewRecord(user) {
		return false, []error{errors.New(NotFound)}
	} else {
		user.ResetPasswordToken = ""
		user.ResetPasswordSentAt = nil
		valid, errs := Save(user)

		return valid, errs
	}
}

func FirstName(user *entities.User) string {
	return strings.Split(user.Name, " ")[0]
}

// local methods

func isLocaleValid(locale string) bool {
	locales := config.App.Locales

	for _, a := range locales {
		if a == locale {
			return true
		}
	}

	return false
}`

var appUserViewContent = `package user

import (
	"test_repository_hub.com/test_account/test_project/commons/app/view"
	"test_repository_hub.com/test_account/test_project/db/entities"
	"time"
)

type PaginationJson struct {
	Pagination view.MainPagination ` + "`" + `json:"pagination"` + "`" + `
	Users      []Json              ` + "`" + `json:"users"` + "`" + `
}

type SuccessfullySavedJson struct {
	SystemMessage view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
	User          Json               ` + "`" + `json:"user"` + "`" + `
}

type Json struct {
	ID        uint       ` + "`" + `json:"id"` + "`" + `
	Name      string     ` + "`" + `json:"name"` + "`" + `
	Email     string     ` + "`" + `json:"email"` + "`" + `
	Admin     bool       ` + "`" + `json:"admin"` + "`" + `
	Locale    string     ` + "`" + `json:"locale"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time  ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt *time.Time ` + "`" + `json:"deleted_at"` + "`" + `
}

func SetJson(user entities.User) Json {
	return Json{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Admin:     user.Admin,
		Locale:    user.Locale,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}`

var appSessionMailerSignUpPtBr = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Bem-vindo {{ .UserFirstName }},</p>

		<p>Obrigado por se cadastrar em <strong>{{ .AppName }}</strong>!</p>

		<p>Qualquer coisa que você precisar, por favor, entre em contato conosco.</p>

		<p>Atenciosamente,</p>

		<p>Nome, Cargo<br>email@example.com</p>
	</body>
</html>`

var appSessionMailerSignUpEn = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Welcome {{ .UserFirstName }},</p>

		<p>Thank you for registering for <strong>{{ .AppName }}</strong>!</p>

		<p>Anything you need from us, please let us know.</p>

		<p>Best,</p>

		<p>Name, Job Position<br>email@example.com</p>
	</body>
</html>`

var appSessionMailerPasswordRecoveryPtBr = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Olá {{ .UserFirstName }}!</p>

		<p>Alguém solicitou o link para alterar sua senha. Você pode fazer isso através do link abaixo.</p>

		<p><a href="{{ .LinkToPasswordRecovery }}">Alterar senha</a></p>

		<p>Se não foi você que solicitou, por favor, apenas ignore este email.</p>
		<p>Sua senha não será alterada até você acessar o link acima e criar uma nova.</p>
	</body>
</html>`

var appSessionMailerPasswordRecoveryEn = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Hello {{ .UserFirstName }}!</p>

		<p>Someone has requested a link to change your password. You can do this through the link below.</p>

		<p><a href="{{ .LinkToPasswordRecovery }}">Change my password</a></p>

		<p>If you didn't request this, please ignore this email.</p>
		<p>Your password won't change until you access the link above and create a new one.</p>
	</body>
</html>`

var appSessionView = `package session

import (
	"bytes"
	"html/template"
	"test_repository_hub.com/test_account/test_project/app/user"
	"test_repository_hub.com/test_account/test_project/commons/app/view"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"test_repository_hub.com/test_account/test_project/config"
	"test_repository_hub.com/test_account/test_project/db/entities"
)

type SignInSuccess struct {
	Message view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
	Token   string             ` + "`" + `json:"token"` + "`" + `
	Expires int                ` + "`" + `json:"expires"` + "`" + `
}

type SignOutSuccess struct {
	Message view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
}

type SignUpSuccess struct {
	UserFirstName string
	AppName       string
}

type PasswordRecoveryInstructions struct {
	UserFirstName          string
	LinkToPasswordRecovery string
}

func SignInSuccessMessage(mType string, content string, token string) SignInSuccess {
	return SignInSuccess{Message: view.SystemMessage{mType, content}, Token: token, Expires: config.App.TokenExpirationSeconds}
}

func SignOutSuccessMessage(mType string, content string) SignOutSuccess {
	return SignOutSuccess{Message: view.SystemMessage{mType, content}}
}

func RefreshSuccessMessage(mType string, content string, token string) SignInSuccess {
	return SignInSuccessMessage(mType, content, token)
}

func SignUpSuccessMessage(mType string, content string, token string) SignInSuccess {
	return SignInSuccessMessage(mType, content, token)
}

func SignUpMailer(currentUser *entities.User) string {
	var content bytes.Buffer

	data := SignUpSuccess{UserFirstName: user.FirstName(currentUser), AppName: config.App.AppName}

	tmpl, err := template.ParseFiles("./app/session/mailer/sign_up." + currentUser.Locale + ".html")
	if err != nil {
		log.Error.Println(err)
	}

	err = tmpl.Execute(&content, &data)

	return content.String()
}

func PasswordRecoveryInstructionsMailer(currentUser *entities.User, token string) string {
	var content bytes.Buffer

	data := PasswordRecoveryInstructions{UserFirstName: user.FirstName(currentUser), LinkToPasswordRecovery: config.App.ResetPasswordUrl + "?token=" + token}

	tmpl, err := template.ParseFiles("./app/session/mailer/password_recovery." + currentUser.Locale + ".html")
	if err != nil {
		log.Error.Println(err)
	}

	err = tmpl.Execute(&content, &data)

	return content.String()
}`

var appSessionModel = `package session

import (
	"errors"
	"test_repository_hub.com/test_account/test_project/commons/app/model"
	"test_repository_hub.com/test_account/test_project/db/entities"
	"time"
)

const NotFound = "session was not found"

func Find(id interface{}) (entities.Session, error) {
	var session entities.Session
	var err error

	model.Db.First(&session, id)
	if model.Db.NewRecord(session) {
		err = errors.New(NotFound)
	}

	return session, err
}

func IsValid(session *entities.Session) (bool, []error) {
	return true, []error{}
}

func Update(session *entities.Session) (bool, []error) {
	var newValue, currentValue interface{}
	var valid bool
	var errs []error

	mapUpdate := make(map[string]interface{})

	currentSession, findErr := Find(session.ID)
	if findErr != nil {
		return false, []error{findErr}
	}

	valid, errs = IsValid(session)

	if valid {
		columns := model.ColumnsFromTable(session, false)
		for _, column := range columns {
			newValue, _ = model.GetColumnValue(session, column)
			currentValue, _ = model.GetColumnValue(currentSession, column)

			if newValue != currentValue {
				mapUpdate[column] = newValue
			}
		}

		if len(mapUpdate) > 0 {
			model.Db.Model(&session).Updates(mapUpdate)
		}

	}

	return valid, errs
}

func Create(session *entities.Session) (bool, []error) {
	valid, errs := IsValid(session)
	if valid && model.Db.NewRecord(session) {
		model.Db.Create(&session)

		if model.Db.NewRecord(session) {
			errs = append(errs, errors.New("database error"))
			return false, errs
		}
	}

	return valid, errs
}

func Save(session *entities.Session) (bool, []error) {
	if model.Db.NewRecord(session) {
		return Create(session)
	} else {
		return Update(session)
	}
}

func Destroy(session *entities.Session) bool {
	if model.Db.NewRecord(session) {
		return false
	} else {
		model.Db.Delete(&session)
		return true
	}
}

func FindByJti(jti string) (entities.Session, error) {
	var session entities.Session
	var err error

	model.Db.Where("jti = ?", jti).First(&session)
	if model.Db.NewRecord(session) {
		session = entities.Session{}
		err = errors.New(NotFound)
	}

	return session, err
}

func Deactivate(session *entities.Session) (bool, []error) {
	session.Active = false
	return Save(session)
}

func IncrementStats(session *entities.Session) (bool, []error) {
	t := time.Now()
	session.LastRequestAt = &t
	session.Requests = session.Requests + 1
	return Save(session)
}`

var appMyselfView = `package myself

import (
	"test_repository_hub.com/test_account/test_project/db/entities"
	"time"
)

type Json struct {
	ID        uint      ` + "`" + `json:"id"` + "`" + `
	Name      string    ` + "`" + `json:"name"` + "`" + `
	Email     string    ` + "`" + `json:"email"` + "`" + `
	Locale    string    ` + "`" + `json:"locale"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func SetJson(userMyself entities.User) Json {
	return Json{
		ID:        userMyself.ID,
		Name:      userMyself.Name,
		Email:     userMyself.Email,
		Locale:    userMyself.Locale,
		CreatedAt: userMyself.CreatedAt,
		UpdatedAt: userMyself.UpdatedAt,
	}
}`

var appHandlersUser = `package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"test_repository_hub.com/test_account/test_project/app/user"
	"test_repository_hub.com/test_account/test_project/commons/app/handler"
	"test_repository_hub.com/test_account/test_project/commons/app/view"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"test_repository_hub.com/test_account/test_project/db/entities"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	type UserPermittedParams struct {
		Name     string ` + "`" + `json:"name"` + "`" + `
		Email    string ` + "`" + `json:"email"` + "`" + `
		Password string ` + "`" + `json:"password"` + "`" + `
		Locale   string ` + "`" + `json:"locale"` + "`" + `
		Admin    bool   ` + "`" + `json:"admin"` + "`" + `
	}

	var userNew = entities.User{}

	log.Info.Println("Handler: UserCreate")
	w.Header().Set("Content-Type", "application/json")

	var userParams UserPermittedParams
	err := json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		log.Error.Println("could not parse JSON")
		handler.Error400(w, r, true)
		return
	}

	handler.SetPermittedParamsToEntity(&userParams, &userNew)

	valid, errs := user.Create(&userNew)

	if valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user was successfully created"), User: user.SetJson(userNew)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not created", errs))
	}
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	type UserPermittedParams struct {
		Name   string ` + "`" + `json:"name"` + "`" + `
		Email  string ` + "`" + `json:"email"` + "`" + `
		Locale string ` + "`" + `json:"locale"` + "`" + `
		Admin  bool   ` + "`" + `json:"admin"` + "`" + `
	}

	log.Info.Println("Handler: UserUpdate")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	userCurrent, err := user.Find(params["id"])
	if err != nil {
		handler.Error404(w, r)
		return
	}

	var userParams UserPermittedParams
	err = json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		log.Error.Println("could not parse JSON")
		handler.Error400(w, r, true)
		return
	}

	handler.SetPermittedParamsToEntity(&userParams, &userCurrent)

	if valid, errs := user.Update(&userCurrent); valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user was successfully updated"), User: user.SetJson(userCurrent)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", errs))
	}
}

func UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
	type UserPermittedParams struct {
		Password string ` + "`" + `json:"password"` + "`" + `
	}

	log.Info.Println("Handler: UserUpdate")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	userCurrent, err := user.Find(params["id"])
	if err != nil {
		handler.Error404(w, r)
		return
	}

	var userParams UserPermittedParams
	err = json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		log.Error.Println("could not parse JSON")
		handler.Error400(w, r, true)
		return
	}

	handler.SetPermittedParamsToEntity(&userParams, &userCurrent)

	if valid, errs := user.Update(&userCurrent); valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user password was successfully updated"), User: user.SetJson(userCurrent)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user password was not updated", errs))
	}
}

func UserDestroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: UserDestroy")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userCurrent, err := user.Find(params["id"])

	if err == nil && user.Destroy(&userCurrent) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully destroyed"))
	} else {
		handler.Error404(w, r)
	}
}

func UserShow(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: UserShow")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userCurrent, err := user.Find(params["id"])

	if err == nil {
		json.NewEncoder(w).Encode(user.SetJson(userCurrent))
	} else {
		handler.Error404(w, r)
	}
}

func UserList(w http.ResponseWriter, r *http.Request) {
	var i, page, entries, pages int
	var userList []entities.User

	userJsons := []user.Json{}

	log.Info.Println("Handler: UserList")
	w.Header().Set("Content-Type", "application/json")

	criteria := handler.QueryParamsToMapCriteria("search", r.URL.Query())
	order := userSanitizeOrder(r.FormValue("order"))

	userList, page, pages, entries = user.Paginate(criteria, order, r.FormValue("page"), r.FormValue("per_page"))

	for i = 0; i < len(userList); i++ {
		userJsons = append(userJsons, user.SetJson(userList[i]))
	}

	pagination := view.MainPagination{CurrentPage: page, TotalPages: pages, TotalEntries: entries}
	json.NewEncoder(w).Encode(user.PaginationJson{Pagination: pagination, Users: userJsons})
}

func userSanitizeOrder(value string) string {
	var allowedParams = []*regexp.Regexp{
		regexp.MustCompile(` + "`" + `(?i)\A\s*id(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*name(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*email(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*password(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*admin(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*locale(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*created_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*updated_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*deleted_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `)}

	for _, allowedParam := range allowedParams {
		if allowedParam.MatchString(value) {
			return value
		}
	}

	return ""
}`

var appHandlersSession = `package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"test_repository_hub.com/test_account/test_project/app/session"
	"test_repository_hub.com/test_account/test_project/app/user"
	"test_repository_hub.com/test_account/test_project/commons/app/handler"
	"test_repository_hub.com/test_account/test_project/commons/app/view"
	"test_repository_hub.com/test_account/test_project/commons/locale"
	"test_repository_hub.com/test_account/test_project/commons/log"
	"test_repository_hub.com/test_account/test_project/commons/mailer"
	"test_repository_hub.com/test_account/test_project/config"
	"test_repository_hub.com/test_account/test_project/db/entities"
	"time"
)

type SessionSignInParams struct {
	Email    string ` + "`" + `json:"email"` + "`" + `
	Password string ` + "`" + `json:"password"` + "`" + `
}

type SessionClaims struct {
	Uid uint   ` + "`" + `json:"uid"` + "`" + `
	Jti string ` + "`" + `json:"jti"` + "`" + `
	jwt.StandardClaims
}

const (
	privateKeyPath = "config/keys/app.key.rsa"
	publicKeyPath  = "config/keys/app.key.rsa.pub"
)

func SessionSignIn(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: SessionSignIn")
	w.Header().Set("Content-Type", "application/json")

	var signInParams SessionSignInParams
	err := json.NewDecoder(r.Body).Decode(&signInParams)
	if err != nil {
		log.Error.Println("could not parse JSON")
		handler.Error400(w, r, true)
		return
	}

	userAuth, err := user.Authenticate(signInParams.Email, signInParams.Password)

	if !user.IsNil(&userAuth) {
		json.NewEncoder(w).Encode(session.SignInSuccessMessage("notice", "signed in successfully", sessionGenerateToken(userAuth, r.RemoteAddr)))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "could not sign in", []error{err}))
	}
}

func SessionSignOut(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: SessionSignOut")
	w.Header().Set("Content-Type", "application/json")

	authToken, _ := sessionAuthToken(r.Header.Get("Authorization"))

	claims, ok := authToken.Claims.(*SessionClaims)
	if !ok || !authToken.Valid {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	sessionSignOut, errorFindByJti := session.FindByJti(claims.Jti)
	if errorFindByJti != nil {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	if deactivated, _ := session.Deactivate(&sessionSignOut); deactivated {
		json.NewEncoder(w).Encode(session.SignOutSuccessMessage("notice", "signed out successfully"))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
	}
}

func SessionRefresh(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: SessionRefresh")
	w.Header().Set("Content-Type", "application/json")

	authToken, _ := sessionAuthToken(r.Header.Get("Authorization"))

	claims, ok := authToken.Claims.(*SessionClaims)
	if !ok || !authToken.Valid {
		log.Error.Println("invalid token")
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	currentSession, errorFindByJti := session.FindByJti(claims.Jti)
	if errorFindByJti != nil {
		log.Error.Printf("could not find session by token %s", claims.Jti)
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	valid, _ := session.Deactivate(&currentSession)
	if !valid {
		log.Error.Println("could not deactivate session")
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	userSession, errorUserNotFound := user.Find(currentSession.UserID)
	if errorUserNotFound != nil {
		log.Error.Printf("could not find user by %d\n", currentSession.UserID)
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "access denied", []error{errors.New("invalid token")}))
		return
	}

	json.NewEncoder(w).Encode(session.RefreshSuccessMessage("notice", "session was successfully refreshed", sessionGenerateToken(userSession, r.RemoteAddr)))
}

func SessionSignUp(w http.ResponseWriter, r *http.Request) {
	type UserPermittedParams struct {
		Name     string ` + "`" + `json:"name"` + "`" + `
		Email    string ` + "`" + `json:"email"` + "`" + `
		Password string ` + "`" + `json:"password"` + "`" + `
		Locale   string ` + "`" + `json:"locale"` + "`" + `
	}

	var userNew = entities.User{}

	log.Info.Println("Handler: SessionSignUp")
	w.Header().Set("Content-Type", "application/json")

	var userParams UserPermittedParams
	err := json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		log.Error.Println("could not parse JSON")
		handler.Error400(w, r, true)
		return
	}

	handler.SetPermittedParamsToEntity(&userParams, &userNew)
	userNew.Admin = false

	if valid, errs := user.Save(&userNew); valid {
		locale.Load(userNew.Locale)

		mailer.AddTo(userNew.Name, userNew.Email)
		subject := locale.I18n.Welcome + " " + user.FirstName(&userNew)
		body := session.SignUpMailer(&userNew)
		go mailer.Send(subject, body, true)

		json.NewEncoder(w).Encode(session.SignUpSuccessMessage("notice", "user was successfully created", sessionGenerateToken(userNew, r.RemoteAddr)))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not created", errs))
	}
}

func SessionPassword(w http.ResponseWriter, r *http.Request) {
	var currentUser, _ = user.FindByEmail(r.FormValue("email"))

	log.Info.Println("Handler: SessionPassword")
	w.Header().Set("Content-Type", "application/json")

	if user.Exists(&currentUser) {
		locale.Load(currentUser.Locale)

		token, _ := user.SetRecovery(&currentUser)
		mailer.AddTo(currentUser.Name, currentUser.Email)
		subject := locale.I18n.PasswordRecoveryInstructions
		body := session.PasswordRecoveryInstructionsMailer(&currentUser, token)
		go mailer.Send(subject, body, true)
	}

	json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user password recovery instructions was successfully sent"))
}

func SessionRecovery(w http.ResponseWriter, r *http.Request) {
	var errs []error
	var valid bool

	log.Info.Println("Handler: SessionRecovery")
	w.Header().Set("Content-Type", "application/json")

	currentUser, _ := user.FindByResetPasswordToken(r.FormValue("token"))
	currentUser.Password = r.FormValue("new_password")

	if !user.Exists(&currentUser) {
		errs = append(errs, errors.New("invalid reset password token"))
	} else if r.FormValue("new_password") != r.FormValue("password_confirmation") {
		errs = append(errs, errors.New("password confirmation does not match new password"))
	} else if valid, errs = user.Save(&currentUser); valid {
		user.ClearRecovery(&currentUser)
		json.NewEncoder(w).Encode(session.SignInSuccessMessage("notice", "password was successfully changed", sessionGenerateToken(currentUser, r.RemoteAddr)))
	}

	if len(errs) > 0 {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "password could not be changed", errs))
	}
}

func SessionCheck(token string) (uint, error) {
	authToken, err := sessionAuthToken(token)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	claims, ok := authToken.Claims.(*SessionClaims)

	if ok && authToken.Valid && sessionUpdateStats(claims.Jti) {
		return claims.Uid, nil
	} else {
		return 0, errors.New("invalid token")
	}
}

func sessionBuildClaims(jti string, userSession *entities.User) (jwt.MapClaims, time.Time) {
	expiresAt := time.Now().Add(time.Second * time.Duration(config.App.TokenExpirationSeconds))

	claims := make(jwt.MapClaims)
	claims["iss"] = config.App.AppRepository
	claims["exp"] = expiresAt.Unix()
	claims["jti"] = jti
	claims["uid"] = userSession.ID

	return claims, expiresAt
}

func sessionUpdateStats(jti string) bool {
	currentSession, _ := session.FindByJti(jti)

	if !currentSession.Active {
		return false
	} else {
		session.IncrementStats(&currentSession)
		return true
	}
}

func sessionAuthToken(token string) (*jwt.Token, error) {
	var publicBytes []byte
	var publicKey *rsa.PublicKey
	var errorReadFile error
	var errorParseRsa error
	var err error
	var authToken *jwt.Token

	if token == "" {
		return authToken, errors.New("invalid token")
	} else {
		removeBearer := regexp.MustCompile(` + "`" + `^\s*Bearer\s+` + "`" + `)
		token = removeBearer.ReplaceAllString(token, "")
	}

	publicBytes, errorReadFile = ioutil.ReadFile(publicKeyPath)
	if errorReadFile != nil {
		log.Error.Println(errorReadFile)
		return authToken, errorReadFile
	}

	publicKey, errorParseRsa = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if errorParseRsa != nil {
		log.Error.Println(errorParseRsa)
		return authToken, errorParseRsa
	}

	authToken, err = jwt.ParseWithClaims(token, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	return authToken, err
}

func sessionGenerateToken(userSession entities.User, remoteAddr string) string {
	var privateBytes []byte
	var privateKey *rsa.PrivateKey
	var err error
	var expiresAt time.Time
	var sessionNew entities.Session

	privateBytes, err = ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal.Println(err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal.Println(err)
	}

	jti := uuid.Must(uuid.NewV4()).String()

	signer := jwt.New(jwt.SigningMethodRS256)
	signer.Claims, expiresAt = sessionBuildClaims(jti, &userSession)

	token, err := signer.SignedString(privateKey)
	if err != nil {
		log.Error.Println(err)
	} else {
		log.Info.Println("Token was successfully created for user " + userSession.Email)
	}

	t := time.Now()
	ip, _, _ := net.SplitHostPort(remoteAddr)
	sessionNew = entities.Session{Jti: jti, App: "Default", Requests: 0, LastRequestAt: &t, UserID: userSession.ID, Address: ip, ExpiresIn: config.App.TokenExpirationSeconds, ExpiresAt: expiresAt}
	session.Save(&sessionNew)

	return token
}`

var appHandlersMyself = `package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"test_repository_hub.com/test_account/test_project/app/myself"
	"test_repository_hub.com/test_account/test_project/app/user"
	"test_repository_hub.com/test_account/test_project/commons/app/handler"
	"test_repository_hub.com/test_account/test_project/commons/app/view"
	"test_repository_hub.com/test_account/test_project/commons/log"
)

func MyselfUpdate(w http.ResponseWriter, r *http.Request) {
	type MyselfPermittedParams struct {
		Name   string ` + "`" + `json:"name"` + "`" + `
		Locale string ` + "`" + `json:"locale"` + "`" + `
	}

	log.Info.Println("Handler: MyselfUpdate")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	var myselfParams MyselfPermittedParams
	err := json.NewDecoder(r.Body).Decode(&myselfParams)
	if err != nil {
		log.Error.Println("could not parser input JSON")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "could not parser input JSON", []error{err}))
		return
	}

	handler.SetPermittedParamsToEntity(&myselfParams, &userMyself)

	if valid, errs := user.Save(&userMyself); valid {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully updated"))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", errs))
	}
}

func MyselfUpdatePassword(w http.ResponseWriter, r *http.Request) {
	type MyselfPasswordParams struct {
		Password             string ` + "`" + `json:"password"` + "`" + `
		PasswordConfirmation string ` + "`" + `json:"password_confirmation"` + "`" + `
	}

	var errs []error
	var valid bool

	log.Info.Println("Handler: MyselfChangePassword")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	var myselfParams MyselfPasswordParams
	err := json.NewDecoder(r.Body).Decode(&myselfParams)
	if err != nil {
		log.Error.Println("could not parser input JSON")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "could not parser input JSON", []error{err}))
		return
	}

	userMyself.Password = myselfParams.Password

	if !user.Exists(&userMyself) {
		errs = append(errs, errors.New("invalid user"))
	} else if myselfParams.Password != myselfParams.PasswordConfirmation {
		errs = append(errs, errors.New("password confirmation does not match new password"))
	} else if valid, errs = user.Save(&userMyself); valid {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "password was successfully changed"))
	}

	if len(errs) > 0 {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "password could not be changed", errs))
	}
}

func MyselfDestroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: MyselfDestroy")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	if user.Destroy(&userMyself) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("alert", "user could not be destroyed"))
	}
}

func MyselfShow(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: MyselfShow")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(myself.SetJson(user.Current))
}`

var dbEntitySessionContent = `package entities

import (
	"time"
)

type Session struct {
	ID            uint       ` + "`" + `gorm:"primary_key"` + "`" + `
	UserID        uint       ` + "`" + `gorm:"index"` + "`" + `
	Jti           string     ` + "`" + `gorm:"type:varchar(255);unique_index"` + "`" + `
	App           string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	Requests      int        ` + "`" + `gorm:"not null;default:0"` + "`" + `
	ExpiresIn     int        ` + "`" + `gorm:"not null;default:0"` + "`" + `
	Address       string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	Active        bool       ` + "`" + `gorm:"default:true"` + "`" + `
	LastRequestAt *time.Time ` + "`" + `gorm:"default:null"` + "`" + `
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ExpiresAt     time.Time
}`

var dbEntityUserContent = `package entities

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name                string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	Email               string     ` + "`" + `gorm:"type:varchar(255);unique_index"` + "`" + `
	Admin               bool       ` + "`" + `gorm:"default:false"` + "`" + `
	Password            string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	ResetPasswordToken  string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	ResetPasswordSentAt *time.Time ` + "`" + `gorm:"default:null"` + "`" + `
	Locale              string     ` + "`" + `gorm:"type:varchar(255);default:'en'"` + "`" + `
	Sessions            []Session
}`

var dbSchemaMigrateContentV1 = `package schema

import (
	"test_repository_hub.com/test_account/test_project/app/user"
	"test_repository_hub.com/test_account/test_project/commons/app/model"
	"test_repository_hub.com/test_account/test_project/commons/crypto"
	"test_repository_hub.com/test_account/test_project/db/entities"
)

func Migrate() {
	model.Db.AutoMigrate(&entities.User{})

	_, err := user.FindByEmail("user@example.com")
	if err != nil {
		model.Db.Create(&entities.User{Name: "User Name", Email: "user@example.com", Password: crypto.SetPassword("Secret123!"), Locale: "en", Admin: true})
	}

	model.Db.AutoMigrate(&entities.Session{})
	model.Db.Model(&entities.Session{}).AddForeignKey("user_id", "users(id)", "NO ACTION", "NO ACTION")
}`
