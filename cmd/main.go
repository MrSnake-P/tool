package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tool"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

func main() {
	a := app.New()
	w := a.NewWindow("tool")
	content := container.NewMax()
	title := widget.NewLabel("welcome")
	setTutorial := func(t tool.Tutorial) {
		title.SetText(t.Title)
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}
	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
	split.Offset = 0.2
	w.SetContent(split)
	w.Resize(fyne.NewSize(580, 480))
	t := &tool.MyTheme{}
	t.SetFonts("bindata", "")
	a.Settings().SetTheme(t)
	w.ShowAndRun()
}

func makeNav(setTutorial func(tutorial tool.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tool.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tool.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := tool.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := tool.Tutorials[uid]; ok {
				if unsupportedTutorial(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func unsupportedTutorial(t tool.Tutorial) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}
