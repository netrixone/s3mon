package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"net"
	"strconv"
	"strings"
)

var (
	toolBarHeight = 2
	itemWidth     = 10
	itemHeight    = 4
	itemPadding   = 1
	maxX          = 0
	maxY          = 0
	plotterX      = 0
	plotterY      = 0
)

func initGui(config *Config) *gocui.Gui {
	gui, err := gocui.NewGui(gocui.Output256, false)
	if err != nil {
		LogPanic(err.Error())
	}

	maxX, maxY = gui.Size()
	gui.Mouse = true

	createToolBarView(gui)
	createMonitoringView(config.Items, gui)

	if err := gui.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		LogPanic(err.Error())
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		LogPanic(err.Error())
	}

	return gui
}

func createToolBarView(gui *gocui.Gui) {
	view, err := gui.SetView("toolbar", plotterX, plotterY, maxX-1, plotterY+toolBarHeight, 0)
	if !gocui.IsUnknownView(err) {
		LogPanic(err.Error())
	}

	fmt.Fprintf(view, "%s v%s | IP: %s", prog, version, getOutboundIP())

	plotterY += toolBarHeight + itemPadding
}

func createMonitoringView(items []*Item, gui *gocui.Gui) {
	for i, item := range items {
		if err := createMonitoringItemView(item, i, maxX, maxY, gui); err != nil {
			LogPanic(err.Error())
		}
	}

	plotterY += itemHeight + itemPadding
}

func createMonitoringItemView(item *Item, col int, maxX int, maxY int, gui *gocui.Gui) error {
	itemsPerLine := maxX / (itemWidth + itemPadding)
	row := 0

	if col >= itemsPerLine {
		row = col / itemsPerLine
		col %= itemsPerLine
	}

	x := plotterX + col*(itemWidth+itemPadding)
	y := plotterY + row*(itemHeight+itemPadding)

	view, err := gui.SetView(item.Label, x, y, x+itemWidth, y+itemHeight, 0)
	if !gocui.IsUnknownView(err) {
		return err
	}

	view.Frame = true
	view.Title = item.Label
	view.Wrap = true
	view.Overwrite = true

	fmt.Fprint(view, "\n\n loading")

	// Bind mouse-click to item refresh.
	gui.SetKeybinding(item.Label, gocui.MouseLeft, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		updateItem(item, gui)
		return nil
	})

	return nil
}

func updateItemView(item *Item, result string, gui *gocui.Gui) {
	gui.Update(func(gui *gocui.Gui) error {
		view, err := gui.View(item.Label)
		if err != nil {
			LogPanic(err.Error())
		}

		view.Clear()
		fmt.Fprint(view, "\n\n ")

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

func getOutboundIP() string {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return "?"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
