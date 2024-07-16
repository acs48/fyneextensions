package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/acs48/fyneextensions"
	"sort"
	"strings"
	"time"
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

type demoAction struct {
	mAction *fyneextensions.ActionItem
	w       fyne.Window
	mForm   *fyneextensions.FormGenUtility

	SetupTime      time.Time `formGenInclude:"true" formGenDescription:"Demo time.Time formgen" formGenDefaultValue:"now"`
	SetupString    string    `formGenInclude:"true" formGenDescription:"Demo string formgen" formGenIsRequired:"true"`
	SetupCheck     bool      `formGenInclude:"true" formGenDescription:"Demo bool formgen" formGenLabel:"Uncheck me" formGenDefaultValue:"true"`
	SetupCheckOpt1 bool      `formGenInclude:"true" formGenDescription:"Demo group formgen" formGenLabel:"Opt 1" formGenDefaultValue:"false" formGenCheckGroup:"group1"`
	SetupCheckOpt2 bool      `formGenInclude:"true" formGenDescription:"Demo group formgen" formGenLabel:"Opt 2" formGenDefaultValue:"false" formGenCheckGroup:"group1"`
	SetupCheckOpt3 bool      `formGenInclude:"true" formGenDescription:"Demo group formgen" formGenLabel:"Opt 3" formGenDefaultValue:"false" formGenCheckGroup:"group1"`
	SetupRadioOpt1 bool      `formGenInclude:"true" formGenDescription:"Demo group formgen" formGenLabel:"Opt 1" formGenDefaultValue:"true" formGenRadioGroup:"group2"`
	SetupRadioOpt2 bool      `formGenInclude:"true" formGenDescription:"Demo group formgen" formGenLabel:"Opt 2" formGenDefaultValue:"false" formGenRadioGroup:"group2"`
	SetupRadioOpt3 bool      `formGenInclude:"true" formGenDescription:"Demo group formgen" formGenLabel:"Opt 3" formGenDefaultValue:"false" formGenRadioGroup:"group2"`
}

