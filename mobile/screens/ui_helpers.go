package screens

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"todo-go/mobile/styles"
)

// fillPageBackground fills the whole screen with the app's standard page
// background color.
func fillPageBackground(gtx layout.Context) {
	paint.Fill(gtx.Ops, styles.Background)
}

// spacer returns a rigid vertical spacer of the given height in dp.
func spacer(height float32) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(height))
		gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(height))
		return layout.Spacer{}.Layout(gtx)
	}
}

// errorBanner draws a soft rounded banner with the given error message.
// The height is fixed (mirroring the fixed-height cards used elsewhere in
// the app) so the background can be sized before the label is measured.
func errorBanner(
	gtx layout.Context,
	theme *material.Theme,
	message string,
) layout.Dimensions {

	gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(40))
	gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(40))

	styles.FillMax(gtx, styles.ErrorSoft, unit.Dp(12))

	return layout.Inset{
		Top:    unit.Dp(9),
		Bottom: unit.Dp(9),
		Left:   unit.Dp(12),
		Right:  unit.Dp(12),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Body2(theme, message)
		label.Color = styles.ErrorText
		return label.Layout(gtx)
	})
}
