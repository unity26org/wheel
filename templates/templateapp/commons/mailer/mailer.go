package mailer

var Path = []string{"commons", "mailer", "mailer.go"}

var Content = `package mailer

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"{{ .AppRepository }}/commons/log"
)

type Config struct {
	User     string
	Name     string
	Password string
	Address  string
	Port     int
}

var (
	from mail.Address
	to   []mail.Address
	cc   []mail.Address
	bcc  []mail.Address
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
	config := loadConfigFile()

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
		log.Fatal.Panic(err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Fatal.Panic(err)
	}

	// Auth
	if err = client.Auth(auth); err != nil {
		log.Fatal.Panic(err)
	}

	// To
	if err = client.Mail(from.Address); err != nil {
		log.Fatal.Panic(err)
	}

	// To, Cc and Bcc
	for i := 0; i < len(receipts); i++ {
		if err = client.Rcpt(receipts[0]); err != nil {
			log.Fatal.Panic(err)
		}
	}

	// Data
	socket, err := client.Data()
	if err != nil {
		log.Fatal.Panic(err)
	}

	_, err = socket.Write([]byte(message))
	if err != nil {
		log.Fatal.Panic(err)
	}

	err = socket.Close()
	if err != nil {
		log.Fatal.Panic(err)
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

func loadConfigFile() Config {
	config := Config{}

	err := yaml.Unmarshal(readConfigFile(), &config)
	if err != nil {
		log.Error.Fatalf("error: %v", err)
	}

	return config
}

func readConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/email.yml")
	if err != nil {
		log.Error.Fatal(err)
	}

	return data
}`
