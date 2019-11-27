package notify

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func FatalIfError(err error) {
	if err != nil {
		Error(err)
		os.Exit(0)
	}
}

func Error(err interface{}) {
	ErrorJustified(err, -1)
}

func Success(message string) {
	GenericColorful("success", message, 32, -1)
}

func Created(message string) {
	GenericColorful("created", message, 32, 15)
}

func Updated(message string) {
	GenericColorful("updated", message, 36, 15)
}

func Skip(message string) {
	GenericColorful("skip", message, 33, 15)
}

func Force(message string) {
	GenericColorful("force", message, 33, 15)
}

func Identical(message string) {
	GenericColorful("identical", message, 34, 15)
}

func Warn(message string) {
	GenericColorful("warn", message, 93, 15)
}

func Patch(message string) {
	GenericColorful("patch", message, 35, 15)
}

func ErrorJustified(err interface{}, leftJustify int) {
	var message string

	switch v := err.(type) {
	case int:
		message = strconv.Itoa(err.(int))
	case float64:
		message = strconv.FormatFloat(err.(float64), 'E', -1, 64)
	case string:
		message = err.(string)
	case error:
		message = err.(error).Error()
	default:
		FatalIfError(errors.New(fmt.Sprintf("invalid type %v", v)))
		return
	}

	GenericColorful("error", message, 91, leftJustify)
}

func GenericColorful(prefix string, message string, color int, leftJustify int) {
	if leftJustify > 0 {
		for {
			if len(prefix) > leftJustify {
				break
			} else {
				prefix = " " + prefix
			}
		}
	}

	fmt.Println("\033["+strconv.Itoa(color)+"m"+prefix+":\033[39m", message)
}

func Simple(message string) {
	fmt.Print(message)
}

func Simpleln(message string) {
	fmt.Println(message)
}

func Fatal(message string) {
	Simpleln(message)
	os.Exit(0)
}

func WarnAppendToRoutes(err error, newCode string) {
	Warn(err.Error())
	fmt.Println("")
	fmt.Println("Edit file \"routes/routes.go\" and append this new code bellow to \"func Routes()\"")
	fmt.Println(newCode)
	fmt.Println("")
}

func WarnAppendToMigrate(err error, newCode string) {
	Warn(err.Error())
	fmt.Println("")
	fmt.Println("Edit file \"db/schema/migrate.go\" and append this new code bellow at \"func Migrate()\"")
	fmt.Println("")
	fmt.Println("\t" + newCode)
	fmt.Println("")
}

func WarnAppendToAuthorize(err error, newCode string) {
	Warn(err.Error())
	fmt.Println("")
	fmt.Println("Edit file \"routes/authorize.go\" and append this new code bellow at \"func init()\"")
	fmt.Println("")
	fmt.Println(newCode)
	fmt.Println("")
}

func NewApp(rootAppPath string) {
	fmt.Println("")
	fmt.Println("\033[32mYour RESTful API was successfully created!\033[39m")
	fmt.Println("")
	fmt.Println("Change to the root directory using the command line below: ")
	fmt.Println("\033[32mcd " + rootAppPath + "\033[39m")
	fmt.Println("")
	fmt.Println("Set up your database connection modifying the file config/database.yml")
	fmt.Println("")
	fmt.Println("For more details call help:")
	fmt.Println("go run main.go --help")
}
