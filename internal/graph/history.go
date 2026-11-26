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

// RenderMultiHistory renders multiple RingBuffer histories using the provided Renderer.
func RenderMultiHistory(buffers []*history.RingBuffer[history.Sample], renderer Renderer, data GraphData) string {
	if len(buffers) > 0 {
		var multiSamples [][]float64
		var timestamps []time.Time

		for _, buf := range buffers {
			if buf != nil {
				slice := buf.Slice()
				samples := make([]float64, len(slice))
				for i, s := range slice {
					samples[i] = s.Value
				}
				multiSamples = append(multiSamples, samples)

				// Keep timestamps from the longest buffer
				if len(slice) > len(timestamps) {
					timestamps = make([]time.Time, len(slice))
					for i, s := range slice {
						timestamps[i] = s.Time
					}
				}
			}
		}
		data.MultiSamples = multiSamples
		data.Timestamps = timestamps
	}
	return renderer.Render(data)
}
