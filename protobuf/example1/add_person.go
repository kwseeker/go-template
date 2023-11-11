package example1

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func promptForAddress(r io.Reader) (*Person, error) {
	// A protocol buffer can be created like any struct.
	p := &Person{}

	rd := bufio.NewReader(r)
	fmt.Print("Enter person ID number: ")
	// An int32 field in the .proto file is represented as an int32 field
	// in the generated Go struct.
	if _, err := fmt.Fscanf(rd, "%d\n", &p.Id); err != nil {
		return p, err
	}

	fmt.Print("Enter name: ")
	name, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	// A string field in the .proto file results in a string field in Go.
	// We trim the whitespace because rd.ReadString includes the trailing
	// newline character in its output.
	p.Name = strings.TrimSpace(name)

	fmt.Print("Enter email address (blank for none): ")
	email, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Email = strings.TrimSpace(email)

	for {
		fmt.Print("Enter a phone number (or leave blank to finish): ")
		phone, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}
		// The PhoneNumber message type is nested within the Person
		// message in the .proto file.  This results in a Go struct
		// named using the name of the parent prefixed to the name of
		// the nested message.  Just as with pb.Person, it can be
		// created like any other struct.
		pn := &Person_PhoneNumber{
			Number: phone,
		}

		fmt.Print("Is this a mobile, home, or work phone? ")
		ptype, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		ptype = strings.TrimSpace(ptype)

		// A proto enum results in a Go constant for each enum value.
		switch ptype {
		case "mobile":
			pn.Type = Person_PHONE_TYPE_MOBILE
		case "home":
			pn.Type = Person_PHONE_TYPE_HOME
		case "work":
			pn.Type = Person_PHONE_TYPE_WORK
		default:
			fmt.Printf("Unknown phone type %q.  Using default.\n", ptype)
		}

		// A repeated proto field maps to a slice field in Go.  We can
		// append to it like any other slice.
		p.Phones = append(p.Phones, pn)
	}

	return p, nil
}
