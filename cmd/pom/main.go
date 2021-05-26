package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
	"github.com/cheggaaa/pb/v3"
	"github.com/gen2brain/beeep"

	"github.com/nuetoban/go-pomodoro/dto"
	"github.com/nuetoban/go-pomodoro/pomodoro"
	"github.com/nuetoban/go-pomodoro/storage"
	"github.com/nuetoban/go-pomodoro/terminal"
)

func main() {
	parser := argparse.NewParser("pom", "Runs pomodoro timer")

	// Define arguments
	nameArgument := parser.String("n", "name", &argparse.Options{Default: "Timer", Help: "Name of the task"})
	minutesArgument := parser.Int("m", "minutes", &argparse.Options{Default: 25, Help: "Number of minutes"})

	// Parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// Create new storage
	sqlite, err := storage.NewSqlite()
	if err != nil {
		panic(err)
	}

	// Create new pomodoro service
	pom := pomodoro.NewService(sqlite)

	// Variables representing state
	// TODO: Replace by FSM
	cont := true
	rest := false

	// Run tasks or rest until stopped
	for cont {
		var (
			minutes  int
			name     = *nameArgument
			progress = make(chan struct{})
			done     = make(chan struct{})
			sigs     = make(chan os.Signal)
			ctx      = context.Background()
		)

		// Create context with cancellation
		ctx, cancel := context.WithCancel(ctx)

		// Handle signals
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			cancel()
		}()

		// New bar template
		template := NewBarTemplate()

		// Check rest mode
		if rest {
			minutes = 5
			name = "Rest"
			template.SetEmoji("💤")
		} else {
			minutes = *minutesArgument
			template.SetEmoji("🍅")
		}

		template.SetName(name)

		// Run pomodoro
		go pom.RunTask(ctx, dto.Task{Name: name, Interval: time.Minute * time.Duration(minutes)}, progress, done)

		// Create new bar
		seconds := minutes * 60
		bar := pb.New(seconds)
		bar.SetTemplateString(template.String())
		bar.SetMaxWidth(120)
		bar.Start()

		// We don't need a cursor in terminal
		terminal.HideCursor()

		// Track time
		func() {
			for {
				select {
				// Progress ticks every second
				case <-progress:
					if bar.Current() < int64(seconds) {
						bar.Increment()
					}

				// The timer is finished or cancelled
				case <-done:
					for bar.Current() < int64(seconds) {
						bar.Increment()
					}
					bar.Finish()

					// Make alarm about task finish and show cursor
					terminal.Bell()
					terminal.ShowCursor()
					err := beeep.Notify("Pomodoro", fmt.Sprintf("The task %s has been finished!", *nameArgument), "assets/information.png")
					if err != nil {
						panic(err)
					}

					// Get next action
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Wanna stop/rest/work? [s/r/w]: ")
					text, err := reader.ReadString('\n')
					if err != nil {
						panic(err)
					}

					// Clear one line so next timer could start right below previous timer
					terminal.ClearLine()

					// Decide on input
					switch strings.TrimSpace(text) {

					// Restart working
					case "w":
						if rest {
							rest = false
						}
						return

					// Take a rest
					case "r":
						rest = true
						return

					// Stop
					default:
						cont = false
						return
					}
				}
			}
		}()
	}
}
