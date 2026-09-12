package screens

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"todo-go/mobile/api"
	"todo-go/mobile/auth"
	"todo-go/mobile/components"
	"todo-go/mobile/models"
)

type apiResult struct {
	Todos []models.Todo
	Err   error
}

type Home struct {
	Client *api.Client
	Window *app.Window
	Theme  *material.Theme

	Todos []models.Todo

	List widget.List

	Editor widget.Editor

	AddButton    widget.Clickable
	LogoutButton widget.Clickable

	DeleteButtons map[int]*widget.Clickable
	CheckBoxes    map[int]*widget.Bool

	Results chan apiResult

	Loading bool
	Error   string

	OnLogout func()
}

func NewHome(
	client *api.Client,
	window *app.Window,
) *Home {

	home := &Home{
		Client: client,
		Window: window,
		Theme:  material.NewTheme(),

		DeleteButtons: make(map[int]*widget.Clickable),
		CheckBoxes:    make(map[int]*widget.Bool),

		Results: make(chan apiResult, 10),
	}

	home.List.Axis = layout.Vertical

	home.Editor.SingleLine = true
	home.Editor.Submit = true
	home.Editor.MaxLen = 255

	return home
}

// ============================================================
// API
// ============================================================

func (h *Home) LoadTodos() {

	if h.Loading {
		return
	}

	h.Loading = true
	h.Error = ""

	go func() {

		todos, err := h.Client.GetTodos()

		h.Results <- apiResult{
			Todos: todos,
			Err:   err,
		}

		h.Window.Invalidate()
	}()
}

func (h *Home) processResults() {

	for {
		select {

		case result := <-h.Results:

			h.Loading = false

			if result.Err != nil {
				h.Error = result.Err.Error()
				continue
			}

			h.Error = ""
			h.Todos = result.Todos

			for _, todo := range h.Todos {

				if _, exists := h.CheckBoxes[todo.ID]; !exists {
					h.CheckBoxes[todo.ID] = &widget.Bool{}
				}

				h.CheckBoxes[todo.ID].Value = todo.Completed

				if _, exists := h.DeleteButtons[todo.ID]; !exists {
					h.DeleteButtons[todo.ID] = &widget.Clickable{}
				}
			}

		default:
			return
		}
	}
}

// ============================================================
// CREATE
// ============================================================

func (h *Home) createTodo() {

	title := strings.TrimSpace(h.Editor.Text())

	if title == "" || h.Loading {
		return
	}

	h.Loading = true
	h.Error = ""

	h.Editor.SetText("")

	go func() {

		_, err := h.Client.CreateTodo(title)

		if err != nil {

			h.Results <- apiResult{
				Err: err,
			}

			h.Window.Invalidate()
			return
		}

		todos, err := h.Client.GetTodos()

		h.Results <- apiResult{
			Todos: todos,
			Err:   err,
		}

		h.Window.Invalidate()
	}()
}

// ============================================================
// UPDATE
// ============================================================

func (h *Home) toggleTodo(
	id int,
	completed bool,
) {

	if h.Loading {
		return
	}

	h.Loading = true
	h.Error = ""

	go func() {

		_, err := h.Client.UpdateTodo(id, completed)

		if err != nil {

			h.Results <- apiResult{
				Err: err,
			}

			h.Window.Invalidate()
			return
		}

		todos, err := h.Client.GetTodos()

		h.Results <- apiResult{
			Todos: todos,
			Err:   err,
		}

		h.Window.Invalidate()
	}()
}

// ============================================================
// DELETE
// ============================================================

func (h *Home) deleteTodo(id int) {

	if h.Loading {
		return
	}

	h.Loading = true
	h.Error = ""

	go func() {

		err := h.Client.DeleteTodo(id)

		if err != nil {

			h.Results <- apiResult{
				Err: err,
			}

			h.Window.Invalidate()
			return
		}

		todos, err := h.Client.GetTodos()

		h.Results <- apiResult{
			Todos: todos,
			Err:   err,
		}

		h.Window.Invalidate()
	}()
}

// ============================================================
// LOGOUT
// ============================================================

func (h *Home) logout() {

	auth.ClearToken()

	h.Todos = nil
	h.Editor.SetText("")
	h.Error = ""

	h.List.Position.First = 0
	h.List.Position.Offset = 0

	if h.OnLogout != nil {
		h.OnLogout()
	}
}

// ============================================================
// HEADER
// ============================================================

func (h *Home) header(
	gtx layout.Context,
) layout.Dimensions {

	gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(72))
	gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(72))

	return layout.Inset{
		Top:    unit.Dp(8),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(24),
		Right:  unit.Dp(24),
	}.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {

			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(
				gtx,

				// TITLE

				layout.Flexed(
					1,
					func(gtx layout.Context) layout.Dimensions {

						return material.H6(
							h.Theme,
							"My Tasks",
						).Layout(gtx)
					},
				),

				// COUNT

				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {

						return layout.Inset{
							Right: unit.Dp(10),
						}.Layout(
							gtx,
							material.Body2(
								h.Theme,
								h.statsText(),
							).Layout,
						)
					},
				),

				// LOGOUT

				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {

						gtx.Constraints.Min.Y =
							gtx.Dp(unit.Dp(44))

						gtx.Constraints.Max.Y =
							gtx.Dp(unit.Dp(44))

						if h.LogoutButton.Clicked(gtx) {
							h.logout()
						}

						return material.Button(
							h.Theme,
							&h.LogoutButton,
							"Logout",
						).Layout(gtx)
					},
				),
			)
		},
	)
}

