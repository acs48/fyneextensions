package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MiniWidget is a Fyne compatible widget. It represents a miniaturized widget structure in the GUI
// which can be adapted to multiple usages, such as menu ribbons, sidebars, etc.
// It's ideal for use-cases like notifications, small data charts, or compact control elements.
// It can be configured with a header or footer area, with close, minimize, move up and move down buttons
// An instance of a MiniWidget can be created with the factory NewMiniWidget
type MiniWidget struct {
	widget.DisableableWidget
	mContainer      *fyne.Container
	widContainer    *fyne.Container
	boxContainer    *fyne.Container
	shadowContainer *fyne.Container
	mHeader         *fyne.Container
	mLabel          *SizableLabel
	mContent        fyne.CanvasObject
	mBackground     *canvas.Rectangle
	mShadow         *canvas.Rectangle
	mSeparator      *widget.Separator
	mSize           float32

	disabler  binding.Bool
	texter    binding.String
	closer    binding.Bool
	minimizer binding.Int

	onMinimize  func(bool)
	minimizeBtn *FlexButton

	onClose  func()
	closeBtn *FlexButton

	onMoveUp  func(object *MiniWidget)
	moveUpBtn *FlexButton

	onMoveDown  func(object *MiniWidget)
	moveDownBtn *FlexButton

	onMore      func(object fyne.CanvasObject)
	mMoreButton *FlexButton
}

// NewMiniWidget is a factory function that creates and initializes a new MiniWidget.
//
// The function accepts the following parameters:
//
//	header: a string that sets the header of the MiniWidget.
//	headerOnTop: a bool that determines if the header should be on top or bottom of the MiniWidget.
//	headerSize: a float32 that sets the size of the header.
//	content: a fyne.CanvasObject which represents the content of the MiniWidget.
//	shadowed: a bool to enable or disable shadow for the MiniWidget.
//	showMinimize: a bool to show or hide minimize functionality for the MiniWidget.
//	onMinimize: a callback function that gets executed when the MiniWidget is minimized.
//	showClose: a bool to show or hide close functionality for the MiniWidget.
//	onClose: a callback function that gets executed when the MiniWidget is closed.
//	showMove: a bool to enable or disable movement functionality for the MiniWidget.
//	onMoveUp: a callback function that gets executed when the MiniWidget is moved upwards.
//	onMoveDown: a callback function that gets executed when the MiniWidget is moved downwards.
//	showMore: a bool to enable or disable more functionality for the MiniWidget.
//	onMore: a callback function that gets executed when the more button on the MiniWidget is clicked.
//	texter: a binding.String instance to bind header text to the MiniWidget.
//	disabler: a binding.Bool instance to disable the MiniWidget.
//	closer: a binding.Bool instance to bind custom close function to the MiniWidget.
//	minimizer: a binding.Int instance to bind custom minimize function to the MiniWidget.
//	mCanvas: a fyne.Canvas instance where the MiniWidget will be drawn.
//
// This function returns a pointer to a fully initialized MiniWidget.
func NewMiniWidget(header string, headerOnTop bool, headerSize float32,
	content fyne.CanvasObject, shadowed bool,
	showMinimize bool, onMinimize func(bool), showClose bool, onClose func(), showMove bool, onMoveUp func(miniWidget *MiniWidget), onMoveDown func(miniWidget *MiniWidget), showMore bool, onMore func(fyne.CanvasObject),
	texter binding.String, disabler binding.Bool, closer binding.Bool, minimizer binding.Int, mCanvas fyne.Canvas) *MiniWidget {
	t := &MiniWidget{
		mSize:       headerSize,
		mContent:    content,
		mHeader:     container.New(&ExpandingFirstUnpaddedHBox{}),
		mBackground: canvas.NewRectangle(theme.BackgroundColor()),
		mShadow:     nil,
		mSeparator:  widget.NewSeparator(),
		mLabel:      NewSizableLabel(header, headerSize, false, false, theme.ForegroundColor(), theme.ButtonColor()),

		texter:    texter,
		disabler:  disabler,
		closer:    closer,
		minimizer: minimizer,

		onClose:    onClose,
		onMinimize: onMinimize,
		onMoveUp:   onMoveUp,
		onMoveDown: onMoveDown,
		onMore:     onMore,
	}
	t.ExtendBaseWidget(t)

	if t.texter != nil {
		t.texter.AddListener(t)
		if tx, err := t.texter.Get(); err == nil {
			header = tx
		}
	}
	if t.closer == nil {
		t.closer = binding.NewBool()
	}
	t.closer.AddListener(t)
	if t.disabler != nil {
		t.disabler.AddListener(t)
	}
	if t.minimizer == nil {
		t.minimizer = binding.NewInt()
	}
	t.minimizer.AddListener(t)

	t.mBackground.StrokeColor = theme.FocusColor()
	t.mBackground.FillColor = theme.BackgroundColor()
	t.mBackground.StrokeWidth = 2.

	t.mHeader.Add(t.mLabel)

	t.moveUpBtn = NewFlexButton("", []fyne.Resource{theme.MoveUpIcon()}, false, true, false, false, false, headerSize, 0., mCanvas, func(int) {
		if t.onMoveUp != nil {
			t.onMoveUp(t)
		}
	}, nil, nil, nil, nil, nil)
	t.moveDownBtn = NewFlexButton("", []fyne.Resource{theme.MoveDownIcon()}, false, true, false, false, false, headerSize, 0., mCanvas, func(int) {
		if t.onMoveDown != nil {
			t.onMoveDown(t)
		}
	}, nil, nil, nil, nil, nil)
	t.mHeader.Add(t.moveUpBtn)
	t.mHeader.Add(t.moveDownBtn)
	if !showMove {
		t.moveUpBtn.Hide()
		t.moveDownBtn.Hide()
	}

	t.mMoreButton = NewFlexButton("", []fyne.Resource{theme.MoreVerticalIcon()}, false, true, false, false, false, headerSize, 0., mCanvas, nil, nil, nil, nil, nil, nil)
	t.mHeader.Add(t.mMoreButton)
	if !showMore {
		t.mMoreButton.Hide()
	}
	t.mMoreButton.OnTapped = func(int) {
		if t.onMore != nil {
			t.onMore(t.mMoreButton)
		}
	}

	t.minimizeBtn = NewFlexButton("", []fyne.Resource{theme.ContentRemoveIcon(), theme.ContentAddIcon()}, false, true, false, false, false, headerSize, 0., mCanvas, nil, nil, nil, nil, t.minimizer, nil)
	t.minimizeBtn.OnTapped = func(state int) {
		if state == 1 {
			//t.mContent.Show()
			//t.mSeparator.Show()
			//t.minimized = false
			t.minimizer.Set(0)
			//miniButton.SetIcon(theme.ContentRemoveIcon())
		} else {
			//t.mContent.Hide()
			//t.mSeparator.Hide()
			//t.minimized = true
			t.minimizer.Set(1)
			//miniButton.SetIcon(theme.ContentAddIcon())
		}
		if t.onMinimize != nil {
			t.onMinimize(state == 1)
		}
		t.Refresh()
	}
	t.mHeader.Add(t.minimizeBtn)
	if !showMinimize {
		t.minimizeBtn.Hide()
	}

	t.closeBtn = NewFlexButton("", []fyne.Resource{theme.ContentClearIcon()}, false, true, false, false, false, headerSize, 0., mCanvas, func(int) {
		//t.Hide()
		t.closer.Set(true)

		if t.onClose != nil {
			t.onClose()
		}
	}, nil, nil, nil, nil, nil)
	t.mHeader.Add(t.closeBtn)
	if !showClose {
		t.closeBtn.Hide()
	}

	if headerOnTop {
		t.widContainer = container.New(&ExpandingLastUnpaddedVBox{}, t.mHeader, t.mSeparator, t.mContent)
	} else {
		t.widContainer = container.New(&ExpandingFirstUnpaddedVBox{}, t.mContent, t.mSeparator, t.mHeader)
	}

	t.boxContainer = container.New(&PaddedBox{}, t.mBackground, t.widContainer)

	if shadowed {
		t.mShadow = canvas.NewRectangle(theme.ShadowColor())
		t.shadowContainer = container.New(&ShadowedLayout{}, t.mShadow, t.boxContainer)
	} else {
		t.shadowContainer = t.boxContainer
	}

	t.mContainer = t.shadowContainer

	t.DataChanged()

	return t
}