func newDemoAction(w fyne.Window) *demoAction {
	rv := &demoAction{
		w: w,
	}
	timeSetting := fyneextensions.NewActionItem("Demo time", true, false, []fyne.Resource{theme.SettingsIcon()}, false, false, false, 0, func(int) {}, nil)
	checkSetting := fyneextensions.NewActionItem("demo check", true, false, []fyne.Resource{theme.CheckButtonIcon(), theme.CheckButtonCheckedIcon()}, false, false, true, 1, nil, nil)
	chk1Setting := fyneextensions.NewActionItem("check 1", true, false, []fyne.Resource{theme.CheckButtonIcon(), theme.CheckButtonCheckedIcon()}, false, false, true, 0, nil, nil)
	chk2Setting := fyneextensions.NewActionItem("check 2", true, false, []fyne.Resource{theme.CheckButtonIcon(), theme.CheckButtonCheckedIcon()}, false, false, true, 0, nil, nil)
	chk3Setting := fyneextensions.NewActionItem("check 3", true, false, []fyne.Resource{theme.CheckButtonIcon(), theme.CheckButtonCheckedIcon()}, false, false, true, 0, nil, nil)
	radio1Setting := fyneextensions.NewActionItem("radio 1", true, false, []fyne.Resource{theme.RadioButtonIcon(), theme.RadioButtonCheckedIcon()}, false, false, true, 1, nil, nil)
	radio2Setting := fyneextensions.NewActionItem("radio 2", true, false, []fyne.Resource{theme.RadioButtonIcon(), theme.RadioButtonCheckedIcon()}, false, false, true, 0, nil, nil)
	radio3Setting := fyneextensions.NewActionItem("radio 3", true, false, []fyne.Resource{theme.RadioButtonIcon(), theme.RadioButtonCheckedIcon()}, false, false, true, 0, nil, nil)

	checkSetting.Triggered = func(s int) {
		if s == 0 {
			checkSetting.Stater.Set(1)
			rv.SetupCheck = true
		} else {
			checkSetting.Stater.Set(0)
			rv.SetupCheck = false
		}
	}
	chk1Setting.Triggered = func(s int) {
		if s == 0 {
			chk1Setting.Stater.Set(1)
			rv.SetupCheckOpt1 = true
		} else {
			chk1Setting.Stater.Set(0)
			rv.SetupCheckOpt1 = false
		}
	}
	chk2Setting.Triggered = func(s int) {
		if s == 0 {
			chk2Setting.Stater.Set(1)
			rv.SetupCheckOpt2 = true
		} else {
			chk2Setting.Stater.Set(0)
			rv.SetupCheckOpt2 = false
		}
	}
	chk3Setting.Triggered = func(s int) {
		if s == 0 {
			chk3Setting.Stater.Set(1)
			rv.SetupCheckOpt3 = true
		} else {
			chk3Setting.Stater.Set(0)
			rv.SetupCheckOpt3 = false
		}
	}
	radio1Setting.Triggered = func(s int) {
		if s == 0 {
			radio1Setting.Stater.Set(1)
			rv.SetupRadioOpt1 = true
			rv.SetupRadioOpt2 = false
			rv.SetupRadioOpt3 = false
		}
	}
	radio2Setting.Triggered = func(s int) {
		if s == 0 {
			radio2Setting.Stater.Set(1)
			rv.SetupRadioOpt1 = false
			rv.SetupRadioOpt2 = true
			rv.SetupRadioOpt3 = false
		}
	}
	radio3Setting.Triggered = func(s int) {
		if s == 0 {
			radio3Setting.Stater.Set(1)
			rv.SetupRadioOpt1 = false
			rv.SetupRadioOpt2 = false
			rv.SetupRadioOpt3 = true
		}
	}

	rv.mAction = fyneextensions.NewActionItem("Demo", false, false, []fyne.Resource{}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
		fyneextensions.NewActionItem("Formgen demo", false, false, []fyne.Resource{}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
			fyneextensions.NewActionItem("Test formgen", true, false, []fyne.Resource{theme.SettingsIcon()}, false, false, false, 0, rv.showDialog, nil),
			fyneextensions.NewActionItem("Formgen settings", true, false, []fyne.Resource{theme.SettingsIcon()}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
				timeSetting,
				checkSetting,
			}),
			fyneextensions.NewActionItem("demo check group", false, true, []fyne.Resource{theme.SettingsIcon()}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
				chk1Setting,
				chk2Setting,
				chk3Setting,
			}),
			fyneextensions.NewActionItem("demo radio group", false, true, []fyne.Resource{theme.SettingsIcon()}, false, false, false, 0, nil, []*fyneextensions.ActionItem{
				radio1Setting,
				radio2Setting,
				radio3Setting,
			}),
		}),
	})

	rv.mForm = fyneextensions.NewFormGenDialog(rv, "App settings", "Ok", "Cancel", func(b bool) {
		if b {
			timeSetting.Name.Set(rv.SetupTime.Format("2006-01-02T15:04:05"))
			checkSetting.Name.Set(rv.SetupString)

			if rv.SetupCheck {
				checkSetting.Stater.Set(1)
			} else {
				checkSetting.Stater.Set(0)
			}

			if rv.SetupCheckOpt1 {
				chk1Setting.Stater.Set(1)
			} else {
				chk1Setting.Stater.Set(0)
			}
			if rv.SetupCheckOpt2 {
				chk2Setting.Stater.Set(1)
			} else {
				chk2Setting.Stater.Set(0)
			}
			if rv.SetupCheckOpt3 {
				chk3Setting.Stater.Set(1)
			} else {
				chk3Setting.Stater.Set(0)
			}

			if rv.SetupRadioOpt1 {
				radio1Setting.Stater.Set(1)
				radio2Setting.Stater.Set(0)
				radio3Setting.Stater.Set(0)
			}
			if rv.SetupRadioOpt2 {
				radio1Setting.Stater.Set(0)
				radio2Setting.Stater.Set(1)
				radio3Setting.Stater.Set(0)
			}
			if rv.SetupRadioOpt3 {
				radio1Setting.Stater.Set(0)
				radio2Setting.Stater.Set(0)
				radio3Setting.Stater.Set(1)
			}
		}
	}, w, fyne.NewSize(750, 550))

	rv.SetupCheck = true
	rv.SetupRadioOpt1 = true
	return rv
}

func (ea *demoAction) GetActions() *fyneextensions.ActionItem {
	return ea.mAction
}

func (ea *demoAction) GetCanvas() fyne.Canvas {
	return ea.w.Canvas()
}

func (ea *demoAction) showDialog(int) {
	ea.mForm.ShowDialog(true)
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
	mDemoAction := newDemoAction(w)
	mRibbon := container.NewAppTabs()
	homeTab, _ := fyneextensions.BuildTabItemRibbon(mHomeAction, 60., 30., messageString)
	homeMenu := fyneextensions.NewActionableMenu(mHomeAction.GetActions())
	mRibbon.Append(homeTab)
	editTab, _ := fyneextensions.BuildTabItemRibbon(mEditAction, 60., 30., messageString)
	editMenu := fyneextensions.NewActionableMenu(mEditAction.GetActions())
	mRibbon.Append(editTab)
	demoTab, _ := fyneextensions.BuildTabItemRibbon(mDemoAction, 60., 30., messageString)
	demoMenu := fyneextensions.NewActionableMenu(mDemoAction.GetActions())
	mRibbon.Append(demoTab)

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
		demoMenu.Menu,
	))
	// Show and run the application
	w.ShowAndRun()

}
