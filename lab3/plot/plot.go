package plot

import (
	"image/color"
	"log"
	"math"

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

func StyleCartesian(p *plot.Plot) {
	SetupCartesian(p)
	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.XOffs = -vg.Points(5)
	p.Legend.YOffs = -vg.Points(5)
}

func SaveLinePlot(x, y []float64, filename, xlabel, ylabel, title string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}

	p.Add(line)
	StyleCartesian(p)

	if err := p.Save(8*vg.Inch, 6*vg.Inch, filename); err != nil {
		log.Fatal(err)
	}
}

func SaveComparisonPlot(x, y1, y2 []float64, label1, label2, filename, xlabel, ylabel, title string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	pts1 := make(plotter.XYs, len(x))
	pts2 := make(plotter.XYs, len(x))
	for i := range x {
		pts1[i].X = x[i]
		pts1[i].Y = y1[i]
		pts2[i].X = x[i]
		pts2[i].Y = y2[i]
	}

	line1, _ := plotter.NewLine(pts1)
	line1.Color = color.NRGBA{B: 255, A: 255}
	line1.Width = vg.Points(2)

	line2, _ := plotter.NewLine(pts2)
	line2.Color = color.NRGBA{R: 255, A: 255}
	line2.Width = vg.Points(2)

	p.Add(line1, line2)
	p.Legend.Add(label1, line1)
	p.Legend.Add(label2, line2)

	StyleCartesian(p)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		log.Fatal(err)
	}
}

func PlotSingleDependency(x []float64, yList [][]float64, names []string, title, xLabel, yLabel, filename string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	colors := []color.Color{
		color.NRGBA{R: 0, G: 0, B: 255, A: 255},
		color.NRGBA{R: 255, G: 0, B: 0, A: 255},
		color.NRGBA{R: 0, G: 128, B: 0, A: 255},
		color.NRGBA{R: 255, G: 165, B: 0, A: 255},
	}

	for idx, y := range yList {
		limit := len(x)
		if len(y) < limit {
			limit = len(y)
		}
		pts := make(plotter.XYs, 0, limit)
		for i := 0; i < limit; i++ {
			if math.IsNaN(x[i]) || math.IsInf(x[i], 0) || math.IsNaN(y[i]) || math.IsInf(y[i], 0) {
				continue
			}
			pts = append(pts, plotter.XY{X: x[i], Y: y[i]})
		}
		if len(pts) < 2 {
			log.Fatalf("недостаточно корректных точек для графика %s", filename)
		}
		line, err := plotter.NewLine(pts)
		if err != nil {
			log.Fatalf("ошибка построения графика %s: %v", filename, err)
		}
		line.Color = colors[idx%len(colors)]
		line.Width = vg.Points(2)
		if idx%2 != 0 {
			line.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		}
		p.Add(line)
		if len(names) > idx {
			p.Legend.Add(names[idx], line)
		}
	}

	StyleCartesian(p)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		log.Fatal(err)
	}
}
