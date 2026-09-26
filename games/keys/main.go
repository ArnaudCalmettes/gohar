// Command keys plays the notes of a MIDI keyboard, and nothing else.
//
// The other half of the latency measurement. games/otolatency covered
// the path from a decision to a sound; this one starts at a finger.
//
// # What to listen for
//
// Play. That is the whole protocol. A keyboard that answers late feels
// soft under the hand, and no figure says it better than the hand does.
// The delay printed is only the part between the event reaching us and
// the sample being written, which is the part we could be blamed for.
// What the keyboard, the bus, the driver and the card add before and
// after is not visible from here.
//
// # What it is not
//
// A synthesiser anyone would want to hear. One sine per key, no
// envelope beyond the ramps that stop it clicking, no velocity curve.
// It plays in tune and on time, which is all that is being tested.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/synth"
)

func main() {
	list := flag.Bool("list", false, "list the MIDI inputs and exit")
	port := flag.String("port", "", "an input's number as -list prints it, "+
		"or part of its name; the first one if empty")
	device := flag.Duration("device", synth.DefaultBuffer, "device buffer")
	player := flag.Duration("player", 0, "player buffer, zero means the same")
	a4 := flag.Float64("a4", 440, "diapason in hertz")
	timbre := flag.String("timbre", "sine", "sine, pulse12, pulse25, square, triangle or noise")
	authentic := flag.Bool("authentic", false, "keep the consoles' aliasing and stepped triangle instead of smoothing them")
	verbose := flag.Bool("v", false, "print every key event received")
	flag.Parse()

	defer keyboard.Shutdown()

	ports := keyboard.Ports()
	if *list {
		for i, p := range ports {
			fmt.Printf("%d: %s\n", i, p)
		}
		if len(ports) == 0 {
			fmt.Println("no MIDI input found.")
			fmt.Println("check that the keyboard is plugged in and powered,")
			fmt.Println("and that this build has rtmidi, which needs cgo and ALSA headers.")
		}
		return
	}

	keys, err := keyboard.OpenMIDI(*port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer keys.Close()

	engine, err := synth.NewEngine(*timbre, !*authentic)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	engine.Tuning = synth.Tuning{A4: *a4}

	out, err := synth.Open(engine, synth.Options{Device: *device, Player: *player})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer out.Close()

	err = keys.Listen(func(e keyboard.Event) {
		if *verbose {
			fmt.Printf("key %3d  velocity %.2f  down %v\n", e.Key, e.Velocity, e.Down)
		}
		if e.Down {
			engine.NoteOn(e.Key, e.Velocity, e.At)
		} else {
			engine.NoteOff(e.Key, e.At)
		}
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("listening to %s\n", keys.Name())
	fmt.Printf("device %v, queued %v, floor about %v\n\n",
		*device, out.Queued(), out.Floor())
	fmt.Println("Play. Ctrl-C to stop.")
	fmt.Println()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	report := time.NewTicker(time.Second)
	defer report.Stop()

	for {
		select {
		case <-interrupt:
			fmt.Println()
			return
		case <-report.C:
			if err := out.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "audio:", err)
				return
			}
			last, worst := engine.Delays()
			if worst == 0 {
				continue
			}
			fmt.Printf("\revent to samples %5.1f ms   worst %5.1f ms   queued %5.1f ms   ",
				ms(last), ms(worst), ms(out.Queued()))
		}
	}
}

func ms(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}
