package graph

import (
	"theWatcher/internal/history"
)

// RenderHistory renders a RingBuffer history using the provided Renderer.
func RenderHistory(buf *history.RingBuffer[float64], renderer Renderer, data GraphData) string {
	if buf != nil {
		data.Samples = buf.Slice()
	}
	return renderer.Render(data)
}
