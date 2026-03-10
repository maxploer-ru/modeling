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

/*
Задача 3.
	u'(x) = x² + u²,  u(0) = 0

Метод Пикара (4 приближения):
	u_0(x) = 0
	u_{n+1}(x) = ∫[0..x] (t² + u_n(t)²) dt

	u_1(x) = x³/3
	u_2(x) = x³/3 + x⁷/63
	u_3(x) = x³/3 + x⁷/63 + 2x¹¹/2079 + x¹⁵/59535
	u_4(x) — вычисляем числено через интегрирование.

Решение имеет вертикальную асимптоту (blow-up) при x ≈ 2.003.
x_max определяем из условия: относительное различие Euler(h) vs Euler(h/2)
по правилу Рунге не превышает 10^{-4} (начиная с области |u| > threshold).
*/

// picard1 — первое приближение Пикара: x³/3
func picard1(x float64) float64 {
	return x * x * x / 3.0
}

// picard2 — второе приближение: x³/3 + x⁷/63
func picard2(x float64) float64 {
	x3 := x * x * x
	x7 := x3 * x * x * x * x
	return x3/3.0 + x7/63.0
}

// picard3 — третье приближение: x³/3 + x⁷/63 + 2x¹¹/2079 + x¹⁵/59535
func picard3(x float64) float64 {
	x2 := x * x
	x3 := x2 * x
	x7 := x3 * x2 * x2
	x11 := x7 * x2 * x2
	x15 := x11 * x2 * x2
	return x3/3.0 + x7/63.0 + 2.0*x11/2079.0 + x15/59535.0
}

// picard4Numerical вычисляет 4-е приближение Пикара числено.
// u_4(x) = ∫[0..x] (t² + u_3(t)²) dt
func picard4Numerical(xs []float64) []float64 {
	n := len(xs)
	vals := make([]float64, n)
	vals[0] = 0
	for i := 1; i < n; i++ {
		dt := xs[i] - xs[i-1]
		t0, t1 := xs[i-1], xs[i]
		p3_0 := picard3(t0)
		p3_1 := picard3(t1)
		f0 := t0*t0 + p3_0*p3_0
		f1 := t1*t1 + p3_1*p3_1
		vals[i] = vals[i-1] + dt*(f0+f1)/2.0
	}
	return vals
}

// findXMax определяет x_max — максимальное x, при котором метод Эйлера
// с шагом h обеспечивает относительную точность eps.
// Эталон: Эйлер с шагом h/2, оценка по правилу Рунге.
// Пропускает начальную область, где |u| < minU.
func findXMax(h float64, eps float64) float64 {
	f := func(x, u float64) float64 { return x*x + u*u }
	xLimit := 2.0

	nH := int(xLimit / h)
	eulerH := methods.Euler(f, 0, 0, h, nH)

	h2 := h / 2.0
	nH2 := int(xLimit / h2)
	eulerH2 := methods.Euler(f, 0, 0, h2, nH2)

	xMax := 0.0
	started := false
	for i := 1; i < len(eulerH); i++ {
		refIdx := i * 2
		if refIdx >= len(eulerH2) {
			break
		}
		uH := eulerH[i].U
		uH2 := eulerH2[refIdx].U

		if math.Abs(uH2) < 0.01 {
			xMax = eulerH[i].X
			continue
		}
		started = true

		relErr := math.Abs(uH-uH2) / math.Abs(uH2)
		if relErr > eps {
			break
		}
		xMax = eulerH[i].X
	}
	if !started {
		return xLimit
	}
	return xMax
}

