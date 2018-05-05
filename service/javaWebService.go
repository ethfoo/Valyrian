package service

import (
	"Valyrian/utils"
)

func InitOutputFiles(shellname string) error {

	err := utils.DirCopy("output/" + shellname + "/", "build-shell/java-web-built")
	// _,err := utils.FileCopy("output/settings.xml", "build-shell/java-web-built/m2/settings.xml")
	if err != nil {
		return err
	}
	return nil
}
