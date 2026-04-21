//go:build js && wasm

package main

import (
	"bytes"
	"syscall/js"

	"github.com/monjiapawne/svg-to-graphml/converter"
)

func main() {
	js.Global().Set("convertSVG", js.FuncOf(convert))
	select {}
}

func convert(_ js.Value, args []js.Value) any {
	data := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(data, args[0])

	node, err := converter.NewNodeFromBytes(data)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}

	var buf bytes.Buffer
	nodes := converter.Nodes{node}
	if err := nodes.RenderTemplate(&buf); err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{"result": buf.String()}
}
