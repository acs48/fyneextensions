package fyneextensions

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

/*
FlexButton is a fyne compatible widget which defines a flexible and customizable button widget.
It provides a consistent look and feel across different parts of our app,
and encapsulates styles and behaviors that are commonly associated with button elements.

the NewFlexButton factory function should be used to create an instance of FlexButton

The main additional functionalities over the standard widget.Button object are the following:
- image and text management: it is possible to define one or more image, and dynamic text via binding.String
- dynamic layout: it is possible to define if text must be visible, and how to layout image and text.
- sizable: it is possible to define image size, text will be displayed or not depending on container size
- encapsulates state: the object can hold an int state value, based on binding.Int object, depending on which it will display different image
- can be hidden or disabled via binding.Bool objects
- when binding.String tool tip is defined, it will push its text to the binding object when mouse is over.
If not, the text will be displayed on a tooltip popup
- it can include a side image when the button triggers a sub-menu

It is the basic object for the MainRibbon widget
*/
type FlexButton struct {
	widget.DisableableWidget

	mSize              fyne.Size
	imageMaxMinSize    fyne.Size
	mState             int
	textIconHorizontal bool
	compressText       bool
	hasMoreIcon        bool
	isMoreIconBelow    bool
	mTextSize          float32
	mTextString        string

	mTextLabel  *SizableLabel
	mPrimImage  []*canvas.Image
	mSideImage  *canvas.Image
	mBackground *canvas.Rectangle
	mContainer  *fyne.Container

	OnTapped func(int)

	tapAnim *fyne.Animation
	tapBG   *canvas.Rectangle

	mCanvas     fyne.Canvas
	mPopUp      *widget.PopUp
	mPupLbl     *SizableLabel
	mPopUpTimer *time.Ticker
	mRelPos     fyne.Position

	Texter     binding.String
	Disabler   binding.Bool
	Hider      binding.Bool
	Stater     binding.Int
	ToolTipper binding.String
}

