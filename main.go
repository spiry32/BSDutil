package main

import (
	"fmt"
	"os/exec"

	"github.com/rivo/tview"
)

var apps = []string{"neofetch", "vim", "nano", "firefox", "firefox-esr"}

func installApplications(selected map[int]string) {
	for _, app := range selected {
		fmt.Printf("Installing %s...\n", app)
		cmd := exec.Command("sudo", "pkg", "install", "-y", app)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to install %s: %v\n", app, err)
		} else {
			fmt.Printf("%s installed successfully.\n", app)
		}
	}
}

func main() {
	app := tview.NewApplication()

	list := tview.NewList().
		AddItem("neofetch", "", '1', nil).
		AddItem("vim", "", '2', nil).
		AddItem("nano", "", '3', nil).
		AddItem("firefox", "", '4', nil).
		AddItem("firefox-esr", "", '5', nil)

	selected := make(map[int]string)

	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		appName := apps[index]
		if _, ok := selected[index]; ok {
			delete(selected, index)
			list.SetSecondaryText(index, "")
		} else {
			selected[index] = appName
			list.SetSecondaryText(index, "[green]+")
		}
	})

	installButton := tview.NewButton("Install").SetSelectedFunc(func() {
		if len(selected) > 0 {
			installApplications(selected)
			app.Stop()
		}
	})

	flex := tview.NewFlex().
		AddItem(list, 0, 1, true).
		AddItem(installButton, 1, 0, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
