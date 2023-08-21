package guilib

import (
	"os"

	"github.com/therecipe/qt/widgets"
)

func InitNewMainWindow(title string, x, y int) *widgets.QMainWindow {
	widgets.NewQApplication(len(os.Args), os.Args)
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(x, y)
	window.SetWindowTitle(title)

	return window
}

func AddCentralWidget(window *widgets.QMainWindow) *widgets.QWidget {

	widget := widgets.NewQWidget(nil, 0)

	widget.SetLayout(widgets.NewQVBoxLayout())

	window.SetCentralWidget(widget)

	return widget
}

func AddButton(text string, widget *widgets.QWidget) *widgets.QPushButton {

	button := widgets.NewQPushButton2(text, nil)

	widget.Layout().AddWidget(button)

	return button
}

func AddEdit(text string, widget *widgets.QWidget) *widgets.QLineEdit {

	edit := widgets.NewQLineEdit(nil)
	edit.SetPlaceholderText(text)

	widget.Layout().AddWidget(edit)

	return edit
}

func AddLabel(text string, widget *widgets.QWidget) *widgets.QLabel {
	label := widgets.NewQLabel2(text, widget, 20)
	widget.Layout().AddWidget(label)

	return label
}
