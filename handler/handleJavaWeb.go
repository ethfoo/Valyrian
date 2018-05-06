package handler

import (
	"Valyrian/service"
	"Valyrian/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
	"os"
)

type DefaultValue struct {
	Repo string
}

type JavaWebGenShell struct {
	ShellName     string
	SshPri        string
	SshHost       string
	MavenSettings string
	Standalone    string
	TargetName    string
	UseSsh        string
	GitAddr       string
	GitUsername   string
	GitPassword   string
	Repository    string
}

type JavaWebRunShell struct {
	ShellName string
	Branch string
	TargetName string
}

func HandleJavaWebGenShell(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleJavaWebGenShell")
	if r.Method == "GET" {
		t, err := template.ParseFiles("static/java-web/genJavaWeb.gtpl")
		if err != nil {
			log.Printf("parse java-web/genJavaWeb.gtpl err: %#v", err)
		}
		defaultVal := DefaultValue{Repo: "hub.c.163.com/ethfoo/"}
		t.Execute(w, defaultVal)
	} else {
		r.ParseForm()
		genDto := JavaWebGenShell{}
		genDto.ShellName = r.FormValue("shellName")
		genDto.SshPri = r.FormValue("sshPri")
		genDto.SshHost = r.FormValue("sshHost")
		genDto.MavenSettings = r.FormValue("mavenSettings")
		genDto.Standalone = r.FormValue("standalone")
		genDto.UseSsh = r.FormValue("useSsh")
		genDto.GitAddr = r.FormValue("gitAddr")
		genDto.GitUsername = r.FormValue("gitUsername")
		genDto.GitPassword = r.FormValue("gitPassword")
		genDto.Repository = r.FormValue("repository")

		log.Printf("genDto: %#v", genDto)

		/*参数检查*/
		if genDto.ShellName == "" {
			utils.Return(utils.Response{Code: "400", Msg: "shellname cannot be null"}, w)
			return
		}

		if genDto.GitAddr == "" {
			utils.Return(utils.Response{"400", "git address cannot be null"}, w)
			return
		}

		if genDto.Repository == "" {
			utils.Return(utils.Response{"400", "repository cannot be null"}, w)
			return
		}

		if genDto.UseSsh == "true" {
			if genDto.SshPri == "" || genDto.SshHost == "" {
				utils.Return(utils.Response{"400", "private secret or known host cannot be null when git clone via ssh"}, w)
				return
			}
			// os.OpenFile("output/" + genDto.ShellName + "/build-shell/")
		} else {
			if genDto.GitPassword == "" || genDto.GitUsername == "" {
				// 可能存在开源的无需git权限的场景，这里就不强制要求username和password不为空
				//utils.Return(utils.Response{"400", "username and password cannot be null when git clone via https"}, w)
			}
			if !strings.HasPrefix(genDto.GitAddr, "https") {
				utils.Return(utils.Response{Code:"400", Msg:"git addr must start with https"}, w)
				return
			}
		}

		// 判断脚本名称是否已经被使用
		dirList, err := utils.ListDir("output")
		if err!=nil {
			utils.ReturnInternalError(w)
			return
		}
		for _, fi:= range dirList {
			if fi.Name() == genDto.ShellName {
				utils.Return(utils.Response{Code:"400", Msg:"shell name is being used"}, w)
				return
			}
		}

		// 初始化脚本文件，将模板file目录下的文件复制到output文件夹里
		err = service.InitOutputFiles(genDto.ShellName)
		if err != nil {
			log.Printf("initOutputFiles err: %#v", err)
			utils.ReturnInternalError(w)
			return
		}
		os.Chmod("output/" + genDto.ShellName + "/build-shell/java-web-built/dockerbuild-buildimage.sh", 0777)
		// os.Chmod("output/" + genDto.ShellName + "/build-shell/java-web-built/run-java-web-images.sh", 0777)

		if genDto.UseSsh == "true" {
			prefix := "output/" + genDto.ShellName + "/build-shell/java-web-built/ssh/";
			_, err := utils.WriteFile(prefix + "id_rsa", genDto.SshPri)
			_, err = utils.WriteFile(prefix + "known_hosts", genDto.SshHost)
			if err != nil {
				log.Printf("write ssh file err:%#v", err)
				utils.ReturnInternalError(w)
				return
			}
		} else {
			if strings.HasPrefix(genDto.GitAddr, "https") {
				if genDto.GitUsername != "" && genDto.GitPassword != "" {
					splits := strings.SplitAfter(genDto.GitAddr, "https://")
					log.Printf("splits[0]: %#v, splits[1]: %#v", splits[0], splits[1])
					genDto.GitAddr = splits[0] + genDto.GitUsername + ":" + genDto.GitPassword + "@" + splits[1]
				}
			} else {
				utils.Return(utils.Response{Code:"400", Msg:"git addr must start with https"}, w)
				return
			}
		}

		if genDto.MavenSettings != "" {
			_, err := utils.WriteFile("output/" + genDto.ShellName + "/build-shell/java-web-built/m2", genDto.MavenSettings)
			if err!=nil {
				log.Printf("write maven settings file err:%#v", err)
				utils.ReturnInternalError(w)
				return
			}
		}

		//修改脚本
		buildShVar, err := os.OpenFile("output/" + genDto.ShellName + "/build-shell/java-web-built/build-var.sh", os.O_RDWR|os.O_APPEND, 0777)
		if err!= nil {
			log.Println(err)
			utils.ReturnInternalError(w)
			return
		}
		defer buildShVar.Close()

		gitAddrStr := "REMOTE_GIT_ADDR=" +"'" + genDto.GitAddr + "'\n"
		repositoryStr := "IMAGE_PRE_NAME=" + "'" + genDto.Repository + "'\n"
		standaloneStr := "STANDALONE=" + "'" + genDto.Standalone + "'\n"
		targetNameStr := "TARGET_MODULE=" + "'" + genDto.ShellName + "'\n"
		shellStr := gitAddrStr + repositoryStr + standaloneStr + targetNameStr
		buildShVar.WriteString(shellStr)

		log.Printf("----------begin to build buildimage--------\n")
		//构建 构建时镜像
		buildimageShell := "output/" + genDto.ShellName + "/build-shell/java-web-built/dockerbuild-buildimage.sh"
		log.Printf("buildimageShell: %s", buildimageShell)
		_, err = utils.RunShellFile(buildimageShell)
		if err!=nil {
			// log.Printf("run shell file err:%#v", err)
			utils.ReturnInternalError(w)
			return
		}
		utils.Return(utils.Response{Code:"200", Msg:"genarate shell success"}, w)
	}
}

