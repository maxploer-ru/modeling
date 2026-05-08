package shooting

import (
	"math"

	"modeling/lab3/task2/core"
	"modeling/lab3/task2/rk4"
)

func integrateSegment(m *core.Model, r0, u0, F0, h, tEnd float64) ([]float64, []float64, []float64) {
	systemFunc := rk4.CreateOdeSystem(m)
	solver := rk4.RK4Solver{Func: systemFunc, Y0: []float64{u0, F0}, H: h, T0: r0}
	rHist, yHist := solver.Solve(tEnd)

	rArr := make([]float64, len(rHist))
	uArr := make([]float64, len(rHist))
	FArr := make([]float64, len(rHist))
	for i := range rHist {
		rArr[i] = rHist[i]
		uArr[i] = yHist[i][0]
		FArr[i] = yHist[i][1]
	}
	return rArr, uArr, FArr
}

func SolveCauchy(m *core.Model, chi float64) ([]float64, []float64, []float64) {
	N := m.N
	h := m.R / float64(N)

	T0 := m.TOfR(0.0)
	up0 := m.UPlanck(T0)
	u0 := chi * up0

	return integrateSegment(m, 0.0, u0, 0.0, h, m.R)
}

func ShootingResidual(m *core.Model, chi float64) float64 {
	_, u, F := SolveCauchy(m, chi)
	last := len(u) - 1
	return F[last] - 0.39*u[last]
}

func Solve(m *core.Model, chi1, chi2, tol float64, maxIter int) (float64, int) {
	r1 := ShootingResidual(m, chi1)
	r2 := ShootingResidual(m, chi2)
	chiPrev := chi2

	for i := 0; i < maxIter; i++ {
		if math.Abs(r2-r1) < 1e-20 {
			return chiPrev, i + 1
		}
		chiNew := chi2 - r2*(chi2-chi1)/(r2-r1)

		if math.Abs(chiNew-chiPrev) <= tol*math.Max(1.0, math.Abs(chiNew)) {
			return chiNew, i + 1
		}

		chi1, r1 = chi2, r2
		chi2 = chiNew
		r2 = ShootingResidual(m, chiNew)
		chiPrev = chiNew
	}

	return chiPrev, maxIter
}

func SolveTwoSided(m *core.Model, chiInit, betaInit, tol float64, maxIter int) ([]float64, []float64, []float64, int) {
	N := m.N
	mid := N / 2
	h := m.R / float64(N)
	if chiInit == 0 {
		chiInit = 1.0
	}
	if betaInit == 0 {
		betaInit = 1.0
	}

	T0 := m.TOfR(0.0)
	up0 := m.UPlanck(T0)
	chi := chiInit
	beta := betaInit

	midR := h * float64(mid)

	for iter := 0; iter < maxIter; iter++ {
		u0 := chi * up0
		_, uF, FF := integrateSegment(m, 0.0, u0, 0.0, h, midR)
		uMidF := uF[len(uF)-1]
		FMidF := FF[len(FF)-1]

		uR := beta / 0.39
		_, uB, FB := integrateSegment(m, m.R, uR, beta, -h, midR)
		uMidB := uB[0]
		FMidB := FB[0]

		diffU := uMidB - uMidF
		diffF := FMidF - FMidB
		maxU := math.Max(1.0, math.Abs(uMidF))
		maxF := math.Max(1.0, math.Abs(FMidF))

		if tol > 0 && math.Abs(diffU) <= tol*maxU && math.Abs(diffF) <= tol*maxF {
			rF, uFfull, FFull := integrateSegment(m, 0.0, u0, 0.0, h, midR)
			rB, uBfull, FBfull := integrateSegment(m, m.R, uR, beta, -h, midR)

			rArr := make([]float64, 0, N+1)
			uArr := make([]float64, 0, N+1)
			FArr := make([]float64, 0, N+1)

			rArr = append(rArr, rF...)
			uArr = append(uArr, uFfull...)
			FArr = append(FArr, FFull...)

			for i := 1; i < len(rB); i++ {
				rArr = append(rArr, rB[i])
				uArr = append(uArr, uBfull[i])
				FArr = append(FArr, FBfull[i])
			}
			return rArr, uArr, FArr, iter + 1
		}

		if uMidF != 0.0 {
			chi = chi * (uMidB / uMidF)
		}
		if FMidB != 0.0 {
			beta = beta * (FMidF / FMidB)
		}
	}

	rArr, uArr, FArr := SolveCauchy(m, chi)
	return rArr, uArr, FArr, maxIter
}
