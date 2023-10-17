package gui

import (
	"WeatherMonitor/database"
	"WeatherMonitor/file_io"
	"WeatherMonitor/plotting"
	"database/sql"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"image/color"
	"slices"
	"strconv"
	"strings"
	"time"
)

type LineElement struct {
	Label        *widget.Label
	Rect         *canvas.Rectangle
	EditButton   *widget.Button
	DeleteButton *widget.Button

	PlotData      *plotting.PlotElem
	LineContainer *fyne.Container
}

const (
	linePlotChosen             = 0
	barPlotChosen              = 1
	effectiveTemperatureChosen = 2
)

func (l *LineElement) GetLineContainer() *fyne.Container {
	return l.LineContainer
}

func DeleteLineElem(slice []*plotting.PlotElem, n int) []*plotting.PlotElem {
	slice[n] = nil
	return slices.Delete(slice, n, n+1)
}

func NewLineElement(label string, color color.Color, EditFunc func(), DeleteFunc func()) *LineElement {
	lineElement := &LineElement{
		Label:        widget.NewLabel(label),
		Rect:         canvas.NewRectangle(color),
		EditButton:   widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), EditFunc),
		DeleteButton: widget.NewButtonWithIcon("", theme.DeleteIcon(), DeleteFunc),
	}

	LineName := container.NewHScroll(lineElement.Label)
	LineName.Resize(fyne.NewSize(32, 100))
	lineElement.Rect.Resize(fyne.NewSize(3, 32))
	lineElement.LineContainer = container.NewBorder(nil,
		nil,
		lineElement.Rect,
		container.NewHBox(lineElement.DeleteButton, lineElement.EditButton),
		LineName)
	//lineElement.EditButton)

	return lineElement
}

func NewDateTimePicker(w *fyne.Window, setTime time.Time, buttonText binding.String, targetTime *time.Time) *fyne.Container {
	hoursEntry := widget.NewEntry()
	hoursEntry.SetText("00")
	minutesEntry := widget.NewEntry()
	minutesEntry.SetText("00")
	secondsEntry := widget.NewEntry()
	secondsEntry.SetText("00")

	timeEntryContainer := NewAdaptiveGridWithRatios([]float32{0.1, 0.2, 0.1, 0.2, 0.1, 0.2, 0.1},
		widget.NewLabel(""),
		hoursEntry,
		widget.NewLabel(":"),
		minutesEntry,
		widget.NewLabel(":"),
		secondsEntry,
		widget.NewLabel(""))

	onClicked := func(t time.Time) {
		hours, err := strconv.Atoi(hoursEntry.Text)
		minutes, err := strconv.Atoi(minutesEntry.Text)
		seconds, err := strconv.Atoi(secondsEntry.Text)

		if err != nil || hours >= 24 || minutes >= 60 || seconds >= 60 {
			err = errors.New("Wrong time format\nCheck your fields")
			dialog.NewError(err, *w).Show()
		} else {
			*targetTime = time.Date(t.Year(), t.Month(), t.Day(), hours, minutes, seconds, 0, time.Local)
			buttonText.Set(targetTime.Format(strings.Replace(time.DateTime, " ", "\n", 1)))
		}
	}

	content := container.NewBorder(nil, timeEntryContainer, nil, nil, xwidget.NewCalendar(setTime, onClicked))

	return content
}

