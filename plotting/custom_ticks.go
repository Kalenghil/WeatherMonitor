package plotting

import (
	"fmt"
	"gonum.org/v1/plot"
	"math"
)

type CustomTimeTicks struct{}

// Ticks CustomTimeTicks provides a slice of ticks for time series plots.
func (CustomTimeTicks) Ticks(xmin, xmax float64) []plot.Tick {
	var ticks []plot.Tick

	// computing order of range (position of least significant digit)
	xorder := int(math.Log10(xmax-xmin)+0.5) - 1
	format := ".0f"
	if xorder < 1 {
		format = fmt.Sprintf(".%df", -xorder)
	}
	// stepping is a power of 10 with integer exponent (xorder)
	xstep := math.Pow10(xorder)

	// make step a multiple of the largest convenient time increment
	// (may want to refine/extend the cases here)
	switch {
	case xstep >= 4*3600: // hours
		xstep = 3600.0 * float64((int(xstep) / 3600))
	case xstep >= 4*600: // 10 minutes
		xstep = 600.0 * float64((int(xstep) / 600))
	case xstep >= 4*60: // minutes
		xstep = 60.0 * float64((int(xstep) / 60))
	case xstep >= 40: // 10 seconds
		xstep = 10.0 * float64((int(xstep) / 10))
	case xstep >= 4: // seconds
		xstep = 1
	}

	// tuning step
	if (xmax-xmin)/xstep > 20 {
		xstep *= 5
	}

	// first big tick is rounded to the correct significant digit
	xoffset := float64(int(xmin/xstep)) * xstep

	// creating big ticks
	for x := xoffset; x <= xmax; x += xstep {
		label := fmt.Sprintf("%"+format, x)
		ticks = append(ticks, plot.Tick{Value: x, Label: label})
	}

	// 5 small ticks for each big tick
	xsub := xstep / 5
	for x := xoffset - xsub; x >= xmin; x -= xsub {
		ticks = append(ticks, plot.Tick{Value: x, Label: ""})
	}
	for x := xoffset + xsub; x <= xmax; x += xsub {
		ticks = append(ticks, plot.Tick{Value: x, Label: ""})
	}

	return ticks
}