/*
NewFlexButton is the factory function for FlexButton object

it requires the following inputs:
- text: The button text.
- images: An array of images to be displayed on the button. Each image is specified as an object that satisfies the fyne.Resource interface.
- isTextAndIconLaidHorizontal: A boolean to determine whether the text and the image are laid out horizontally. If true, the text and image are placed side by side; otherwise, they are stacked vertically.
- compressText: A boolean to determine whether the text must be displayed. If true and FlexButton container size is not big enough, text is omitted. This implies that MinSize() of the FlexButton is the size of the image only
- isImagePadded: A boolean to determine whether the image is padded. If true, padding is added around the image within the button's bounds.
- hasMoreIcon: A boolean to determine whether the button has a "more" icon. If true, a "more" icon is displayed side or below the button to indicate the presence of additional functionality or content.
- isMoreIconBelow: A boolean to determine the placement of the "more" icon. If true, the "more" icon is placed below the button content.
- fullHeight: The full height of the button. For text and image laid out vertically with compressText=false, this is the image height + text height. On other cases it is the image height only
- textHeight: The height of the button text. The combination of fullHeight, textHeight and compressText define the FlexButton MinSize() and layout
- mCanvas: The fyne.Canvas where the button is rendered.
- onTapped: The function to be called when the button is tapped. OnTapped requires an int as input, which is the state of the button when OnTapped is called
- texter: A binding that determines the text displayed on the button.
- disabler: A binding that determines whether the button is disabled.
- hider: A binding that determines whether the button is visible.
- stater: A binding that determines the state of the button.
- toolTipper: A binding that determines the tooltip text for the button.
Returns:
- A pointer to the created FlexButton instance.
*/
func NewFlexButton(text string, images []fyne.Resource, isTextAndIconLaidHorizontal bool, compressText bool, isImagePadded bool, hasMoreIcon bool, isMoreIconBelow bool, fullHeight float32, textHeight float32, mCanvas fyne.Canvas, onTapped func(int), texter binding.String, disabler binding.Bool, hider binding.Bool, stater binding.Int, toolTipper binding.String) *FlexButton {
	t := &FlexButton{
		mBackground:        canvas.NewRectangle(theme.ButtonColor()),
		tapBG:              canvas.NewRectangle(color.Transparent),
		OnTapped:           onTapped,
		mSize:              fyne.NewSize(0., fullHeight),
		mCanvas:            mCanvas,
		textIconHorizontal: isTextAndIconLaidHorizontal,
		compressText:       compressText,
		hasMoreIcon:        hasMoreIcon,
		isMoreIconBelow:    isMoreIconBelow,
		mTextSize:          textHeight,
		mTextString:        text,

		Texter:     texter,
		Disabler:   disabler,
		Hider:      hider,
		Stater:     stater,
		ToolTipper: toolTipper,
	}
	t.ExtendBaseWidget(t)

	t.tapAnim = newButtonTapAnimation(t.tapBG, t)
	t.tapAnim.Curve = fyne.AnimationEaseOut

	imMinSize := fyne.NewSize(0., 0.)
	imMinSize.Height = fullHeight
	if text != "" || texter != nil {
		if !compressText {
			if !isTextAndIconLaidHorizontal {
				imMinSize.Height = fullHeight - textHeight
			}
		}
	}
	if hasMoreIcon {
		if isMoreIconBelow {
			imMinSize.Height -= textHeight
		}
	}

	for i, o := range images {
		t.mPrimImage = append(t.mPrimImage, canvas.NewImageFromResource(o))
		t.mPrimImage[i].FillMode = canvas.ImageFillContain
		t.mPrimImage[i].ScaleMode = canvas.ImageScaleSmooth

		imMinSize.Width = imMinSize.Height * t.mPrimImage[i].Aspect()
		if t.mSize.Width < imMinSize.Width {
			t.mSize.Width = imMinSize.Width
		}
		t.mPrimImage[i].SetMinSize(imMinSize)
		if imMinSize.Height <= 0 {
			t.mPrimImage[i].Hide()
		}
	}
	t.imageMaxMinSize = t.mSize

	if text != "" || texter != nil {
		t.mTextLabel = NewSizableLabel(text, t.mTextSize, true, true, theme.ForegroundColor(), color.Transparent)
		txtSz := t.mTextLabel.MinSize()

		if isTextAndIconLaidHorizontal {
			if !compressText {
				t.mSize.Width += txtSz.Width
			}
		} else {
			if !compressText {
				if txtSz.Width > t.mSize.Width {
					t.mSize.Width = txtSz.Width
				}
			}
		}
	}

	var mItm *fyne.Container
	var innerLayout, outerLayout fyne.Layout
	if isImagePadded {
		innerLayout = &StackFixedRatioPadded{} //&StackPadded{}
	} else {
		innerLayout = &StackFixedRatioUnpadded{} //layout.NewStackLayout()
	}
	if isTextAndIconLaidHorizontal {
		if compressText {
			outerLayout = &GiveUpLastExpandingUnpaddedHBox{}
		} else {
			outerLayout = &ExpandingLastUnpaddedHBox{}
		}
	} else {
		if compressText {
			outerLayout = &GiveUpLastExpandingUnpaddedVBox{}
		} else {
			outerLayout = &ExpandingFirstUnpaddedVBox{}
		}
	}
	if t.mTextLabel != nil && t.mPrimImage != nil {
		innerContainer := container.New(innerLayout)
		for _, o := range t.mPrimImage {
			innerContainer.Add(o)
		}

		mItm = container.New(
			outerLayout,
			innerContainer,
			//container.New(&StackFixedRatioUnpadded{}, innerContainer),
			t.mTextLabel,
		)
	} else if t.mTextLabel != nil {
		mItm = container.NewStack(t.mTextLabel)
	} else if t.mPrimImage != nil {
		innerContainer := container.New(innerLayout)
		for _, o := range t.mPrimImage {
			innerContainer.Add(o)
		}

		mItm = container.New(&StackFixedRatioUnpadded{}, innerContainer)
	} else {
		panic(fmt.Errorf("FlexButton cannot have both text and resource empty"))

	}

	var mItm2 *fyne.Container
	if hasMoreIcon {
		if isTextAndIconLaidHorizontal {
			t.mSideImage = canvas.NewImageFromResource(theme.MenuExpandIcon())
		} else {
			t.mSideImage = canvas.NewImageFromResource(theme.MenuDropDownIcon())
		}
		t.mSideImage.FillMode = canvas.ImageFillOriginal
		t.mSideImage.ScaleMode = canvas.ImageScaleSmooth
		t.mSideImage.SetMinSize(fyne.NewSize(t.mTextSize*t.mSideImage.Aspect(), t.mTextSize))
		if isMoreIconBelow {
			if t.mTextSize*t.mSideImage.Aspect() > t.mSize.Width {
				t.mSize.Width = t.mTextSize * t.mSideImage.Aspect()
			}
			mItm2 = container.New(&ExpandingFirstUnpaddedVBox{}, mItm, container.New(&StackFixedRatioUnpadded{}, t.mSideImage))
		} else {
			t.mSize.Width += t.mTextSize * t.mSideImage.Aspect()
			mItm2 = container.New(&ExpandingFirstUnpaddedHBox{}, mItm, container.New(&StackFixedRatioUnpadded{}, t.mSideImage))
		}
	} else {
		mItm2 = mItm
	}
	t.mContainer = container.NewStack(t.mBackground, t.tapBG, mItm2)

	if t.ToolTipper == nil && mCanvas != nil {
		if text != "" || texter != nil {
			t.mPupLbl = NewSizableLabel(text, 20., false, false, theme.ForegroundColor(), color.Transparent)
			t.mPopUp = widget.NewPopUp(t.mPupLbl, mCanvas)
			t.mPopUp.Hide()
			t.mPopUpTimer = time.NewTicker(2 * time.Second)
			t.mPopUpTimer.Stop()
			go func() {
				for {
					<-t.mPopUpTimer.C
					nPos := fyne.NewPos(t.mRelPos.X, t.mRelPos.Y)
					nPos.Y -= t.mPopUp.MinSize().Height
					t.mPopUp.ShowAtRelativePosition(nPos, t)
					t.mPopUpTimer.Stop()
				}
			}()
		}
	}

	if texter != nil {
		texter.AddListener(t)
		t.mTextString, _ = texter.Get()
	}
	if disabler != nil {
		disabler.AddListener(t)
	}
	if hider != nil {
		hider.AddListener(t)
	}
	if stater != nil {
		stater.AddListener(t)
	}

	return t
}

