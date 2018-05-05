package utils

import (
	"log"
	"os/exec"
)

func RunShellFile(fileName string) (output string, err error) {
	cmd := exec.Command("/bin/bash", fileName)
	bytes, err := cmd.Output()
	if err!=nil {
		log.Printf("exec shell err:%s", err)
		return "", err
	}
	log.Printf("exec shell output: %s", string(bytes))
	return string(bytes), err
}