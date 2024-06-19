package fyneextensions

import (
	"fyne.io/fyne/v2"
)

// ActionableMenu represents a Menu that has
// actions associated with it.
// Implements data binding listener interfaces to enable
// reactive behaviours to configuration changes such as
// name, being disabled, state changes and visibility.
type ActionableMenu struct {
	mActionItem *ActionItem

	mActionableMenuItem *ActionableMenuItem

	// Menu represents the actual fyne menu
	Menu *fyne.Menu
}

// NewActionableMenu creates a new ActionableMenu from a given ActionItem.
// It initializes a new Fyne menu and an actionable MenuItem for the provided Action.
// If any of the listener interfaces (Name, Disabler, Stater, Hider)
// for the associated ActionItem are non-nil, the newly created ActionableMenu
// is added as a listener to these.
// It returns a pointer to the constructed ActionableMenu.
// Input: an ActionItem to be associated with the new ActionableMenu.
// Output: a pointer to the newly created ActionableMenu.
func NewActionableMenu(item *ActionItem) *ActionableMenu {
	am3 := &ActionableMenu{
		mActionItem: item,
		Menu:        fyne.NewMenu(""),
	}

	am3.mActionableMenuItem = NewActionableMenuItem(item, am3, am3)

	am3.Menu.Items = make([]*fyne.MenuItem, MaxMenuItems)[0:0]

	if am3.mActionItem.Name != nil {
		am3.mActionItem.Name.AddListener(am3)
	}
	if am3.mActionItem.Disabler != nil {
		am3.mActionItem.Disabler.AddListener(am3)
	}
	if am3.mActionItem.Stater != nil {
		am3.mActionItem.Stater.AddListener(am3)
	}
	if am3.mActionItem.Hider != nil {
		am3.mActionItem.Hider.AddListener(am3)
	}

	return am3
}

// NewActionableMenu2 creates a new ActionableMenu from a given slice of ActionItem.
func NewActionableMenu2(items ...*ActionItem) *ActionableMenu {
	item := NewActionItem("", false, false, []fyne.Resource{}, false, false, false, 0, nil, items)

	am := &ActionableMenu{
		mActionItem: item,
		Menu:        fyne.NewMenu(""),
	}

	am.mActionableMenuItem = NewActionableMenuItem(item, am, am)

	am.Menu.Items = make([]*fyne.MenuItem, MaxMenuItems)[0:0]

	if am.mActionItem.Name != nil {
		am.mActionItem.Name.AddListener(am)
	}
	if am.mActionItem.Disabler != nil {
		am.mActionItem.Disabler.AddListener(am)
	}
	if am.mActionItem.Stater != nil {
		am.mActionItem.Stater.AddListener(am)
	}
	if am.mActionItem.Hider != nil {
		am.mActionItem.Hider.AddListener(am)
	}

	return am
}

func (am *ActionableMenu) DataChanged() {
	if am.mActionItem.Name != nil {
		if name, err := am.mActionItem.Name.Get(); err == nil {
			am.Menu.Label = name
		}
	}

	if am.mActionableMenuItem != nil {
		newItems := am.mActionableMenuItem.getItems()
		if len(newItems) > 0 {
			if newItems[0].IsSeparator {
				newItems = newItems[1:]
			}
		}
		if len(newItems) > 0 {
			if newItems[len(newItems)-1].IsSeparator {
				newItems = newItems[:len(newItems)-1]
			}
		}
		if len(newItems) < cap(am.Menu.Items) {
			am.Menu.Items = am.Menu.Items[:len(newItems)]
			for i, o := range newItems {
				am.Menu.Items[i] = o
			}
			am.Menu.Refresh()
		} else {
			//panic(fmt.Errorf("max number of child items reached"))
		}
	}
}
