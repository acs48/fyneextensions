package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// PaddedBox is a fyne compatible layout. It must have 2 fyne.CanvasObjects and will lay them out
// filling the container size for the first one, and centered with a theme.Padding() border the second one
type PaddedBox struct{}

// MaxMinBox is a fyne compatible layout. It must have 2 fyne.CanvasObjects and will lay them out
// filling the container size for the first one, and centered at its MinSize() the second one
type MaxMinBox struct{}

func (d *PaddedBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	if len(objects) != 2 {
		w, h = 0., 0.
		return fyne.NewSize(w, h)
	}
	w = objects[1].MinSize().Width + theme.Padding()
	h = objects[1].MinSize().Height + theme.Padding()
	return fyne.NewSize(w, h)
}
func (d *PaddedBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	w, h := containerSize.Width-theme.Padding(), containerSize.Height-theme.Padding()
	if len(objects) != 2 {
		for _, o := range objects {
			o.Resize(fyne.NewSize(0., 0.))
			o.Move(fyne.NewPos(0., 0.))
		}
		return
	}
	objects[1].Resize(fyne.NewSize(w, h))
	objects[1].Move(fyne.NewPos(theme.Padding()/2., theme.Padding()/2.))
	objects[0].Resize(containerSize)
	objects[0].Move(fyne.NewPos(0., 0.))
}

func (d *MaxMinBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	if len(objects) != 2 {
		w, h = 0., 0.
		return fyne.NewSize(w, h)
	}
	w = objects[1].MinSize().Width
	h = objects[1].MinSize().Height
	return fyne.NewSize(w, h)
}
func (d *MaxMinBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 2 {
		for _, o := range objects {
			o.Resize(fyne.NewSize(0., 0.))
			o.Move(fyne.NewPos(0., 0.))
		}
		return
	}
	o1Sz := objects[1].MinSize()
	objects[1].Resize(o1Sz)
	objects[1].Move(fyne.NewPos((containerSize.Width-o1Sz.Width)/2., (containerSize.Height-o1Sz.Height)/2.))
	objects[0].Resize(containerSize)
	objects[0].Move(fyne.NewPos(0., 0.))
}
