package errorpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowErrorWindow(a fyne.App, message string) {
	w := a.NewWindow("Error")
	w.Resize(fyne.NewSize(400, 200))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	// Red icon and message
	icon := widget.NewIcon(theme.ErrorIcon())
	label := widget.NewLabelWithStyle("Oops! Something went wrong.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	label.TextStyle.Bold = true

	errorMessage := widget.NewLabel(message)
	errorMessage.Wrapping = fyne.TextWrapWord
	errorMessage.Alignment = fyne.TextAlignCenter

	closeBtn := widget.NewButton("Close", func() {
		a.Quit()
	})

	content := container.NewVBox(
		icon,
		label,
		errorMessage,
		closeBtn,
	)
	content = container.NewCenter(content)
	w.SetContent(content)

	w.Show()
}