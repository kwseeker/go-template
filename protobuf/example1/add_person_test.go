package example1

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	//"google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/proto"
)

func TestPromptForAddress(t *testing.T) {
	in := `12345
Example Name
name@example.com
123-456-7890
home
222-222-2222
mobile
111-111-1111
work
777-777-7777
unknown

`
	got, err := promptForAddress(strings.NewReader(in))
	if err != nil {
		t.Fatalf("promptForAddress(%q) had unexpected error: %s", in, err.Error())
	}
	if got.Id != 12345 {
		t.Errorf("promptForAddress(%q) got %d, want ID %d", in, got.Id, 12345)
	}
	if got.Name != "Example Name" {
		t.Errorf("promptForAddress(%q) => want name %q, got %q", in, "Example Name", got.Name)
	}
	if got.Email != "name@example.com" {
		t.Errorf("promptForAddress(%q) => want email %q, got %q", in, "name@example.com", got.Email)
	}

	want := []*Person_PhoneNumber{
		{Number: "123-456-7890", Type: Person_PHONE_TYPE_HOME},
		{Number: "222-222-2222", Type: Person_PHONE_TYPE_MOBILE},
		{Number: "111-111-1111", Type: Person_PHONE_TYPE_WORK},
		{Number: "777-777-7777", Type: Person_PHONE_TYPE_UNSPECIFIED},
	}
	if len(got.Phones) != len(want) {
		t.Errorf("want %d phone numbers, got %d", len(want), len(got.Phones))
	}
	phones := len(got.Phones)
	if phones > len(want) {
		phones = len(want)
	}
	for i := 0; i < phones; i++ {
		if !proto.Equal(got.Phones[i], want[i]) {
			t.Errorf("want phone %q, got %q", want[i], got.Phones[i])
		}
	}

	bytes, err := proto.Marshal(got)
	cloned := &Person{}
	err = proto.Unmarshal(bytes, cloned)

	fn := "address_book.txt"
	if err := ioutil.WriteFile(fn, bytes, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}

	fi, err := ioutil.ReadFile(fn)
	personUnmarshal := &Person{}
	err = proto.Unmarshal(fi, personUnmarshal)

	descriptor := proto.MessageV2(personUnmarshal).ProtoReflect().Descriptor()
	_ = descriptor.Options().(*descriptorpb.MessageOptions)

}
