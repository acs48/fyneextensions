package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// EquallySpacedUnpaddedVBox is a fyne compatible layout.
// It defines vertical box without padding, where all object have the same height.
// Width is the container width
type EquallySpacedUnpaddedVBox struct{}

// EquallySpacedUnpaddedHBox is a fyne compatible layout.
// It defines horizontal box without padding, where all object have the same width.
// Height is the container height
type EquallySpacedUnpaddedHBox struct{}

// EquallySpacedPaddedVBox is a fyne compatible layout.
// It defines vertical box with padding, where all object have the same height.
// Width is the container width
type EquallySpacedPaddedVBox struct{}

// EquallySpacedPaddedHBox is a fyne compatible layout.
// It defines horizontal box without padding, where all object have the same width.
// Height is the container height
type EquallySpacedPaddedHBox struct{}

func (d *EquallySpacedUnpaddedVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		if o.Visible() {
			childSize := o.MinSize()
			if childSize.Width > w {
				w = childSize.Width
			}
			h += childSize.Height
		}
	}
	return fyne.NewSize(w, h)
}
func (d *EquallySpacedUnpaddedVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)

	numObj := 0
	for _, o := range objects {
		if o.Visible() {
			numObj++
		}
	}

	commonHeight := float32(0.)
	if numObj > 0 {
		commonHeight = containerSize.Height / float32(numObj)
	}

	for _, o := range objects {
		if o.Visible() {
			size := fyne.NewSize(containerSize.Width, commonHeight)
			o.Resize(size)
			o.Move(pos)

			pos.Y += commonHeight
		} else {
			o.Resize(fyne.NewSize(0, 0))
			o.Move(pos)
		}
	}
}

func (d *EquallySpacedUnpaddedHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		if o.Visible() {
			childSize := o.MinSize()
			if childSize.Height > h {
				h = childSize.Height
			}
			w += childSize.Width
		}
	}
	return fyne.NewSize(w, h)
}
func (d *EquallySpacedUnpaddedHBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)

	numObj := 0
	for _, o := range objects {
		if o.Visible() {
			numObj++
		}
	}

	commonWidth := float32(0.)
	if numObj > 0 {
		commonWidth = containerSize.Width / float32(numObj)
	}

	for _, o := range objects {
		if o.Visible() {
			size := fyne.NewSize(commonWidth, containerSize.Height)
			o.Resize(size)
			o.Move(pos)

			pos.X += commonWidth
		} else {
			o.Resize(fyne.NewSize(0, 0))
			o.Move(pos)
		}
	}
}

func (d *EquallySpacedPaddedVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	numChild := 0
	for _, o := range objects {
		if o.Visible() {
			numChild++
			childSize := o.MinSize()
			if childSize.Width > w {
				w = childSize.Width
			}
			h += childSize.Height
		}
	}
	h += theme.Padding() * float32(numChild-1)
	if numChild == 0 {
		h = 0.
	}
	return fyne.NewSize(w, h)
}
func (d *EquallySpacedPaddedVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)

	numObj := 0
	for _, o := range objects {
		if o.Visible() {
			numObj++
		}
	}

	commonHeight := float32(0.)
	if numObj > 0 {
		commonHeight = (containerSize.Height - theme.Padding()*float32(numObj-1)) / float32(numObj)
	}

	for _, o := range objects {
		if o.Visible() {
			size := fyne.NewSize(containerSize.Width, commonHeight)
			o.Resize(size)
			o.Move(pos)

			pos.Y += commonHeight + theme.Padding()
		} else {
			o.Resize(fyne.NewSize(0, 0))
			o.Move(pos)
		}
	}
}

func (d *EquallySpacedPaddedHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	numChild := 0
	for _, o := range objects {
		if o.Visible() {
			numChild++
			childSize := o.MinSize()
			if childSize.Height > h {
				h = childSize.Height
			}
			w += childSize.Width
		}
	}
	w += theme.Padding() * float32(numChild-1)
	if numChild == 0 {
		w = 0.
	}
	return fyne.NewSize(w, h)
}
func (d *EquallySpacedPaddedHBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)

	numObj := 0
	for _, o := range objects {
		if o.Visible() {
			numObj++
		}
	}

	commonWidth := float32(0.)
	if numObj > 0 {
		commonWidth = (containerSize.Width - theme.Padding()*float32(numObj-1)) / float32(numObj)
	}

	for _, o := range objects {
		if o.Visible() {
			size := fyne.NewSize(commonWidth, containerSize.Height)
			o.Resize(size)
			o.Move(pos)

			pos.X += commonWidth + theme.Padding()
		} else {
			o.Resize(fyne.NewSize(0, 0))
			o.Move(pos)
		}
	}
}
