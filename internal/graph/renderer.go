package graph

import "time"

// GraphData defines input data payload for all graph renderers.
type GraphData struct {
	Samples      []float64
	Timestamps   []time.Time
	Width        int
	Height       int
	Min          float64
	Max          float64
	Title        string
	Unit         string
	ShowAxis     bool
	ShowLegend   bool
	Color        string
	ShowRefLines bool
	AvgVal       float64
	MinVal       float64
	MaxVal       float64
	MultiSamples [][]float64
	MultiColors  []string
	AvgColor     string
	MinColor     string
	MaxColor     string
}

// Renderer interface for graph visualization components.
type Renderer interface {
	Render(data GraphData) string
}
