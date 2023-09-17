package canvashttp

import "net/http"

type canvasClientOpts func(*CanvasClient)

func WithClient(c *http.Client) canvasClientOpts {
	return func(cc *CanvasClient) {
		cc.client = c
	}
}
