package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func PrettyPrint(v interface{}, title ...string) (string, error) {
	name := "log"
	if len(title) > 0 {
		name = strings.Join(title, ", ")
	}
	marshalData, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Printf("================ error : %s ================\n", name)
		fmt.Println(err.Error())
		fmt.Printf("================ end : %s ================\n", name)
		return "", err
	} else {
		fmt.Printf("================ start : %s ================\n", name)
		fmt.Println(string(marshalData))
		fmt.Printf("================ end : %s ================\n", name)
		return string(marshalData), nil
	}
}

const (
	MODE_DEBUG = iota
	MODE_RELEASE
)

func RenderError(err error, replace string, mode int) error {
	if mode == MODE_DEBUG {
		log.Printf("ERROR : %s", err.Error())
		return err
	}
	return errors.New(replace)
}

type Action struct {
}

type JsonResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Trace   interface{} `json:"trace,omitempty"`
}

func setHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/sjon")
	return w
}

func (a *Action) WriteJson(w http.ResponseWriter, v interface{}, statusCode ...int) {
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

func (a *Action) ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
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

func (a *Action) PrintLog(msg string, v interface{}) {
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

var (
	ErrorColor   = "\033[31m"
	SuccessColor = "\033[32m"
	InfoColor    = "\033[36m"
	ResetCOlor   = "\033[0m"
)

type Logger struct {
}

func init() {
	if runtime.GOOS == "windows" {
		ErrorColor = ""
		SuccessColor = ""
		InfoColor = ""
		ResetCOlor = ""
	}
}

func (l *Logger) LogInfo(args ...interface{}) {
	fmt.Println(InfoColor, args, ResetCOlor)
}

func (l *Logger) LogError(args ...interface{}) {
	fmt.Println(ErrorColor, args, ResetCOlor)
}
func (l *Logger) LogSuccess(args ...interface{}) {
	fmt.Println(SuccessColor, args, ResetCOlor)
}

type RandomOption struct {
	LowerCase   bool
	UpperCase   bool
	Numeric     bool
	SpecialChar bool
}

var (
	LowerCase        = "abcdefghijklmnopqrstuvwxyz"
	UpperCase        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numeric          = "1234567890"
	SpecialCharacter = "!@#$%^&*()_+-."
)

func (a *RandomOption) GenerateRandomString(length int) (string, error) {
	if length == 0 {
		return "", errors.New("minimum 1 character's length")
	}
	minChar := 0
	rawSeed := ""
	result := ""

	if a.LowerCase {
		minChar++
		rawSeed += LowerCase
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(LowerCase)[rand.Intn(len(LowerCase))])
	}
	if a.UpperCase {
		minChar++
		rawSeed += UpperCase
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(UpperCase)[rand.Intn(len(UpperCase))])
	}
	if a.Numeric {
		minChar++
		rawSeed += Numeric
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(Numeric)[rand.Intn(len(Numeric))])
	}
	if a.SpecialChar {
		minChar++
		rawSeed += SpecialCharacter
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(SpecialCharacter)[rand.Intn(len(SpecialCharacter))])
	}
	if minChar == 0 {
		b := make([]rune, length)
		rand.Seed(time.Now().UnixNano())
		for i := range b {
			b[i] = []rune(SpecialCharacter)[rand.Intn(len(SpecialCharacter))]
		}
		return string(b), nil
	}

	if minChar > length {
		return "", errors.New("min char is " + strconv.Itoa(minChar))
	}
	for i := minChar; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		result += string([]rune(rawSeed)[rand.Intn(len(rawSeed))])
	}
	return string(ShuffleBytes([]byte(result))), nil

}

func ShuffleBytes(src []byte) []byte {
	final := make([]byte, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}
