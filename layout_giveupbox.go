package fyneextensions

import (
	"fmt"
	"fyne.io/fyne/v2"
)

type GiveUpLastExpandingUnpaddedVBox struct {
}
type GiveUpLastExpandingUnpaddedHBox struct {
}

func (d *GiveUpLastExpandingUnpaddedVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) != 2 {
		fmt.Println("GiveUpLayout must have 2 objects")
		return fyne.NewSize(0., 0.)
	}
	return objects[0].MinSize()
}
func (d *GiveUpLastExpandingUnpaddedVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 2 {
		fmt.Println("GiveUpLayout must have 2 objects")
		for _, o := range objects {
			o.Resize(fyne.NewSize(0., 0.))
			o.Move(fyne.NewPos(0., 0.))
		}
		return
	}

	realMinSize := fyne.NewSize(0, 0)
	realMinSize.Height, realMinSize.Width = objects[0].MinSize().Height+objects[1].MinSize().Height, objects[0].MinSize().Width
	if objects[1].MinSize().Width > objects[0].MinSize().Width {
		realMinSize.Width = objects[1].MinSize().Width
	}

	pos := fyne.NewPos(0, containerSize.Height)
	if containerSize.Height < realMinSize.Height || containerSize.Width < realMinSize.Width {
		o := objects[1]
		o.Hide()

		o = objects[0]
		size := containerSize
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	} else {
		o := objects[1]
		size := o.MinSize()
		size.Width = containerSize.Width
		pos.Y -= size.Height

		o.Show()
		o.Resize(size)
		o.Move(pos)

		o = objects[0]

		size = fyne.NewSize(containerSize.Width, pos.Y)
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (d *GiveUpLastExpandingUnpaddedHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) != 2 {
		fmt.Println("GiveUpLayout must have 2 objects")
		return fyne.NewSize(0., 0.)
	}
	return objects[0].MinSize()
}
func (d *GiveUpLastExpandingUnpaddedHBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 2 {
		fmt.Println("GiveUpLayout must have 2 objects")
		for _, o := range objects {
			o.Resize(fyne.NewSize(0., 0.))
			o.Move(fyne.NewPos(0., 0.))
		}
		return
	}

	realMinSize := fyne.NewSize(0, 0)
	realMinSize.Height, realMinSize.Width = objects[0].MinSize().Height, objects[0].MinSize().Width+objects[1].MinSize().Width
	if objects[0].MinSize().Height > realMinSize.Height {
		realMinSize.Height = objects[0].MinSize().Height
	}

	pos := fyne.NewPos(containerSize.Width, 0.)
	if containerSize.Width < realMinSize.Width || containerSize.Height < realMinSize.Height {
		o := objects[1]
		o.Hide()

		o = objects[0]
		size := containerSize
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	} else {
		o := objects[1]
		size := o.MinSize()
		size.Height = containerSize.Height
		pos.X -= size.Width

		o.Show()
		o.Resize(size)
		o.Move(pos)

		o = objects[0]
		size = fyne.NewSize(pos.X, containerSize.Height)
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}