type RunShellDefaultVal struct {
	ShellName string
}

func HandleJavaWebRunShell(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleJavaWebRunShell")
	if r.Method == "GET" {
		r.ParseForm()
		dfVal := RunShellDefaultVal{}
		dfVal.ShellName = r.FormValue("sn")

		// 判断脚本是否存在
		exist := isShellExist(dfVal.ShellName, w)
		if !exist {
			return
		}

		t, err := template.ParseFiles("static/java-web/runJavaWeb.gtpl")
		if err != nil {
			log.Printf("parse java-web/runJavaWeb.gtpl err: %#v", err)
		}
		//defaultVal := DefaultValue{Repo: "hub.c.163.com/ethfoo/"}
		t.Execute(w, dfVal)
	} else {
		r.ParseForm()
		runDto := JavaWebRunShell{}
		runDto.ShellName = r.FormValue("ns")
		runDto.Branch = r.FormValue("branch")
		runDto.TargetName = r.FormValue("targetName")
		log.Printf("runDto: %#v", runDto)

		// 参数检查
		if runDto.Branch == "" {
			utils.Return(utils.Response{Code:"400", Msg:"branch cannot be null"}, w)
			return
		}
		if runDto.TargetName == "" {
			utils.Return(utils.Response{Code:"400", Msg:"target name cannot be null", }, w)
			return
		}
		// 判断脚本是否存在
		exist := isShellExist(runDto.ShellName, w)
		if !exist {
			return
		}

		// 生成并运行构建最终镜像脚本
		runFileName := "output/" + runDto.ShellName + "/build-shell/java-web-built/run-build-image.sh"
		runFile, err := os.OpenFile( runFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
		if err!=nil {
			utils.ReturnInternalError(w)
			return
		}
		err = os.Chmod(runFileName, 0777)
		if err!=nil {
			utils.ReturnInternalError(w)
			return
		}
		defer runFile.Close()

		isStandalone, err := utils.CheckLine("output/" + runDto.ShellName + "/build-shell/java-web-built/build-var.sh", "STANDALONE='-s'\n")
		if err!=nil {
			utils.ReturnInternalError(w)
			return
		}
		var standaloneStr string
		if isStandalone {
			standaloneStr=" -s "
		}
		branchStr := " -b " + runDto.Branch
		argsStr := branchStr + standaloneStr + "\n"
		
		shellStr := "#!/bin/bash\n" + "docker run -v ~/.m2/repository:/root/.m2/repository -v /var/run/docker.sock:/var/run/docker.sock -v $(which docker):/usr/bin/docker valyrian:buildv1 "
		rmStr := "rm " + runFileName + "\n"
		shellStr = shellStr + argsStr + rmStr
		runFile.WriteString(shellStr)

		log.Printf("-----------执行run-build-image.sh---------")
		runbuildimageShell := "output/" + runDto.ShellName + "/build-shell/java-web-built/run-build-image.sh"
		log.Printf("runbuildimageShell: %s", runbuildimageShell)
		_, err = utils.RunShellFile(runbuildimageShell)
		if err!=nil {
			utils.ReturnInternalError(w)
			return
		}
		utils.Return(utils.Response{Code:"200", Msg:"run build image shell success"}, w)
	}
}

func isShellExist(shellName string, w http.ResponseWriter) bool {
	log.Printf("isShellExist, input shellname: %s", shellName)
	dirList, err := utils.ListDir("output")
		if err!=nil {
			utils.ReturnInternalError(w)
			return false
		}
		isExist := false
		for _, fi:= range dirList {
			log.Printf("isShellExist,list fileName: %s", fi.Name())
			if fi.Name() == shellName {
				isExist = true
			}
		}
		if !isExist {
			utils.Return(utils.Response{Code:"400", Msg:"shell name is not exist"}, w)
			return false
		}
		return true
}