func (t *FlexButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.mContainer)
}

func (t *FlexButton) MinSize() fyne.Size {
	return t.mSize
}

func (t *FlexButton) Refresh() {
	if t.mSideImage != nil {
		t.mSideImage.Refresh()
	}

	if !t.Disabled() {
		t.mBackground.FillColor = theme.ButtonColor()
	} else {
		t.mBackground.FillColor = theme.DisabledColor()
	}
	t.mBackground.Refresh()

	if t.Stater != nil {
		for i, o := range t.mPrimImage {
			if i == t.mState {
				o.Show()
			} else {
				o.Hide()
			}
		}
	}

	if t.mTextLabel != nil {
		//t.mTextLabel.mBackgroundColor = theme.ButtonColor()
		t.mTextLabel.mTextColor = theme.ForegroundColor()
		t.mTextLabel.Refresh()
	}
	if t.mPopUp != nil {
		t.mPupLbl.mTextColor = theme.ForegroundColor()
		t.mPopUp.Refresh()
	}

}

func (t *FlexButton) Tapped(*fyne.PointEvent) {
	if !t.Disabled() {
		if t.OnTapped != nil {
			t.OnTapped(t.mState)
			if t.mPopUpTimer != nil {
				t.mPopUpTimer.Stop()
			}
		}
		t.tapAnim.Stop()

		if fyne.CurrentApp().Settings().ShowAnimations() {
			t.tapAnim.Start()
		}
	}
}

