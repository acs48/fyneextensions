package fyneextensions

import (
	"fyne.io/fyne/v2"
	"image/color"
)

func DarkTheme(fallback fyne.Theme) fyne.Theme {
	return &forcedVariant{Theme: fallback, variant: 0} // avoid import loops
}

func LightTheme(fallback fyne.Theme) fyne.Theme {
	return &forcedVariant{Theme: fallback, variant: 1} // avoid import loops
}

type forcedVariant struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.variant)
}

func blendColor(under, over color.Color) color.Color {
	// This alpha blends with the over operator, and accounts for RGBA() returning alpha-premultiplied values
	dstR, dstG, dstB, dstA := under.RGBA()
	srcR, srcG, srcB, srcA := over.RGBA()

	srcAlpha := float32(srcA) / 0xFFFF
	dstAlpha := float32(dstA) / 0xFFFF

	outAlpha := srcAlpha + dstAlpha*(1-srcAlpha)
	outR := srcR + uint32(float32(dstR)*(1-srcAlpha))
	outG := srcG + uint32(float32(dstG)*(1-srcAlpha))
	outB := srcB + uint32(float32(dstB)*(1-srcAlpha))
	// We create an RGBA64 here because the color components are already alpha-premultiplied 16-bit values (they're just stored in uint32s).
	return color.RGBA64{R: uint16(outR), G: uint16(outG), B: uint16(outB), A: uint16(outAlpha * 0xFFFF)}
}

// ToNRGBA converts a color to RGBA values which are not premultiplied, unlike color.RGBA().
func ToNRGBA(c color.Color) (r, g, b, a int) {
	// We use UnmultiplyAlpha with RGBA, RGBA64, and unrecognized implementations of Color.
	// It works for all Colors whose RGBA() method is implemented according to spec, but is only necessary for those.
	// Only RGBA and RGBA64 have components which are already premultiplied.
	switch col := c.(type) {
	// NRGBA and NRGBA64 are not premultiplied
	case color.NRGBA:
		r = int(col.R)
		g = int(col.G)
		b = int(col.B)
		a = int(col.A)
	case *color.NRGBA:
		r = int(col.R)
		g = int(col.G)
		b = int(col.B)
		a = int(col.A)
	case color.NRGBA64:
		r = int(col.R) >> 8
		g = int(col.G) >> 8
		b = int(col.B) >> 8
		a = int(col.A) >> 8
	case *color.NRGBA64:
		r = int(col.R) >> 8
		g = int(col.G) >> 8
		b = int(col.B) >> 8
		a = int(col.A) >> 8
	// Gray and Gray16 have no alpha component
	case *color.Gray:
		r = int(col.Y)
		g = int(col.Y)
		b = int(col.Y)
		a = 0xff
	case color.Gray:
		r = int(col.Y)
		g = int(col.Y)
		b = int(col.Y)
		a = 0xff
	case *color.Gray16:
		r = int(col.Y) >> 8
		g = int(col.Y) >> 8
		b = int(col.Y) >> 8
		a = 0xff
	case color.Gray16:
		r = int(col.Y) >> 8
		g = int(col.Y) >> 8
		b = int(col.Y) >> 8
		a = 0xff
	// Alpha and Alpha16 contain only an alpha component.
	case color.Alpha:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A)
	case *color.Alpha:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A)
	case color.Alpha16:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A) >> 8
	case *color.Alpha16:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A) >> 8
	default: // RGBA, RGBA64, and unknown implementations of Color
		r, g, b, a = unmultiplyAlpha(c)
	}
	return
}

// unmultiplyAlpha returns a color's RGBA components as 8-bit integers by calling c.RGBA() and then removing the alpha premultiplication.
// It is only used by ToRGBA.
func unmultiplyAlpha(c color.Color) (r, g, b, a int) {
	red, green, blue, alpha := c.RGBA()
	if alpha != 0 && alpha != 0xffff {
		red = (red * 0xffff) / alpha
		green = (green * 0xffff) / alpha
		blue = (blue * 0xffff) / alpha
	}
	// Convert from range 0-65535 to range 0-255
	r = int(red >> 8)
	g = int(green >> 8)
	b = int(blue >> 8)
	a = int(alpha >> 8)
	return
}
