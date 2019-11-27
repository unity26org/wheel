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
