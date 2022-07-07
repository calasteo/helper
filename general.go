package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
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
