package diff

import (
	"bytes"
	"fmt"
	"github.com/unity26org/wheel/commons/fileutil"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Diff(content string, filePath string, fileName string, removeTempFile bool) (string, error) {
	var out bytes.Buffer

	tmpfile, err := ioutil.TempFile(filePath, strings.Replace(fileName, ".go", "_*.go", -1))
	if err != nil {
		return "", err
	}

	_, err = tmpfile.WriteString(content)
	if err != nil {
		return "", err
	}

	err = tmpfile.Close()
	if err != nil {
		return "", err
	}

	fileutil.GoFmtFile(tmpfile.Name())

	cmd := exec.Command("diff", "-u", filepath.Join(filePath, fileName), tmpfile.Name())
	cmd.Stdout = &out
	err = cmd.Run()
	// notify.Error(err)

	if removeTempFile {
		os.Remove(tmpfile.Name())
	}

	return out.String(), nil
}

func Patch(content string, filePath string, fileName string) error {
	var out bytes.Buffer
	var in bytes.Buffer

	tmpfile, err := ioutil.TempFile(filePath, strings.Replace(fileName, ".go", "_*.patch", -1))
	if err != nil {
		return err
	}

	diff, err := Diff(content, filePath, fileName, false)
	if err != nil {
		return err
	}

	removeFile := regexp.MustCompile(`\n`).Split(diff, -1)[1]
	removeFile = regexp.MustCompile(`(\s|\t)+`).Split(removeFile, -1)[1]

	defer os.Remove(tmpfile.Name())
	defer os.Remove(removeFile)

	_, err = tmpfile.WriteString(diff)
	if err != nil {
		return err
	}

	err = tmpfile.Close()
	if err != nil {
		return err
	}

	b, err := fileutil.ReadTextFile("", tmpfile.Name())
	if err != nil {
		return err
	}

	in.Write([]byte(b))

	cmd := exec.Command("patch", "-p0")
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Stdin = &in
	_ = cmd.Run()

	return nil
}

func Print(content string, filePath string, fileName string) {
	fmt.Println(Diff(content, filePath, fileName, true))
}
