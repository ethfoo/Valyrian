package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Code string
	Msg  string
}

func Return(resp Response, w http.ResponseWriter) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(respBytes))
}

func ReturnMap(response map[string]interface{}, w http.ResponseWriter) {
	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(respBytes))
}

func ReturnInternalError(w http.ResponseWriter) {
	resp := Response{Code: "500", Msg: "Internal Server Error"}
	Return(resp, w)
}
