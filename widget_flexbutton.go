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

// SetMinSize fff
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
