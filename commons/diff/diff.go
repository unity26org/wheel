package diff

import (
	"bytes"
	"fmt"
	"github.com/unity26org/wheel/commons/fileutil"
	"github.com/unity26org/wheel/commons/notify"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Diff(content string, filePath string, fileName string, removeTempFile bool) string {
	var out bytes.Buffer

	tmpfile, err := ioutil.TempFile(filePath, strings.Replace(fileName, ".go", "_*.go", -1))
	notify.FatalIfError(err)

	_, err = tmpfile.WriteString(content)
	notify.FatalIfError(err)

	err = tmpfile.Close()
	notify.FatalIfError(err)

	fileutil.GoFmtFile(tmpfile.Name())

	cmd := exec.Command("diff", "-u", filepath.Join(filePath, fileName), tmpfile.Name())
	cmd.Stdout = &out
	err = cmd.Run()
	// notify.Error(err)

	if removeTempFile {
		os.Remove(tmpfile.Name())
	}

	return out.String()
}

func Patch(content string, filePath string, fileName string) {
	var out bytes.Buffer
	var in bytes.Buffer

	tmpfile, err := ioutil.TempFile(filePath, strings.Replace(fileName, ".go", "_*.patch", -1))
	notify.FatalIfError(err)

	diff := Diff(content, filePath, fileName, false)

	removeFile := regexp.MustCompile(`\n`).Split(diff, -1)[1]
	removeFile = regexp.MustCompile(`(\s|\t)+`).Split(removeFile, -1)[1]

	defer os.Remove(tmpfile.Name())
	defer os.Remove(removeFile)

	_, err = tmpfile.WriteString(diff)
	notify.FatalIfError(err)

	err = tmpfile.Close()
	notify.FatalIfError(err)

	in.Write([]byte(fileutil.ReadTextFile("", tmpfile.Name())))

	cmd := exec.Command("patch", "-p0")
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Stdin = &in
	_ = cmd.Run()
}

func Print(content string, filePath string, fileName string) {
	fmt.Println(Diff(content, filePath, fileName, true))
}
