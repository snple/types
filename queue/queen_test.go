package queue

import "testing"

func TestQueueSimple(t *testing.T) {
	q := New[int]()

	for i := 0; i < minQueueLen; i++ {
		q.Push(i)
	}
	for i := 0; i < minQueueLen; i++ {
		n := q.Peek()
		if n.Get() != i {
			t.Error("peek", i, "had value", q.Peek())
		}
		x := q.Pop()
		if x.Get() != i {
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestQueueWrapping(t *testing.T) {
	q := New[int]()

	for i := 0; i < minQueueLen; i++ {
		q.Push(i)
	}
	for i := 0; i < 3; i++ {
		q.Pop()
		q.Push(minQueueLen + i)
	}

	for i := 0; i < minQueueLen; i++ {
		n := q.Peek()
		if n.Get() != i+3 {
			t.Error("peek", i, "had value", q.Peek())
		}
		q.Pop()
	}
}

func TestQueueLength(t *testing.T) {
	q := New[int]()

	if q.Length() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		q.Push(i)
		if q.Length() != i+1 {
			t.Error("adding: queue with", i, "elements has length", q.Length())
		}
	}
	for i := 0; i < 1000; i++ {
		q.Pop()
		if q.Length() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", q.Length())
		}
	}
}

func TestQueueGet(t *testing.T) {
	q := New[int]()

	for i := 0; i < 1000; i++ {
		q.Push(i)
		for j := 0; j < q.Length(); j++ {
			n := q.Get(j)
			if n.Get() != j {
				t.Errorf("index %d doesn't contain %d", j, j)
			}
		}
	}
}

func TestQueueGetNegative(t *testing.T) {
	q := New[int]()

	for i := 0; i < 1000; i++ {
		q.Push(i)
		for j := 1; j <= q.Length(); j++ {
			n := q.Get(-j)
			if n.Get() != q.Length()-j {
				t.Errorf("index %d doesn't contain %d", -j, q.Length()-j)
			}
		}
	}
}

func TestQueueGetOutOfRangePanics(t *testing.T) {
	q := New[int]()

	q.Push(1)
	q.Push(2)
	q.Push(3)

	assertPanics(t, "should panic when negative index", func() {
		n := q.Get(-4)
		n.Get()
	})

	assertPanics(t, "should panic when index greater than length", func() {
		n := q.Get(4)
		n.Get()
	})
}

func TestQueuePeekOutOfRangePanics(t *testing.T) {
	q := New[int]()

	assertPanics(t, "should panic when peeking empty queue", func() {
		n := q.Peek()
		n.Get()
	})

	q.Push(1)
	q.Pop()

	assertPanics(t, "should panic when peeking emptied queue", func() {
		n := q.Peek()
		n.Get()
	})
}

func TestQueueRemoveOutOfRangePanics(t *testing.T) {
	q := New[int]()

	assertPanics(t, "should panic when removing empty queue", func() {
		n := q.Pop()
		n.Get()
	})

	q.Push(1)
	q.Pop()

	assertPanics(t, "should panic when removing emptied queue", func() {
		n := q.Pop()
		n.Get()
	})
}

func assertPanics(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: didn't panic as expected", name)
		}
	}()

	f()
}

// General warning: Go's benchmark utility (go test -bench .) increases the number of
// iterations until the benchmarks take a reasonable amount of time to run; memory usage
// is *NOT* considered. On my machine, these benchmarks hit around ~1GB before they've had
// enough, but if you have less than that available and start swapping, then all bets are off.

func BenchmarkQueueSerial(b *testing.B) {
	q := New[interface{}]()
	for i := 0; i < b.N; i++ {
		q.Push(nil)
	}
	for i := 0; i < b.N; i++ {
		q.Peek()
		q.Pop()
	}
}

func BenchmarkQueueGet(b *testing.B) {
	q := New[int]()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Get(i)
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := New[interface{}]()
	for i := 0; i < b.N; i++ {
		q.Push(nil)
		q.Peek()
		q.Pop()
	}
}
