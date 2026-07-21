package graph

// GraphData defines input data payload for all graph renderers.
type GraphData struct {
	Samples    []float64
	Width      int
	Height     int
	Min        float64
	Max        float64
	Title      string
	Unit       string
	ShowAxis   bool
	ShowLegend bool
	Color      string
}

// Renderer interface for graph visualization components.
type Renderer interface {
	Render(data GraphData) string
}
