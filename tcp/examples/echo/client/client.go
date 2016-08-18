package main

import (
	"fmt"
	"log"
	"net"
	"github.com/ihaiker/gokit/tcp/examples/echo"
	"time"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)


	// ping <--> pong
	for i := 0; i < 20; i++ {
		// write
		bytes := echo.NewEchoPacket([]byte("hello")).Serialize();
		conn.Write(bytes)
		// read

		for {
			lengthBytes := make([]byte, 10)
			readed, err := conn.Read(lengthBytes)
			if err != nil || readed == 0 {
				fmt.Println("over .....")
				return
			}
			fmt.Println(string(lengthBytes[0:readed]))
		}

		time.Sleep(time.Second)
	}
	time.Sleep(time.Minute)
	fmt.Println("over")
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
