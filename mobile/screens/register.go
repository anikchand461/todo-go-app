package screens

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"todo-go/mobile/api"
	"todo-go/mobile/auth"
)

type RegisterScreen struct {
	Theme *material.Theme
	API   *api.Client

	Username widget.Editor
	Password widget.Editor
	Confirm  widget.Editor

	RegisterButton widget.Clickable
	LoginButton    widget.Clickable

	Error string

	OnRegisterSuccess func()
	OnLoginRequest    func()
}

func NewRegister(
	theme *material.Theme,
	client *api.Client,
) *RegisterScreen {

	screen := &RegisterScreen{
		Theme: theme,
		API:   client,
	}

	screen.Username.SingleLine = true
	screen.Username.Submit = true
	screen.Username.MaxLen = 100

	screen.Password.SingleLine = true
	screen.Password.Submit = true
	screen.Password.MaxLen = 100
	screen.Password.Mask = '•'

	screen.Confirm.SingleLine = true
	screen.Confirm.Submit = true
	screen.Confirm.MaxLen = 100
	screen.Confirm.Mask = '•'

	return screen
}

func (s *RegisterScreen) Layout(gtx layout.Context) layout.Dimensions {

	paint.Fill(
		gtx.Ops,
		color.NRGBA{
			R: 245,
			G: 247,
			B: 250,
			A: 255,
		},
	)

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		gtx.Constraints.Min.X = gtx.Dp(unit.Dp(340))
		gtx.Constraints.Max.X = gtx.Dp(unit.Dp(380))

		return layout.Inset{
			Top:    unit.Dp(25),
			Bottom: unit.Dp(25),
			Left:   unit.Dp(24),
			Right:  unit.Dp(24),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(
				gtx,

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.H4(
							s.Theme,
							"Create Account",
						).Layout(gtx)
					})
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(28))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(28))
					return layout.Spacer{}.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Username",
						&s.Username,
					)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(12))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(12))
					return layout.Spacer{}.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Password",
						&s.Password,
					)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(12))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(12))
					return layout.Spacer{}.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Confirm password",
						&s.Confirm,
					)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					if s.Error == "" {
						gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(18))
						gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(18))
						return layout.Spacer{}.Layout(gtx)
					}

					return material.Body2(
						s.Theme,
						s.Error,
					).Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					return material.Button(
						s.Theme,
						&s.RegisterButton,
						"Register",
					).Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(15))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(15))
					return layout.Spacer{}.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

						return material.Button(
							s.Theme,
							&s.LoginButton,
							"Back to Login",
						).Layout(gtx)
					})
				}),
			)
		})
	})
}

func (s *RegisterScreen) editor(
	gtx layout.Context,
	label string,
	editor *widget.Editor,
) layout.Dimensions {

	gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(58))
	gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(58))

	return material.Editor(
		s.Theme,
		editor,
		label,
	).Layout(gtx)
}

func (s *RegisterScreen) HandleEvents(gtx layout.Context) {

	if s.RegisterButton.Clicked(gtx) {

		username := strings.TrimSpace(
			s.Username.Text(),
		)

		password := s.Password.Text()
		confirm := s.Confirm.Text()

		if username == "" ||
			password == "" ||
			confirm == "" {

			s.Error = "All fields are required"
			return
		}

		if password != confirm {
			s.Error = "Passwords do not match"
			return
		}

		response, err := s.API.Register(
			username,
			password,
		)

		if err != nil {
			s.Error = err.Error()
			return
		}

		auth.SetToken(response.Token)

		s.Error = ""

		if s.OnRegisterSuccess != nil {
			s.OnRegisterSuccess()
		}
	}

	if s.LoginButton.Clicked(gtx) {

		if s.OnLoginRequest != nil {
			s.OnLoginRequest()
		}
	}
}

// Keep image import used by older Gio versions/build setups.
var _ = image.Point{}
var _ = clip.Rect{}
var _ = paint.FillShape
