package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andlabs/ui"
)

func main() {
	//TODO: make file configurable and triggered by a ui action
	sshCfg, err := readSSHConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err != nil {
		panic(err)
	}

	fmt.Println(sshCfg.String())

	err = ui.Main(func() {
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
