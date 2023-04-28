package tool

import "fyne.io/fyne/v2"

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"input": {"切分文件",
			"",
			makeInputTab,
			true,
		},
	}
	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"": {"input"},
	}
)