// ============================================================
// STATS
// ============================================================

func (h *Home) statsText() string {

	total := len(h.Todos)

	if total == 1 {
		return "1 task"
	}

	return strconv.Itoa(total) + " tasks"
}

// ============================================================
// ADD AREA
// ============================================================

func (h *Home) inputArea(
	gtx layout.Context,
) layout.Dimensions {

	gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(88))
	gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(88))

	return layout.Inset{
		Left:  unit.Dp(24),
		Right: unit.Dp(24),
	}.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {

			radius := gtx.Dp(unit.Dp(16))

			rect := clip.UniformRRect(
				imageRect(gtx),
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
				rect.Op(gtx.Ops),
			)

			return layout.Inset{
				Top:    unit.Dp(10),
				Bottom: unit.Dp(10),
				Left:   unit.Dp(14),
				Right:  unit.Dp(14),
			}.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {

					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(
						gtx,

						// EDITOR

						layout.Flexed(
							1,
							func(gtx layout.Context) layout.Dimensions {

								gtx.Constraints.Min.Y =
									gtx.Dp(unit.Dp(50))

								gtx.Constraints.Max.Y =
									gtx.Dp(unit.Dp(50))

								return material.Editor(
									h.Theme,
									&h.Editor,
									"What needs to be done?",
								).Layout(gtx)
							},
						),

						// ADD BUTTON

						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {

								return layout.Inset{
									Left: unit.Dp(12),
								}.Layout(
									gtx,
									func(gtx layout.Context) layout.Dimensions {

										gtx.Constraints.Min.Y =
											gtx.Dp(unit.Dp(48))

										gtx.Constraints.Max.Y =
											gtx.Dp(unit.Dp(48))

										if h.AddButton.Clicked(gtx) {
											h.createTodo()
										}

										return material.Button(
											h.Theme,
											&h.AddButton,
											"Add",
										).Layout(gtx)
									},
								)
							},
						),
					)
				},
			)
		},
	)
}

// ============================================================
// TODO LIST
// ============================================================

func (h *Home) todoList(
	gtx layout.Context,
) layout.Dimensions {

	if h.Loading && len(h.Todos) == 0 {

		return layout.Center.Layout(
			gtx,
			material.Body1(
				h.Theme,
				"Loading your tasks...",
			).Layout,
		)
	}

	if h.Error != "" && len(h.Todos) == 0 {

		return layout.Center.Layout(
			gtx,
			material.Body1(
				h.Theme,
				"Unable to load tasks",
			).Layout,
		)
	}

	if len(h.Todos) == 0 {

		return layout.Center.Layout(
			gtx,
			func(gtx layout.Context) layout.Dimensions {

				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(
					gtx,

					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {

							return material.H6(
								h.Theme,
								"You're all caught up!",
							).Layout(gtx)
						},
					),

					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {

							return layout.Inset{
								Top: unit.Dp(8),
							}.Layout(
								gtx,
								material.Body2(
									h.Theme,
									"Add a task above to get started.",
								).Layout,
							)
						},
					),
				)
			},
		)
	}

	return material.List(
		h.Theme,
		&h.List,
	).Layout(
		gtx,
		len(h.Todos),

		func(
			gtx layout.Context,
			index int,
		) layout.Dimensions {

			todo := h.Todos[index]

			check := h.CheckBoxes[todo.ID]
			deleteButton := h.DeleteButtons[todo.ID]

			// 76dp card + 4dp top + 4dp bottom.
			itemHeight := gtx.Dp(unit.Dp(84))

			gtx.Constraints.Min.Y = itemHeight
			gtx.Constraints.Max.Y = itemHeight

			if deleteButton.Clicked(gtx) && !h.Loading {
				h.deleteTodo(todo.ID)
			}

			before := check.Value

			card := components.NewTodoCard(
				todo,
				check,
				deleteButton,
			)

			// Same left/right alignment as input box.
			dimensions := layout.Inset{
				Top:    unit.Dp(4),
				Bottom: unit.Dp(4),
				Left:   unit.Dp(24),
				Right:  unit.Dp(4),
			}.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {

					return card.Layout(
						gtx,
						h.Theme,
					)
				},
			)

			after := check.Value

			if before != after && !h.Loading {

				h.toggleTodo(
					todo.ID,
					after,
				)
			}

			return dimensions
		},
	)
}

// ============================================================
// MAIN LAYOUT
// ============================================================

func (h *Home) Layout(
	gtx layout.Context,
) layout.Dimensions {

	h.processResults()

	paint.Fill(
		gtx.Ops,
		color.NRGBA{
			R: 246,
			G: 247,
			B: 249,
			A: 255,
		},
	)

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,

		layout.Rigid(
			h.header,
		),

		layout.Rigid(
			h.inputArea,
		),

		layout.Flexed(
			1,
			h.todoList,
		),
	)
}

// ============================================================
// HELPER
// ============================================================

func imageRect(gtx layout.Context) image.Rectangle {

	return image.Rectangle{
		Max: gtx.Constraints.Max,
	}
}
