package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// TopUI represents the main application window.
type TopUI struct {
	window    *widgets.QMainWindow
	listBox   *widgets.QListWidget
	cmdBuffer bytes.Buffer
}

func (t *TopUI) updateUI() {
	fmt.Printf("\n\nS T A R T\n\n")
	fmt.Printf("%q\n", t.cmdBuffer.String())
	t.cmdBuffer.Reset()
}

func (t *TopUI) scanSTDOUT(scanner *bufio.Scanner) {
	for scanner.Scan() {
		t.cmdBuffer.Write(scanner.Bytes())
	}
}

func (t TopUI) execTop() (*bufio.Scanner, error) {
	command := exec.Command("top")
	stdout, err := command.StdoutPipe()

	if err != nil {
		return nil, err
	}
	if err := command.Start(); err != nil {
		return nil, err
	}
	return bufio.NewScanner(stdout), nil
}

// RunApp initializes the GUI and all associated GUI types.
func (t *TopUI) RunApp() {
	scanner, err := t.execTop()
	if err != nil {
		panic(err)
	}

	go t.scanSTDOUT(scanner)

	ui := widgets.NewQApplication(len(os.Args), os.Args)

	t.window = widgets.NewQMainWindow(nil, 0)
	t.window.SetMinimumSize2(450, 675)
	t.window.SetMaximumSize2(450, 675)
	t.window.SetWindowTitle("Top Process Monitor")

	h1 := widgets.NewQHBoxLayout()
	h2 := widgets.NewQHBoxLayout()
	h3 := widgets.NewQHBoxLayout()
	v := widgets.NewQVBoxLayout()

	timer1 := core.NewQTimer(t.window)
	timer1.ConnectTimeout(func() { t.updateUI() })
	timer1.Start(1000)

	heading1 := widgets.NewQLabel(t.window, 0)
	heading1.SetText("PID") // 0

	heading2 := widgets.NewQLabel(t.window, 0)
	heading2.SetText("Program%") // 1

	heading3 := widgets.NewQLabel(t.window, 0)
	heading3.SetText("CPU%") // 2

	heading4 := widgets.NewQLabel(t.window, 0)
	heading4.SetText("MEMORY") // 7

	heading5 := widgets.NewQLabel(t.window, 0)
	heading5.SetText("UID") // 16

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
	h1.Layout().AddWidget(heading5)
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Errorf("an error occurred: %v", err))
			os.Exit(1)
		}
	}()

	var topUI = &TopUI{}
	topUI.RunApp()
}
