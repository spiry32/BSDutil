package main

import (
	"fmt"
	"os/exec"

	"github.com/gdamore/tcell/v2"
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

	list := tview.NewList()

	for i, appName := range apps {
		item := fmt.Sprintf("[%d] %s", i+1, appName)
		if err := list.AddItem(item, "", '0'+rune(i), nil); err != nil {
			fmt.Printf("Failed to add item %s: %v\n", appName, err)
		}
	}

	selected := make(map[int]string)

	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		appName := apps[index]
		if _, ok := selected[index]; ok {
			delete(selected, index)
			list.RemoveItem(index)
			list.InsertItem(index, mainText, "", '0'+rune(index), nil)
		} else {
			selected[index] = appName
			list.RemoveItem(index)
			list.InsertItem(index, mainText, "[green]+", '0'+rune(index), nil)
		}
	})

	okButton := tview.NewButton("OK").SetSelectedFunc(func() {
		if len(selected) > 0 {
			installApplications(selected)
			app.Stop()
		}
	})

	cancelButton := tview.NewButton("CANCEL").SetSelectedFunc(func() {
