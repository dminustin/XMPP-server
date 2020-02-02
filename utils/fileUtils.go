package utils

import (
	"encoding/base64"
	"io/ioutil"
	"log"
)

func Base64ReadFile(filename string) string {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println(err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(content)

}
