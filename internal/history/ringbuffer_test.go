package history

import (
	"reflect"
	"testing"
)

func TestRingBufferBasic(t *testing.T) {
	buf := NewRingBuffer[int](3)
	if buf.Size() != 0 {
		t.Errorf("expected size 0, got %d", buf.Size())
	}

	buf.Push(1)
	buf.Push(2)
	if buf.Size() != 2 {
		t.Errorf("expected size 2, got %d", buf.Size())
	}

	latest, ok := buf.Latest()
	if !ok || latest != 2 {
		t.Errorf("expected latest 2, got %d", latest)
	}

	expected := []int{1, 2}
	if !reflect.DeepEqual(buf.Slice(), expected) {
		t.Errorf("expected slice %v, got %v", expected, buf.Slice())
	}
}

func TestRingBufferOverwrite(t *testing.T) {
	buf := NewRingBuffer[int](3)
	buf.Push(10)
	buf.Push(20)
	buf.Push(30)
	buf.Push(40) // Overwrites 10

	if buf.Size() != 3 {
		t.Errorf("expected size 3 after overflow, got %d", buf.Size())
	}

	expected := []int{20, 30, 40}
	if !reflect.DeepEqual(buf.Slice(), expected) {
		t.Errorf("expected slice %v, got %v", expected, buf.Slice())
	}
}
