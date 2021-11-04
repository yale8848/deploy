package util

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetUserPassWord(fp string) (string, string, error) {

	by, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", "", err
	}
	user := ""
	pass := ""
	str := strings.TrimSpace(string(by))

	ps := strings.Split(str, "\n")
	if len(ps) == 2 {
		user = strings.Trim(ps[0], "\r")
		pass = strings.Trim(ps[1], "\r")
		return user, pass, nil
	}

	ps = strings.Split(str, "@")
	if len(ps) > 0 {
		user = ps[0]
	}
	if len(ps) > 1 {
		pass = ps[1]
	}
	return user, pass, nil

}
