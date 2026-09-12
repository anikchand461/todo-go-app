package screens

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"todo-go/mobile/api"
	"todo-go/mobile/auth"
	"todo-go/mobile/styles"
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
				Axis: layout.Vertical,
			}.Layout(
				gtx,

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						title := material.H4(s.Theme, "Create Account")
						title.Color = styles.TextPrimary
						return title.Layout(gtx)
					})
				}),

				layout.Rigid(spacer(6)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						subtitle := material.Body2(s.Theme, "Start tracking your tasks in seconds")
						subtitle.Color = styles.TextSecondary
						return subtitle.Layout(gtx)
					})
				}),

				layout.Rigid(spacer(26)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Username",
						&s.Username,
					)
				}),

				layout.Rigid(spacer(12)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Password",
						&s.Password,
					)
				}),

				layout.Rigid(spacer(12)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.editor(
						gtx,
						"Confirm password",
						&s.Confirm,
					)
				}),

				layout.Rigid(spacer(18)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					if s.Error == "" {
						gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(18))
						gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(18))
						return layout.Spacer{}.Layout(gtx)
					}

					return errorBanner(gtx, s.Theme, s.Error)
				}),

				layout.Rigid(spacer(8)),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

					button := material.Button(
						s.Theme,
						&s.RegisterButton,
						"Register",
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
						&s.LoginButton,
						"Back to Login",
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

func (s *RegisterScreen) editor(
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
