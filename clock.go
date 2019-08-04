// Provides a midi master clock and various other utils for working with midi
package midi

import (
	"io"
	"time"
)

// Midi standard commands
const (
	Start    = 0xFA
	Stop     = 0xFC
	Tick     = 0xF8
	Continue = 0xFB
)

var (
	midiStart    = []byte{Start}
	midiStop     = []byte{Stop}
	midiTick     = []byte{Tick}
	midiContinue = []byte{Continue}
)

// Pulses per quarter note
const ppqn = 24
const uSecInMin = 6000000

type Bpm float32

// Converts bpm to a 24ppqn pulse interval in microseconds
func tempoToPulseInterval(t Bpm) time.Duration {

	return time.Duration((uSecInMin/(t/10.00))/ppqn) * time.Microsecond
}

type Clock struct {
	midiOut   chan []byte
	pulseRate chan time.Duration
	device    io.Writer
}

// Create a new midi clock. The clock starts to send tick events as soon as it is created.
// device, err := os.OpenFile("/dev/snd/midiC1D0", os.O_WRONLY, 0664)
// if err != nil {
//	log.Fatal(err)
// }
//
// clk := midi.NewClock(device)
// clk.SetTempo(120.00)
// clk.Start()
func NewClock(device io.Writer) *Clock {
	clk := &Clock{
		device:    device,
		midiOut:   make(chan []byte),
		pulseRate: make(chan time.Duration),
	}

	go clk.run()

	return clk
}

// Change the BPM of the clock
func (clk *Clock) SetTempo(t Bpm) {
	clk.pulseRate <- tempoToPulseInterval(t)
}

// Send MIDI sequencer start event
func (clk *Clock) Start() {
	clk.midiOut <- midiStart
}

// Send MIDI sequencer stop event
func (clk *Clock) Stop() {
	clk.midiOut <- midiStop
}

// Send MIDI sequencer stop event
func (clk *Clock) Continue() {
	clk.midiOut <- midiContinue
}

func (clk *Clock) run() {

	pulseRate := tempoToPulseInterval(120.00)
	var t VarTicker
	t.SetDuration(pulseRate)

	go func() {
		for {
			select {
			case <-t.C:
				clk.device.Write(midiTick)
			case newPulseRate := <-clk.pulseRate:
				t.SetDuration(newPulseRate)
			case midiCmd := <-clk.midiOut:
				clk.device.Write(midiCmd)
			}
		}
	}()
}
