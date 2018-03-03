// Create by Yale 2018/3/2 13:49
package main

import (
	"testing"
	"strings"
	"fmt"
)

func Test_getServers(t * testing.T)  {

	ss:=strings.Split("aaa",",")
	fmt.Println(len(ss))

	ss=strings.Split("aaac,bbb",",")
	fmt.Println(len(ss))

}