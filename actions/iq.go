package actions

import (
	"encoding/xml"
	"fmt"
	"os"
)

type dataFormat_Iq struct {
	IQ []struct {
		Bind struct {

		}`xml:"Bind" json:"bind"`
	} `xml:"Iq" json:"iq"`
}



func ActinIQ() {
	s := "<iq type=\"set\" id=\"5bbf3932-701f-4f82-89c0-3e180758e53b\">" +
		"<bind xmlns=\"urn:ietf:params:xml:ns:xmpp-bind\"><resource>gajim.5ET60J5M</resource></bind></iq>"
	var inData = []byte(s)
	data := &dataFormat_Iq{}
	err := xml.Unmarshal(inData, data)
	if nil != err {
		fmt.Println("Error unmarshalling from XML", err)
	}
	fmt.Println(data, err)
	os.Exit(1)

}
