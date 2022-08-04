package api

import (
	"concurrent-http-server/pkg"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UploadFileJsonRequest struct {
	File string `json:"file"`
}

type UploadFileJsonResponse struct {
	FileId string `json:"file_id"`
}

type DownloadFileJsonRequest struct {
	FileId string `json:"file_id"`
}

type ErrorJsonResponse struct {
	Error string `json:"error"`
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: /uploadFile")

	isJson, err := IsJsonRequest(r.Header.Get("Content-Type"))
	if err != nil {
		ReturnError(w, err)
		fmt.Println(err)
		return
	}

	var fileData []byte
	var fileName string
	if isJson {
		fileData, fileName, err = GetFileFromJson(r)
	} else {
		fileData, fileName, err = GetFileFromForm(r)
	}
	if err != nil {
		ReturnError(w, err)
		fmt.Println(err)
		return
	}

	hash := GetFileHash(fileData) + ":" + base64.URLEncoding.EncodeToString([]byte(fileName))
	err = pkg.ConcurrentWrite("./files/"+hash, fileData)
	if err != nil {
		ReturnError(w, err)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadFileJsonResponse{FileId: hash})
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: /downloadFile")

	isJson, err := IsJsonRequest(r.Header.Get("Content-Type"))
	if err != nil {
		ReturnError(w, err)
		fmt.Println(err)
		return
	}

	var fileId string
	if isJson {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ReturnError(w, err)
			fmt.Println(err)
			return
		}

		var jsonData DownloadFileJsonRequest
		err = json.Unmarshal(reqBody, &jsonData)
		if err != nil {
			ReturnError(w, err)
			fmt.Println(err)
			return
		}
		fileId = jsonData.FileId
	} else {
		err = r.ParseMultipartForm(0)
		if err != nil {
			ReturnError(w, err)
			fmt.Println(err)
			return
		}

		fileId = r.FormValue("file_id")
	}

	http.ServeFile(w, r, "./files/"+fileId)
}
