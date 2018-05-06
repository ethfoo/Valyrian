package utils

import (
	"io"
	"bufio"
	"log"
	"os/exec"
)

type Handler func(line string)

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

func RealtimeRunShell(fileName string, handle Handler) error {
	cmd := exec.Command("/bin/bash", fileName)
	stdout, err := cmd.StdoutPipe()
	if err!=nil {
		log.Printf("RealtimeRunShell err:%#v", err)
		return err
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err!=nil || err == io.EOF {
			break
		}
		handle(line)
	}
	cmd.Wait()
	return nil
}