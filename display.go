package main

import (
	"fmt"
	"net/netip"
	"sort"
	"time"

	"github.com/rivo/tview"
)

type ArpKeyValue struct {
	Key   string
	Value Host
}

func Display() {
	app := tview.NewApplication()

	box := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignLeft)

	// Function to update the box content
	go func() {
		for {
			time.Sleep(1 * time.Second)

			cacheList := make([]ArpKeyValue, 0)
			for key, host := range arpCache {
				cacheList = append(cacheList, ArpKeyValue{Key: key, Value: host})
			}
			sort.Slice(cacheList, func(i, j int) bool {
				addr1, _ := netip.ParseAddr(cacheList[i].Key)
				addr2, _ := netip.ParseAddr(cacheList[j].Key)
				return addr1.Compare(addr2) == -1

			})
			cacheStr := ""
			for _, record := range cacheList {
				cacheStr += fmt.Sprintf("%-15s: %17s\n", record.Key, record.Value.MAC)
			}
			box.SetText(fmt.Sprintf("cache size: %d\n%s", len(arpCache), cacheStr))
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
