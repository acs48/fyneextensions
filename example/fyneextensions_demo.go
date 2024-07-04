package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/acs48/fyneextensions"
	"sort"
	"strings"
)

type homeAction struct {
	mAction *fyneextensions.ActionItem
	w       fyne.Window
}

func newHomeAction(w fyne.Window) *homeAction {
	rv := &homeAction{
		w: w,
		mAction: fyneextensions.NewActionItem("Home", false, false, []fyne.Resource{theme.FileApplicationIcon()}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
			fyneextensions.NewActionItem("File", false, false, []fyne.Resource{theme.FileApplicationIcon()}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
				fyneextensions.NewActionItem("New", false, false, []fyne.Resource{theme.DocumentCreateIcon()}, false, false, false, 0, func(int) {}, nil),
				fyneextensions.NewActionItem("Open", false, false, []fyne.Resource{theme.FolderOpenIcon()}, false, false, false, 0, func(int) {}, nil),
				fyneextensions.NewActionItem("Save", false, false, []fyne.Resource{theme.DocumentSaveIcon()}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
					fyneextensions.NewActionItem("Save", true, false, []fyne.Resource{theme.DocumentSaveIcon()}, false, false, false, 0, func(int) {}, nil),
					fyneextensions.NewActionItem("Save as", true, false, []fyne.Resource{theme.DocumentSaveIcon()}, false, false, false, 0, func(int) {}, nil),
				}),
			}),
		}),
	}
	return rv
}

func (ha *homeAction) GetActions() *fyneextensions.ActionItem {
	return ha.mAction
}

func (ha *homeAction) GetCanvas() fyne.Canvas {
	return ha.w.Canvas()
}

type editAction struct {
	mAction *fyneextensions.ActionItem
	w       fyne.Window
}

func newEditAction(w fyne.Window) *editAction {
	checkerItem := fyneextensions.NewActionItem("Check me!", true, false, []fyne.Resource{theme.CheckButtonIcon(), theme.CheckButtonCheckedIcon()}, false, false, true, 0, func(int) {}, nil)
	checkerItem.Triggered = func(i int) {
		if i == 0 {
			checkerItem.Stater.Set(1)
		} else {
			checkerItem.Stater.Set(0)
		}
	}
	rv := &editAction{
		w: w,
		mAction: fyneextensions.NewActionItem("Edit", false, false, []fyne.Resource{}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
			fyneextensions.NewActionItem("Clipboard", false, false, []fyne.Resource{}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
				fyneextensions.NewActionItem("Copy", false, false, []fyne.Resource{theme.ContentCopyIcon()}, false, false, false, 0, func(int) {}, nil),
				fyneextensions.NewActionItem("Cut", false, false, []fyne.Resource{theme.ContentCutIcon()}, false, false, false, 0, func(int) {}, nil),
				fyneextensions.NewActionItem("Paste", false, false, []fyne.Resource{theme.ContentPasteIcon()}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
					fyneextensions.NewActionItem("Paste", true, false, []fyne.Resource{theme.ContentPasteIcon()}, false, false, false, 0, func(int) {}, nil),
					fyneextensions.NewActionItem("Paste Special", true, false, []fyne.Resource{theme.ContentPasteIcon()}, false, false, false, 0, func(int) {}, nil),
				}),
			}),
			fyneextensions.NewActionItem("Enable", false, false, []fyne.Resource{}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
				checkerItem,
			}),
		}),
	}
	return rv
}

func (ea *editAction) GetActions() *fyneextensions.ActionItem {
	return ea.mAction
}

func (ea *editAction) GetCanvas() fyne.Canvas {
	return ea.w.Canvas()
}

type mStringList []string

func (lst mStringList) Len() int {
	return len(lst)
}
func (lst mStringList) Less(i, j int) bool {
	return lst[i] < lst[j]
}
func (lst mStringList) Swap(i, j int) {
	lst[i], lst[j] = lst[j], lst[i]
}

type sampleList struct {
	allItems     mStringList
	visibleItems mStringList
	selectedItem string

	mWindow fyne.Window
}

func (dm *sampleList) GetCanvas() fyne.Canvas {
	return dm.mWindow.Canvas()
}

func (dm *sampleList) ListLen() int {
	return len(dm.visibleItems)
}

func (dm *sampleList) CreateListItem() fyne.CanvasObject {
	name := widget.NewLabel("")
	return name
}

