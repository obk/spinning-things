// Spinning shapes rendered as ASCII art in the terminal.
//
//	go run . cube   (default)
//	go run . donut
//	go run . lemniscate
//
// Press Ctrl+C to quit.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"time"
)

const (
	width      = 80
	height     = 40
	frameDelay = 16 * time.Millisecond
)

// shape draws one frame into out and advances its animation.
type shape interface {
	frame(out *bufio.Writer)
}

func main() {
	var s shape = &cube{}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "cube":
		case "donut":
			s = &donut{}
		case "lemniscate":
			s = &lemniscate{}
		default:
			fmt.Fprintf(os.Stderr, "unknown shape %q (want cube, donut or lemniscate)\n", os.Args[1])
			os.Exit(2)
		}
	}

	out := bufio.NewWriterSize(os.Stdout, width*height*8)

	// Clear the screen and hide the cursor.
	fmt.Print("\x1b[2J\x1b[?25l")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(frameDelay)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			// Restore the cursor and colors before exiting.
			fmt.Print("\x1b[0m\x1b[?25h\n")
			return
		case <-ticker.C:
			// Move the cursor home and redraw over the previous frame.
			out.WriteString("\x1b[H")
			s.frame(out)
			out.WriteString("\x1b[0m")
			out.Flush()
		}
	}
}
