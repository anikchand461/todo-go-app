package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"todo-go/mobile/models"
)

type TodoCard struct {
	Todo   models.Todo
	Check  *widget.Bool
	Delete *widget.Clickable
}

func NewTodoCard(
	todo models.Todo,
	check *widget.Bool,
	deleteButton *widget.Clickable,
) TodoCard {

	return TodoCard{
		Todo:   todo,
		Check:  check,
		Delete: deleteButton,
	}
}

// ============================================================
// TODO CARD
// ============================================================

func (c TodoCard) Layout(
	gtx layout.Context,
	theme *material.Theme,
) layout.Dimensions {

	// ----------------------------------------------------------
	// CARD
	// ----------------------------------------------------------

	cardHeight := gtx.Dp(unit.Dp(76))

	gtx.Constraints.Min.Y = cardHeight
	gtx.Constraints.Max.Y = cardHeight

	radius := gtx.Dp(unit.Dp(14))

	cardRect := image.Rectangle{
		Max: gtx.Constraints.Max,
	}

	cardShape := clip.UniformRRect(
		cardRect,
		radius,
	)

	paint.FillShape(
		gtx.Ops,
		color.NRGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
		cardShape.Op(gtx.Ops),
	)

	// ----------------------------------------------------------
	// TITLE
	// ----------------------------------------------------------

	title := material.Body1(
		theme,
		c.Todo.Title,
	)

	if c.Todo.Completed {
		title.Color = color.NRGBA{
			R: 145,
			G: 145,
			B: 150,
			A: 255,
		}
	}

	// ----------------------------------------------------------
	// CONTENT
	// ----------------------------------------------------------

	return layout.Inset{
		Left:  unit.Dp(14),
		Right: unit.Dp(14),
	}.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {

			// Give the content a predictable height.
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(76))
			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(76))

			return layout.Center.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {

					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(
						gtx,

						// --------------------------------------------------
						// CHECKBOX
						// --------------------------------------------------

						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {

								gtx.Constraints.Min.X =
									gtx.Dp(unit.Dp(48))

								gtx.Constraints.Max.X =
									gtx.Dp(unit.Dp(48))

								gtx.Constraints.Min.Y =
									gtx.Dp(unit.Dp(48))

								gtx.Constraints.Max.Y =
									gtx.Dp(unit.Dp(48))

								return layout.Center.Layout(
									gtx,
									material.CheckBox(
										theme,
										c.Check,
										"",
									).Layout,
								)
							},
						),

						// --------------------------------------------------
						// TITLE
						// --------------------------------------------------

						layout.Flexed(
							1,
							func(gtx layout.Context) layout.Dimensions {

								return layout.Center.Layout(
									gtx,
									func(gtx layout.Context) layout.Dimensions {

										return layout.Inset{
											Left:  unit.Dp(8),
											Right: unit.Dp(8),
										}.Layout(
											gtx,
											title.Layout,
										)
									},
								)
							},
						),

						// --------------------------------------------------
						// DELETE BUTTON
						// --------------------------------------------------

						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {

								gtx.Constraints.Min.X =
									gtx.Dp(unit.Dp(50))

								gtx.Constraints.Max.X =
									gtx.Dp(unit.Dp(50))

								gtx.Constraints.Min.Y =
									gtx.Dp(unit.Dp(48))

								gtx.Constraints.Max.Y =
									gtx.Dp(unit.Dp(48))

								button := material.Button(
									theme,
									c.Delete,
									"×",
								)

								button.Background = color.NRGBA{
									R: 245,
									G: 245,
									B: 247,
									A: 255,
								}

								button.Color = color.NRGBA{
									R: 100,
									G: 100,
									B: 105,
									A: 255,
								}

								button.CornerRadius = unit.Dp(10)

								return layout.Center.Layout(
									gtx,
									button.Layout,
								)
							},
						),
					)
				},
			)
		},
	)
}
