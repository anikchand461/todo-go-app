package screens

import (
	"image"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"todo-go/mobile/api"
	"todo-go/mobile/auth"
	"todo-go/mobile/styles"
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

	fillPageBackground(gtx)

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		gtx.Constraints.Min.X = gtx.Dp(unit.Dp(340))
		gtx.Constraints.Max.X = gtx.Dp(unit.Dp(380))

		return layout.Inset{
			Top:    unit.Dp(24),
			Bottom: unit.Dp(24),
			Left:   unit.Dp(24),
			Right:  unit.Dp(24),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(
				gtx,

				layout.Rigid(s.brandMark),

				layout.Rigid(spacer(16)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						title := material.H4(s.Theme, "Go Todo")
						title.Color = styles.TextPrimary
						return title.Layout(gtx)
					})
				}),

				layout.Rigid(spacer(6)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						subtitle := material.Body2(s.Theme, "Sign in to manage your tasks")
						subtitle.Color = styles.TextSecondary
						return subtitle.Layout(gtx)
					})
				}),

				layout.Rigid(spacer(28)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Username",
						&s.Username,
					)
				}),

				layout.Rigid(spacer(14)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Password",
						&s.Password,
					)
				}),

				layout.Rigid(spacer(18)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if s.Error == "" {
						gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(20))
						gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(20))
						return layout.Spacer{}.Layout(gtx)
					}

					return errorBanner(gtx, s.Theme, s.Error)
				}),

				layout.Rigid(spacer(8)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					button := material.Button(
						s.Theme,
						&s.LoginButton,
						"Login",
					)

					button.Background = styles.Primary
					button.Color = styles.OnPrimary
					button.CornerRadius = unit.Dp(14)

					return button.Layout(gtx)
				}),

				layout.Rigid(spacer(12)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					button := material.Button(
						s.Theme,
						&s.RegisterButton,
						"Create account",
					)

					button.Background = styles.PrimarySoft
					button.Color = styles.Primary
					button.CornerRadius = unit.Dp(14)

					return button.Layout(gtx)
				}),
			)
		})
	})
}

// brandMark draws a small circular accent badge above the title.
func (s *LoginScreen) brandMark(gtx layout.Context) layout.Dimensions {

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		size := gtx.Dp(unit.Dp(64))

		gtx.Constraints.Min = image.Point{X: size, Y: size}
		gtx.Constraints.Max = gtx.Constraints.Min

		styles.FillMax(gtx, styles.Primary, unit.Dp(32))

		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			mark := material.H5(s.Theme, "✓")
			mark.Color = styles.OnPrimary
			return mark.Layout(gtx)
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

	styles.FillMax(gtx, styles.Surface, unit.Dp(14))

	return layout.Inset{
		Left:  unit.Dp(16),
		Right: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Editor(
				s.Theme,
				editor,
				label,
			).Layout(gtx)
		})
	})
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
