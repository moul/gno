// Package txbench is a tiny explicit stage timer for the transaction
// benchmark harness. It is off by default and costs one bool check per call
// site when off. It assumes a single goroutine drives DeliverTx/Commit, which
// is what the harness does; it is not safe for concurrent use.
//
// Usage at a call site:
//
//	defer txbench.End(txbench.Begin("stage.name"))
//
// Stages nest. Each stage records inclusive time and self time (inclusive
// minus the inclusive time of its direct children), and, when Allocs is on,
// inclusive and self heap allocation bytes and objects.
package txbench

import (
	"runtime"
	"time"
)

var (
	// On enables timing.
	On bool
	// Allocs enables per-stage allocation accounting through
	// runtime.ReadMemStats. It stops the world at every boundary, so timings
	// taken with Allocs on are perturbed and must not be reported as latency.
	Allocs bool
)

type Stat struct {
	N                    int64
	Incl, Self           time.Duration
	InclBytes, SelfBytes uint64
	InclObjs, SelfObjs   uint64
}

type frame struct {
	name                 string
	start                time.Time
	childDur             time.Duration
	bytes0, objs0        uint64
	childBytes, childObj uint64
}

var (
	stack []frame
	stats = map[string]*Stat{}
	ms    runtime.MemStats
)

func readAllocs() (uint64, uint64) {
	runtime.ReadMemStats(&ms)
	return ms.TotalAlloc, ms.Mallocs
}

// Begin opens a stage and returns a token for End.
func Begin(name string) int {
	if !On {
		return -1
	}
	f := frame{name: name}
	if Allocs {
		f.bytes0, f.objs0 = readAllocs()
	}
	f.start = time.Now()
	stack = append(stack, f)
	return len(stack) - 1
}

// End closes the stage opened by the Begin that returned tok, and any stage
// left open above it (a panic that skipped an End).
func End(tok int) {
	if !On || tok < 0 || tok >= len(stack) {
		return
	}
	now := time.Now()
	var b1, o1 uint64
	if Allocs {
		b1, o1 = readAllocs()
	}
	for len(stack) > tok {
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		d := now.Sub(f.start)
		st := stats[f.name]
		if st == nil {
			st = &Stat{}
			stats[f.name] = st
		}
		st.N++
		// A recursive stage (an object load that loads its children) counts
		// inclusive time once, at its outermost frame.
		if !open(f.name) {
			st.Incl += d
		}
		st.Self += d - f.childDur
		var db, do uint64
		if Allocs {
			db, do = b1-f.bytes0, o1-f.objs0
			if !open(f.name) {
				st.InclBytes += db
				st.InclObjs += do
			}
			st.SelfBytes += db - f.childBytes
			st.SelfObjs += do - f.childObj
		}
		if len(stack) > 0 {
			p := &stack[len(stack)-1]
			p.childDur += d
			p.childBytes += db
			p.childObj += do
		}
	}
}

// Take returns the stats accumulated since the last Take and resets them.
func Take() map[string]Stat {
	out := make(map[string]Stat, len(stats))
	for k, v := range stats {
		out[k] = *v
	}
	clear(stats)
	stack = stack[:0]
	return out
}

// open reports whether a frame named name is still on the stack.
func open(name string) bool {
	for i := range stack {
		if stack[i].name == name {
			return true
		}
	}
	return false
}
