package graph

import (
	"time"
	"theWatcher/internal/history"
)

// RenderHistory renders a RingBuffer history using the provided Renderer.
func RenderHistory(buf *history.RingBuffer[history.Sample], renderer Renderer, data GraphData) string {
	if buf != nil {
		slice := buf.Slice()
		samples := make([]float64, len(slice))
		timestamps := make([]time.Time, len(slice))
		for i, s := range slice {
			samples[i] = s.Value
			timestamps[i] = s.Time
		}
		data.Samples = samples
		data.Timestamps = timestamps
	}
	return renderer.Render(data)
}
