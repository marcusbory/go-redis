// go-redis
package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port 6379")
	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	// will close when main returns
	defer l.Close()

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	// will close when main returns - LIFO
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		_ = value

		// this writer will write the response back to the client
		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