func (t *MiniWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.mContainer)
}

func (t *MiniWidget) Minimize(state int) {

}

func (t *MiniWidget) MinSize() fyne.Size {
	return t.mContainer.MinSize()
}

func (t *MiniWidget) Refresh() {
	if t.mShadow != nil {
		t.mShadow.StrokeColor = theme.ShadowColor()
		t.mShadow.FillColor = theme.ShadowColor()
		t.mShadow.Refresh()
	}

	//t.mBackground.SetMinSize(t.widContainer.MinSize().Add(fyne.NewDelta(6., 6.)))
	t.mBackground.StrokeColor = theme.FocusColor()
	t.mBackground.FillColor = theme.BackgroundColor()
	t.mBackground.StrokeWidth = 2.

	t.mSeparator.Refresh()
	t.mContent.Refresh()
	t.mHeader.Refresh()
	t.mBackground.Refresh()

	t.mLabel.mTextColor = theme.ForegroundColor()
	t.mLabel.mBackgroundColor = theme.ButtonColor()
	t.mLabel.Refresh()
}

func (t *MiniWidget) setShowMore(showMore bool) {
	if showMore {
		t.mMoreButton.Show()
	} else {
		t.mMoreButton.Hide()
	}
}

func (t *MiniWidget) DataChanged() {
	if t.disabler != nil {
		d, err := t.disabler.Get()
		if err == nil {
			if d {
				t.Disable()
			} else {
				t.Enable()
			}
		}
	}

	if t.closer != nil {
		if d, err := t.closer.Get(); err == nil {
			if d {
				t.Hide()
			} else {
				t.Show()
			}
		}
	}

	if t.minimizer != nil {
		if d, err := t.minimizer.Get(); err == nil {
			if d == 1 {
				t.mContent.Hide()
				t.mSeparator.Hide()
			} else {
				t.mContent.Show()
				t.mSeparator.Show()
			}
		}
	}

	if t.texter != nil {
		if t.texter != nil {
			if d, err := t.texter.Get(); err == nil {
				if t.mLabel.mText.Text != d {
					t.mLabel.mText.Text = d
					t.Refresh()
				}
			}
		}
	}

	t.Refresh()
}
