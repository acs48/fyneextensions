package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type StackFixedRatioUnpadded struct{}
type StackFixedRatioPadded struct{}
type StackPadded struct{}

func (d *StackFixedRatioUnpadded) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		if o.Visible() {
			if childSize.Width > w {
				w = childSize.Width
			}
			if childSize.Height > h {
				h = childSize.Height
			}
		}
	}
	return fyne.NewSize(w, h)
}
func (d *StackFixedRatioUnpadded) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for _, o := range objects {
		baseSize := o.MinSize()
		rath := containerSize.Height / baseSize.Height
		ratw := containerSize.Width / baseSize.Width
		minrat := rath
		if ratw < rath {
			minrat = ratw
		}
		fullSize := fyne.NewSize(baseSize.Width*minrat, baseSize.Height*minrat)
		o.Resize(fullSize)

		fullPos := fyne.NewPos(0., 0.)
		fullPos.Y = (containerSize.Height - fullSize.Height) / 2.
		fullPos.X = (containerSize.Width - fullSize.Width) / 2.
		o.Move(fullPos)
	}
}

func (d *StackFixedRatioPadded) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		if o.Visible() {
			if childSize.Width > w {
				w = childSize.Width
			}
			if childSize.Height > h {
				h = childSize.Height
			}
		}
	}
	return fyne.NewSize(w+theme.Padding()*2., h+theme.Padding()*2.)
}
func (d *StackFixedRatioPadded) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	unpaddedSize := fyne.NewSize(containerSize.Width-theme.Padding()*2., containerSize.Height-theme.Padding()*2.)

	for _, o := range objects {
		baseSize := o.MinSize()
		rath := unpaddedSize.Height / baseSize.Height
		ratw := unpaddedSize.Width / baseSize.Width
		minrat := rath
		if ratw < rath {
			minrat = ratw
		}
		fullSize := fyne.NewSize(baseSize.Width*minrat, baseSize.Height*minrat)
		o.Resize(fullSize)

		fullPos := fyne.NewPos(0., 0.)
		fullPos.Y = (unpaddedSize.Height-fullSize.Height)/2. + theme.Padding()
		fullPos.X = (unpaddedSize.Width-fullSize.Width)/2. + theme.Padding()
		o.Move(fullPos)
	}
}

func (d *StackPadded) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		if o.Visible() {
			if childSize.Width > w {
				w = childSize.Width
			}
			if childSize.Height > h {
				h = childSize.Height
			}
		}
	}
	return fyne.NewSize(w+theme.Padding()*2., h+theme.Padding()*2.)
}
func (d *StackPadded) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for _, o := range objects {
		nSz := fyne.NewSize(containerSize.Width-theme.Padding()*2., containerSize.Height-theme.Padding()*2.)
		o.Resize(nSz)
		o.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	}
}
