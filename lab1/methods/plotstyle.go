package methods

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func SetupCartesian(p *plot.Plot) {

	p.Add(plotter.NewGrid())

	transparent := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	p.X.Color = transparent
	p.Y.Color = transparent

	p.X.Tick.Color = color.Black
	p.Y.Tick.Color = color.Black
	p.X.Tick.Label.Color = color.Black
	p.Y.Tick.Label.Color = color.Black
}

func AddAxis(p *plot.Plot) {
	axisColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	xMin, xMax := p.X.Min, p.X.Max
	if xMin == xMax {
		xMin, xMax = -1, 1
	}
	hLine := plotter.XYs{{X: xMin, Y: 0}, {X: xMax, Y: 0}}
	hAxis, _ := plotter.NewLine(hLine)
	hAxis.LineStyle.Color = axisColor
	hAxis.LineStyle.Width = vg.Points(1.5)
	hAxis.LineStyle.Dashes = nil
	p.Add(hAxis)

	yMin, yMax := p.Y.Min, p.Y.Max
	if yMin == yMax {
		yMin, yMax = -1, 1
	}
	vLine := plotter.XYs{{X: 0, Y: yMin}, {X: 0, Y: yMax}}
	vAxis, _ := plotter.NewLine(vLine)
	vAxis.LineStyle.Color = axisColor
	vAxis.LineStyle.Width = vg.Points(1.5)
	vAxis.LineStyle.Dashes = nil
	p.Add(vAxis)
}

func StyleCartesian(p *plot.Plot) {
	SetupCartesian(p)

	if p.Y.Min > 0 {
		p.Y.Min = -0.1 * p.Y.Max
	}
	if p.Y.Max < 0 {
		p.Y.Max = -0.1 * p.Y.Min
	}
	if p.X.Min > 0 {
		p.X.Min = -0.05 * p.X.Max
	}
	if p.X.Max < 0 {
		p.X.Max = -0.05 * p.X.Min
	}

	AddAxis(p)

	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.XOffs = -vg.Points(5)
	p.Legend.YOffs = -vg.Points(5)
}
