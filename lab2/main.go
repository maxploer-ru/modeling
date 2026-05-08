package main

import (
	"fmt"
	"math"
)

var (
	T0Table            []float64
	ITable, mTable     []float64
	TTable, sigmaTable []float64

	R  = 0.35
	lP = 12.0
	LK = 187e-6
	CK = 268e-6
	RK = 0.25
	U0 = 1400.0
	I0 = 0.5
	TW = 2000.0
)

func binSearch(arr []float64, target float64) int {
	low, high := -1, len(arr)
	for low+1 < high {
		mid := (low + high) / 2
		if arr[mid] < target {
			low = mid
		} else {
			high = mid
		}
	}
	return high
}

func newtonCoeffs(x, y []float64) []float64 {
	n := len(x)
	a := make([]float64, n)
	copy(a, y)
	for j := 1; j < n; j++ {
		for i := n - 1; i >= j; i-- {
			a[i] = (a[i] - a[i-1]) / (x[i] - x[i-j])
		}
	}
	return a
}

func newtonEval(coeffs, x []float64, target float64) float64 {
	res := 0.0
	for i := len(coeffs) - 1; i >= 0; i-- {
		res = coeffs[i] + (target-x[i])*res
	}
	return res
}

func newtonInterp(xTab, yTab []float64, target float64, degree int) float64 {
	if degree > len(xTab) {
		degree = len(xTab)
	}
	idx := binSearch(xTab, target)
	left := idx - degree/2
	right := left + degree
	if left < 0 {
		left = 0
		right = degree
	}
	if right > len(xTab) {
		right = len(xTab)
		left = right - degree
	}
	xSub := xTab[left:right]
	ySub := yTab[left:right]
	coeffs := newtonCoeffs(xSub, ySub)
	return newtonEval(coeffs, xSub, target)
}

func T0FromI(I float64) float64 {
	return newtonInterp(ITable, T0Table, I, 4)
}

func mFromI(I float64) float64 {
	return newtonInterp(ITable, mTable, I, 4)
}

func sigmaFromT(T float64) float64 {
	return newtonInterp(TTable, sigmaTable, T, 4)
}

func TZ(I, z float64) float64 {
	T0 := T0FromI(I)
	m := mFromI(I)
	return T0 + (TW-T0)*math.Pow(z, m)
}

func trapezoid(x, y []float64) float64 {
	sum := 0.0
	for i := 1; i < len(x); i++ {
		h := x[i] - x[i-1]
		sum += (y[i] + y[i-1]) * h / 2.0
	}
	return sum
}

func sigmaIntegral(I float64, steps int) float64 {
	z := make([]float64, steps+1)
	f := make([]float64, steps+1)
	dz := 1.0 / float64(steps)
	for i := 0; i <= steps; i++ {
		zi := float64(i) * dz
		z[i] = zi
		T := TZ(I, zi)
		f[i] = sigmaFromT(T) * zi
	}
	return trapezoid(z, f)
}

func Rp(I float64) float64 {
	integral := sigmaIntegral(I, 100)
	return lP / (2.0 * math.Pi * R * R * integral)
}

func dIdt(U, I, RExtra float64, useRp bool) float64 {
	var RSum float64
	if useRp {
		RSum = RK + Rp(I)
	} else {
		RSum = RExtra
	}
	return (U - RSum*I) / LK
}

func dUdt(I float64) float64 {
	return -I / CK
}

func rk2Step(t, U, I, h, RExtra float64, useRp bool) (float64, float64, float64) {
	alpha := 0.5
	k1 := dIdt(U, I, RExtra, useRp)
	q1 := dUdt(I)

	Umid := U + q1*(h/(2.0*alpha))
	Imid := I + k1*(h/(2.0*alpha))

	k2 := dIdt(Umid, Imid, RExtra, useRp)
	q2 := dUdt(Imid)

	Inew := I + h*((1.0-alpha)*k1+alpha*k2)
	Unew := U + h*((1.0-alpha)*q1+alpha*q2)
	tnew := t + h

	return tnew, Unew, Inew
}

func solve(h float64, tEnd float64, RExtra float64, useRp bool) ([]float64, []float64, []float64) {
	n := int(math.Ceil(tEnd/h)) + 1
	t := make([]float64, n)
	U := make([]float64, n)
	I := make([]float64, n)
	t[0] = 0.0
	U[0] = U0
	I[0] = I0
	for i := 0; i < n-1; i++ {
		t[i+1], U[i+1], I[i+1] = rk2Step(t[i], U[i], I[i], h, RExtra, useRp)
	}
	return t, U, I
}

func solveBackward(tStart, tEnd, h float64, RExtra float64, useRp bool, UInit, IInit float64) ([]float64, []float64, []float64) {
	if h > 0 {
		h = -h
	}
	n := int(math.Ceil(math.Abs(tEnd-tStart)/math.Abs(h))) + 1
	t := make([]float64, n)
	U := make([]float64, n)
	I := make([]float64, n)

	t[0] = tStart
	U[0] = UInit
	I[0] = IInit

	for i := 0; i < n-1; i++ {
		t[i+1], U[i+1], I[i+1] = rk2Step(t[i], U[i], I[i], h, RExtra, useRp)
	}
	return t, U, I
}

