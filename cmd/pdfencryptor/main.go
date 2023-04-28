package main

import (
	"strings"

	"github.com/dcapitajp/pdfencryptor/pkg/crypto"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func hasPDFExtension(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".pdf")
}

func confForAlgorithm(aes bool, keyLength int) *model.Configuration {
	c := model.NewDefaultConfiguration()
	c.EncryptUsingAES = aes
	c.EncryptKeyLength = keyLength
	c.Cmd = model.ENCRYPT
	return c
}

type App struct {
	*walk.MainWindow
	UserPW     *walk.LineEdit
	UserPWVal  *walk.LineEdit
	OwnerPW    *walk.LineEdit
	OwnerPWVal *walk.LineEdit
}

type AppDialog struct {
	dlg      *walk.Dialog
	acceptPB *walk.PushButton
}

func (app *App) validateUserPassword() bool {
	return app.UserPW.Text() == app.UserPW.Text()
}

func (app *App) validateOwnerPassword() bool {
	return app.OwnerPW.Text() == app.OwnerPWVal.Text()
}

func (app *App) openDialog(title, msg string) (int, error) {
	dialog := new(AppDialog)
	return Dialog{
		AssignTo:      &dialog.dlg,
		Title:         title,
		DefaultButton: &dialog.acceptPB,
		MinSize:       Size{Width: 200, Height: 100},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: msg,
					},
				},
			},
		},
	}.Run(app)
}

func main() {
	app := &App{}
	MainWindow{
		AssignTo: &app.MainWindow,
		Title:    "PDFパスワード保護",
		Size:     Size{Width: 320, Height: 180},
		Layout:   VBox{},

		OnDropFiles: func(files []string) {
			if !app.validateUserPassword() {
				app.openDialog("エラー", "閲覧パスワードが一致しません")
				return
			}

			if !app.validateOwnerPassword() {
				app.openDialog("エラー", "権限パスワードが一致しません")
				return
			}
			if app.UserPW.Text() == "" {
				// 何もしない
				return
			}
			if app.OwnerPW.Text() == "" {
				app.OwnerPW.SetText(app.UserPW.Text())
			}

			conf := confForAlgorithm(true, 256)
			conf.UserPW = app.UserPW.Text()
			conf.OwnerPW = app.OwnerPW.Text()
			conf.Permissions = model.PermissionsAll
			for _, f := range files {
				if !hasPDFExtension(f) {
					continue
				}
				if err := crypto.EncryptInplace(f, conf); err != nil {
					app.openDialog("Error", err.Error())
					return
				}

			}
		},
		Children: []Widget{
			Composite{Layout: Grid{Columns: 3},
				Alignment: AlignHCenterVNear,
				Children: []Widget{
					Label{Text: "閲覧パスワード*"},
					LineEdit{AssignTo: &app.UserPW, ColumnSpan: 2},
					Label{Text: "閲覧パスワード(確認)*"},
					LineEdit{AssignTo: &app.UserPWVal, ColumnSpan: 2},
					Label{Text: "権限パスワード"},
					LineEdit{AssignTo: &app.OwnerPW, ColumnSpan: 2},
					Label{Text: "権限パスワード(確認)"},
					LineEdit{AssignTo: &app.OwnerPWVal, ColumnSpan: 2},
				},
			},
			TextLabel{Text: "権限パスワードを指定しない場合は、閲覧パスワードがセットされます"},
		},
	}.Run()
}
