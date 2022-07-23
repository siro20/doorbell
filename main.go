/* SPDX-License-Identifier: MIT License */
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Wrong Arguments. Expected <filename.mp3> <volume>")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Could not open mp3: %v", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatalf("Could not decode mp3: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/20))
	time.Sleep(time.Second / 10)
	volumeReq, _ := strconv.Atoi(os.Args[2])
	wait := make(chan bool)
	if len(os.Args) == 4 && os.Args[3] == "-t" {
		streamer.Seek(0)
		speaker.Play(beep.Seq(&effects.Volume{
			Streamer: streamer,
			Base:     2,
			Volume:   float64(volumeReq),
			Silent:   false,
		}, beep.Callback(func() { wait <- true })))
	} else {
		c, err := gpiod.NewChip("gpiochip0")
		if err != nil {
			log.Fatalf("Could not open gpiochip0: %v", err)
		}
		l, err := c.RequestLine(rpi.GPIO17, gpiod.WithEventHandler(func(evt gpiod.LineEvent) {
			go func() {
				streamer.Seek(0)
				speaker.Play(beep.Seq(&effects.Volume{
					Streamer: streamer,
					Base:     2,
					Volume:   float64(volumeReq),
					Silent:   false,
				}, beep.Callback(func() {})))
				fmt.Printf("Ding Dong\n")
			}()
		}), gpiod.WithFallingEdge, gpiod.WithDebounce(10*time.Millisecond), gpiod.WithPullUp)
		if err != nil {
			log.Fatalf("Could not request GPIO line: %v", err)
		}
		val, _ := l.Value()
		fmt.Printf("GPIO17 is %d\n", val)
		defer l.Close()
	}
	<-wait
}
