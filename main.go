package main

import (
	"fmt"
	"os/exec"

	"github.com/rivo/tview"
)

var apps = []string{"neofetch", "vim", "nano", "firefox"}

func installApplications(selected []int) {
	for _, index := range selected {
		app := apps[index]
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
	selected := make(map[int]struct{})

	for i, app := range apps {
		i, app := i, app // capture loop variables
		list.AddItem(fmt.Sprintf("[ ] %s", app), "", 0, func() {
			if _, ok := selected[i]; ok {
				delete(selected, i)
			} else {
				selected[i] = struct{}{}
			}
			updateList(list, selected)
		})
	}

	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if _, ok := selected[index]; ok {
			delete(selected, index)
		} else {
			selected[index] = struct{}{}
		}
		updateList(list, selected)
	})

	installButton := tview.NewButton("Install").SetSelectedFunc(func() {
		selectedIndexes := make([]int, 0, len(selected))
		for index := range selected {
			selectedIndexes = append(selectedIndexes, index)
		}
		installApplications(selectedIndexes)
		app.Stop()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(installButton, 1, 0, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func updateList(list *tview.List, selected map[int]struct{}) {
	list.Clear()
	for i, app := range apps {
		prefix := "[ ] "
		if _, ok := selected[i]; ok {
			prefix = "[*] "
		}
		list.AddItem(fmt.Sprintf("%s%s", prefix, app), "", 0, nil)
	}
}