func NewPlot() {
	myApp := app.NewWithID("com.kirilmak.weathermonitor.preferences")

	mainWindow := myApp.NewWindow("Weather Monitor")
	mainWindow.Resize(fyne.NewSize(1920, 1080))
	myApp.Settings().SetTheme(theme.DarkTheme())

	var db *sql.DB
	plotChosen := linePlotChosen
	mainWindow.SetMaster()
	var plotName string
	var DevicesInfo file_io.DevicesInfoMap
	// var plotPath string
	loadJSONItem := fyne.NewMenuItem("Load .json", func() {
		openJsonDialog := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			db, err = database.CreateDB()
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			var openPath string
			if dir != nil {
				openPath = dir.URI().Path()
			}
			if openPath == "" {
				dialog.ShowError(errors.New("JSON file hasnt been chosen"), mainWindow)
			}
			wait1Dialog := dialog.NewProgressInfinite("Unmarshalling JSON", "Unmarshalling JSON\nPlease wait...", mainWindow)
			wait1Dialog.Show()
			responseMap, err := file_io.UnmarshalJSONIntoResponseMap(openPath)
			if err != nil {
				wait1Dialog.Hide()
				dialog.ShowError(err, mainWindow)
				return
			}
			wait1Dialog.Hide()

			wait2Dialog := dialog.NewProgressInfinite("Inserting Data", "Inserting data into database\nPlease wait...", mainWindow)
			wait2Dialog.Show()
			DevicesInfo, err = database.InsertDataIntoDB(db, responseMap)
			if err != nil {
				wait2Dialog.Hide()
				dialog.ShowError(err, mainWindow)
				return
			}
			wait2Dialog.Hide()

			wait3Dialog := dialog.NewProgressInfinite("Saving Device Info", "Saving device info\nPlease wait...", mainWindow)
			wait3Dialog.Show()
			err = file_io.MergeJSONFromSource(file_io.DefaultDevicesInfoPath, DevicesInfo)
			if err != nil {
				wait3Dialog.Hide()
				dialog.ShowError(err, mainWindow)
				return
			}
			wait3Dialog.Hide()

		}, mainWindow)
		openJsonDialog.SetOnClosed(func() {})
		openJsonDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		defaultLocation, _ := storage.ListerForURI(storage.NewFileURI("C:\\Users\\user\\GolandProjects\\WeatherMonitor"))
		openJsonDialog.SetLocation(defaultLocation)
		openJsonDialog.Show()
	})
	loadDBItem := fyne.NewMenuItem("Load .db", func() {
		openDBDialog := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			var openPath string
			if dir != nil {
				openPath = dir.URI().Path()
			}
			wait1Dialog := dialog.NewProgressInfinite("Wait...", "", mainWindow)
			wait1Dialog.Show()
			db, err = database.OpenDataBaseFromFile(openPath)
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			DevicesInfo, err = file_io.GetDeviceInfo(file_io.DefaultDevicesInfoPath)
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			wait1Dialog.Hide()
		}, mainWindow)
		openDBDialog.SetFilter(storage.NewExtensionFileFilter([]string{".db"}))
		defaultLocation, _ := storage.ListerForURI(storage.NewFileURI("C:\\Users\\user\\GolandProjects\\WeatherMonitor"))
		openDBDialog.SetLocation(defaultLocation)
		openDBDialog.Show()
	})

	beginDateButtonText := binding.NewString()
	beginDateButtonText.Set("PickDate")

	endDateButtonText := binding.NewString()
	endDateButtonText.Set("PickDate")

	beginDate := time.Time{}
	endDate := time.Time{}

	var pickBeginDataButton, pickEndDataButton, refreshDateButton *widget.Button
	Lines := make([]*plotting.PlotElem, 0)

	dateBeginPicker := dialog.NewCustom("Pick Date",
		"Submit",
		NewDateTimePicker(&mainWindow,
			time.Date(time.Now().Year(), time.February, time.Now().Day(), 0, 0, 0, 0, time.Local),
			beginDateButtonText,
			&beginDate),
		mainWindow)
	pickBeginDataButton = widget.NewButton("Pick Date", dateBeginPicker.Show)
	dateEndPicker := dialog.NewCustom("Pick Date",
		"Submit",
		NewDateTimePicker(&mainWindow,
			time.Date(time.Now().Year(), time.February, time.Now().Day(), 0, 0, 0, 0, time.Local),
			endDateButtonText,
			&endDate),
		mainWindow)
	pickEndDataButton = widget.NewButton("Pick Date", dateEndPicker.Show)

	refreshDateButton = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		textBegin, _ := beginDateButtonText.Get()
		pickBeginDataButton.SetText(textBegin)
		pickBeginDataButton.Refresh()

		textEnd, _ := endDateButtonText.Get()
		pickEndDataButton.SetText(textEnd)
		pickEndDataButton.Refresh()

		for i, _ := range Lines {
			Lines[i].Request.BeginDateTime = beginDate.Format(time.DateTime)
			Lines[i].Request.EndDateTime = endDate.Format(time.DateTime)
			fmt.Println(*Lines[i])
		}
	})

	IsDarkTheme := false
	fileMenu := fyne.NewMenu("File", loadJSONItem, fyne.NewMenuItemSeparator(), loadDBItem)
	themeMenu := fyne.NewMenu("Theme", fyne.NewMenuItem("Change Theme", func() {
		if IsDarkTheme {
			myApp.Settings().SetTheme(theme.LightTheme())
			IsDarkTheme = false
		} else {
			myApp.Settings().SetTheme(theme.DarkTheme())
			IsDarkTheme = true
		}
	}))
	plotNameEntry := widget.NewEntry()
	plotNameEntry.PlaceHolder = "Name your plot!"
	plotNameContainer := container.NewVBox(plotNameEntry)
	nameOfPlotMenu := fyne.NewMenuItem("Name of plot", func() {
		dialog.NewCustomConfirm("Name Your Plot", "Confirm", "Dismiss", plotNameContainer, func(confirmed bool) {
			if confirmed {
				plotName = plotNameEntry.Text
			}

		}, mainWindow).Show()
	})
	typeOfPlotMenu := fyne.NewMenuItem("Type of plot", nil)
	typeOfPlotMenu.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Line plot", func() { plotChosen = linePlotChosen }),
		fyne.NewMenuItem("Bar plot", func() { plotChosen = barPlotChosen }),
		fyne.NewMenuItem("Effective temperature plot", func() { plotChosen = effectiveTemperatureChosen }))
	plotMenu := fyne.NewMenu("Plot", typeOfPlotMenu, nameOfPlotMenu)
	mainMenu := fyne.NewMainMenu(fileMenu, themeMenu, plotMenu)

	imgUrl, _ := fyne.LoadResourceFromPath("test_plot.png")

	img := canvas.NewImageFromResource(imgUrl)
	img.FillMode = canvas.ImageFillContain

	mainWindow.SetMainMenu(mainMenu)

	var graphList *widget.List
	graphList = widget.NewList(
		func() int { return len(Lines) },
		func() fyne.CanvasObject {
			return NewLineElement("Test", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, nil, nil).GetLineContainer()
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*widget.Label).SetText(Lines[id].ElemName)
			obj.(*fyne.Container).Objects[1].(*canvas.Rectangle).FillColor = Lines[id].Color
			obj.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				Lines = DeleteLineElem(Lines, id)
				graphList.Refresh()
			}
			obj.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				newWindow := *AddNewLineWindow(&myApp, db, DevicesInfo, Lines[id], graphList.Refresh, true)
				newWindow.Show()
			}
		},
	)

	addElementButton := widget.NewButton("Add New Element", func() {
		newLine := plotting.PlotElem{}
		newLine.Request.BeginDateTime = beginDate.Format(time.DateTime)
		newLine.Request.EndDateTime = endDate.Format(time.DateTime)
		newWindow := *AddNewLineWindow(&myApp, db, DevicesInfo, &newLine, func() { graphList.Refresh() }, false)
		newWindow.Show()
		newWindow.SetOnClosed(func() {
			if newLine.Color != nil {
				Lines = append(Lines, &newLine)
			}
		})
	})

	imgContainer := container.NewMax(img)
	refreshPlotButton := widget.NewButtonWithIcon("Refresh Plot", theme.ViewRefreshIcon(), func() {
		for i, _ := range Lines {
			progressBar := dialog.NewProgressInfinite("Fetching Data For"+Lines[i].ElemName, "Fetching data from database.\nPlease wait.", mainWindow)
			progressBar.Show()
			sliceOfDots, err := database.GetDataFromDB(db, Lines[i].Request)
			if err != nil {
				progressBar.Hide()
				errorDialog := dialog.NewError(err, mainWindow)
				errorDialog.SetOnClosed(mainWindow.Close)
			}

			progressBar.Hide()
			Lines[i].SliceOfDots = sliceOfDots
		}
		progressBar := dialog.NewProgressInfinite("Making Plot", "In process of making plot.\nPlease wait...", mainWindow)
		progressBar.Show()
		switch plotChosen {
		case linePlotChosen:
			err := plotting.CreateLinePlot(Lines, plotName)
			if err != nil {
				errorDialog := dialog.NewError(err, mainWindow)
				errorDialog.Show()
				errorDialog.SetOnClosed(myApp.Quit)
			}
			imgUrl, _ := fyne.LoadResourceFromPath(plotting.DefaultLinePlotPath)
			img := canvas.NewImageFromResource(imgUrl)
			imgContainer.RemoveAll()
			imgContainer.Add(img)
		case barPlotChosen:
			err := plotting.CreateBarChart(Lines, plotName)
			if err != nil {
				errorDialog := dialog.NewError(err, mainWindow)
				errorDialog.Show()
				errorDialog.SetOnClosed(myApp.Quit)
			}
			imgUrl, _ := fyne.LoadResourceFromPath(plotting.DefaultBarPlotPath)
			img := canvas.NewImageFromResource(imgUrl)
			imgContainer.RemoveAll()
			imgContainer.Add(img)
		default:
			err := plotting.CreateLinePlot(Lines, plotName)
			if err != nil {
				errorDialog := dialog.NewError(err, mainWindow)
				errorDialog.Show()
				errorDialog.SetOnClosed(myApp.Quit)
			}
			imgUrl, _ := fyne.LoadResourceFromPath(plotting.DefaultLinePlotPath)
			img := canvas.NewImageFromResource(imgUrl)
			imgContainer.RemoveAll()
			imgContainer.Add(img)
		}

		progressBar.Hide()

	})

	graphElem := container.NewBorder(
		container.NewHBox(widget.NewLabel("From"), pickBeginDataButton, widget.NewLabel("To"), pickEndDataButton, refreshDateButton),
		container.NewHBox(addElementButton, refreshPlotButton),
		nil,
		nil,
		graphList)

	content := container.NewHSplit(graphElem, imgContainer)
	content.SetOffset(0.15)
	mainWindow.SetContent(content)

	mainWindow.SetOnClosed(func() {
		db.Close()
	})
	mainWindow.Show()
	myApp.Run()
}
