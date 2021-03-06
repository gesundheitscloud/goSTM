package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/andlabs/ui"
	"github.com/kevinburke/ssh_config"
)

// Tunnel ..
type Tunnel struct {
	SSHConfig *ssh_config.Config
	Host      string
	Context   context.Context
	Cancel    context.CancelFunc
}

// GoSTMUI ..
type GoSTMUI struct {
	TunnelListParent *ui.Box
	TunnelBox        *ui.Box
	TunnelList       []TunnelUIItem
	ConfigField      *ui.Entry
}

// TunnelUIItem ..
type TunnelUIItem struct {
	Item   *ui.Checkbox
	Icon   *ui.Label
	Tunnel Tunnel
}

var goSTMUI GoSTMUI

func main() {
	err := ui.Main(func() {
		startTunnelButton := ui.NewButton("Start")
		stopTunnelButton := ui.NewButton("Stop")
		readConfigButton := ui.NewButton("Get Config")
		box := ui.NewVerticalBox()
		tunnelListParentBox := ui.NewVerticalBox()
		configEntry := ui.NewEntry()
		box.Append(tunnelListParentBox, false)

		goSTMUI = GoSTMUI{tunnelListParentBox, nil, nil, configEntry}

		configEntry.SetText(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
		getConfig()

		box.Append(configEntry, false)
		box.Append(readConfigButton, false)
		box.Append(startTunnelButton, false)
		box.Append(stopTunnelButton, false)
		window := ui.NewWindow("goSTM", 200, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		startTunnelButton.OnClicked(func(*ui.Button) {
			startSelectedTunnels()
		})
		stopTunnelButton.OnClicked(func(*ui.Button) {
			stopSelectedTunnels()
		})
		readConfigButton.OnClicked(func(*ui.Button) {
			getConfig()
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

func getConfig() {
	var sshCfg *ssh_config.Config
	if goSTMUI.ConfigField.Text() == "" {
		var serr error
		sshCfg, serr = getDefaultConfig()
		if serr != nil {
			panic(serr)
		}
	} else {
		if isValidURL(goSTMUI.ConfigField.Text()) {
			// This seems to be an URL
			fmt.Println("Get config from URL")
			var serr error
			sshCfg, serr = getConfigFromURL()
			if serr != nil {
				panic(serr)
			}
		} else {
			// This should be a local path
			fmt.Println("Get config from local path")
			var serr error
			sshCfg, serr = getConfigFromLocalPath()
			if serr != nil {
				panic(serr)
			}
		}
	}
	displayTunnels(sshCfg)
}

func displayTunnels(sshCfg *ssh_config.Config) {
	fmt.Println("Refresh Tunnel List")
	tunnelListBox := ui.NewVerticalBox()
	var tunnelList []TunnelUIItem
	for _, host := range sshCfg.Hosts {
		if len(host.String()) > 0 {
			ctx, cancel := context.WithCancel(context.Background())
			tunnel := Tunnel{sshCfg, host.Patterns[0].String(), ctx, cancel}
			tunnelUIItem := TunnelUIItem{ui.NewCheckbox(host.String()), ui.NewLabel("Disabled"), tunnel}
			tunnelList = append(tunnelList, tunnelUIItem)
			tunnelBox := ui.NewHorizontalBox()
			tunnelBox.Append(tunnelUIItem.Item, false)
			tunnelBox.Append(tunnelUIItem.Icon, false)
			tunnelListBox.Append(tunnelBox, false)
		}
	}

	if goSTMUI.TunnelBox != nil {
		goSTMUI.TunnelListParent.Delete(0)
		goSTMUI.TunnelBox.Destroy()
	}

	goSTMUI.TunnelList = tunnelList
	goSTMUI.TunnelBox = tunnelListBox
	fmt.Println("Set new List")
	goSTMUI.TunnelListParent.Append(tunnelListBox, false)
}

func startSelectedTunnels() {
	for _, tunnel := range goSTMUI.TunnelList {
		if tunnel.Item.Checked() {
			fmt.Printf("%s\n", tunnel.Tunnel.Host)
			start(tunnel.Tunnel.Context, tunnel.Tunnel.SSHConfig, tunnel.Tunnel.Host)
			tunnel.Icon.SetText("Active")
		}
	}
}

func stopSelectedTunnels() {
	for _, tunnel := range goSTMUI.TunnelList {
		if tunnel.Item.Checked() {
			tunnel.Tunnel.Cancel()
			tunnel.Icon.SetText("Disabled")
		}
	}
}
