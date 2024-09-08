package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

func Display() {
	app := tview.NewApplication()
	box := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("0")

	// Function to update the box content
	go func() {
		for {
			time.Sleep(1 * time.Second)
			box.SetText(fmt.Sprintf("cache size: %d", len(arpCache)))
			app.Draw()
		}
	}()

	// Create a flex layout and add the box
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(box, 0, 1, true)

	// Set up the application and run
	if err := app.SetRoot(flex, true).Run(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
	}
}
