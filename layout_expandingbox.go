package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type ExpandingLastUnpaddedVBox struct {
}
type ExpandingFirstUnpaddedVBox struct {
}
type ExpandingLastPaddedVBox struct {
}
type ExpandingFirstPaddedVBox struct {
}

type ExpandingFirstUnpaddedHBox struct {
}
type ExpandingLastUnpaddedHBox struct {
}

type ExpandingFirstPaddedHBox struct {
}

type ExpandingAllProportionallyPaddedHBox struct{}

func (d *ExpandingLastUnpaddedVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
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
func (d *ExpandingLastUnpaddedVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)

	visibleObjects := make([]fyne.CanvasObject, 0)
	for _, o := range objects {
		if o.Visible() {
			visibleObjects = append(visibleObjects, o)
		}
	}
	if len(objects) > 0 {
		if !objects[len(objects)-1].Visible() {
			visibleObjects = append(visibleObjects, objects[len(objects)-1])
		}
	}

	for i, o := range visibleObjects {
		if i != len(visibleObjects)-1 {
			size := o.MinSize()
			size.Width = containerSize.Width
			o.Resize(size)
			o.Move(pos)
			pos.Y += size.Height
		} else {
			size := containerSize.Subtract(pos)
			o.Resize(size)
			o.Move(pos)
		}
	}
}

func (d *ExpandingFirstUnpaddedVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
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
func (d *ExpandingFirstUnpaddedVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, containerSize.Height)

	visibleObjects := make([]fyne.CanvasObject, 0)
	if len(objects) > 0 {
		visibleObjects = append(visibleObjects, objects[0])
	}
	for i, o := range objects {
		if i > 0 {
			if o.Visible() {
				visibleObjects = append(visibleObjects, o)
			}
		}
	}

	for i := len(visibleObjects) - 1; i > 0; i-- {
		o := visibleObjects[i]
		size := o.MinSize()
		size.Width = containerSize.Width
		pos = pos.Subtract(fyne.NewPos(0, size.Height))

		o.Resize(size)
		o.Move(pos)
	}
	if len(visibleObjects) > 0 {
		o := visibleObjects[0]
		size := fyne.NewSize(containerSize.Width, pos.Y)
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (d *ExpandingLastPaddedVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	visibleObj := 0
	for _, o := range objects {
		if o.Visible() {
			visibleObj++
			childSize := o.MinSize()
			if childSize.Width > w {
				w = childSize.Width
			}
			h += childSize.Height
		}
	}
	if visibleObj > 0 {
		h += float32(visibleObj-1) * theme.Padding()
	}
	return fyne.NewSize(w, h)
}
func (d *ExpandingLastPaddedVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)

	visibleObjects := make([]fyne.CanvasObject, 0)
	for _, o := range objects {
		if o.Visible() {
			visibleObjects = append(visibleObjects, o)
		}
	}
	if len(objects) > 0 {
		if !objects[len(objects)-1].Visible() {
			visibleObjects = append(visibleObjects, objects[len(objects)-1])
		}
	}

	for i, o := range visibleObjects {
		if i != len(visibleObjects)-1 {
			size := o.MinSize()
			size.Width = containerSize.Width
			o.Resize(size)
			o.Move(pos)
			pos.Y += size.Height + theme.Padding()
		} else {
			size := containerSize.Subtract(pos)
			o.Resize(size)
			o.Move(pos)
		}
	}
}

func (d *ExpandingFirstPaddedVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	visibleObj := 0
	for _, o := range objects {
		if o.Visible() {
			visibleObj++
			childSize := o.MinSize()
			if childSize.Width > w {
				w = childSize.Width
			}
			h += childSize.Height
		}
	}
	if visibleObj > 0 {
		h += float32(visibleObj-1) * theme.Padding()
	}
	return fyne.NewSize(w, h)
}
func (d *ExpandingFirstPaddedVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	visibleObjects := make([]fyne.CanvasObject, 0)
	if len(objects) > 0 {
		visibleObjects = append(visibleObjects, objects[0])
	}
	for i, o := range objects {
		if i > 0 {
			if o.Visible() {
				visibleObjects = append(visibleObjects, o)
			}
		}
	}

	pos := fyne.NewPos(0, containerSize.Height+theme.Padding())
	for i := len(visibleObjects) - 1; i > 0; i-- {
		o := visibleObjects[i]
		size := o.MinSize()
		size.Width = containerSize.Width
		if i == len(visibleObjects)-1 {
			pos.Y -= size.Height
		} else {
			pos.Y -= size.Height + theme.Padding()
		}

		o.Resize(size)
		o.Move(pos)
	}
	if len(visibleObjects) > 0 {
		o := visibleObjects[0]
		size := fyne.NewSize(containerSize.Width, pos.Y /*-theme.Padding()*/)
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (d *ExpandingLastUnpaddedHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
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
func (d *ExpandingLastUnpaddedHBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)

	visibleObjects := make([]fyne.CanvasObject, 0)
	for _, o := range objects {
		if o.Visible() {
			visibleObjects = append(visibleObjects, o)
		}
	}
	if len(objects) > 0 {
		if !objects[len(objects)-1].Visible() {
			visibleObjects = append(visibleObjects, objects[len(objects)-1])
		}
	}

	for i, o := range visibleObjects {
		if i != len(visibleObjects)-1 {
			size := o.MinSize()
			size.Height = containerSize.Height
			o.Resize(size)
			o.Move(pos)

			pos.X += size.Width
		} else {
			size := containerSize.Subtract(pos)

			o.Resize(size)
			o.Move(pos)
		}
	}
}

func (d *ExpandingFirstUnpaddedHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	visibleObj := make([]fyne.CanvasObject, 0)

	for _, o := range objects {
		if o.Visible() {
			visibleObj = append(visibleObj, o)
		}
	}

	w, h := float32(0), float32(0)
	for _, o := range visibleObj {
		childSize := o.MinSize()
		if childSize.Height > h {
			h = childSize.Height
		}
		w += childSize.Width
	}
	return fyne.NewSize(w, h)
}
func (d *ExpandingFirstUnpaddedHBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(containerSize.Width, 0)

	visibleObj := make([]fyne.CanvasObject, 0)
	if len(objects) > 0 {
		visibleObj = append(visibleObj, objects[0])
	}
	for i, o := range objects {
		if i > 0 {
			if o.Visible() {
				visibleObj = append(visibleObj, o)
			}
		}
	}

	for i := len(visibleObj) - 1; i > 0; i-- {
		o := visibleObj[i]
		size := o.MinSize()
		size.Height = containerSize.Height
		pos = pos.Subtract(fyne.NewPos(size.Width, 0))

		o.Resize(size)
		o.Move(pos)
	}
	if len(visibleObj) > 0 {
		o := visibleObj[0]

		size := fyne.NewSize(pos.X, containerSize.Height)
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (d *ExpandingFirstPaddedHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)

	visibleObj := 0
	for _, o := range objects {
		if o.Visible() {
			visibleObj++
			childSize := o.MinSize()
			if childSize.Height > h {
				h = childSize.Height
			}
			w += childSize.Width
		}
	}
	if visibleObj > 0 {
		w += float32(visibleObj-1) * theme.Padding()
	}
	return fyne.NewSize(w, h)
}
func (d *ExpandingFirstPaddedHBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	visibleObj := make([]fyne.CanvasObject, 0)
	if len(objects) > 0 {
		visibleObj = append(visibleObj, objects[0])
	}
	for i, o := range objects {
		if i > 0 {
			if o.Visible() {
				visibleObj = append(visibleObj, o)
			}
		}
	}

	pos := fyne.NewPos(containerSize.Width, 0)
	for i := len(visibleObj) - 1; i > 0; i-- {
		o := visibleObj[i]
		size := o.MinSize()
		size.Height = containerSize.Height
		if i == len(visibleObj)-1 {
			pos.X -= size.Width
		} else {
			pos.X -= size.Width + theme.Padding()
		}

		o.Resize(size)
		o.Move(pos)
	}
	if len(visibleObj) > 0 {
		o := visibleObj[0]

		size := fyne.NewSize(pos.X /*-theme.Padding()*/, containerSize.Height)
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (d *ExpandingAllProportionallyPaddedHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	visibleObj := make([]fyne.CanvasObject, 0)

	for _, o := range objects {
		if o.Visible() {
			visibleObj = append(visibleObj, o)
		}
	}

	w, h := float32(0), float32(0)
	for _, o := range visibleObj {
		childSize := o.MinSize()
		if childSize.Height > h {
			h = childSize.Height
		}
		w += childSize.Width
	}
	if len(visibleObj) > 0 {
		w += float32(len(visibleObj)-1) * theme.Padding()
	}
	return fyne.NewSize(w, h)
}
func (d *ExpandingAllProportionallyPaddedHBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	visibleObj := make([]fyne.CanvasObject, 0)

	for _, o := range objects {
		if o.Visible() {
			visibleObj = append(visibleObj, o)
		}
	}

	wMin := float32(0)
	for _, o := range visibleObj {
		childSize := o.MinSize()
		wMin += childSize.Width
	}
	wMinPad := wMin
	if len(visibleObj) > 0 {
		wMinPad += float32(len(visibleObj)-1) * theme.Padding()
	}

	multiplier := float32(1.)
	if containerSize.Width > wMinPad {
		diff := containerSize.Width - wMinPad
		multiplier = (wMin + diff) / wMin
	}

	pos := fyne.NewPos(0, 0)
	for _, o := range visibleObj {
		o.Move(pos)
		oMs := o.MinSize()
		o.Resize(fyne.NewSize(oMs.Width*multiplier, containerSize.Height))
		pos.X += oMs.Width*multiplier + theme.Padding()
	}
}

func NewExpandingFirstVBox(objects ...fyne.CanvasObject) *fyne.Container {
	a := container.New(&ExpandingFirstUnpaddedVBox{})
	for i := range objects {
		a.Add(objects[i])
	}
	return a
}
func NewExpandingLastVBox(objects ...fyne.CanvasObject) *fyne.Container {
	a := container.New(&ExpandingLastUnpaddedVBox{})
	for i := range objects {
		a.Add(objects[i])
	}
	return a
}
