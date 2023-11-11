package jsonpb

import (
	"github.com/golang/protobuf/jsonpb"
	"io/ioutil"
	"testing"
)

func parseJsonFile(fn string) (*MyMessage, error) {
	bytesContent, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	myMessage := &MyMessage{}
	if err = jsonpb.UnmarshalString(string(bytesContent), myMessage); err != nil {
		return nil, err
	}
	return myMessage, nil
}

func TestJsonParse(t *testing.T) {
	myMessage, err := parseJsonFile("config.json")
	if err != nil {
		t.Error("failed to parse json file")
	}

	if myMessage.Name != "Arvin" || myMessage.Age != 18 || len(myMessage.Contacts) != 2 {
		t.Error("failed to parse json file")
	}
	email := myMessage.Contacts[1]
	if email.Type != int32(ContactType_Email) || email.Addr != "xxx@gmail.com" {
		t.Error("failed to parse json file")
	}
}
