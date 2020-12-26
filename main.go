package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// TopUI represents the main application window.
type TopUI struct {
	window        *widgets.QMainWindow
	intervalLabel *widgets.QLabel
	listBox       *widgets.QListWidget
}

// reset the gui and variables to inital/empty values.
func (t *TopUI) reset() {

	// t.inputTextBox.SetText(t.PDFIn)
	t.window.SetWindowTitle("")
}

func (t *TopUI) updateUI() {

}

// RunApp initializes the GUI and all associated GUI types.
func (t *TopUI) RunApp() {

	ui := widgets.NewQApplication(len(os.Args), os.Args)

	t.window = widgets.NewQMainWindow(nil, 0)
	t.window.SetMinimumSize2(450, 675)
	t.window.SetMaximumSize2(450, 675)
	t.window.SetWindowTitle("Top UI")

	h1 := widgets.NewQHBoxLayout()
	h2 := widgets.NewQHBoxLayout()
	h3 := widgets.NewQHBoxLayout()

	v := widgets.NewQVBoxLayout()

	timer1 := core.NewQTimer(t.window)
	timer1.ConnectTimeout(func() { t.updateUI() })
	timer1.Start(1000)

	title := widgets.NewQGraphicsScene(t.window)
	title.AddText("t e m p", gui.NewQFont2("Menlo", 15, 1, false))
	titleView := widgets.NewQGraphicsView(t.window)
	titleView.SetScene(title)

	t.listBox = widgets.NewQListWidget(t.window)
	t.listBox.SetFixedHeight(525)

	t.intervalLabel = widgets.NewQLabel(t.window, 0)
	t.intervalLabel.SetText("g r e e t i n g s")

	intervalButton := widgets.NewQPushButton2("Browse", nil)
	intervalButton.ConnectClicked(func(bool) { t.updateUI() })

	h1.Layout().AddWidget(titleView)
	h2.Layout().AddWidget(t.listBox)
	h3.Layout().AddWidget(intervalButton)
	h3.Layout().AddWidget(t.intervalLabel)

	for _, layout := range []*widgets.QHBoxLayout{h1, h2, h3} {
		v.AddLayout(layout, 0)
	}

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(v)
	t.window.SetCentralWidget(widget)
	t.window.Show()
	ui.Exec()
}

func main() {

	var ui = &TopUI{}
	ui.RunApp()
}
