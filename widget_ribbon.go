package fyneextensions

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"math"
)

/*
The MainRibbon struct is a fyne.Widget. Its main purpose is to provide a user interface
element that organizes other widgets into a ribbon-like structure.

It should be added to a container.AppTabs widget. The BuildTabItemRibbon function is a factory function
for creating new container.TabItem to be added to the container.AppTabs in Fyne.

It is based on ActionItem, which dictates the ribbon layout. Items can be laid out horizontally
or vertically, or in context menus, depending on ActionItem depth and length
*/
type MainRibbon struct {
	widget.BaseWidget

	items              []*ActionItem
	rems               []int
	canvas             fyne.Canvas
	maxSize, blockSize float32
	toolTipper         binding.String

	mContainer *fyne.Container
	mMasterCnt fyne.CanvasObject

	mMiniWidgets  []*MiniWidget
	sContainer    []*fyne.Container
	sAllObj       [][]fyne.CanvasObject
	sMenu         []*ActionableMenu
	sAllMenuItems [][]*ActionableMenuItem

	minSize      fyne.Size
	lastRenderer *mainRibbonRenderer
}

func (mr *MainRibbon) DataChanged() {
	mr.canvas.Refresh(mr)
	if mr.lastRenderer != nil {
		mr.lastRenderer.Layout(mr.Size())
	}
	mr.canvas.Refresh(mr)
}

func (mr *MainRibbon) CreateRenderer() fyne.WidgetRenderer {
	mr.lastRenderer = &mainRibbonRenderer{
		mRibbon: mr,
	}
	return mr.lastRenderer
}

func (mr *MainRibbon) MinSize() fyne.Size {
	return mr.minSize
}

type mainRibbonRenderer struct {
	mRibbon *MainRibbon
}

func (mrr *mainRibbonRenderer) Destroy() {}

func (mrr *mainRibbonRenderer) Layout(containerSize fyne.Size) {
	if containerSize.Width <= 0. {
		mrr.mRibbon.mMasterCnt.Resize(containerSize)
		mrr.mRibbon.mMasterCnt.Move(fyne.NewPos(0, 0))
		mrr.mRibbon.mContainer.Resize(containerSize)
		mrr.mRibbon.mContainer.Move(fyne.NewPos(0, 0))
		return
	}

	locRems := make([]int, len(mrr.mRibbon.rems))

	copy(locRems, mrr.mRibbon.rems)

	if containerSize.Width == mrr.mRibbon.mContainer.MinSize().Width {
		mrr.mRibbon.mMasterCnt.Resize(containerSize)
		mrr.mRibbon.mMasterCnt.Move(fyne.NewPos(0, 0))
		mrr.mRibbon.mContainer.Resize(containerSize)
		mrr.mRibbon.mContainer.Move(fyne.NewPos(0, 0))
		return
	}

	if containerSize.Width > mrr.mRibbon.mContainer.MinSize().Width {
		sOldZs := make([]fyne.Size, len(mrr.mRibbon.mMiniWidgets))
		for i, o := range mrr.mRibbon.mMiniWidgets {
			sOldZs[i] = o.MinSize()
		}

		locSkip := make([]bool, len(locRems))
		for i, o := range mrr.mRibbon.mMiniWidgets {
			locSkip[i] = !o.Visible()
		}

		for {
			i := 0
			for i = 0; i < len(locRems); i++ {
				if !locSkip[i] {
					if locRems[i] > 0 {
						locRems[i] -= 1
						break
					}
				}
			}
			if i == len(locRems) {
				break
			}

			mrr.mRibbon.sContainer[i].Objects = mrr.mRibbon.sAllObj[i][:len(mrr.mRibbon.sAllObj[i])-locRems[i]]
			mrr.mRibbon.sMenu[i].mActionableMenuItem.subActionableMenuItems = mrr.mRibbon.sAllMenuItems[i][len(mrr.mRibbon.sAllObj[i])-locRems[i]:]
			mrr.mRibbon.canvas.Refresh(mrr.mRibbon.sContainer[i])
			mrr.mRibbon.sMenu[i].DataChanged()

			sOldZs[i] = mrr.mRibbon.mMiniWidgets[i].MinSize()
			if mrr.mRibbon.mContainer.MinSize().Width > containerSize.Width {
				break
			}
		}
	}
	if mrr.mRibbon.mContainer.MinSize().Width > containerSize.Width {
		sOldZs := make([]fyne.Size, len(mrr.mRibbon.mMiniWidgets))
		for i, o := range mrr.mRibbon.mMiniWidgets {
			sOldZs[i] = o.MinSize()
		}

		locSkip := make([]bool, len(locRems))
		for i, o := range mrr.mRibbon.mMiniWidgets {
			locSkip[i] = !o.Visible()
		}

		for {
			i := 0
			for i = len(locRems) - 1; i >= 0; i-- {
				if !locSkip[i] {
					if locRems[i] < len(mrr.mRibbon.items[i].SubActions)-1 {
						locRems[i] += 1
						break
					}
				}
			}
			if i < 0 {
				break
			}

			mrr.mRibbon.sContainer[i].Objects = mrr.mRibbon.sAllObj[i][:len(mrr.mRibbon.sAllObj[i])-locRems[i]]
			mrr.mRibbon.sMenu[i].mActionableMenuItem.subActionableMenuItems = mrr.mRibbon.sAllMenuItems[i][len(mrr.mRibbon.sAllObj[i])-locRems[i]:]
			mrr.mRibbon.canvas.Refresh(mrr.mRibbon.sContainer[i])
			mrr.mRibbon.sMenu[i].DataChanged()

			if mrr.mRibbon.sAllObj[i][len(mrr.mRibbon.sContainer[i].Objects)-1].Visible() {
				if mrr.mRibbon.mMiniWidgets[i].MinSize().Width < sOldZs[i].Width {
					sOldZs[i] = mrr.mRibbon.mMiniWidgets[i].MinSize()
					if mrr.mRibbon.mContainer.MinSize().Width < containerSize.Width {
						break
					}
				} else {
					locRems[i] -= 1
					locSkip[i] = true
					mrr.mRibbon.sContainer[i].Objects = mrr.mRibbon.sAllObj[i][:len(mrr.mRibbon.sAllObj[i])-locRems[i]]
					mrr.mRibbon.sMenu[i].mActionableMenuItem.subActionableMenuItems = mrr.mRibbon.sAllMenuItems[i][len(mrr.mRibbon.sAllObj[i])-locRems[i]:]
					mrr.mRibbon.canvas.Refresh(mrr.mRibbon.sContainer[i])
					mrr.mRibbon.sMenu[i].DataChanged()
				}
			}
		}
	}

	copy(mrr.mRibbon.rems, locRems)

	for i, o := range mrr.mRibbon.sMenu {
		mrr.mRibbon.mMiniWidgets[i].setShowMore(len(o.mActionableMenuItem.subActionableMenuItems) != 0)
	}

	mrr.mRibbon.mMasterCnt.Resize(containerSize)
	mrr.mRibbon.mMasterCnt.Move(fyne.NewPos(0, 0))
	mrr.mRibbon.mContainer.Resize(containerSize)
	mrr.mRibbon.mContainer.Move(fyne.NewPos(0, 0))
}

