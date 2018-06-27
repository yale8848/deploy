// Create by Yale 2018/3/2 18:03
package zipfile

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"regexp"
)

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
func ZipFile(zipFiles []string, zipName string, zipRegexp []string) string {

	if len(zipFiles) == 0 || len(zipName) == 0 {
		return zipName
	}
	files := []*os.File{}
	for _, z := range zipFiles {
		f, e := os.Open(z)
		if e != nil {
			checkError(e)
		}
		files = append(files, f)
	}
	err := Compress(files, zipRegexp, zipName)
	fmt.Println("[zip] zip " + zipName + " finish")
	checkError(err)
	return zipName
}

func Compress(files []*os.File, zipRegexp []string, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w, zipRegexp)
		if err != nil {
			return err
		}
	}
	return nil
}

func isFilter(file *os.File, zipRegexp []string) bool {

	if len(zipRegexp) == 0 {
		return false
	}

	n := file.Name()
	for _, r := range zipRegexp {

		if len(r) == 0 {
			continue
		}

		m, e := regexp.MatchString(r, n)
		if e != nil {
			continue
		}
		if m {
			return true
		}

	}
	return false
}

func compress(file *os.File, prefix string, zw *zip.Writer, zipRegexp []string) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if isFilter(file, zipRegexp) {
		return nil
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw, zipRegexp)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
