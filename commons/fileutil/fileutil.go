package fileutil

import (
	"bytes"
	"github.com/unity26org/wheel/commons/notify"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func GoFmtFile(fullpath string) {
	var out bytes.Buffer

	cmd := exec.Command("go", "fmt", "-x", fullpath)
	cmd.Stdout = &out
	_ = cmd.Run()
}

func ReadBytesFile(filePath string, fileName string) []byte {
	var fullPath string

	if filePath == "" {
		fullPath = fileName
	} else {
		fullPath = filepath.Join(filePath, fileName)
	}

	file, err := os.Open(fullPath)
	notify.FatalIfError(err)

	defer file.Close()

	b, err := ioutil.ReadAll(file)

	return b
}

func ReadTextFile(filePath string, fileName string) string {
	return string(ReadBytesFile(filePath, fileName))
}

func DirOrFileExists(fullPath string) bool {
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

func DestroyDirOrFile(fullPath string) {
	err := os.Remove(fullPath)
	notify.FatalIfError(err)
}

func DestroyAllDirOrFile(fullPath string) {
	err := os.RemoveAll(fullPath)
	notify.FatalIfError(err)
}

func UpdateTextFile(content string, filePath string, fileName string) {
	PersistFile(content, filePath, fileName, "a")
}

func SaveTextFile(content string, filePath string, fileName string) {
	PersistFile(content, filePath, fileName, "w")
}

func PersistFile(content string, filePath string, fileName string, pseudoMode string) {
	err := os.MkdirAll(filePath, 0775)
	notify.FatalIfError(err)

	fullPath := filepath.Join(filePath, fileName)

	f, err := os.Create(fullPath)
	notify.FatalIfError(err)

	defer f.Close()

	_, err = f.WriteString(content)
	notify.FatalIfError(err)

	f.Sync()

	GoFmtFile(fullPath)

	switch pseudoMode {
	case "w":
		notify.Created(fullPath)
	case "a":
		notify.Updated(fullPath)
	case "i":
		notify.Identical(fullPath)
	case "f":
		notify.Force(fullPath)
	case "s":
		notify.Skip(fullPath)
	}
}
