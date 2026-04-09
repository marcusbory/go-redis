package main

// map of command name (string) to handler function (func([]Value) Value)
var Handlers = map[string]func([]Value) Value{
	"PING": ping,
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}
