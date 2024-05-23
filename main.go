package main

import (
	"fmt"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var apps = []string{"neofetch", "vim", "nano", "firefox"}

func installApplications(selected map[int]struct{}) {
	selectedIndexes := make([]int, 0, len(selected))
	for index := range selected {
		selectedIndexes = append(selectedIndexes, index)
	}
	for _, index := range selectedIndexes {
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

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case ' ':
				index := list.GetCurrentItem()
				if _, ok := selected[index]; ok {
					delete(selected, index)
				} else {
					selected[index] = struct{}{}
				}
				updateList(list, selected)
			case 'j', 'J', 's', 'S':
				// tasta "j" sau "s" pentru deplasare în jos
				index := list.GetCurrentItem()
				if index < len(apps)-1 {
					index++
					list.SetCurrentItem(index)
				}
			case 'k', 'K', 'w', 'W':
				// tasta "k" sau "w" pentru deplasare în sus
				index := list.GetCurrentItem()
				if index > 0 {
					index--
					list.SetCurrentItem(index)
				}
			case '\n':
				if len(selected) > 0 {
					installApplications(selected)
					app.Stop()
				}
			}
		}
		return event
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true)

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
