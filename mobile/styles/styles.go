// Package styles centralizes the app's visual design tokens (colors) and a
// couple of small drawing helpers so every screen shares one consistent
// look. It intentionally contains no business logic.
package styles

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ------------------------------------------------------------------
// PALETTE
// ------------------------------------------------------------------

var (
	// Page background.
	Background = color.NRGBA{R: 244, G: 245, B: 250, A: 255}

	// Card / input surfaces.
	Surface = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	// Brand accent.
	Primary     = color.NRGBA{R: 88, G: 74, B: 234, A: 255}
	PrimaryDark = color.NRGBA{R: 67, G: 55, B: 199, A: 255}
	PrimarySoft = color.NRGBA{R: 231, G: 228, B: 253, A: 255}

	// Text.
	TextPrimary   = color.NRGBA{R: 32, G: 34, B: 46, A: 255}
	TextSecondary = color.NRGBA{R: 121, G: 125, B: 140, A: 255}
	TextMuted     = color.NRGBA{R: 168, G: 171, B: 182, A: 255}
	OnPrimary     = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	// Borders / dividers.
	Border = color.NRGBA{R: 232, G: 234, B: 241, A: 255}

	// Status.
	ErrorText = color.NRGBA{R: 194, G: 51, B: 51, A: 255}
	ErrorSoft = color.NRGBA{R: 252, G: 231, B: 231, A: 255}
	Success   = color.NRGBA{R: 39, G: 163, B: 114, A: 255}

	// Neutral fills.
	Muted = color.NRGBA{R: 243, G: 244, B: 248, A: 255}
)

// ------------------------------------------------------------------
// DRAWING HELPERS
// ------------------------------------------------------------------

// FillRect paints a filled, uniformly-rounded rectangle of the given size
// at the widget's current drawing origin.
func FillRect(gtx layout.Context, size image.Point, c color.NRGBA, radius unit.Dp) {

	r := gtx.Dp(radius)

	rect := clip.UniformRRect(
		image.Rectangle{Max: size},
		r,
	)

	paint.FillShape(
		gtx.Ops,
		c,
		rect.Op(gtx.Ops),
	)
}

// FillMax paints a filled, uniformly-rounded rectangle across the widget's
// current maximum constraints. This is the same "fill then inset content"
// idiom already used throughout the app for cards and input backgrounds.
func FillMax(gtx layout.Context, c color.NRGBA, radius unit.Dp) {
	FillRect(gtx, gtx.Constraints.Max, c, radius)
}
