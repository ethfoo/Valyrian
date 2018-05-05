package handler

import (
	"Valyrian/service"
	"Valyrian/utils"
	"fmt"
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
	UseSsh        string
	GitAddr       string
	GitUsername   string
	GitPassword   string
	Repository    string
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
		buildSh, err := os.OpenFile("output/" + genDto.ShellName + "/build-shell/java-web-built/build.sh", os.O_RDWR, 0666)
		if err!= nil {
			log.Println(err)
			utils.ReturnInternalError(w)
			return
		}
		defer buildSh.Close()

		curOffset,_ := buildSh.Seek(30, os.SEEK_SET)
		// log.Printf("curOffset: %d\n", curOffset)
		gitAddrStr := "REMOTE_GIT_ADDR=" +"'" + genDto.GitAddr + "'\n"
		repositoryStr := "IMAGE_PRE_NAME=" + "'" + genDto.Repository + "'\n"
		shellStr := gitAddrStr + repositoryStr
		buildSh.WriteAt([]byte(shellStr) ,curOffset)

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

func HandleUploadSSH(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleUploadSSH")
	if r.Method == "POST" {
		copiedBytes, err := utils.Upload("uploadSshPriFile", "java-web-built/ssh/", 0600, "", r)

		if err != nil {
			fmt.Fprintf(w, "upload err: %#v", err)
			return
		}
		fmt.Fprintf(w, "upload success %#v bytes", copiedBytes)
		return
	}
}

func HandleUploadMaven(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleUploadMaven")
	if r.Method == "POST" {
		copiedBytes, err := utils.Upload("uploadMavenFile", "java-web-built/m2/", 0666, "settings.xml", r)
		if err != nil {
			fmt.Fprintf(w, "upload err: %#v", err)
			return
		}
		fmt.Fprintf(w, "upload success %#v bytes", copiedBytes)
		return
	}
}
