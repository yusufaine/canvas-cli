package main

import (
	"github.com/yusufaine/canvas-cli/internal/app/canvas"
)

func main() {
	config := canvas.NewConfig()
	canvas.Start(config)
}
