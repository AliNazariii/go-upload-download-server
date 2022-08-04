package http_server

import (
	"concurrent-http-server/api"
	"fmt"
	"net/http"
)

func Main() {
	http.HandleFunc("/uploadFile", api.UploadFile)
	http.HandleFunc("/downloadFile", api.DownloadFile)
	fmt.Println("Running server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
