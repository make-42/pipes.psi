package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/shawnridgeway/wfc"
	"golang.org/x/term"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	checkError(err)
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func main() {
	// Params
	inputFilePath := "input.png"
	// Code
	// Termbox Init
	err := termbox.Init()
	checkError(err)
	defer termbox.Close()
	// Make a go routine for catching inputs
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	// Get terminal width and height
	width, height, err := term.GetSize(0)
	checkError(err)
	// Load input pattern from file
	inputImg, err := getImageFromFilePath(inputFilePath)
	checkError(err)
	// Create a new model
	model := wfc.NewOverlappingModel(inputImg, 3, width, height, true, true, 2, false)

mainloop:
	for {
		i := 0
		model.Clear()
		// Run the algorithm iteratively, stopping at success or contradiction
	generationloop:
		for {
			select {
			// If there's a new key press event read it
			case ev := <-event_queue:
				if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
					break mainloop
				}
				break
			// If not keep drawing
			default:
				// Iterate the model once
				outputImg, finished, success := model.Iterate(1)
				// Display to screen the result
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						r, _, _, _ := outputImg.At(x, y).RGBA()
						if r == 65535 {
							tbprint(x, y, termbox.ColorWhite, termbox.ColorDefault, "█")
						} else if r > 0 {
							tbprint(x, y, termbox.ColorDarkGray, termbox.ColorDefault, "░")
						} else {
							tbprint(x, y, termbox.ColorWhite, termbox.ColorDefault, " ")
						}
					}

				}
				// Show info to the user
				tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, fmt.Sprintf("Iteration: %d", i))
				tbprint(0, height-1, termbox.ColorWhite, termbox.ColorDefault, "Press ESC to exit.")
				// On finish or contradiction restart drawing
				if finished {
					if !success {
						tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, "Algorithm failed because of a contradiction!")
						break generationloop
					}
					tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, "Algorithm finished")
					break generationloop
				}
				// Update screen
				termbox.Flush()
				// Increment the iteration counter
				i++
			}
		}
	}
}
