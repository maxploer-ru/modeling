package main

import (
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func setupCartesian(p *plot.Plot) {
	p.Add(plotter.NewGrid())
	transparent := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	p.X.Color = transparent
	p.Y.Color = transparent
	p.X.Tick.Color = color.Black
	p.Y.Tick.Color = color.Black
	p.X.Tick.Label.Color = color.Black
	p.Y.Tick.Label.Color = color.Black
}

func styleCartesian(p *plot.Plot) {
	setupCartesian(p)
	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.XOffs = -vg.Points(5)
	p.Legend.YOffs = -vg.Points(5)
}

func saveLinePlot(x, y []float64, filename, xlabel, ylabel, title string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	line, _ := plotter.NewLine(pts)
	p.Add(line)
	styleCartesian(p)
	if err := p.Save(8*vg.Inch, 6*vg.Inch, filename); err != nil {
		log.Fatal(err)
	}
}