func (mrr *mainRibbonRenderer) MinSize() fyne.Size {
	return mrr.mRibbon.MinSize()
}

func (mrr *mainRibbonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{mrr.mRibbon.mMasterCnt}
}

func (mrr *mainRibbonRenderer) Refresh() {
	mrr.mRibbon.mContainer.Refresh()
}

/*
BuildTabItemRibbon is a function that constructs a `container.TabItem` and a `MainRibbon` from an actionable item, dimensions (maxSize and blockSize), and a string for tooltips.

The function signature is:
- act is the object implementing Actionable interface which will indicate all ribbon functionalities.
- maxSize and blockSize are the maximum size and block size. These are used for laying out vertically components within the MainRibbon. For example, a maxSize of 90 and blockSize of 30 will allow up to 3 lines of objects in the ribbon
- toolTipper is a binding.String that serves for adding tooltips to components. This is optional. if set to nil, the tip will be displayed on a context popup instead when passing over a button with the mouse
*/
func BuildTabItemRibbon(act Actionable, maxSize, blockSize float32, toolTipper binding.String) (*container.TabItem, *MainRibbon) {
	item := act.GetActions()
	mCanvas := act.GetCanvas()
	var mContainer *MainRibbon

	if item.Triggered != nil {
		mContainer = newMainRibbon([]*ActionItem{item}, mCanvas, maxSize, blockSize, toolTipper)
	} else if len(item.SubActions) > 0 {
		mContainer = newMainRibbon(item.SubActions, mCanvas, maxSize, blockSize, toolTipper)
	}

	ribName, _ := item.Name.Get()
	retV := container.NewTabItem(ribName, container.NewStack(mContainer))
	item.Name.AddListener(mContainer)

	return retV, mContainer
}

