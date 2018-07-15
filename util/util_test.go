package util

import (
	"testing"
	"fmt"
	"path/filepath"
)

func TestGetCurrentPath(t *testing.T) {

	p,e:=GetCurrentPath()
	if e!=nil {
		panic(e)
	}
	p1:=filepath.Join(p,"../")
	p2:=filepath.Join(p,"aa")
	if len(p1)>0 ||len(p2)>0{
		
	}
	fmt.Printf(p)

}
