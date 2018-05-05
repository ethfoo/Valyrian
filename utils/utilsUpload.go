package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Upload(formName string, dirName string, perm os.FileMode, overrideFileName string, r *http.Request) (int64, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(formName)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer file.Close()
	var fileName string
	var f *os.File
	if overrideFileName != "" {
		fileName = "file/" + dirName + overrideFileName
		f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm) //覆盖默认settings

	} else {
		fileName = "file/" + dirName + handler.Filename
		f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, perm)
	}
	os.Chmod(fileName, perm) //预防万一，改文件权限，否则可能因为权限太高导致git拉不下来代码
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer f.Close()
	copiedBytes, err := io.Copy(f, file)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return copiedBytes, nil
}
