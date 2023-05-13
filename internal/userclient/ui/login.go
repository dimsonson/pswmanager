package ui

import (
	"context"
	"os"
	"syscall"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/rs/zerolog"
)

type UsersServicesProvider interface {
	CreateUser(ctx context.Context) error
	ReadUser(ctx context.Context) (config.UserConfig, error)
	CheckUser(ctx context.Context, login string, passwHex string) error
}

type LoginServicesProvider interface {
	ProcessingLogin(ctx context.Context, record models.LoginRecord) error
	SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error)
}

type CardServicesProvider interface {
	ProcessingCard(ctx context.Context, record models.CardRecord) error
	SearchCard(ctx context.Context, searchInput string) ([]models.CardRecord, error)
}

type BinaryServicesProvider interface {
	ProcessingBinary(ctx context.Context, record models.BinaryRecord) error
	SearchBinary(ctx context.Context, searchInput string) ([]models.BinaryRecord, error)
}


func (ui *UI) FlexLogin() {
	ui.flexLogin = ui.NewCustomFlex(ui.listLogin, 10)
	ui.flexLog = ui.NewCustomFlex(ui.loginform, 10)
	ui.flexReg = ui.NewCustomFlex(ui.regform, 10)
}

func (ui *UI) ListLogin() {
	ui.listLogin = tview.NewList().
		AddItem("Login", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.loginFrm()
			ui.pages.SwitchToPage(LoginForm)
		}).
		AddItem("Registration", "", 'b', func() {
			ui.regform.Clear(true)
			ui.registerFrm()
			if ui.cfg.UserLogin != "" {
				ui.ShowOk("App already registered for user. Please login to App.", func() {
					ui.pages.SwitchToPage(LoginPage)
				})
			}
			if ui.cfg.UserLogin == "" {
				ui.pages.SwitchToPage(RegisterForm)
			}
		}).
		AddItem("Quit", "", 'q', func() {
			log.Logg = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			ui.MainApp.Stop()
			err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			if err != nil {
				log.Print("stop programm error")
				return
			}
		})
	ui.listLogin.SetBorder(true)
	ui.listLogin.SetTitle("Login menu")
	ui.listLogin.SetTitleAlign(tview.AlignLeft)
	ui.listLogin.SetWrapAround(true)
	ui.listLogin.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listLogin)
}

func (ui *UI) loginFrm() *tview.Form {
	tmpUserCfg := config.UserConfig{}
	ui.loginform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
		tmpUserCfg.UserLogin = ulogin
	})
	ui.loginform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
		tmpUserCfg.UserPsw = upsw
	})
	ui.loginform.AddButton("Login", func() {
		err := ui.u.CheckUser(ui.ctx, tmpUserCfg.UserLogin, tmpUserCfg.UserPsw)
		if err != nil {
			ui.ShowConfirm("Wrong password or username", "Do you like try again?",
				func() {
					ui.loginform.SetFocus(0)
					ui.pages.ShowPage(LoginForm)
				},
				func() {
					ui.pages.SwitchToPage(LoginPage)
				})
		}
		if err == nil {
			ui.pages.SwitchToPage(MainPage)
		}
	})
	ui.loginform.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(LoginPage)
	})
	return ui.loginform
}

// + проверяем, зарегистрировано ли приложение на пользователя
// + если да, выводим сообщение
// + если нет, идем дальше к вводу логина и пароля
// + принимаем логин и пароль
// - проверяем пароль на длинну
// - отправляем структуру конфигурации пользователя в сервис создания пользователей
// - если err == nil возвращаем пользователю сообщение об успешной регистрации и переводим в меню входа

func (ui *UI) registerFrm() *tview.Form {
	ui.regform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
		ui.cfg.UserLogin = ulogin
	})
	ui.regform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
		if len(upsw) >= 72 {
			ui.ShowOk("Maximun password leght is 71 character. \nPlease enter shorter password", func() {
				ui.pages.SwitchToPage(RegisterForm)
			})
		}
		ui.cfg.UserPsw = upsw
	})
	ui.regform.AddButton("Register", func() {
		err := ui.u.CreateUser(ui.ctx)
		if err != nil {
			log.Print("registration error:", err)
			ui.ShowOk("Registration error.", func() {
				ui.pages.SwitchToPage(RegisterForm)
			})
		}
		if err == nil {
			ui.ShowOk("Registration successful. \n Please, keep your credentials safe.", func() {
				ui.pages.SwitchToPage(LoginPage)
			})
		}
	})
	ui.regform.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(LoginPage)
	})
	return ui.regform
}
