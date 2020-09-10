package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/urfave/cli"
)

//TODO: add authors, clean up code, timer ui

func main() {

	app := &cli.App{
		Name:    "Go Time",
		Usage:   "A Golang CLI Countdown",
		Version: "v1.0.0",
		Commands: []*cli.Command{
			{
				Name:    "start timer using seconds",
				Aliases: []string{"st"},
				Usage:   "Start the timer, providing seconds, default timer is 10 seconds.",
				Action: func(c *cli.Context) error {

					fmt.Printf("Hello: arg = %q\n", c.Args().Get(0))

					var timerSetting int = 10
					var err error
					var timeSeconds = int(time.Second)

					// Set Default timer to 10 if no argument given
					if c.Args().Get(0) != "" {
						timerSetting, err = strconv.Atoi(c.Args().Get(0))
						if err != nil {
							log.Fatal(err)
						}
					}

					ticker := time.NewTicker(time.Second)
					defer ticker.Stop()
					done := make(chan bool)
					go func() {
						time.Sleep(time.Duration(timerSetting * timeSeconds))
						done <- true
					}()
					for {
						select {
						case <-done:
							fmt.Println("Done!")

							//TODO: change path to find beep mp3
							//TODO: Extract function from here!
							f, err := os.Open("beep-06.mp3")
							if err != nil {
								log.Fatal(err)
							}

							streamer, format, err := mp3.Decode(f)
							if err != nil {
								log.Fatal(err)
							}
							defer streamer.Close()

							//Initialize the speaker to play the sound. We only want to initialize this once
							speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

							done := make(chan bool)
							speaker.Play(beep.Seq(streamer, beep.Callback(func() {
								done <- true
							})))

							<-done

							return nil
						case t := <-ticker.C:
							fmt.Println("Current time: ", t)
						}
					}
				},
			},
			{
				Name:    "start timer using minutes",
				Aliases: []string{"mt"},
				Usage:   "Start the timer, providing multiplier for seconds, default timer is 10 seconds.",
				Action: func(c *cli.Context) error {
					//TODO: add minutes logic here
					fmt.Println("It's GoTime (with Minutes)!")
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("It's GoTime!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
