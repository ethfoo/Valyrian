package router

import (
	"net/http"
	"Valyrian/handler"
)

func Route() {
	http.HandleFunc("/", handler.HandleMainPage)
	routeJavaWeb()
}

func routeJavaWeb() {
	http.HandleFunc("/java-web/gen", handler.HandleJavaWebGenShell)
	http.HandleFunc("/java-web/build", handler.HandleJavaWebRunShell)
	// http.HandleFunc("/java-web/upload/ssh", handler.HandleUploadSSH)
	// http.HandleFunc("/java-web/upload/maven", handler.HandleUploadMaven)
}