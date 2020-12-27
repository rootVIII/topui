package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// TopUI represents the main application window.
type TopUI struct {
	window  *widgets.QMainWindow
	listBox *widgets.QListWidget
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

	heading1 := widgets.NewQLabel(t.window, 0)
	heading1.SetText("PID")
	heading2 := widgets.NewQLabel(t.window, 0)
	heading2.SetText("CPU%")
	heading3 := widgets.NewQLabel(t.window, 0)
	heading3.SetText("MEMORY")
	heading4 := widgets.NewQLabel(t.window, 0)
	heading4.SetText("UID")

	divider := widgets.NewQGraphicsScene(t.window)
	titleView := widgets.NewQGraphicsView(t.window)
	titleView.SetScene(divider)
	titleView.SetFixedHeight(2)

	t.listBox = widgets.NewQListWidget(t.window)
	t.listBox.SetFixedHeight(600)

	h1.Layout().AddWidget(heading1)
	h1.Layout().AddWidget(heading2)
	h1.Layout().AddWidget(heading3)
	h1.Layout().AddWidget(heading4)
	h2.Layout().AddWidget(titleView)
	h3.Layout().AddWidget(t.listBox)

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
