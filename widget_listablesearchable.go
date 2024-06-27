package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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
