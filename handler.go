package main

// map of command name (string) to handler function (func([]Value) Value)
var Handlers = map[string]func([]Value) Value{
	"PING": ping,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}
