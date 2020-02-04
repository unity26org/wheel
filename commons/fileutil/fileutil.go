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

func ReadBytesFile(filePath string, fileName string) ([]byte, error) {
	var fullPath string

	if filePath == "" {
		fullPath = fileName
	} else {
		fullPath = filepath.Join(filePath, fileName)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()

	return ioutil.ReadAll(file)
}

func ReadTextFile(filePath string, fileName string) (string, error) {
	b, err := ReadBytesFile(filePath, fileName)
	if err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func DirOrFileExists(fullPath string) bool {
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

func DestroyDirOrFile(fullPath string) error {
	return os.Remove(fullPath)
}

func DestroyAllDirOrFile(fullPath string) error {
	return os.RemoveAll(fullPath)
}

func UpdateTextFile(content string, filePath string, fileName string) error {
	return PersistFile(content, filePath, fileName, "a")
}

func SaveTextFile(content string, filePath string, fileName string) error {
	return PersistFile(content, filePath, fileName, "w")
}

func PersistFile(content string, filePath string, fileName string, pseudoMode string) error {
	err := os.MkdirAll(filePath, 0775)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(filePath, fileName)

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

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

	return nil
}
