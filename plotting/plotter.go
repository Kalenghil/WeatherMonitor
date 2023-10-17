package plotting

import (
	"WeatherMonitor/database"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"time"
)

const (
	defaultWidth        = 1610
	defaultHeigth       = 1025
	DefaultLinePlotPath = "./plots/line_plot.png"
	DefaultBarPlotPath  = "./plots/bar_plot.png"
	DefaultHeatMapPath  = "./plots/heapmap.png"
)

type PlotElem struct {
	ElemName       string
	Color          color.Color
	DeviceName     string
	Serial         string
	SensorName     string
	IsGrouping     string
	TypeOfGrouping string
	TypeOfFunc     string

	Request database.SensorDataRequest

	SliceOfDots *database.PlotDataArray
}

func castToPoints(array *database.PlotDataArray) plotter.XYs {
	points := make(plotter.XYs, len(*array))
	fmt.Println(len(*array))
	for i := range *array {
		date, _ := time.Parse(time.DateTime, (*array)[i].Datetime)
		points[i].X = float64(date.Unix())
		points[i].Y = float64((*array)[i].Value)
	}
	return points
}

func CreateLinePlot(linesData []*PlotElem, plotName string) error {
	p := plot.New()
	p.Add(plotter.NewGrid())
	p.Title.Text = plotName
	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04:02", Ticker: CustomTimeTicks{}}

	lines := make([]plotter.Line, len(linesData))
	for i, _ := range lines {
		lineXYs := castToPoints(linesData[i].SliceOfDots)
		line, err := plotter.NewLine(lineXYs)
		if err != nil {
			return err
		}
		line.Color = linesData[i].Color
		p.Add(line)
		p.Legend.Add(linesData[i].ElemName, line)
	}

	p.X.Tick.Marker = xticks

	err := p.Save(defaultWidth, defaultHeigth, DefaultLinePlotPath)
	return err
}

func CreateBarChart(linesData []*PlotElem, plotName string) error {
	p := plot.New()
	p.Add(plotter.NewGrid())
	p.Title.Text = plotName
	// xticks := plot.TimeTicks{Format: "2006-01-02\n15:04:02", Ticker: CustomTimeTicks{}}
	castToValues := func(array *database.PlotDataArray) plotter.Values {
		values := make(plotter.Values, 0)
		for i := range *array {
			values = append(values, float64((*array)[i].Value))
		}
		return values
	}
	castToStrings := func(array *database.PlotDataArray) []string {
		values := make([]string, 0)
		for i := range *array {
			values = append(values, (*array)[i].Datetime)
		}
		return values
	}
	for i, _ := range linesData {
		fmt.Println(*linesData[i].SliceOfDots)
		values := castToValues(linesData[i].SliceOfDots)
		fmt.Println(values)
		barChart, err := plotter.NewBarChart(values, vg.Points(10))
		if err != nil {
			return err
		}
		barChart.Color = linesData[i].Color
		p.Add(barChart)
		p.Legend.Add(linesData[i].ElemName, barChart)
		p.NominalX(castToStrings(linesData[i].SliceOfDots)...)
	}
	/*p.X.Tick.Marker = xticks*/

	/*date, _ := time.Parse(time.DateTime, linesData[0].Request.BeginDateTime)
	p.X.Min = float64(date.Unix())

	date, _ = time.Parse(time.DateTime, linesData[0].Request.EndDateTime)
	p.X.Max = float64(date.Unix())*/

	err := p.Save(defaultWidth, defaultHeigth, DefaultBarPlotPath)
	return err
}
