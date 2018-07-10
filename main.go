package main

import (
	"github.com/andlabs/ui"
)

func main() {
	err := ui.Main(func() {
		input := ui.NewEntry()
		button := ui.NewButton("Greet")
		greeting := ui.NewLabel("")
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("Enter your name:"), false)
		box.Append(input, false)
		box.Append(button, false)
		box.Append(greeting, false)
		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			addTunnelWindow()
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func addTunnelWindow() {
	box := ui.NewVerticalBox()
	box.Append(ui.NewLabel("Add SSH Tunnel"), false)
	window := ui.NewWindow("Add SSH Tunnel", 200, 100, false)
	window.SetMargined(true)
	window.SetChild(box)
	window.OnClosing(func(*ui.Window) bool {
		return true
	})
	window.Show()
}
