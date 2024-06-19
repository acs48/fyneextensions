package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type SizableLabel struct {
	widget.BaseWidget
	mText       *canvas.Text
	mBackground *canvas.Rectangle
	mContainer  *fyne.Container

	mSize            float32
	mTextColor       color.Color
	mBackgroundColor color.Color
	isPadded         bool
}

func NewSizableLabel(text string, verticalSize float32, padded bool, alignCenter bool, textColor color.Color, backgroundColor color.Color) *SizableLabel {
	t := &SizableLabel{
		mText:            canvas.NewText(text, textColor),
		mBackground:      canvas.NewRectangle(backgroundColor),
		mSize:            verticalSize,
		mTextColor:       textColor,
		mBackgroundColor: backgroundColor,
		isPadded:         padded,
	}
	t.ExtendBaseWidget(t)

	if alignCenter {
		t.mContainer = container.New(&MaxMinBox{})
	} else {
		t.mContainer = container.NewStack()
	}

	textSize := verticalSize
	if padded {
		textSize -= theme.Padding() * 2.
	}

	t.mText.TextSize = t.mSize
	t.mText.TextStyle.Bold = true

	fontSz1 := textSize / 2.
	fontSz2 := textSize * 2.
	fontSz3 := textSize

	for i := 0; i < 20; i++ {
		txtSz1 := fyne.MeasureText(text, fontSz1, fyne.TextStyle{Bold: true})
		txtSz2 := fyne.MeasureText(text, fontSz2, fyne.TextStyle{Bold: true})
		h1 := txtSz1.Height
		h2 := txtSz2.Height

		r3 := (textSize - h2) / (h1 - h2)
		fontSz3 = fontSz1 + r3*(fontSz2-fontSz1)
		txtSz3 := fyne.MeasureText(text, fontSz3, fyne.TextStyle{Bold: true})
		h3 := txtSz3.Height
		if h3 > textSize {
			fontSz2 = fontSz3
		} else if h3 < textSize {
			fontSz1 = fontSz3
		} else {
			break
		}
	}

	t.mText.TextSize = fontSz1

	t.mContainer.Add(t.mBackground)
	t.mContainer.Add(t.mText)

	t.Refresh()

	return t
}

func (t *SizableLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.mContainer)
}

func (t *SizableLabel) SetTextSize(size float32) {
	t.mSize = size
	t.Refresh()
}

func (t *SizableLabel) GetFontSize() float32 {
	return t.mText.TextSize
}

func (t *SizableLabel) MinSize() fyne.Size {
	txtSz := fyne.MeasureText(t.mText.Text, t.mText.TextSize, t.mText.TextStyle)
	wSz := txtSz.Width
	if t.isPadded {
		wSz += 4. * theme.Padding()
	} else {
		wSz += 2. * theme.Padding()
	}
	return fyne.NewSize(wSz, t.mSize)
}

func (t *SizableLabel) Refresh() {
	t.mText.Color = t.mTextColor

	//t.mBackground.SetMinSize(t.mTextLabel.MinSize())
	//t.mBackground.Resize(t.mTextLabel.Size())
	t.mBackground.FillColor = t.mBackgroundColor
	if t.mBackgroundColor == color.Transparent {
		t.mBackground.Hide()
	}

	t.mText.Refresh()
	t.mBackground.Refresh()
}
