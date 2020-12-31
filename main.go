package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// 2020 rootVIII - List top processes in a simple QT GUI

// TopUI represents the main application window.
type TopUI struct {
	window    *widgets.QMainWindow
	tableBox  *widgets.QTableWidget
	cmdBuffer []string
}

func (t *TopUI) updateUI() {
	// if len(t.cmdBuffer) > 0 {
	// 	t.tableBox.Clear()
	// 	t.tableBox.AddItems(t.cmdBuffer)
	// 	t.cmdBuffer = nil
	// }
}

func (t *TopUI) scanSTDOUT(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) > 0 {
			if _, err := strconv.Atoi(line[0]); err == nil {
				if line[1] != "top" && line[1] != "topui" {
					fields := fmt.Sprintf("%-10s%-20s%-10s", line[0], line[1], line[2])
					t.cmdBuffer = append(t.cmdBuffer, fields)
				}
			}
		}
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
	t.window.SetMinimumSize2(480, 675)
	t.window.SetMaximumSize2(480, 675)
	t.window.SetWindowTitle("Top 20 Processes")

	h1 := widgets.NewQHBoxLayout()
	h2 := widgets.NewQHBoxLayout()
	v := widgets.NewQVBoxLayout()

	timer1 := core.NewQTimer(t.window)
	timer1.ConnectTimeout(func() { t.updateUI() })
	timer1.Start(1000)

	divider := widgets.NewQGraphicsScene(t.window)
	titleView := widgets.NewQGraphicsView(t.window)
	titleView.SetScene(divider)
	titleView.SetFixedHeight(2)

	t.tableBox = widgets.NewQTableWidget(t.window)
	t.tableBox.SetFixedHeight(625)
	t.tableBox.SetColumnCount(3)
	t.tableBox.SetRowCount(500)
	t.tableBox.SetHorizontalHeaderLabels([]string{"PID", "CPU%", "APP"})
	t.tableBox.VerticalHeader().Hide()
	t.tableBox.SetColumnWidth(0, 100)
	t.tableBox.SetColumnWidth(1, 100)
	t.tableBox.SetColumnWidth(2, 250)

	h1.Layout().AddWidget(titleView)
	h2.Layout().AddWidget(t.tableBox)

	for _, layout := range []*widgets.QHBoxLayout{h1, h2} {
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
