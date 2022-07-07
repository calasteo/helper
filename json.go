package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Action struct {
}
type JsonResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    any            `json:"data,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}
type ResponseError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Trace   any    `json:"trace,omitempty"`
}

func setHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/sjon")
	return w
}

func (a *Action) WriteJson(w http.ResponseWriter, v any, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	mData, err := json.Marshal(v)
	if err != nil {
		code = http.StatusBadRequest
		write := setHeader(w)
		write.WriteHeader(code)
		var jResponse JsonResponse
		jResponse.Success = false
		jResponse.Message = "Error writing message"
		jResponse.Error.Code = http.StatusBadRequest
		jResponse.Error.Message = err.Error()
		jResponse.Error.Trace = fmt.Sprintf("%+v", err)
		res, _ := json.Marshal(jResponse)
		write.Write(res)
		return
	}
	write := setHeader(w)
	write.WriteHeader(code)
	write.Write(mData)
	return
}

func (a *Action) ErrJson(w http.ResponseWriter, message string, err error, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	write := setHeader(w)
	write.WriteHeader(code)
	var jResponse JsonResponse
	jResponse.Success = false
	jResponse.Message = message
	jResponse.Error.Code = http.StatusBadRequest
	jResponse.Error.Message = err.Error()
	jResponse.Error.Trace = fmt.Sprintf("%+v", err)
	res, _ := json.Marshal(jResponse)
	write.Write(res)
	return
}

func (a *Action) ReadJSON(w http.ResponseWriter, r *http.Request, v any) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("invalid request data")
		}
		return err
	}
	return nil
}

func (a *Action) PrintLog(msg string, v any) {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		now := time.Now().Format("2006-01-02 15:04:05")
		log.Printf("%s | error marshalling data: %+v\n", now, v)
	} else {
		fmt.Println("===================================")
		fmt.Println(msg)
		fmt.Println(string(data))
		fmt.Println("===================================")
	}
}
