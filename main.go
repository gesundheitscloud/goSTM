package main

import (
	"os"
	"path/filepath"

	"github.com/andlabs/ui"
)

// Tunnel ...
type Tunnel struct {
	UIItem *ui.Checkbox
	UIIcon *ui.Label
}

func main() {
	//TODO: make file configurable and triggered by a ui action
	sshCfg, err := readSSHConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err != nil {
		panic(err)
	}

	err = ui.Main(func() {
		var UIList []Tunnel
		startTunnelButton := ui.NewButton("Start")
		stopTunnelButton := ui.NewButton("Stop")
		box := ui.NewVerticalBox()
		var tunnelBox *ui.Box

		// Display tunnels
		for _, host := range sshCfg.Hosts {
			if len(host.String()) > 0 {
				tunnel := Tunnel{ui.NewCheckbox(host.String()), ui.NewLabel("Disabled")}
				UIList = append(UIList, tunnel)
				tunnelBox = ui.NewHorizontalBox()
				tunnelBox.Append(tunnel.UIItem, false)
				tunnelBox.Append(tunnel.UIIcon, false)
				box.Append(tunnelBox, false)
			}
		}

		box.Append(startTunnelButton, false)
		box.Append(stopTunnelButton, false)
		window := ui.NewWindow("goSTM", 200, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		startTunnelButton.OnClicked(func(*ui.Button) {
			startSelectedTunnels(UIList)
		})
		stopTunnelButton.OnClicked(func(*ui.Button) {
			stopSelectedTunnels(UIList)
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

func startSelectedTunnels(UIList []Tunnel) {
	for _, tunnel := range UIList {
		if tunnel.UIItem.Checked() {
			// TODO: activate
			tunnel.UIIcon.SetText("Active")
		}
	}
}

func stopSelectedTunnels(UIList []Tunnel) {
	for _, tunnel := range UIList {
		if tunnel.UIItem.Checked() {
			// TODO: disable
			tunnel.UIIcon.SetText("Disabled")
		}
	}
}