func (t *FlexButton) MouseIn(*desktop.MouseEvent) {
	if t.Disabled() {
		t.mBackground.FillColor = theme.DisabledColor()
	} else {
		t.mBackground.FillColor = blendColor(theme.ButtonColor(), theme.HoverColor())
	}
	t.mBackground.Refresh()

	if t.mTextLabel != nil {
		t.mTextLabel.mBackgroundColor = t.mBackground.FillColor
		t.mTextLabel.Refresh()
	}

	if t.ToolTipper != nil {
		t.ToolTipper.Set(t.mTextString)
	}
}

func (t *FlexButton) MouseMoved(me *desktop.MouseEvent) {
	if t.mPopUp != nil {
		t.mRelPos = me.Position
		if !t.mTextLabel.Visible() && !t.mPopUp.Visible() {
			t.mPopUpTimer.Reset(2 * time.Second)
		}
	}
}

func (t *FlexButton) MouseOut() {
	if t.Disabled() {
		t.mBackground.FillColor = theme.DisabledColor()
	} else {
		t.mBackground.FillColor = theme.ButtonColor()
	}
	t.mBackground.Refresh()

	if t.mTextLabel != nil {
		t.mTextLabel.mBackgroundColor = t.mBackground.FillColor
		t.mTextLabel.Refresh()
	}

	if t.mPopUp != nil {
		t.mPopUpTimer.Stop()
		t.mPopUp.Hide()
	}

	if t.ToolTipper != nil {
		t.ToolTipper.Set("")
	}
}

func (t *FlexButton) SetMinSize(size fyne.Size) {
	for _, o := range t.mPrimImage {
		o.SetMinSize(size)
	}
}

func newButtonTapAnimation(bg *canvas.Rectangle, w fyne.Widget) *fyne.Animation {
	return fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
		mid := w.Size().Width / 2
		size := mid * done
		bg.Resize(fyne.NewSize(size*2, w.Size().Height))
		bg.Move(fyne.NewPos(mid-size, 0))

		r, g, bb, a := ToNRGBA(theme.PressedColor())
		aa := uint8(a)
		fade := aa - uint8(float32(aa)*done)
		if fade > 0 {
			bg.FillColor = &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(bb), A: fade}
		} else {
			bg.FillColor = color.Transparent
		}
		canvas.Refresh(bg)
	})
}

func (t *FlexButton) DataChanged() {
	if t.Texter != nil {
		if tx, err := t.Texter.Get(); err == nil {
			t.mTextString = tx

			if t.mPupLbl != nil {
				t.mPupLbl.mText.Text = tx
				t.mPupLbl.Refresh()
			}

			if t.mTextLabel != nil {
				if t.mTextLabel.mText.Text != tx {
					t.mTextLabel.mText.Text = tx
					txtSz := t.mTextLabel.MinSize()

					if t.textIconHorizontal {
						if !t.compressText {
							t.mSize.Width = t.imageMaxMinSize.Width + txtSz.Width
						}
					} else {
						if !t.compressText {
							if txtSz.Width > t.imageMaxMinSize.Width {
								t.mSize.Width = txtSz.Width
							}
						}
					}
					if t.hasMoreIcon {
						if t.mSideImage != nil {
							sideImSize := t.mSideImage.MinSize()
							if t.isMoreIconBelow {
								if sideImSize.Width+theme.Padding()*2. > t.mSize.Width {
									t.mSize.Width = sideImSize.Width + theme.Padding()*2.
								}
							} else {
								t.mSize.Width += sideImSize.Width + theme.Padding()*2.
							}
						}
					}
					t.Refresh()
				}
			}
		}
	}

	if t.Disabler != nil {
		mb, err := t.Disabler.Get()
		if err == nil {
			if mb {
				t.Disable()
			} else {
				t.Enable()
			}
			t.Refresh()
		}
	}

	if t.Hider != nil {
		mb, err := t.Hider.Get()
		if err == nil {
			if mb {
				t.Hide()
			} else {
				t.Show()
			}
			t.Refresh()
		}
	}

	if t.Stater != nil {
		mi, err := t.Stater.Get()
		if err == nil {
			t.mState = mi
			t.Refresh()
		}
	}
}
