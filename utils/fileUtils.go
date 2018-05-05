package utils

import (
	"log"
	"os"
	"io"
	"path/filepath"
	"io/ioutil"
)

func WriteFile(fileName, content string) (written int, err error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err!=nil {
		return 0, err
	}
	defer file.Close()
	n, err := file.WriteString(content)
	if err!=nil {
		return 0, err
	}
	return n, nil
}


func FileCopy(dstName, srcName string) (copiedBytes int64, err error) {
	src, err := os.Open(srcName)
	// src, err := os.OpenFile(dstName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("file copy err:%#v", err)
		return 0, err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("file copy err:%#v", err)
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func DirCopy(dstDir string, srcDir string) error {
	err := filepath.Walk(srcDir, func(filename string, file os.FileInfo, err error) error {
		if err!=nil {
			return err
		}
		if file.IsDir() {
			// log.Printf("file is dir, filename:%s", filename)
			os.MkdirAll(dstDir + filename, os.ModeDir)
			return nil
		}
		// log.Printf("walk file: %#v", filename)
		_, err = FileCopy(dstDir + filename ,filename)
		if err != nil {
			log.Fatalf("copy file to output dir err: %#v", err)
			return err
		}
		return nil
	})
	return err
}

func ListDir(dirName string) ([]os.FileInfo, error) {
	fileList, err := ioutil.ReadDir(dirName)
	if err!=nil {
		return nil, err
	}
	return fileList, nil
}