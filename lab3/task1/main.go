package main

import (
	"math"
	"modeling/lab3/plot"
)

func linspace(start, end float64, num int) []float64 {
	res := make([]float64, num)
	if num == 1 {
		res[0] = start
		return res
	}
	step := (end - start) / float64(num-1)
	for i := range res {
		res[i] = start + float64(i)*step
	}
	return res
}

//func galerkin(x []float64) []float64 {
//	c1 := 2331.0 / 34724.0
//	c2 := -28.0 / 8681.0
//	c3 := 1316.0 / 8681.0
//
//	y := make([]float64, len(x))
//	for i, val := range x {
//		val2 := val * val
//		val3 := val2 * val
//		val4 := val3 * val
//		y[i] = val + c1*(val2-2*val) + c2*(val3-3*val) + c3*(val4-4*val)
//	}
//	return y
//}

func galerkin(x []float64) []float64 {
	c1 := -0.49853
	c2 := 0.015273
	c3 := -0.0024697

	y := make([]float64, len(x))
	for i, val := range x {
		y[i] = val + c1*math.Sin(math.Pi/2*val) + c2*math.Sin(3*math.Pi/2*val) + c3*math.Sin(5*math.Pi/2*val)
	}
	return y
}

func progonka(x []float64) []float64 {
	nodes := len(x)
	h := x[1] - x[0]

	An := make([]float64, nodes)
	Bn := make([]float64, nodes)
	Cn := make([]float64, nodes)
	Fn := make([]float64, nodes)

	for i, val := range x {
		An[i] = 1.0 + val*h
		Bn[i] = 2.0 - 2.0*h*h
		Cn[i] = 1.0 - val*h
		Fn[i] = -val * h * h
	}

	eta := make([]float64, nodes)
	ksi := make([]float64, nodes)

	for i := 1; i < nodes; i++ {
		denom := Bn[i-1] - An[i-1]*ksi[i-1]
		ksi[i] = Cn[i-1] / denom
		eta[i] = (An[i-1]*eta[i-1] + Fn[i-1]) / denom
	}

	y := make([]float64, nodes)
	y[nodes-1] = (3*h*h - 2*h - 2*eta[nodes-1]) / (2*ksi[nodes-1] + 2*h*h - 2)

	for i := nodes - 2; i >= 0; i-- {
		y[i] = ksi[i+1]*y[i+1] + eta[i+1]
	}

	return y
}

func estimateOptimalStep(tolerance float64) (int, float64) {
	N := 11

	for {
		x1 := linspace(0, 1, N)
		u1 := progonka(x1)

		N2 := 2*N - 1
		x2 := linspace(0, 1, N2)
		u2 := progonka(x2)

		maxErr := 0.0
		for i := 0; i < N; i++ {
			err := math.Abs(u1[i]-u2[2*i]) / 3.0
			if err > maxErr {
				maxErr = err
			}
		}

		if maxErr <= tolerance {
			h := 1.0 / float64(N2-1)
			return N2, h
		}

		N = N2

		if N > 100000 {
			h := 1.0 / float64(N-1)
			return N, h
		}
	}
}

func main() {
	//tolerance := 1e-3
	//N, _ := estimateOptimalStep(tolerance)
	N := 100
	x := linspace(0, 1, N)

	uPolynom := galerkin(x)
	uProgonka := progonka(x)

	plot.SaveComparisonPlot(
		x,
		uPolynom,
		uProgonka,
		"Метод Галёркина",
		"Прогонка",
		"result.png",
		"x",
		"u(x)",
		"Решение краевой задачи",
	)
}
