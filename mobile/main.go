package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"todo-go/mobile/api"
	"todo-go/mobile/auth"
	"todo-go/mobile/screens"
)

const defaultAPIURL = "https://todo-go-app-kbsx.onrender.com/"

type AppScreen int

const (
	ScreenLogin AppScreen = iota
	ScreenRegister
	ScreenHome
)

func main() {

	go func() {

		window := new(app.Window)

		window.Option(
			app.Title("Go Todo"),
			app.Size(
				unit.Dp(400),
				unit.Dp(800),
			),
		)

		// ------------------------------------------------------
		// API URL
		// ------------------------------------------------------

		apiURL := os.Getenv("TODO_API_URL")

		if apiURL == "" {
			apiURL = defaultAPIURL
		}

		client := api.NewClient(apiURL)

		// ------------------------------------------------------
		// THEME
		// ------------------------------------------------------

		theme := material.NewTheme()

		// ------------------------------------------------------
		// SCREENS
		// ------------------------------------------------------

		login := screens.NewLogin(
			theme,
			client,
		)

		register := screens.NewRegister(
			theme,
			client,
		)

		home := screens.NewHome(
			client,
			window,
		)

		currentScreen := ScreenLogin

		// ------------------------------------------------------
		// LOGIN → HOME
		// ------------------------------------------------------

		login.OnLoginSuccess = func() {

			home.LoadTodos()

			currentScreen = ScreenHome

			window.Invalidate()
		}

		// ------------------------------------------------------
		// LOGIN → REGISTER
		// ------------------------------------------------------

		login.OnRegisterRequest = func() {

			currentScreen = ScreenRegister

			window.Invalidate()
		}

		// ------------------------------------------------------
		// REGISTER → HOME
		// ------------------------------------------------------

		register.OnRegisterSuccess = func() {

			home.LoadTodos()

			currentScreen = ScreenHome

			window.Invalidate()
		}

		// ------------------------------------------------------
		// REGISTER → LOGIN
		// ------------------------------------------------------

		register.OnLoginRequest = func() {

			currentScreen = ScreenLogin

			window.Invalidate()
		}

		// ------------------------------------------------------
		// HOME → LOGOUT → LOGIN
		// ------------------------------------------------------

		home.OnLogout = func() {

			currentScreen = ScreenLogin

			window.Invalidate()
		}

		// ------------------------------------------------------
		// INITIAL SCREEN
		// ------------------------------------------------------

		if auth.IsLoggedIn() {

			home.LoadTodos()

			currentScreen = ScreenHome
		}

		// ------------------------------------------------------
		// EVENT LOOP
		// ------------------------------------------------------

		var ops op.Ops

		for {

			event := window.Event()

			switch e := event.(type) {

			case app.DestroyEvent:

				if e.Err != nil {
					log.Println(e.Err)
				}

				return

			case app.FrameEvent:

				gtx := app.NewContext(
					&ops,
					e,
				)

				switch currentScreen {

				case ScreenLogin:

					login.HandleEvents(gtx)
					login.Layout(gtx)

				case ScreenRegister:

					register.HandleEvents(gtx)
					register.Layout(gtx)

				case ScreenHome:

					home.Layout(gtx)
				}

				e.Frame(gtx.Ops)
			}
		}

	}()

	app.Main()
}
