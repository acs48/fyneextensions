package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

const MaxMenuItems int = 100

type ActionableMenuItem struct {
	mActionItem *ActionItem

	subActionableMenuItems []*ActionableMenuItem

	mItem *fyne.MenuItem
}

/*
NewActionableMenuItem is a constructor function that takes an `*ActionItem`,
a `parentItem` and a `rootItem` both as `binding.DataListener` and produces an `*ActionableMenuItem`.
ParentItem and RootItem are typically nil.

In a nutshell, NewActionableMenuItem provides a way to create a MenuItem which is actionable,
and is associated with other data listeners to react to changes. This allows for dynamically
adjusting menu items based on the program state
*/
func NewActionableMenuItem(item *ActionItem, parentItem binding.DataListener, rootItem binding.DataListener) *ActionableMenuItem {
	ami := &ActionableMenuItem{
		mActionItem: item,
	}

	ami.mItem = fyne.NewMenuItem("dummy", nil)
	if len(ami.mActionItem.Resources) > 0 {
		ami.mItem.Icon = ami.mActionItem.Resources[0]
	}
	ami.mItem.Action = func() {
		if item.Triggered != nil {
			ms := 0
			if item.Stater != nil {
				ms, _ = item.Stater.Get()
			}
			item.Triggered(ms)
		}
	}

	for _, o := range item.SubActions {
		ami.subActionableMenuItems = append(ami.subActionableMenuItems, NewActionableMenuItem(o, ami, rootItem))
	}

	if len(ami.subActionableMenuItems) > 0 {
		ami.mItem.ChildMenu = fyne.NewMenu("dummy")
		ami.mItem.ChildMenu.Items = make([]*fyne.MenuItem, MaxMenuItems)[0:0]
	}

	if ami.mActionItem.Name != nil {
		ami.mActionItem.Name.AddListener(ami)
		if rootItem != nil {
			ami.mActionItem.Name.AddListener(rootItem)
		}
	}
	if ami.mActionItem.Disabler != nil {
		ami.mActionItem.Disabler.AddListener(ami)
		//if parentItem != nil {
		//	ami.mActionItem.Disabler.AddListener(parentItem)
		//}
	}
	if ami.mActionItem.Stater != nil {
		ami.mActionItem.Stater.AddListener(ami)
		//if parentItem != nil {
		//	ami.mActionItem.Stater.AddListener(parentItem)
		//}
	}
	if ami.mActionItem.Hider != nil {
		if parentItem != nil {
			ami.mActionItem.Hider.AddListener(parentItem)
			if parentItem != rootItem && rootItem != nil {
				ami.mActionItem.Hider.AddListener(rootItem)
			}
		}
	}

	return ami
}

func (ami *ActionableMenuItem) DataChanged() {
	if ami.mActionItem.Name != nil {
		if name, err := ami.mActionItem.Name.Get(); err == nil {
			ami.mItem.Label = name
		}
	}
	if ami.mActionItem.Disabler != nil {
		if disabled, err := ami.mActionItem.Disabler.Get(); err == nil {
			ami.mItem.Disabled = disabled
		}
	}
	if ami.mActionItem.Stater != nil {
		if state, err := ami.mActionItem.Stater.Get(); err == nil {
			if len(ami.mActionItem.Resources) > state {
				ami.mItem.Icon = ami.mActionItem.Resources[state]
			}
		}
	}
	if ami.mItem.ChildMenu != nil {
		var newItems []*fyne.MenuItem
		for _, o := range ami.subActionableMenuItems {
			newItems = append(newItems, o.getItems()...)
		}
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
		if len(newItems) < cap(ami.mItem.ChildMenu.Items) {
			ami.mItem.ChildMenu.Items = ami.mItem.ChildMenu.Items[:len(newItems)]
			for i, o := range newItems {
				ami.mItem.ChildMenu.Items[i] = o
			}
			ami.mItem.ChildMenu.Refresh()
		} else {
			//panic(fmt.Errorf("max num of sub menu items reached"))
		}
	}
}

func (ami *ActionableMenuItem) getItems() (items []*fyne.MenuItem) {
	item := ami.mActionItem
	if item.Hider != nil {
		if hidden, err := item.Hider.Get(); err == nil {
			if hidden {
				return nil
			}
		}
	}

	hasMore := len(ami.subActionableMenuItems) > 0

	if !hasMore {
		items = append(items, ami.mItem)
		return
	}

	if item.AlwaysShowAsContainer {
		items = append(items, ami.mItem)
		return
	}

	items = append(items, fyne.NewMenuItemSeparator())
	for _, o := range ami.subActionableMenuItems {
		items = append(items, o.getItems()...)
	}

	return
}

/*
The GetMenu method is a part of the `ActionableMenuItem` struct and it returns a pointer to a `fyne.Menu` or `nil`.

This function would be called when you want to obtain the associated menu of an ActionableMenuItem,
if it exists, updating its data if there have been changes to the associated ActionItem.
*/
func (ami *ActionableMenuItem) GetMenu() *fyne.Menu {

	if ami.mItem.ChildMenu != nil {
		ami.DataChanged()
		return ami.mItem.ChildMenu
	}
	return nil
}
