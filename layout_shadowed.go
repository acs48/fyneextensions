package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ShadowedLayout struct{}

func (d *ShadowedLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	if len(objects) != 2 {
		w, h = 0., 0.
		return fyne.NewSize(w, h)
	}
	w = objects[1].MinSize().Width + theme.Padding()
	h = objects[1].MinSize().Height + theme.Padding()
	return fyne.NewSize(w, h)
}
func (d *ShadowedLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	w, h := containerSize.Width-theme.Padding(), containerSize.Height-theme.Padding()
	if len(objects) != 2 {
		for _, o := range objects {
			o.Resize(fyne.NewSize(0., 0.))
			o.Move(fyne.NewPos(0., 0.))
		}
		return
	}
	objects[1].Resize(fyne.NewSize(w, h))
	objects[1].Move(fyne.NewPos(0., 0.))
	objects[0].Resize(fyne.NewSize(w, h))
	objects[0].Move(fyne.NewPos(theme.Padding(), theme.Padding()))
}

type ShadowedWidget struct {
	widget.BaseWidget
	mCanvasObj  fyne.CanvasObject
	mBackground *canvas.Rectangle
	mContainer  *fyne.Container
}

func (t *ShadowedWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.mContainer)
}

func (t *ShadowedWidget) MinSize() fyne.Size {
	w, h := t.mCanvasObj.MinSize().Width, t.mCanvasObj.MinSize().Height
	return fyne.NewSize(w+theme.Padding(), h+theme.Padding())
}

func (t *ShadowedWidget) Refresh() {
	t.mBackground.FillColor = theme.ShadowColor()
	t.mBackground.StrokeColor = theme.ForegroundColor()

	t.mBackground.Refresh()
	t.mCanvasObj.Refresh()
}