func pulseDuration(t, I []float64) float64 {
	Imax := I[0]
	for _, v := range I {
		if v > Imax {
			Imax = v
		}
	}
	thresh := 0.35 * Imax
	var t1, t2 float64
	for i := 0; i < len(I); i++ {
		if I[i] >= thresh {
			if i == 0 {
				t1 = t[i]
			} else {
				frac := (thresh - I[i-1]) / (I[i] - I[i-1])
				t1 = t[i-1] + frac*(t[i]-t[i-1])
			}
			break
		}
	}
	for i := len(I) - 1; i >= 0; i-- {
		if I[i] >= thresh {
			if i == len(I)-1 {
				t2 = t[i]
			} else {
				frac := (thresh - I[i]) / (I[i+1] - I[i])
				t2 = t[i] + frac*(t[i+1]-t[i])
			}
			break
		}
	}
	return t2 - t1
}

func estimateStep(tEnd, RExtra float64, useRp bool) float64 {
	h := 1e-1
	tolerance := 1e-3
	for {
		_, U1, I1 := solve(h, tEnd, RExtra, useRp)
		_, U2, I2 := solve(h/2, tEnd, RExtra, useRp)
		errI := math.Abs(I1[len(I1)-1]-I2[len(I2)-1]) / math.Abs(I1[len(I1)-1])
		errU := math.Abs(U1[len(U1)-1]-U2[len(U2)-1]) / math.Abs(U1[len(U1)-1])
		if errI < tolerance && errU < tolerance {
			return h
		}
		h /= 2
		if h < 1e-8 {
			return h
		}
	}
}

func main() {
	ITable, T0Table, mTable = readTable1("lab2/table1.txt")
	TTable, sigmaTable = readTable2("lab2/table2.txt")

	h := 1e-6
	tEnd := 600e-6
	fmt.Println(estimateStep(tEnd, 0.0, true))

	t, U, I := solve(h, tEnd, 0.0, true)
	//t, U, I := solveBackward(tEnd, 0.0, h, 0.0, true, 107.38206305692975, 305.8902099384868)
	//fmt.Println("I(tEnd) =", I[len(I)-1])
	//fmt.Println("U(tEnd) =", U[len(U)-1])
	//I(tEnd) = 0.5334032574671053
	//U(tEnd) = 1400.0000146796363

	fmt.Printf("Шаг сетки h = %.0e с\n", h)

	saveLinePlot(t, I, "I_t.png", "t, с", "I, А", "Ток I(t)")

	saveLinePlot(t, U, "U_t.png", "t, с", "U, В", "Напряжение U(t)")

	RpVals := make([]float64, len(t))
	IrpVals := make([]float64, len(t))
	T0Vals := make([]float64, len(t))
	for i, curI := range I {
		RpVals[i] = Rp(curI)
		IrpVals[i] = curI * RpVals[i]
		T0Vals[i] = T0FromI(curI)
	}
	saveLinePlot(t, RpVals, "Rp_t.png", "t, с", "Rp, Ом", "Сопротивление Rp(t)")
	saveLinePlot(t, IrpVals, "I_Rp_t.png", "t, с", "I*Rp, В", "Произведение I(t)*Rp(t)")
	saveLinePlot(t, T0Vals, "T0_t.png", "t, с", "T0, K", "Температура T0(t)")

	t2, _, I2 := solve(1e-6, 2000e-6, 0.0, false)
	saveLinePlot(t2, I2, "I_t_R0.png", "t, с", "I, А", "Ток при Rk+Rp=0")

	t3, _, I3 := solve(1e-8, 20e-6, 200.0, false)
	saveLinePlot(t3, I3, "I_t_R200.png", "t, с", "I, А", "Ток при R=200 Ом")

	CVals := []float64{150e-6, 180e-6, 210e-6, 240e-6, 270e-6, 300e-6, 330e-6, 360e-6, 390e-6, 420e-6}
	var durC []float64
	origC := CK
	for _, c := range CVals {
		CK = c
		_, _, Icur := solve(h, tEnd, 0.0, true)
		durC = append(durC, pulseDuration(t, Icur))
	}
	CK = origC

	LVals := []float64{50e-6, 80e-6, 110e-6, 140e-6, 170e-6, 200e-6, 230e-6, 260e-6, 290e-6, 320e-6}
	var durL []float64
	origL := LK
	for _, l := range LVals {
		LK = l
		_, _, Icur := solve(h, tEnd, 0.0, true)
		durL = append(durL, pulseDuration(t, Icur))
	}
	LK = origL

	RVals := []float64{0.10, 0.15, 0.20, 0.25, 0.30, 0.35, 0.40, 0.45, 0.50, 0.55}
	var durR []float64
	origR := RK
	for _, r := range RVals {
		RK = r
		_, _, Icur := solve(h, tEnd, 0.0, true)
		durR = append(durR, pulseDuration(t, Icur))
	}
	RK = origR

	saveLinePlot(CVals, durC, "C_duration.png", "Ck, Ф", "tимп, с", "Влияние ёмкости на длительность")
	saveLinePlot(LVals, durL, "L_duration.png", "Lk, Гн", "tимп, с", "Влияние индуктивности на длительность")
	saveLinePlot(RVals, durR, "R_duration.png", "Rk, Ом", "tимп, с", "Влияние сопротивления на длительность")
}
