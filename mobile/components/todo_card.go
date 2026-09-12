package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"todo-go/mobile/models"
	"todo-go/mobile/styles"
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

	styles.FillMax(gtx, styles.Surface, unit.Dp(14))

	// ----------------------------------------------------------
	// TITLE
	// ----------------------------------------------------------

	title := material.Body1(
		theme,
		c.Todo.Title,
	)

	if c.Todo.Completed {
		title.Color = styles.TextMuted
	} else {
		title.Color = styles.TextPrimary
	}

	// ----------------------------------------------------------
	// CONTENT
	// ----------------------------------------------------------

	return layout.Inset{
		Left:  unit.Dp(10),
		Right: unit.Dp(10),
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
						// STATUS DOT + CHECKBOX
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
									gtx.Dp(unit.Dp(44))

								gtx.Constraints.Max.X =
									gtx.Dp(unit.Dp(44))

								gtx.Constraints.Min.Y =
									gtx.Dp(unit.Dp(44))

								gtx.Constraints.Max.Y =
									gtx.Dp(unit.Dp(44))

								button := material.Button(
									theme,
									c.Delete,
									"×",
								)

								button.Background = styles.ErrorSoft
								button.Color = styles.ErrorText
								button.CornerRadius = unit.Dp(22)

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
