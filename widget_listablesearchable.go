package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*
ListableSearchable defines the basic operations and behavior required for a ListableSearchableWidget.

Methods:
- ListLen() int: This method returns the length of the list.
- CreateListItem() fyne.CanvasObject: This method creates a list item and returns a fyne.CanvasObject. The concrete implementation of this object can vary. Note that this constructor does not need actual data to be displayed, but only the constructors fof the CanvasObjects which will hold the data
- UpdateListItem(int, fyne.CanvasObject): This method updates a list item at a given index position with a fyne.CanvasObject. This will introduce actual data into the list
- GetSelected() int: This method returns the index position of the currently selected item in the list.
- OnSelectedItem(int): This method defines what to do when an item in the list is selected. The specific action is determined by the integer parameter, which represents the index of the item.
- GetListableSearchableActions() *ActionItem: This method returns the actions related to ListableSearchable widget. If the returned *ActionItem is not nil, the widget will display a customized menu
- GetCanvas() fyne.Canvas: This method must return the canvas object where the ListableSearchableWidget will be displayed.
- ClearSearch(): Clears the current search results.
- StartSearch(string): This method starts a search operation on the list using the input string as the query.
*/
type ListableSearchable interface {
	ListLen() int
	CreateListItem() fyne.CanvasObject
	UpdateListItem(int, fyne.CanvasObject)

	GetSelected() int
	OnSelectedItem(int)

	GetListableSearchableActions() *ActionItem
	GetCanvas() fyne.Canvas
	ClearSearch()
	StartSearch(string)
}

/*
The ListableSearchableWidget struct is a Fyne compatible widget.
The widget can display a list of items with a search bar with a default Search and Clear button.
The search bar can also include a customized menu.

it is possible to associate to a ListableSearchableWidget any object which implements ListableSearchable interface

Constructor:
- NewListableSearchableWidget: Returns a new instance of the ListableSearchableWidget struct.
*/
type ListableSearchableWidget struct {
	widget.BaseWidget

	list ListableSearchable

	mainContainer *fyne.Container
	searchEdit    *widget.Entry
	moreButton    *widget.Button
	searchButton  *widget.Button
	clearButton   *widget.Button
	listView      *widget.List
}

func NewListableSearchableWidget(iList ListableSearchable) *ListableSearchableWidget {
	t := &ListableSearchableWidget{}
	t.ExtendBaseWidget(t)

	t.list = iList

	t.searchEdit = widget.NewEntry()
	t.searchEdit.SetPlaceHolder("Search")

	t.listView = widget.NewList(
		func() int {
			if t.list == nil {
				return 0
			}
			return t.list.ListLen()
		},
		func() fyne.CanvasObject {
			if t.list != nil {
				return t.list.CreateListItem()
			}
			return container.NewWithoutLayout()
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if t.list != nil {
				t.list.UpdateListItem(i, o)
			}
		},
	)

	t.listView.OnSelected = func(id widget.ListItemID) {
		if t.list != nil {
			t.list.OnSelectedItem(id)
			//window.Canvas().Refresh(masterContainer)
		}
	}

	t.moreButton = widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), nil)
	t.moreButton.Hide()
	if iList != nil {
		mActionItem := iList.GetListableSearchableActions()
		mCanvas := iList.GetCanvas()
		if mActionItem != nil {
			sMenu := NewActionableMenu2(mActionItem.SubActions...).Menu
			t.moreButton.Show()
			t.moreButton.OnTapped = func() {
				sPopup := widget.NewPopUpMenu(sMenu, mCanvas)
				sPopup.ShowAtRelativePosition(fyne.NewPos(t.moreButton.Size().Width, 0.), t.moreButton)
			}

		}
	}

	t.searchButton = widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		if t.list != nil {
			t.list.StartSearch(t.searchEdit.Text)
		}
		//parent.Canvas().Refresh(masterContainer)
		t.searchEdit.Refresh()
		t.listView.Refresh()
	})

	t.clearButton = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		if t.list != nil {
			t.list.ClearSearch()
		}
		t.searchEdit.Text = ""
		t.searchEdit.Refresh()
		t.listView.Refresh()
		//parent.Canvas().Refresh(masterContainer)
	})

	toolBrContainer := container.New(&ExpandingFirstUnpaddedHBox{})
	toolBrContainer.Add(t.searchEdit)
	toolBrContainer.Add(t.moreButton)
	toolBrContainer.Add(t.searchButton)
	toolBrContainer.Add(t.clearButton)

	t.mainContainer = container.New(&ExpandingLastUnpaddedVBox{})
	t.mainContainer.Add(toolBrContainer)
	t.mainContainer.Add(t.listView)

	return t
}

func (t *ListableSearchableWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.mainContainer)
}

func (t *ListableSearchableWidget) DataChanged() {
	t.searchEdit.Text = ""

	t.searchEdit.Refresh()
	t.listView.Refresh()
	if t.list != nil {
		if selId := t.list.GetSelected(); selId != -1 {
			t.listView.Select(selId)
		}
	}
}

func (t *ListableSearchableWidget) ResetListableSearchableAndRefresh(mList ListableSearchable) {
	if t.list == mList {
		if t.list != nil {
			if selId := t.list.GetSelected(); selId != -1 {
				t.listView.Select(selId)
			}
		}
		t.listView.Refresh()
	} else {
		t.list = mList

		t.searchEdit.Text = ""

		t.searchEdit.Refresh()
		t.listView.Refresh()
		if t.list != nil {
			if selId := t.list.GetSelected(); selId != -1 {
				t.listView.Select(selId)
			}
		}
	}
}
