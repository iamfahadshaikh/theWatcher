# theWatcher - Architecture & Technical Specification

`theWatcher` is a single-executable, production-grade Linux terminal observability platform built with Go, BubbleTea, Lipgloss, and Bubbles.

---

## 1. Data Pipeline & System Architecture

```text
/proc, ss, docker stats
         │
         ▼
    Collectors (internal/collector)
         │
         ▼
     Metrics (internal/metrics)
         │
         ▼
   Ring Buffer (internal/history)
         │
         ▼
  Graph Renderer (internal/graph)
         │
         ▼
      Widgets (internal/widgets)
         │
         ▼
      Panels (internal/panel)
         │
         ▼
 Layout Engine & Focus Manager (internal/layout)
         │
         ▼
    Dashboard View (internal/app)
```

1. **Collectors**: Poll raw Linux `/proc` files (`/proc/stat`, `/proc/meminfo`, `/proc/net/dev`, `/proc/diskstats`, `/proc/[pid]/`), socket tables (`/proc/net/tcp`), and `docker stats`.
2. **Metrics**: Typed domain structs holding snapshot values.
3. **Ring Buffer**: Generic in-memory circular buffer storing up to 300 samples (150s at 500ms refresh).
4. **Graph Renderer**: Abstract renderer interface accepting `GraphData` and rendering high-resolution Braille continuous curves (`⣿`), sparklines, or block bars.
5. **Widgets**: Standardized TUI elements (Header, Footer, Tables, Badges, Progress Bars).
6. **Panels**: Encapsulated UI components with titles, borders, and active focus indicators.
7. **Layout & Focus Manager**: Computes 2D grid dimensions and manages active panel focus + full-screen zoom toggle.

---

## 2. Directory Structure

```text
theWatcher/
cmd/
    theWatcher/
        main.go           # Entry point (< 50 lines)
docs/
    architecture.md       # Architecture specification
internal/
    app/                  # Application runner and tea.Model
    collector/            # Metric collectors implementing Collector interface
    config/               # Constants (RefreshRate, HistorySize, Version)
    dashboard/            # Central dashboard orchestrator
    graph/                # Graph rendering engine (Renderer, Braille, Spark, Bar, Axis)
    history/              # Generic RingBuffer[T]
    layout/               # 2D Grid layout and FocusManager
    metrics/              # Strongly typed metric structs
    panel/                # Panel containers and views
    theme/                # Colors, Palette, Icons, Styles, Spacing, Borders
    widgets/              # Table, Header, Footer, Badge, Progress
    util/                 # Helper utilities
```

---

## 3. Interfaces & Design Patterns

### Collector Interface
```go
type Collector interface {
    Name() string
    Collect() error
    Snapshot() any
}
```

### Graph Renderer Interface
```go
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

type Renderer interface {
    Render(data GraphData) string
}
```

---

## 4. Keybindings & Interactions

| Key | Action |
| --- | --- |
| `q` / `Ctrl+C` | Quit application cleanly |
| `Tab` / `Shift+Tab` | Cycle panel focus forward/backward |
| `Space` | Toggle full-screen zoom on focused panel |
| `Esc` | Restore standard multi-panel grid view |
| `1` - `5` | Switch timeframe zoom (30s, 1m, 5m, 15m, 1h) |
| `?` / `h` | Toggle help modal overlay |

---

## 5. Development Guidelines & Constraints
- Maximum file size: **300 lines**.
- Entry point `main.go`: **< 50 lines**.
- Zero persistent database or daemon requirement.
- Pure in-memory state; disappearing completely upon exit.
