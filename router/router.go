package router

import (
	"net/http"
	"Valyrian/handler"
	"golang.org/x/net/websocket"
)

func Route() {
	http.HandleFunc("/", handler.HandleMainPage)
	routeJavaWeb()
}

func routeJavaWeb() {
	http.HandleFunc("/java-web/gen", handler.HandleJavaWebGenShell)
	http.HandleFunc("/java-web/build", handler.HandleJavaWebRunShell)
	// http.HandleFunc("/ws", websocket.Handler(handleWs))
	http.HandleFunc("/java-web/result", handler.HandleRunResult)
	http.Handle("/ws", websocket.Handler(handler.HandleWebsock))
	// http.HandleFunc("/java-web/upload/ssh", handler.HandleUploadSSH)
	// http.HandleFunc("/java-web/upload/maven", handler.HandleUploadMaven)
}