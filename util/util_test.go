package util

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestGetCurrentPath(t *testing.T) {

	p, e := GetCurrentPath()
	if e != nil {
		panic(e)
	}
	p1 := filepath.Join(p, "../")
	p2 := filepath.Join(p, "aa")
	if len(p1) > 0 || len(p2) > 0 {
	}
	fmt.Printf(p)
}

func TestGetUserPassWord(t *testing.T) {
	u, p, e := GetUserPassWord("usr.password")
	fmt.Printf("%v %v %v", u, p, e)
}
