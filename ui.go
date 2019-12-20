package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strconv"
	"strings"
)

var (
	itemWidth   = 10
	itemHeight  = 4
	itemPadding = 1
)

func initGui(config *Config) *gocui.Gui {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	maxX, maxY := gui.Size()
	for i, item := range config.Items {
		if err := createItemView(item, i, maxX, maxY, gui); err != nil {
			panic(err)
		}
	}

	if err := gui.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		LogPanic(err.Error())
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		LogPanic(err.Error())
	}

	return gui
}

func createItemView(item *Item, col int, maxX int, maxY int, gui *gocui.Gui) error {
	itemsPerLine := maxX / (itemWidth + itemPadding)
	row := 0

	if col >= itemsPerLine {
		row = col / itemsPerLine
		col %= itemsPerLine
	}

	x := col * (itemWidth + itemPadding)
	y := row * (itemHeight + itemPadding)

	view, err := gui.SetView(item.Label, x, y, x+itemWidth, y+itemHeight)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	view.Frame = true
	view.Title = item.Label
	view.Wrap = true
	view.Overwrite = true

	fmt.Fprintln(view, " loading")
	return nil
}

func updateItemView(item *Item, result string, gui *gocui.Gui) {
	gui.Update(func(gui *gocui.Gui) error {
		view, err := gui.View(item.Label)
		if err != nil {
			LogPanic(err.Error())
		}

		view.Clear()

		if result == ResultErr {
			// Something bad happened, the service may not be available.
			view.BgColor = ErrBgColor
			view.FgColor = ErrFgColor
			fmt.Fprintln(view, result)
		} else {
			// Got some response.
			view.BgColor = OkBgColor
			view.FgColor = OkFgColor

			// Try to parse the result as float.
			result = strings.Replace(result, ",", ".", 1)
			floatResult, err := strconv.ParseFloat(result, 32)
			if err == nil {
				// We can apply threshold if available for this item.
				if item.Threshold > 0 && float32(floatResult) >= item.Threshold {
					view.BgColor = WarnBgColor
					view.FgColor = WarnFgColor
				}
				fmt.Fprintf(view, "%.3f", floatResult)
			} else {
				fmt.Fprint(view, result)
			}

			if item.Unit != "" {
				fmt.Fprintf(view, " %s\n", item.Unit)
			}
		}

		return nil
	})
}

func closeGui(gui *gocui.Gui) {
	defer gui.Close()
}

func runMainLoop(gui *gocui.Gui) {
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		LogPanic(err.Error())
	}
}

func quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
