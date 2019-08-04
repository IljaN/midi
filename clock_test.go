package midi

import (
	"bytes"
	"testing"
	"time"
)

func NewTestClock() (clk *Clock, protocolLog *bytes.Buffer) {
	var buf = new(bytes.Buffer)
	buf.Grow(128)
	return NewClock(buf), buf
}

func TestStop(t *testing.T) {
	clk, protoLog := NewTestClock()
	clk.Stop()
	midiStream := protoLog.Bytes()

	if numOfOccurences(midiStream, Stop) == 0 {
		t.Error("No stop command was sent")
	}

	if numOfOccurences(midiStream, Stop) > 1 {
		t.Error("More than one stop command was sent")
	}
}

func TestContinue(t *testing.T) {
	clk, protoLog := NewTestClock()
	clk.Continue()
	midiStream := protoLog.Bytes()

	if numOfOccurences(midiStream, Continue) == 0 {
		t.Error("No continue command was sent")
	}

	if numOfOccurences(midiStream, Stop) > 1 {
		t.Error("More than one continue command was sent")
	}
}

func TestTimings(t *testing.T) {
	t.Run("120 BPM in 1s   ", func(t *testing.T) { assertClockRateFor(120.00, 47, time.Duration(1)*time.Second, t) })
	t.Run("60  BPM in 1s   ", func(t *testing.T) { assertClockRateFor(60.00, 24, time.Duration(1)*time.Second, t) })
	t.Run("60  BPM in 500ms", func(t *testing.T) { assertClockRateFor(60.00, 11, time.Duration(500)*time.Millisecond, t) })
	t.Run("60  BPM in 100ms", func(t *testing.T) { assertClockRateFor(60.00, 2, time.Duration(100)*time.Millisecond, t) })
	t.Run("270 BPM in 1s   ", func(t *testing.T) { assertClockRateFor(270.00, 107, time.Duration(1)*time.Second, t) })
	t.Run("270 BPM in 436ms", func(t *testing.T) { assertClockRateFor(270.00, 47, time.Duration(436)*time.Millisecond, t) })
}

func assertClockRateFor(bpm Bpm, expectedTickCount int, d time.Duration, t *testing.T) {
	clk, protoLog := NewTestClock()
	clk.Stop()
	clk.SetTempo(bpm)
	protoLog.Reset()
	time.Sleep(d)

	var tickCount = len(protoLog.Bytes())

	//  tick results vary between test round, either bug, rounding problem or test timing problem
	if abs(expectedTickCount-tickCount) > 1 {
		t.Errorf("Failed to assert that the tick-count with %v bpm equals %v (actual: %v) ", bpm, expectedTickCount, tickCount)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func numOfOccurences(buf []byte, v byte) int {
	var stopCount = 0
	for _, b := range buf {
		if b == v {
			stopCount++
		}
	}

	return stopCount

}
