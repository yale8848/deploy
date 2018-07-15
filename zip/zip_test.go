// Create by Yale 2018/6/27 15:38
package zipfile

import "testing"

func TestZipFile(t *testing.T) {
	f := make([]string, 3)
	f[0] = "./test"
	f[1] = "./build.cmd"
	f[2] = "G:\\tmp\\mylog.txt"

	reg := make([]string, 1)
	reg[0] = "test/.*"

	ZipFile(f, "aaa.zip", reg)

}