func (dm *sampleList) UpdateListItem(i int, o fyne.CanvasObject) {
	if i >= 0 && i < len(dm.visibleItems) {
		o.(*widget.Label).SetText(dm.visibleItems[i])
	}
}

func (dm *sampleList) GetSelected() (id int) {
	id = -1
	for i, ptr := range dm.visibleItems {
		if ptr == dm.selectedItem {
			id = i
		}
	}
	return
}

func (dm *sampleList) OnSelectedItem(id int) {
	if id >= 0 && id < len(dm.visibleItems) {
		if dm.selectedItem != dm.visibleItems[id] {
			dm.selectedItem = dm.visibleItems[id]
		}
	}
}

func (dm *sampleList) ClearSearch() {
	dm.visibleItems = dm.visibleItems[0:len(dm.allItems)]
	i := 0
	for _, dtsPtr := range dm.allItems {
		dm.visibleItems[i] = dtsPtr
		i++
	}
	dm.visibleItems = dm.visibleItems[0:i]
	sort.Sort(dm.visibleItems)
}

func (dm *sampleList) GetListableSearchableActions() *fyneextensions.ActionItem {
	return nil
}

func (dm *sampleList) StartSearch(text string) {
	if text != "" {
		j := 0
		dm.visibleItems = dm.visibleItems[0:cap(dm.visibleItems)]
		for _, dst := range dm.allItems {
			if strings.Contains(strings.ToUpper(dst), strings.ToUpper(text)) || strings.Contains(strings.ToUpper(dst), strings.ToUpper(text)) {
				dm.visibleItems[j] = dst
				j++
			}
		}
		dm.visibleItems = dm.visibleItems[0:j]
		sort.Sort(dm.visibleItems)
	}
}

// Example function demonstrating the fyneextensions widgets
func main() {
	// Instantiate the Fyne application
	a := app.New()

	// Create a new window
	w := a.NewWindow("Fyne Window")

	// Set the main window content
	messageString := binding.NewString()
	messageLabel := widget.NewLabelWithData(messageString)

	mHomeAction := newHomeAction(w)
	mEditAction := newEditAction(w)
	mRibbon := container.NewAppTabs()
	homeTab, _ := fyneextensions.BuildTabItemRibbon(mHomeAction, 60., 30., messageString)
	homeMenu := fyneextensions.NewActionableMenu(mHomeAction.GetActions())
	mRibbon.Append(homeTab)
	editTab, _ := fyneextensions.BuildTabItemRibbon(mEditAction, 60., 30., messageString)
	editMenu := fyneextensions.NewActionableMenu(mEditAction.GetActions())
	mRibbon.Append(editTab)

	projectTree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			switch id {
			case "":
				return []widget.TreeNodeID{"a", "b", "c"}
			case "a":
				return []widget.TreeNodeID{"a1", "a2"}
			}
			return []string{}
		},
		func(id widget.TreeNodeID) bool {
			return id == "" || id == "a"
		},
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("Branch template")
			}
			return widget.NewLabel("Leaf template")
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			text := id
			if branch {
				text += " (branch)"
			}
			o.(*widget.Label).SetText(text)
		})
	widgetTree := fyneextensions.NewMiniWidget("PROJECT", true, 20., projectTree, false, true, nil, false, nil, true, nil, nil, false, nil, nil, nil, nil, nil, w.Canvas())
	listItem := &sampleList{
		allItems:     mStringList{"alfa", "beta", "gamma", "delta"},
		visibleItems: make(mStringList, 4),
		mWindow:      w,
	}
	listItem.ClearSearch()
	searchList := fyneextensions.NewListableSearchableWidget(listItem)
	searchWidget := fyneextensions.NewMiniWidget("ITEMS", true, 20., searchList, false, true, nil, false, nil, true, nil, nil, false, nil, nil, nil, nil, nil, w.Canvas())

	mainContent := container.NewStack()
	sideContent := fyneextensions.NewSideBar(widgetTree, searchWidget)
	split := container.NewHSplit(sideContent, mainContent)

	mainContainer := container.NewBorder(mRibbon, messageLabel, nil, nil, mRibbon, messageLabel, split)

	w.SetContent(mainContainer)
	w.SetMainMenu(fyne.NewMainMenu(
		homeMenu.Menu,
		editMenu.Menu,
	))
	// Show and run the application
	w.ShowAndRun()

	// The application run in a different goroutine and could not be represented in console output.
	// Let's print something instead
	fmt.Println("A new Fyne window has been created and shown.")

}
