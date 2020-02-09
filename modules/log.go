package modules

import (
	"fmt"
	"log"
	"os"
)

type AppLogStruct struct {
	LogType    string
	LogMessage string
	LogData    interface{}
}

func WriteQueChan(id string, data string) {
	//Do nothing
	if true {
		return
	}
	filename := "./tmp/" + id + ".txt"

	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		os.Remove(filename)
	} else {
		appendFile("./tmp/"+id+".txt", data+"\n")
	}
}

func appendFile(filename string, s string) {
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	f.Write([]byte(s))
}

func (a *AppLogStruct) WriteAppLog() {
	return
	f, err := os.OpenFile("./text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("[%s]: %s\t%s\n", a.LogType, a.LogMessage, a.LogData)); err != nil {
		log.Println(err)
	}
}

func init() {
}