func newMainRibbon(items []*ActionItem, mCanvas fyne.Canvas, maxSize, blockSize float32, toolTipper binding.String) *MainRibbon {
	mr := &MainRibbon{
		items:      items,
		mContainer: container.NewHBox(),
		canvas:     mCanvas,

		maxSize:    maxSize,
		blockSize:  blockSize,
		toolTipper: toolTipper,

		sContainer:    make([]*fyne.Container, 0),
		sAllObj:       make([][]fyne.CanvasObject, 0),
		sMenu:         make([]*ActionableMenu, 0),
		sAllMenuItems: make([][]*ActionableMenuItem, 0),
	}
	mr.ExtendBaseWidget(mr)

	mr.rems = make([]int, len(items))
	for i := range items {
		mr.rems[i] = len(items[i].SubActions) - 1
	}

	for i, o := range items {
		rb, sc, sm := buildL1Ribbon(o, mCanvas, maxSize, blockSize, toolTipper)
		mr.mContainer.Add(rb)
		mr.mMiniWidgets = append(mr.mMiniWidgets, rb)
		mr.sContainer = append(mr.sContainer, sc)
		mr.sAllObj = append(mr.sAllObj, sc.Objects)
		mr.sMenu = append(mr.sMenu, sm)
		mr.sAllMenuItems = append(mr.sAllMenuItems, sm.mActionableMenuItem.subActionableMenuItems)

		o.Name.AddListener(mr)
		o.Hider.AddListener(mr)
		o.Disabler.AddListener(mr)

		mr.sContainer[i].Objects = mr.sAllObj[i][:len(mr.sAllObj[i])-mr.rems[i]]
		mr.sMenu[i].mActionableMenuItem.subActionableMenuItems = mr.sAllMenuItems[i][len(mr.sAllObj[i])-mr.rems[i]:]
		mr.sContainer[i].Refresh()
		mr.sMenu[i].DataChanged()
	}

	mr.mMasterCnt = container.NewHScroll(mr.mContainer)

	mr.minSize = mr.mMasterCnt.MinSize()

	return mr
}

func (mr *MainRibbon) AddItems(items ...*ActionItem) {
	for i := range items {
		mr.rems = append(mr.rems, len(items[i].SubActions)-1)
	}

	oldLen := len(mr.items)

	for j, o := range items {
		mr.items = append(mr.items, o)
		rb, sc, sm := buildL1Ribbon(o, mr.canvas, mr.maxSize, mr.blockSize, mr.toolTipper)
		mr.mContainer.Add(rb)
		mr.mMiniWidgets = append(mr.mMiniWidgets, rb)
		mr.sContainer = append(mr.sContainer, sc)
		mr.sAllObj = append(mr.sAllObj, sc.Objects)
		mr.sMenu = append(mr.sMenu, sm)
		mr.sAllMenuItems = append(mr.sAllMenuItems, sm.mActionableMenuItem.subActionableMenuItems)

		o.Name.AddListener(mr)
		o.Hider.AddListener(mr)
		o.Disabler.AddListener(mr)

		i := oldLen + j
		mr.sContainer[i].Objects = mr.sAllObj[i][:len(mr.sAllObj[i])-mr.rems[i]]
		mr.sMenu[i].mActionableMenuItem.subActionableMenuItems = mr.sAllMenuItems[i][len(mr.sAllObj[i])-mr.rems[i]:]
		mr.sContainer[i].Refresh()
		mr.sMenu[i].DataChanged()
	}
	//mr.Refresh()
}

