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

type LoginScreen struct {
	Theme *material.Theme
	API   *api.Client

	Username widget.Editor
	Password widget.Editor

	LoginButton    widget.Clickable
	RegisterButton widget.Clickable

	Error string

	OnLoginSuccess    func()
	OnRegisterRequest func()
}

func NewLogin(
	theme *material.Theme,
	client *api.Client,
) *LoginScreen {

	screen := &LoginScreen{
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

	return screen
}

func (s *LoginScreen) Layout(gtx layout.Context) layout.Dimensions {

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
			Top:    unit.Dp(30),
			Bottom: unit.Dp(30),
			Left:   unit.Dp(24),
			Right:  unit.Dp(24),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(
				gtx,

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						title := material.H4(
							s.Theme,
							"Go Todo",
						)
						return title.Layout(gtx)
					})
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(12))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(12))
					return layout.Spacer{}.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						title := material.H6(
							s.Theme,
							"Welcome back",
						)
						return title.Layout(gtx)
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
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(14))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(14))
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
					if s.Error == "" {
						gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(20))
						gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(20))
						return layout.Spacer{}.Layout(gtx)
					}

					label := material.Body2(
						s.Theme,
						s.Error,
					)

					return label.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					button := material.Button(
						s.Theme,
						&s.LoginButton,
						"Login",
					)

					return button.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(18))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(18))
					return layout.Spacer{}.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

						button := material.Button(
							s.Theme,
							&s.RegisterButton,
							"Create account",
						)

						return button.Layout(gtx)
					})
				}),
			)
		})
	})
}

func (s *LoginScreen) editor(
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

func (s *LoginScreen) HandleEvents(gtx layout.Context) {

	if s.LoginButton.Clicked(gtx) {

		username := strings.TrimSpace(
			s.Username.Text(),
		)

		password := s.Password.Text()

		if username == "" || password == "" {
			s.Error = "Username and password are required"
			return
		}

		response, err := s.API.Login(
			username,
			password,
		)

		if err != nil {
			s.Error = err.Error()
			return
		}

		auth.SetToken(response.Token)

		s.Error = ""

		if s.OnLoginSuccess != nil {
			s.OnLoginSuccess()
		}
	}

	if s.RegisterButton.Clicked(gtx) {

		if s.OnRegisterRequest != nil {
			s.OnRegisterRequest()
		}
	}
}

// Keep image import used by older Gio versions/build setups.
var _ = image.Point{}
var _ = clip.Rect{}
var _ = paint.FillShape
