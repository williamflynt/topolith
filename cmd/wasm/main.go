//go:build js && wasm

package main

import (
	"encoding/json"
	"github.com/williamflynt/topolith/pkg/app"
	"github.com/williamflynt/topolith/pkg/grammar"
	"github.com/williamflynt/topolith/pkg/world"
	"syscall/js"
)

type JavaScriptReply struct {
	Status int               `json:"status"`
	Data   grammar.Response  `json:"data"`
	Error  map[string]string `json:"error"`
	Raw    string            `json:"raw"`
}

func main() {
	app, err := app.NewApp(world.CreateWorld("default-world"))
	if err != nil {
		panic(err)
	}
	// Create a JavaScript function that can be called from JavaScript to send commands to the App.
	js.Global().Set("sendCommand", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Expected one argument"
		}
		command := args[0].String()
		response := app.Exec(command)
		return jsReplyToJsValue(toReply(response))
	}))
	// Prevent the Go program from exiting.
	select {}
}

func toReply(response string) JavaScriptReply {
	p, err := grammar.Parse(response)
	// TODO: Structure the error as a TopolithError.
	if err != nil {
		return JavaScriptReply{
			Status: 500,
			Data:   p.Response,
			Error:  map[string]string{"error": err.Error()},
			Raw:    response,
		}
	}
	// TODO: Parse the response object Repr and return objects in JSON.
	//  This will mean creating a structure for our World.Tree.
	return JavaScriptReply{
		Status: p.Response.Status.Code,
		Data:   p.Response,
		Error:  map[string]string{"error": ""},
		Raw:    response,
	}
}

func jsReplyToJsValue(reply JavaScriptReply) js.Value {
	data, _ := json.Marshal(reply)
	return js.ValueOf(string(data))
}
