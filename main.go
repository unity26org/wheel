package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/unity26org/wheel/commons/notify"
	"github.com/unity26org/wheel/generator"
	"github.com/unity26org/wheel/help"
	"github.com/unity26org/wheel/version"
	"os"
	"os/exec"
	"strings"
)

func IsGoInstalled() bool {
	var out bytes.Buffer

	cmd := exec.Command("go", "version")
	cmd.Stdout = &out
	err := cmd.Run()

	return err == nil
}

func CheckDependences() {
	var out bytes.Buffer
	var hasDependence bool

	notify.Simpleln("Checking dependences...")

	cmd := exec.Command("go", "list", "...")
	cmd.Stdout = &out

	cmd.Run()
	// if err := cmd.Run(); err != nil {
	//   fmt.Println("error:", err)
	// }

	installedDependences := strings.Split(out.String(), "\n")

	requiredDependences := []string{"github.com/jinzhu/gorm", "gopkg.in/yaml.v2", "github.com/gorilla/mux", "github.com/dgrijalva/jwt-go", "github.com/satori/go.uuid", "github.com/lib/pq", "golang.org/x/crypto/bcrypt"}
	for _, requiredDependence := range requiredDependences {
		hasDependence = false
		for _, installedDependence := range installedDependences {
			hasDependence = (requiredDependence == installedDependence)
			if hasDependence {
				break
			}
		}

		if !hasDependence {
			notify.Simple(fmt.Sprintf("         package %s was not found, installing...", requiredDependence))
			cmd := exec.Command("go", "get", requiredDependence)
			cmd.Stdout = &out
			err := cmd.Run()
			notify.FatalIfError(err)

			notify.Simpleln(fmt.Sprintf("         package %s was successfully installed", requiredDependence))
		} else {
			notify.Simpleln(fmt.Sprintf("         package %s was found", requiredDependence))
		}
	}
}

func optionsAreValid(args []string) bool {
	b := true

	for index, value := range args {
		if (index > 2) && (value != "-G" && value != "--skip-git") {
			b = false
			break
		}
	}

	return b
}

func checkGitIgnore(args []string) bool {
	b := true

	for index, value := range args {
		if index > 2 && value == "-G" || value == "--skip-git" {
			b = false
			break
		}
	}

	return b
}

func handleNewApp(args []string) {
	var options = make(map[string]interface{})

	if !optionsAreValid(args) {
		err := errors.New("invalid option. Run \"wheel --help\" for details")
		notify.FatalIfError(err)
	}

	preOptions := strings.Split(os.Args[2], "/")

	options["app_name"] = preOptions[len(preOptions)-1]
	options["app_repository"] = os.Args[2]
	options["git_ignore"] = checkGitIgnore(os.Args)

	notify.Simpleln("Generating new app...")
	generator.NewApp(options)
}

func buildGenerateOptions(args []string) (map[string]bool, error) {
	var options = make(map[string]bool)
	var subject string
	var err error

	subject = args[2]

	options["model"] = false
	options["entity"] = false
	options["view"] = false
	options["handler"] = false
	options["routes"] = false
	options["migrate"] = false
	options["authorize"] = false

	switch subject {
	case "scaffold":
		options["model"] = true
		options["entity"] = true
		options["view"] = true
		options["handler"] = true
		options["routes"] = true
		options["migrate"] = true
		options["authorize"] = true
		if len(args) < 4 {
			err = errors.New("invalid scaffold name")
		}
	case "model":
		options["model"] = true
		options["entity"] = true
		options["migrate"] = true
		if len(args) < 4 {
			err = errors.New("invalid model name")
		}
	case "handler":
		options["handler"] = true
		options["routes"] = true
		options["authorize"] = true
	case "entity":
		options["entity"] = true
		options["migrate"] = true
		if len(args) < 4 {
			err = errors.New("invalid entity name")
		}
	default:
		err = errors.New("invalid generate subject. Run \"wheel --help\" for details")
	}

	return options, err
}

func handleGenerateNewCrud(args []string, options map[string]bool) {
	var columns []string

	for index, value := range args {
		if index <= 3 {
			continue
		} else {
			columns = append(columns, value)
		}
	}

	notify.Simpleln("Generating new CRUD...")
	generator.NewCrud(args[3], columns, options)
}

func handleGenerate(args []string) {
	var options map[string]bool
	var err error

	options, err = buildGenerateOptions(args)
	notify.FatalIfError(err)

	handleGenerateNewCrud(args, options)
}

func handleHelp() {
	notify.Simpleln(help.Content)
}

func handleVersion() {
	notify.Simpleln(version.Content)
}

func checkIsGoInstalled() {
	if !IsGoInstalled() {
		notify.FatalIfError(errors.New("\"Go\" seems not installed"))
	} else {
		notify.Simpleln("\"Go\" seems installed")
		CheckDependences()
	}
}

func main() {
	command := ""

	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	if command == "new" || command == "n" {
		checkIsGoInstalled()
		handleNewApp(os.Args)
	} else if command == "generate" || command == "g" {
		checkIsGoInstalled()
		handleGenerate(os.Args)
	} else if command == "--help" || command == "-h" {
		handleHelp()
	} else if command == "--version" || command == "-v" {
		handleVersion()
	} else {
		notify.ErrorJustified("invalid argument. Use \"new\" or \"generate\". Run \"wheel --help\" for details", 0)
		handleHelp()
		notify.Fatal("")
	}
}
