package app

import (
	"context"
	"testing"
)

func BenchmarkSequential(b *testing.B) {
	c := NewCollectors()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.System.Collect(ctx)
		_ = c.CPU.Collect(ctx)
		_ = c.Memory.Collect(ctx)
		_ = c.Network.Collect(ctx)
		_ = c.Disk.Collect(ctx)
		_ = c.Process.Collect(ctx)
		_ = c.Ports.Collect(ctx)
		_ = c.Conn.Collect(ctx)
		_ = c.Login.Collect(ctx)

		_ = c.System.Snapshot()
		_ = c.CPU.Snapshot()
		_ = c.Memory.Snapshot()
		_ = c.Network.Snapshot()
		_ = c.Disk.Snapshot()
		_ = c.Process.Snapshot()
		_ = c.Ports.Snapshot()
		_ = c.Conn.Snapshot()
		_ = c.Login.Snapshot()
	}
}

func BenchmarkSchedulerStore(b *testing.B) {
	store := NewMetricsStore()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.GetSystem()
		_ = store.GetCPU()
		_ = store.GetMemory()
		_ = store.GetNetwork()
		_ = store.GetDisk()
		_ = store.GetProcess()
		_ = store.GetPorts()
		_ = store.GetConn()
		_ = store.GetLogin()
	}
}
