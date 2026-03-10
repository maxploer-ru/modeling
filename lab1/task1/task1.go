package task1

import (
	"fmt"
	"math"
	"modeling/lab1/methods"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

//func seriesCoefficients(n int) []float64 {
//	a := make([]float64, n)
//	a[0] = 1.0
//	a[1] = -0.5
//	for k := 1; k < n-1; k++ {
//		a[k+1] = -a[k] / (float64(k+1) * float64(4*k+2))
//	}
//	return a
//}

func evalSeries(h float64, n int) []float64 {
	y := make([]float64, n+1)
	x := 0.0
	for i := 0; i <= n; i++ {
		y[i] = 1 - x/2 + math.Pow(x, 2)/24.0 - math.Pow(x, 3)/720.0
		x += h
	}
	return y
}

func EulerSystem(h float64, n int) []float64 {
	f2 := func(x, u, v float64) float64 {
		if x == 0.0 {
			return 1.0 / 12.0
		}
		return -(2*v + u) / (4 * x)
	}
	x := make([]float64, n+1)
	u0 := make([]float64, n+1)
	v0 := make([]float64, n+1)
	x[0] = 0.0
	u0[0] = 1.0
	v0[0] = -0.5
	for i := 0; i < n; i++ {
		u2 := f2(x[i], u0[i], v0[i])
		x[i+1] = x[i] + h
		u0[i+1] = u0[i] + h*v0[i]
		v0[i+1] = v0[i] + h*u2
	}
	return u0
}

func Run() {
	x0, x1 := 0.0, 10.0
	n := 100
	h := (x1 - x0) / float64(n)

	y := evalSeries(h, n)
	yEuler := EulerSystem(h, n)

	p := plot.New()
	p.Title.Text = "Задача 1: 4xu'' + 2u' + u = 0"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "u(x)"

	seriesXY := make(plotter.XYs, n+1)
	for i := 0; i <= n; i++ {
		seriesXY[i].X = x0 + float64(i)*h
		seriesXY[i].Y = y[i]
	}

	eulerXY := make(plotter.XYs, n+1)
	for i := 0; i <= n; i++ {
		eulerXY[i].X = x0 + float64(i)*h
		eulerXY[i].Y = yEuler[i]
	}

	err := plotutil.AddLinePoints(p, "Степенной ряд", seriesXY, "Метод Эйлера", eulerXY)
	if err != nil {
		fmt.Println("Ошибка построения графика:", err)
		return
	}
	methods.StyleCartesian(p)
	if err := p.Save(8*vg.Inch, 5*vg.Inch, "task1.png"); err != nil {
		fmt.Println("Ошибка сохранения графика:", err)
		return
	}
	fmt.Println("График сохранён: task1.png")
	fmt.Println()
}
