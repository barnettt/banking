package util

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

const ContentTypeJson string = "application/json"
const ContentTypeXml string = "application/xml"

func WriteResponse(writer http.ResponseWriter, code int, data interface{}, contentType string) {
	writer.Header().Add("Content-Type", contentType)
	writer.WriteHeader(code)
	if contentType == ContentTypeXml {
		err := xml.NewEncoder(writer).Encode(data)
		if err != nil {
			panic(err)
		}
		return
	}
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		panic(err)
	}
}
