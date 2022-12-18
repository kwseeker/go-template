package tcp

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestTCPClient(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("err dialing:", err.Error())
		return
	}
	defer conn.Close()

	//inputReader := bufio.NewReader(os.Stdin)
	for {
		//input, _ := inputReader.ReadString('\n')
		input := time.Now().String()
		time.Sleep(2 * time.Second)
		trimmedInput := strings.Trim(input, "\r\n")
		if trimmedInput == "Q" {
			return
		}
		_, err := conn.Write([]byte(trimmedInput))
		if err != nil {
			fmt.Println("err conn.write:", err)
			return
		}
	}
}
