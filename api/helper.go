package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func IsJsonRequest(contentType string) (bool, error) {
	if strings.Contains(contentType, "form") {
		return false, nil
	} else if strings.Contains(contentType, "json") {
		return true, nil
	}
	return false, errors.New("unsupported content type")
}

func ReturnError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorJsonResponse{Error: err.Error()})
}

func GetFileFromJson(r *http.Request) ([]byte, string, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, "nil", err
	}

	var jsonData UploadFileJsonRequest
	err = json.Unmarshal(reqBody, &jsonData)
	if err != nil {
		return nil, "nil", err
	}

	getResponse, err := http.Get(jsonData.File)
	if err != nil {
		return nil, "nil", err
	}

	fileData, err := ioutil.ReadAll(getResponse.Body)
	if err != nil {
		return nil, "nil", err
	}

	tSplit := strings.Split(jsonData.File, "/")
	fileName := tSplit[len(tSplit)-1]

	return fileData, fileName, nil
}

func GetFileFromForm(r *http.Request) ([]byte, string, error) {
	err := r.ParseMultipartForm(0)
	if err != nil {
		return nil, "nil", err
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, "nil", err
	}
	fileName := fileHeader.Filename

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, "nil", err
	}

	return fileData, fileName, nil
}

