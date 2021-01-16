package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// 2021 rootVIII - List top processes in a simple QT GUI

// TopUI represents the main application type
// and controls/updates the GUI in real time.
type TopUI struct {
	window         *widgets.QMainWindow
	tableBox       *widgets.QTableWidget
	titleView      *widgets.QGraphicsView
	verticalLayout *widgets.QVBoxLayout
	cmdBuffer      [][]string
	rowMax         int
}

func (t *TopUI) scanSTDOUT(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) > 0 {
			if _, err := strconv.Atoi(line[0]); err == nil {
				if line[1] != "top" && line[1] != "topui" && line[1] != "qtbox" {
					t.cmdBuffer = append(t.cmdBuffer, []string{line[0], line[2], line[1]})
				}
			}
		}
	}
}

func (t *TopUI) unloadBuffer() {
	if len(t.cmdBuffer) > 0 {
		for index, line := range t.cmdBuffer {
			if index > t.rowMax {
				break
			}
			for col := 0; col < 3; col++ {
				t.tableBox.SetItem(index, col, widgets.NewQTableWidgetItem2(line[col], 0))
			}
		}
		t.cmdBuffer = nil
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

func (t *TopUI) setWindow() {
	t.window = widgets.NewQMainWindow(nil, 0)
	t.window.SetMinimumSize2(480, 685)
	t.window.SetMaximumSize2(480, 685)
	t.window.SetWindowTitle("💻")
}

func (t *TopUI) setTitle() {
	title := widgets.NewQGraphicsScene(t.window)
	title.AddText("Top Processes", gui.NewQFont2("Menlo", 12, 1, false))
	t.titleView = widgets.NewQGraphicsView(t.window)
	t.titleView.SetScene(title)
	t.titleView.SetFixedHeight(30)
}

func (t *TopUI) setTable() {
	t.tableBox = widgets.NewQTableWidget(t.window)
	t.tableBox.SetFixedHeight(620)
	t.tableBox.SetColumnCount(3)
	t.rowMax = 500
	t.tableBox.SetRowCount(t.rowMax + 1)
	t.tableBox.SetHorizontalHeaderLabels([]string{"PID", "CPU%", "APP"})
	t.tableBox.VerticalHeader().Hide()
	t.tableBox.SetColumnWidth(0, 80)
	t.tableBox.SetColumnWidth(1, 80)
	t.tableBox.SetColumnWidth(2, 270)
}

func (t *TopUI) showWindow() {
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(t.verticalLayout)
	t.window.SetCentralWidget(widget)
	t.window.Show()
}

func (t *TopUI) setVLayout() {
	h1 := widgets.NewQHBoxLayout()
	h2 := widgets.NewQHBoxLayout()
	t.verticalLayout = widgets.NewQVBoxLayout()
	h1.Layout().AddWidget(t.titleView)
	h2.Layout().AddWidget(t.tableBox)

	for _, layout := range []*widgets.QHBoxLayout{h1, h2} {
		t.verticalLayout.AddLayout(layout, 0)
	}
}

func (t *TopUI) buildApp() {
	t.setWindow()
	t.setTitle()
	t.setTable()
	t.setVLayout()
	t.showWindow()

	timer := core.NewQTimer(t.window)
	timer.ConnectTimeout(func() { t.unloadBuffer() })
	timer.Start(1000)
}

// RunApp initializes the GUI and all associated GUI types.
// Stdout from top bin is scanned in goroutine, written
// to a buffer, and then the buffer is emptied and used
// to populate the GUI within a QTimed method.
func (t *TopUI) RunApp() {
	ui := widgets.NewQApplication(len(os.Args), os.Args)
	t.buildApp()

	scanner, err := t.execTop()
	if err != nil {
		panic(err)
	}

	go t.scanSTDOUT(scanner)

	ui.Exec()
}

func init() {
	pathVar := os.Getenv("PATH")
	if len(pathVar) < 1 {
		log.Fatal("no $PATH variables found")
	}

	foundTop := false
	for _, path := range strings.Split(pathVar, ":") {
		_, err := os.Stat(filepath.Join(path, "top"))
		if err == nil {
			foundTop = true
			break
		}
	}

	if !foundTop {
		log.Fatal("top executable not found in $PATH")
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(fmt.Errorf("a fatal error occurred: %v", err))
		}
	}()

	var topUI = &TopUI{}
	topUI.RunApp()
}
