package gui

import (
	"WeatherMonitor/database"
	"WeatherMonitor/file_io"
	"WeatherMonitor/plotting"
	"database/sql"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/maps"
	"image/color"
	"math/rand"
	"strings"
)

var defaultRatio = []float32{0.8, 0.2}

func getOnlyUsefulData(sensorsSlice []string) []string {
	usefulSensorSlice := make([]string, 0)
	for _, sensor := range sensorsSlice {
		if !strings.Contains(sensor, "System") && !strings.Contains(sensor, "RTC") {
			usefulSensorSlice = append(usefulSensorSlice, sensor)
		}
	}

	return usefulSensorSlice
}

func AddNewLineWindow(app *fyne.App, db *sql.DB, infoMap file_io.DevicesInfoMap, elem *plotting.PlotElem, parentWidgetRefresh func(), redacting bool) *fyne.Window {
	NewLineWindow := (*app).NewWindow("Add New Element")
	NewLineWindow.Resize(fyne.NewSize(500, 500))

	var SensorSelect, DeviceSelect, SerialSelect, GroupTypeSelect *widget.Select
	var FuncDataRadio *widget.RadioGroup
	firstContainer := NewAdaptiveGridWithRatios(defaultRatio, widget.NewLabel("Line Name"), widget.NewLabel("Color"))

	elemNameEnry := widget.NewEntry()
	colorRectangle := canvas.NewRectangle(color.NRGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255,
	})
	colorPicker := dialog.NewColorPicker("Color Picker", "Chose element color", func(c color.Color) {
		colorRectangle.FillColor = c
		colorRectangle.Refresh()
	}, NewLineWindow)
	colorPicker.Advanced = true
	colorButton := widget.NewButton("", func() {
		colorPicker.Show()
	})
	colorButtonContainer := container.NewMax(colorButton, colorRectangle)
	secondContainer := NewAdaptiveGridWithRatios(defaultRatio, elemNameEnry, colorButtonContainer)

	thirdContainer := NewAdaptiveGridWithRatios(defaultRatio, widget.NewLabel("Device Name"), widget.NewLabel("Serial"))

	DeviceSelect = widget.NewSelect(maps.Keys(infoMap), func(deviceName string) {
		SensorSelect.Selected = ""
		SerialSelect.Selected = ""
		SensorSelect.Options = getOnlyUsefulData(infoMap[deviceName].DeviceSensors)
		if len(infoMap[deviceName].DeviceSerials) > 1 {
			SerialSelect.Options = infoMap[deviceName].DeviceSerials
			if SerialSelect.Disabled() {
				SerialSelect.Enable()
			}
		} else {
			SerialSelect.Disable()
		}
		SensorSelect.Enable()
		SerialSelect.Refresh()
		SensorSelect.Refresh()
	})

	SerialSelect = widget.NewSelect([]string{}, func(serial string) {
	})

	SerialSelect.Disable()

	fourthContainer := NewAdaptiveGridWithRatios(defaultRatio, DeviceSelect, SerialSelect)

	Label3 := widget.NewLabel("Sensor Name")
	SensorSelect = widget.NewSelect([]string{}, func(sensorName string) {
	})
	SensorSelect.Disable()

	GroupTypeSelect = widget.NewSelect([]string{"By Hour", "By 3 Hours", "By Day"}, func(s string) {})
	GroupTypeSelect.Disable()
	GroupDataRadio := widget.NewRadioGroup([]string{"Raw", "Grouped"}, func(s string) {
		switch s {
		case "Raw":
			GroupTypeSelect.Selected = ""
			GroupTypeSelect.Refresh()
			GroupTypeSelect.Disable()
			FuncDataRadio.Selected = ""
			GroupTypeSelect.Refresh()
			FuncDataRadio.Disable()
		case "Grouped":
			GroupTypeSelect.Selected = GroupTypeSelect.Options[0]
			FuncDataRadio.Enable()
			GroupTypeSelect.Enable()
		}
	})

	GroupDataRadio.Selected = "Raw"

	fifthContainer := NewAdaptiveGridWithRatios(defaultRatio, GroupDataRadio, GroupTypeSelect)

	FuncDataRadio = widget.NewRadioGroup([]string{"Max", "Min", "Average"}, func(TypeOfFunc string) {})
	FuncDataRadio.Disable()

	formContainer := container.NewVBox(firstContainer, secondContainer, thirdContainer, fourthContainer, Label3, SensorSelect, fifthContainer, FuncDataRadio)

	if redacting {
		elemNameEnry.Text = elem.ElemName
		colorRectangle.FillColor = elem.Color
		DeviceSelect.SetSelected(elem.DeviceName)
		SerialSelect.SetSelected(elem.Serial)
		SensorSelect.SetSelected(elem.SensorName)
		GroupDataRadio.SetSelected(elem.IsGrouping)
		GroupTypeSelect.SetSelected(elem.TypeOfGrouping)
		FuncDataRadio.SetSelected(elem.TypeOfFunc)

		formContainer.Refresh()
	}
	saveButton := widget.NewButton("Save", func() {
		// check if request is ok
		if DeviceSelect.Selected == "" ||
			SensorSelect.Selected == "" ||
			(SerialSelect.Selected == "" && len(infoMap[DeviceSelect.Selected].DeviceSerials) > 1) ||
			(GroupDataRadio.Selected == GroupDataRadio.Options[0] && FuncDataRadio.Selected != "") ||
			(GroupDataRadio.Selected == GroupDataRadio.Options[1] && FuncDataRadio.Selected == "") {
			dialog.NewInformation("Error", "Wrong request data\nPlease check your data", NewLineWindow).Show()
		} else {
			request := database.SensorDataRequest{
				TableName:    database.ToSnakeCase(infoMap[DeviceSelect.Selected].DeviceName) + "s",
				SensorName:   database.ToSnakeCase(SensorSelect.Selected),
				SensorSerial: SerialSelect.Selected,
			}
			switch GroupDataRadio.Selected {
			case GroupDataRadio.Options[0]:
				request.TypeOfDataGrouping = database.NoGrouping
				request.TypeOfDataFunc = database.IsRaw
			case GroupDataRadio.Options[1]:
				switch GroupTypeSelect.Selected {
				case GroupTypeSelect.Options[0]:
					request.TypeOfDataGrouping = database.PerHour
				case GroupTypeSelect.Options[1]:
					request.TypeOfDataGrouping = database.Per3Hours
				case GroupTypeSelect.Options[2]:
					request.TypeOfDataGrouping = database.PerDay
				}

				switch FuncDataRadio.Selected {
				case FuncDataRadio.Options[0]:
					request.TypeOfDataFunc = database.IsMax
				case FuncDataRadio.Options[1]:
					request.TypeOfDataFunc = database.IsMin
				case FuncDataRadio.Options[2]:
					request.TypeOfDataFunc = database.IsAvg
				}
			}

			elem.ElemName = elemNameEnry.Text
			elem.Color = colorRectangle.FillColor
			elem.DeviceName = DeviceSelect.Selected
			elem.Serial = SerialSelect.Selected
			elem.SensorName = SensorSelect.Selected
			elem.IsGrouping = GroupDataRadio.Selected
			elem.TypeOfGrouping = GroupTypeSelect.Selected
			elem.TypeOfFunc = FuncDataRadio.Selected

			request.BeginDateTime = elem.Request.BeginDateTime
			request.EndDateTime = elem.Request.EndDateTime

			elem.Request = request

			parentWidgetRefresh()
			NewLineWindow.Close()
		}

	})
	cancelButton := widget.NewButton("Cancel", NewLineWindow.Close)
	buttons := NewAdaptiveGridWithRatios([]float32{0.5, 0.25, 0.25}, widget.NewLabel(""), cancelButton, saveButton)
	content := container.NewBorder(nil, buttons, nil, nil, formContainer)

	NewLineWindow.SetContent(content)
	return &NewLineWindow
}
