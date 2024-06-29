package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SideBar is a Fyne compatible widget which acts as a container for MiniWidget components.
// It can manage basic functionalities such as adding content to the sidebar and moving child objects up and down
type SideBar struct {
	widget.BaseWidget

	mContent *fyne.Container
	mObjects []*MiniWidget
}

func NewSideBar(widgets ...*MiniWidget) *SideBar {
	sb := &SideBar{
		mContent: container.NewStack(),
	}

	sb.ExtendBaseWidget(sb)

	sb.AddWidget(widgets...)

	return sb
}

func (sb *SideBar) DataChanged() {
	sb.Refresh()
}

func (sb *SideBar) AddWidget(widgets ...*MiniWidget) {
	for _, o := range widgets {
		sb.mObjects = append(sb.mObjects, o)

		o.minimizeBtn.Show()
		o.minimizer.AddListener(sb)
		o.closer.AddListener(sb)

		o.moveUpBtn.Show()
		o.onMoveUp = sb.moveUpChild

		o.moveDownBtn.Show()
		o.onMoveDown = sb.moveDownChild
	}
	//sb.Refresh()
}

func (sb *SideBar) moveUpChild(mw *MiniWidget) {
	itmNum := -1

	for i, o := range sb.mObjects {
		if o == mw {
			itmNum = i
			break
		}
	}

	if itmNum > 0 {
		sb.mObjects[itmNum-1], sb.mObjects[itmNum] = sb.mObjects[itmNum], sb.mObjects[itmNum-1]
	}

	sb.Refresh()
}

func (sb *SideBar) moveDownChild(mw *MiniWidget) {
	itmNum := -1

	for i, o := range sb.mObjects {
		if o == mw {
			itmNum = i
			break
		}
	}

	if itmNum >= 0 && itmNum < len(sb.mObjects)-1 {
		sb.mObjects[itmNum+1], sb.mObjects[itmNum] = sb.mObjects[itmNum], sb.mObjects[itmNum+1]
	}
	sb.Refresh()
}

func (sb *SideBar) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(sb.mContent)
}

func (sb *SideBar) MinSize() fyne.Size {
	return sb.mContent.MinSize()
}

func (sb *SideBar) Refresh() {
	sb.mContent.RemoveAll()
	visibleObj := make([]*MiniWidget, 0)
	for _, o := range sb.mObjects {
		if closed, _ := o.closer.Get(); !closed {
			visibleObj = append(visibleObj, o)
		}
	}

	groups := make([][]*MiniWidget, 0)
	groups = append(groups, make([]*MiniWidget, 0))
	cGp := 0
	for _, o := range visibleObj {
		groups[cGp] = append(groups[cGp], o)
		if isMin, _ := o.minimizer.Get(); isMin == 0 {
			groups = append(groups, make([]*MiniWidget, 0))
			cGp++
		}
	}

	lastGpAllMin := true
	var top fyne.CanvasObject
	var topCnt *fyne.Container
	var btm fyne.CanvasObject
	var btmCnt *fyne.Container

	if len(groups) > 0 {
		if len(groups[len(groups)-1]) == 0 {
			groups = groups[:len(groups)-1]
		}

		for _, o := range groups[len(groups)-1] {
			if isMin, _ := o.minimizer.Get(); isMin == 0 {
				lastGpAllMin = false
				break
			}
		}

		if lastGpAllMin {
			btmCnt = container.NewVBox()
		} else {
			btmCnt = container.New(&ExpandingLastPaddedVBox{})
		}
		for _, o := range groups[len(groups)-1] {
			btmCnt.Add(o)
		}
		btm = btmCnt
	}

	if len(groups) > 1 {
		topCnt = container.New(&ExpandingLastPaddedVBox{})
		for _, o := range groups[len(groups)-2] {
			topCnt.Add(o)
		}
		top = topCnt

		if lastGpAllMin {
			btmCnt = container.New(&ExpandingFirstPaddedVBox{}, top, btm)
			btm = btmCnt
		} else {
			btm = container.NewVSplit(top, btm)
		}
	}

	if len(groups) > 2 {
		for i := len(groups) - 3; i >= 0; i-- {
			topCnt = container.New(&ExpandingLastPaddedVBox{})
			for _, o := range groups[i] {
				topCnt.Add(o)
			}
			top = topCnt
			btm = container.NewVSplit(top, btm)
		}
	}
	if btm != nil {
		sb.mContent.Add(btm)
	}
}