// Run выполняет задачу 3.
func Run() {
	fmt.Println("=== ЗАДАЧА 3 ===")
	fmt.Println("u'(x) = x² + u²,  u(0) = 0")
	fmt.Println()

	eps := 1e-4

	// Определяем x_max для разных шагов
	fmt.Println("Определение x_max (правило Рунге, относительная точность 1e-4):")
	for _, hh := range []float64{0.1, 0.05, 0.01, 0.001} {
		xm := findXMax(hh, eps)
		fmt.Printf("  h = %.4f → x_max = %.4f\n", hh, xm)
	}
	fmt.Println()

	hEuler := 0.01
	xMax := findXMax(hEuler, eps)
	fmt.Printf("Выбран шаг h = %g, x_max = %.4f (Эйлер точен до этого x)\n", hEuler, xMax)

	// Для таблицы и графика показываем до x=2.0, чтобы видеть расхождение
	xTable := 2.0
	fmt.Printf("Таблица строится до x = %.1f для наглядности расхождения методов.\n", xTable)
	xMax = xTable
	fmt.Println()

	// Сетка
	n := int(math.Round(xMax / hEuler))
	if n < 1 {
		n = 1
	}

	xs := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		xs[i] = float64(i) * hEuler
	}

	// Пикар 4 (численно, на мелкой сетке)
	hFine := 0.0001
	nFine := int(xMax / hFine)
	xsFine := make([]float64, nFine+1)
	for i := 0; i <= nFine; i++ {
		xsFine[i] = float64(i) * hFine
	}
	p4Fine := picard4Numerical(xsFine)

	picard4At := func(x float64) float64 {
		idx := int(math.Round(x / hFine))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(p4Fine) {
			idx = len(p4Fine) - 1
		}
		return p4Fine[idx]
	}

	// Эйлер
	f := func(x, u float64) float64 { return x*x + u*u }
	eulerPts := methods.Euler(f, 0, 0, hEuler, n)

	// Таблица
	fmt.Println("Таблица значений u(x):")
	header := []string{"x", "Пикар 1", "Пикар 2", "Пикар 3", "Пикар 4", "Эйлер"}
	var rows [][]string
	dispStep := n / 25
	if dispStep < 1 {
		dispStep = 1
	}
	for i := 0; i <= n; i += dispStep {
		x := xs[i]
		rows = append(rows, []string{
			fmt.Sprintf("%.4f", x),
			fmt.Sprintf("%.8f", picard1(x)),
			fmt.Sprintf("%.8f", picard2(x)),
			fmt.Sprintf("%.8f", picard3(x)),
			fmt.Sprintf("%.8f", picard4At(x)),
			fmt.Sprintf("%.8f", eulerPts[i].U),
		})
	}
	if n%dispStep != 0 {
		x := xs[n]
		rows = append(rows, []string{
			fmt.Sprintf("%.4f", x),
			fmt.Sprintf("%.8f", picard1(x)),
			fmt.Sprintf("%.8f", picard2(x)),
			fmt.Sprintf("%.8f", picard3(x)),
			fmt.Sprintf("%.8f", picard4At(x)),
			fmt.Sprintf("%.8f", eulerPts[n].U),
		})
	}
	methods.PrintTable(header, rows)
	fmt.Println()

	// График
	pl := plot.New()
	pl.Title.Text = "Задача 3: u' = x² + u²"
	pl.X.Label.Text = "x"
	pl.Y.Label.Text = "u(x)"

	var p1XY, p2XY, p3XY, p4XY, eulerXY plotter.XYs
	for i := 0; i <= n; i++ {
		x := xs[i]
		p1XY = append(p1XY, plotter.XY{X: x, Y: picard1(x)})
		p2XY = append(p2XY, plotter.XY{X: x, Y: picard2(x)})
		p3XY = append(p3XY, plotter.XY{X: x, Y: picard3(x)})
		p4XY = append(p4XY, plotter.XY{X: x, Y: picard4At(x)})
		eulerXY = append(eulerXY, plotter.XY{X: x, Y: eulerPts[i].U})
	}

	err := plotutil.AddLinePoints(pl,
		"Picard 1", p1XY,
		"Picard 2", p2XY,
		"Picard 3", p3XY,
		"Picard 4", p4XY,
		"Euler", eulerXY,
	)
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
