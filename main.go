// go-redis
package main

import (
	"fmt"
	"net"
	"strings"
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

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	// To sync in memory DB from AOF
	fmt.Println("Reading from AOF")
	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	// Listen for connections - as of now: single connection only
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
		fmt.Println("value: ", value)

		// ensure we have an array, and it's not empty
		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// command is the first element of the array, the rest are the arguments
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		// ok is an idiom to check the success of a function
		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: "Invalid command: " + command})
			continue
		}

		// write to AOF if the command is SET or HSET - GET, HGET, HGETALL are read only (clog up aof)
		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
