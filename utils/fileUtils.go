package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
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

func Base64ToSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}