func buildL1Ribbon(item *ActionItem, mCanvas fyne.Canvas, maxSize, blockSize float32, toolTipper binding.String) (*MiniWidget, *fyne.Container, *ActionableMenu) {
	mContent := container.New(&ExpandingAllProportionallyPaddedHBox{})

	var moreMenu *ActionableMenu
	var moreFunc func(object fyne.CanvasObject)

	if item.Triggered != nil {
		nb := NewFlexButton("", item.Resources, false, !item.CriticalName, true, false, true, maxSize, blockSize, mCanvas, item.Triggered, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
		mContent.Add(nb)
	} else if len(item.SubActions) > 0 {
		for _, o := range item.SubActions {
			mContent.Add(buildL2Ribbon(o, mCanvas, maxSize, blockSize, toolTipper))
		}

		moreActItm := NewActionItem("internal Menu, bug if visible", false, false, item.Resources, false, false, false, 0, nil, item.SubActions)
		moreMenu = NewActionableMenu(moreActItm)
		moreFunc = func(co fyne.CanvasObject) {
			moreMenu.DataChanged()
			widget.ShowPopUpMenuAtRelativePosition(moreMenu.Menu, mCanvas, fyne.NewPos(0., co.Size().Height), co)
		}

	} else {
		panic(fmt.Errorf("error in ActionItem: nor func nor container"))
		return nil, nil, nil
	}

	mwName := ""
	if item.Name != nil {
		mwName, _ = item.Name.Get()
	}

	mw := NewMiniWidget(
		mwName,
		false,
		20.,
		mContent,
		true,
		false, nil,
		false, nil,
		false, nil, nil,
		moreFunc != nil, moreFunc,
		item.Name,
		item.Disabler,
		item.Hider,
		nil,
		mCanvas,
	)

	if item.Name != nil {
		item.Name.AddListener(mw)
	}
	if item.Disabler != nil {
		item.Disabler.AddListener(mw)
	}
	if item.Hider != nil {
		item.Hider.AddListener(mw)
		if hidden, _ := item.Hider.Get(); hidden {
			mw.Hide()
		}
	}

	return mw, mContent, moreMenu
}

func buildL2Ribbon(item *ActionItem, mCanvas fyne.Canvas, maxSize, blockSize float32, toolTipper binding.String) (mContent fyne.CanvasObject) {
	if item.Triggered != nil {
		nb := NewFlexButton("", item.Resources, false, !item.CriticalName, false, false, false, maxSize, blockSize, mCanvas, item.Triggered, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
		mContent = nb
	} else if item.AlwaysShowAsContainer {
		nb := NewFlexButton("", item.Resources, false, false, false, true, true, maxSize, blockSize, mCanvas, func(int) {}, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
		mContent = nb

		sMenu := NewActionableMenu2(item.SubActions...).Menu
		nb.OnTapped = func(int) {
			widget.ShowPopUpMenuAtRelativePosition(sMenu, mCanvas, fyne.NewPos(0., nb.Size().Height), nb)
		}
	} else if len(item.SubActions) > 0 {
		maxItems := int(math.Floor(float64(maxSize / blockSize)))
		if len(item.SubActions) <= maxItems {
			mContainer := container.New(&EquallySpacedUnpaddedVBox{})
			for _, o := range item.SubActions {
				mContainer.Add(buildL3Ribbon(o, mCanvas, maxSize/float32(len(item.SubActions)), blockSize, toolTipper))
			}
			mContent = mContainer
		} else {
			nb := NewFlexButton("", item.Resources, false, !item.CriticalName, false, true, true, maxSize, blockSize, mCanvas, item.Triggered, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
			mContent = nb

			sMenu := NewActionableMenu2(item.SubActions...).Menu
			nb.OnTapped = func(int) {
				widget.ShowPopUpMenuAtRelativePosition(sMenu, mCanvas, fyne.NewPos(0., nb.Size().Height), nb)
			}
		}
	} else {
		panic(fmt.Errorf("error in ActionItem: nor func nor container"))
	}
	/*
		if item.Hider != nil {
			if hidden, err := item.Hider.Get(); err == nil {
				if hidden {
					mContent.Hide()
				}
			}
		}
	*/
	return
}

func buildL3Ribbon(item *ActionItem, mCanvas fyne.Canvas, maxSize, blockSize float32, toolTipper binding.String) (mContent fyne.CanvasObject) {
	if item.Triggered != nil {
		nb := NewFlexButton("", item.Resources, true, !item.CriticalName, false, false, false, maxSize, blockSize, mCanvas, item.Triggered, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
		mContent = nb
	} else if len(item.SubActions) > 0 {
		if len(item.SubActions) < 4 {
			mContainer := container.New(&ExpandingFirstPaddedHBox{})
			for _, o := range item.SubActions {
				mContainer.Add(buildL4Ribbon(o, mCanvas, maxSize, blockSize, toolTipper))
			}
			mContent = mContainer
		} else {
			nb := NewFlexButton("", item.Resources, true, !item.CriticalName, false, true, false, maxSize, blockSize, mCanvas, item.Triggered, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
			mContent = nb

			sMenu := NewActionableMenu2(item.SubActions...).Menu

			nb.OnTapped = func(int) {
				sPopup := widget.NewPopUpMenu(sMenu, mCanvas)
				sPopup.ShowAtRelativePosition(fyne.NewPos(nb.Size().Width, 0.), nb)
			}
		}
	} else {
		panic(fmt.Errorf("error in ActionItem: nor func nor container"))
	}
	return
}

func buildL4Ribbon(item *ActionItem, mCanvas fyne.Canvas, maxSize, blockSize float32, toolTipper binding.String) (mContent fyne.CanvasObject) {
	if item.Triggered != nil {
		mContent = NewFlexButton("", item.Resources, true, !item.CriticalName, false, false, false, maxSize, blockSize, mCanvas, item.Triggered, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
	} else if len(item.SubActions) > 0 {
		nb := NewFlexButton("", item.Resources, true, !item.CriticalName, false, true, false, maxSize, blockSize, mCanvas, item.Triggered, item.Name, item.Disabler, item.Hider, item.Stater, toolTipper)
		mContent = nb

		sMenu := NewActionableMenu2(item.SubActions...).Menu

		nb.OnTapped = func(int) {
			sPopup := widget.NewPopUpMenu(sMenu, mCanvas)
			sPopup.ShowAtRelativePosition(fyne.NewPos(nb.Size().Width, 0.), nb)
		}
	} else {
		panic(fmt.Errorf("error in ActionItem: nor func nor container"))
	}
	return
}
