package canvashttp

import "net/http"

type canvasClientOpts func(*CanvasClient)

func WithClient(c *http.Client) canvasClientOpts {
	return func(cc *CanvasClient) {
		cc.client = c
	}
}

// Default: "canvas.nus.edu.sg"
func WithHost(host string) canvasClientOpts {
	return func(cc *CanvasClient) {
		cc.canvasHost = host
	}
}

// Default: "/api/v1"
func WithPathApiPrefix(api string) canvasClientOpts {
	return func(cc *CanvasClient) {
		cc.apiPath = api
	}
}
