package task3

import (
	"fmt"
	"math"
	"modeling/lab1/methods"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func picard1(x float64) float64 {
	return math.Pow(x, 3.0) / 3.0
}

func picard2(x float64) float64 {
	x3 := math.Pow(x, 3.0)
	x7 := math.Pow(x, 7.0)
	return x3/3.0 + x7/63.0
}

func picard3(x float64) float64 {
	x3 := math.Pow(x, 3.0)
	x7 := math.Pow(x, 7.0)
	x11 := math.Pow(x, 11.0)
	x15 := math.Pow(x, 15.0)
	return x3/3.0 + x7/63.0 + 2.0*x11/2079.0 + x15/59535.0
}

func picard4(x float64) float64 {
	x3 := math.Pow(x, 3)
	x7 := math.Pow(x, 7)
	x11 := math.Pow(x, 11)
	x15 := math.Pow(x, 15)
	x19 := math.Pow(x, 19)
	x23 := math.Pow(x, 23)
	x27 := math.Pow(x, 27)
	x31 := math.Pow(x, 31)
	return x3/3 + x7/63 + (2*x11)/2079 + (13*x15)/218295 +
		(82*x19)/37328445 + (662*x23)/10438212015 +
		(4*x27)/3341878155 + x31/109876902975
}

type result struct {
	x  float64
	p1 float64
	p2 float64
	p3 float64
	p4 float64
	eu float64
}

func computeAll(h float64, epsilon float64, sampleEvery int) ([]result, float64) {
	var results []result
	results = append(results, result{0, 0, 0, 0, 0, 0})

	x := h
	yEuler := 0.0
	yRunge := 0.0
	cnt := 1

	for {
		p1 := picard1(x)
		p2 := picard2(x)
		p3 := picard3(x)
		p4 := picard4(x)

		yRunge = yEuler + h/2*(x*x+yRunge*yRunge)
		xHalf := x + h/2
		yRunge = yRunge + h/2*(xHalf*xHalf+yRunge*yRunge)

		yEuler = yEuler + h*(x*x+yEuler*yEuler)

		var relError float64
		if math.Abs(yRunge) > 1e-5 {
			relError = math.Abs(yEuler-yRunge) / yRunge
		}

		if relError > epsilon {
			break
		}

		if cnt%sampleEvery == 0 {
			results = append(results, result{x, p1, p2, p3, p4, yEuler})
		}
		cnt++
		x += h
	}
	xMax := x - h
	return results, xMax
}

func Run() {
	h := 1e-7
	epsilon := 1e-4
	sampleEvery := 10000

	results, xMax := computeAll(h, epsilon, sampleEvery)
	fmt.Printf("x_max = %.4f\n", xMax)
	fmt.Println()

	fmt.Println("Таблица значений u(x):")
	header := []string{"x", "Пикар 1", "Пикар 2", "Пикар 3", "Пикар 4", "Эйлер"}
	var rows [][]string
	for _, r := range results {
		rows = append(rows, []string{
			fmt.Sprintf("%.15f", r.x),
			fmt.Sprintf("%.8f", r.p1),
			fmt.Sprintf("%.8f", r.p2),
			fmt.Sprintf("%.8f", r.p3),
			fmt.Sprintf("%.8f", r.p4),
			fmt.Sprintf("%.8f", r.eu),
		})
	}
	methods.PrintTable(header, rows)
	fmt.Println()

	epsApprox := 0.01
	bounds := [4]float64{}

	for _, r := range results {
		diffs := [4]float64{
			math.Abs(r.p1 - r.p2),
			math.Abs(r.p2 - r.p3),
			math.Abs(r.p3 - r.p4),
			math.Abs(r.p4 - r.eu),
		}
		for j := 0; j < 4; j++ {
			if diffs[j] <= epsApprox {
				bounds[j] = r.x
			}
		}
	}

	fmt.Println("Границы применимости приближений Пикара (|u_i - эталон| ≤ 0.01):")
	for i := 0; i < 4; i++ {
		fmt.Printf("  Пикар %d  применим на [0, %.3f]\n", i+1, bounds[i])
	}
	fmt.Println()

	pl := plot.New()
	pl.Title.Text = "u' = x² + u²"
	pl.X.Label.Text = "x"
	pl.Y.Label.Text = "u(x)"

	var p1XY, p2XY, p3XY, p4XY, eulerXY plotter.XYs
	for _, r := range results {
		p1XY = append(p1XY, plotter.XY{X: r.x, Y: r.p1})
		p2XY = append(p2XY, plotter.XY{X: r.x, Y: r.p2})
		p3XY = append(p3XY, plotter.XY{X: r.x, Y: r.p3})
		p4XY = append(p4XY, plotter.XY{X: r.x, Y: r.p4})
		eulerXY = append(eulerXY, plotter.XY{X: r.x, Y: r.eu})
	}
	err := plotutil.AddLines(pl,
		"Picard 1", p1XY,
		"Picard 2", p2XY,
		"Picard 3", p3XY,
		"Picard 4", p4XY,
		"Euler", eulerXY,
	)
	pl.Y.Min = 0
	pl.Y.Max = 20
	if err != nil {
		fmt.Println("Ошибка построения графика:", err)
		return
	}
	if err := pl.Save(8*vg.Inch, 5*vg.Inch, "task3.png"); err != nil {
		fmt.Println("Ошибка сохранения графика:", err)
		return
	}
	fmt.Println("График сохранён: task3.png")
	fmt.Println()
}
