package rk4

import (
	"math"
	"modeling/lab3/task2/core"
)

type RK4Solver struct {
	Func func(t float64, y []float64) []float64
	Y0   []float64
	H    float64
	T0   float64
}

func (s *RK4Solver) Solve(tEnd float64) ([]float64, [][]float64) {
	direction := 1.0
	if tEnd < s.T0 {
		direction = -1.0
	}
	h := direction * math.Abs(s.H)
	t := s.T0
	y := make([]float64, len(s.Y0))
	copy(y, s.Y0)

	var tHist []float64
	var yHist [][]float64
	tHist = append(tHist, t)
	yHist = append(yHist, append([]float64(nil), y...))

	for (tEnd-t)*direction > 1e-10 {
		k1 := s.Func(t, y)

		y1 := make([]float64, len(y))
		for i := range y {
			y1[i] = y[i] + h/2.0*k1[i]
		}
		k2 := s.Func(t+h/2.0, y1)

		y2 := make([]float64, len(y))
		for i := range y {
			y2[i] = y[i] + h/2.0*k2[i]
		}
		k3 := s.Func(t+h/2.0, y2)

		y3 := make([]float64, len(y))
		for i := range y {
			y3[i] = y[i] + h*k3[i]
		}
		k4 := s.Func(t+h, y3)

		for i := range y {
			y[i] += (h / 6.0) * (k1[i] + 2.0*k2[i] + 2.0*k3[i] + k4[i])
		}
		t += h

		tHist = append(tHist, t)
		yHist = append(yHist, append([]float64(nil), y...))
	}

	if direction == -1 {
		for i, j := 0, len(tHist)-1; i < j; i, j = i+1, j-1 {
			tHist[i], tHist[j] = tHist[j], tHist[i]
			yHist[i], yHist[j] = yHist[j], yHist[i]
		}
	}
	return tHist, yHist
}

func RootScalarSecant(f func(float64) float64, x0, x1, tol float64, maxIter int) float64 {
	for i := 0; i < maxIter; i++ {
		f1 := f(x1)
		f0 := f(x0)
		if math.Abs(f1-f0) < 1e-20 {
			break
		}
		xNew := x1 - f1*(x1-x0)/(f1-f0)
		if math.Abs((xNew-x1)/xNew) < tol {
			return xNew
		}
		x0, x1 = x1, xNew
	}
	return x1
}

func CreateOdeSystem(m *core.Model) func(r float64, y []float64) []float64 {
	return func(r float64, y []float64) []float64 {
		u, F := y[0], y[1]
		du, dF := m.EvalDerivatives(r, u, F)
		return []float64{du, dF}
	}
}
