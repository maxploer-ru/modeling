package main

import (
	"fmt"
	"log"
	"math"

	"modeling/lab3/plot"
	"modeling/lab3/task2/core"
	"modeling/lab3/task2/raznost"
	"modeling/lab3/task2/shooting"
)

func main() {
	// 1. Чтение конфигурации
	cfg, err := core.ReadConfig("lab3/task2/config.json")
	if err != nil {
		log.Fatalf("Ошибка чтения config.json: %v", err)
	}

	// Базовая модель
	mBase := core.BuildModel(cfg, nil)

	shootingMaxIter := cfg.ShootingMaxIter
	shootingTol := cfg.ShootingTol

	// 2. Метод стрельбы (базовый)
	chiOpt, shotsUsed := shooting.Solve(mBase, 0.00, 1.0, shootingTol, shootingMaxIter)
	fmt.Printf("\nОптимальное chi = %.6f\n", chiOpt)
	fmt.Printf("Количество выстрелов: %d (tol=%.1e)\n", shotsUsed, shootingTol)

	rStrelba, uStrelba, FStrelba := shooting.SolveCauchy(mBase, chiOpt)
	uPStrelba := make([]float64, len(rStrelba))
	for i, rVal := range rStrelba {
		uPStrelba[i] = mBase.UPlanck(mBase.TOfR(rVal))
	}

	plot.PlotSingleDependency(rStrelba, [][]float64{FStrelba}, []string{"F(r)"}, "F(r)", "Радиус r, см", "F(r)", "F_strelba.png")
	plot.PlotSingleDependency(rStrelba, [][]float64{uStrelba}, []string{"u(r)"}, "u(r)", "Радиус r, см", "u(r)", "u_strelba.png")
	plot.PlotSingleDependency(rStrelba, [][]float64{uPStrelba}, []string{"u_p(r)"}, "u_p(r)", "Радиус r, см", "u_p(r)", "up_strelba.png")

	// 3. Исследование параметров (как в Python)
	studyRData := [][]float64{rStrelba}
	studyUData := [][]float64{uStrelba}
	studyFData := [][]float64{FStrelba}
	studyNames := []string{"Базовый (V1)"}

	for name, sp := range cfg.Studies {
		mStudy := core.BuildModel(cfg, &sp)
		chiStudy, _ := shooting.Solve(mStudy, 0.00, 1.00, shootingTol, shootingMaxIter)
		rS, uS, FS := shooting.SolveCauchy(mStudy, chiStudy)

		studyRData = append(studyRData, rS)
		studyUData = append(studyUData, uS)
		studyFData = append(studyFData, FS)
		studyNames = append(studyNames, name)
	}

	// Так как в вашей функции PlotSingleDependency общий X, передадим самую длинную или первую сетку
	// (В идеале для разных X лучше использовать PlotComparisonPlot или переписать плоттер под массивы X)
	plot.PlotSingleDependency(studyRData[0], studyUData, studyNames, "Исследование u(r)", "Радиус r, см", "u(r)", "studies_u.png")
	plot.PlotSingleDependency(studyRData[0], studyFData, studyNames, "Исследование F(r)", "Радиус r, см", "F(r)", "studies_F.png")

	// 4. Конечно-разностный метод (Прогонка)
	uRaznost, rRaznost, _, kRaznost, uPRaznost, h := raznost.Solve(mBase)
	N := mBase.N

	FDiff := make([]float64, N+1)
	FDiff[0] = 0.0
	for i := 1; i < N; i++ {
		FDiff[i] = -1.0 / (3.0 * kRaznost[i]) * (uRaznost[i+1] - uRaznost[i-1]) / (2.0 * h)
	}
	FDiff[N] = -1.0 / (3.0 * kRaznost[N]) * (3.0*uRaznost[N] - 4.0*uRaznost[N-1] + uRaznost[N-2]) / (2.0 * h)

	FInt := make([]float64, N+1)
	integrand := make([]float64, N+1)
	integral := make([]float64, N+1)
	for i := 0; i <= N; i++ {
		integrand[i] = kRaznost[i] * (uPRaznost[i] - uRaznost[i]) * rRaznost[i]
	}
	for i := 1; i <= N; i++ {
		integral[i] = integral[i-1] + (integrand[i-1]+integrand[i])*h/2.0
		FInt[i] = integral[i] / mBase.R
	}

	divF := make([]float64, N+1)
	for i := 0; i <= N; i++ {
		divF[i] = -kRaznost[i] * (uRaznost[i] - uPRaznost[i])
	}

	FR1 := -1.0 / (3.0 * kRaznost[N]) * (uRaznost[N] - uRaznost[N-1]) / h
	FR2 := FDiff[N]
	FR3 := FInt[N]
	FR4 := 0.39 * uRaznost[N]

	fmt.Printf("\nСравнение F(R) четырьмя способами:\n")
	fmt.Printf("1) Односторонняя разность 1-го порядка: %e\n", FR1)
	fmt.Printf("2) Формула 2-го порядка:                %e\n", FR2)
	fmt.Printf("3) Интегрирование:                      %e\n", FR3)
	fmt.Printf("4) Из краевого условия:                 %e\n", FR4)
	fmt.Printf("5) Метод стрельбы:                      %e\n", FStrelba[N])

	plot.PlotSingleDependency(rStrelba, [][]float64{FDiff, FStrelba}, []string{"Прогонка F(r)", "Стрельба F(r)"}, "F(r)", "Радиус r, см", "F(r)", "F_compare.png")
	plot.PlotSingleDependency(rStrelba, [][]float64{uRaznost, uStrelba}, []string{"Прогонка u(r)", "Стрельба u(r)"}, "u(r)", "Радиус r, см", "u(r)", "u_compare.png")
	plot.PlotSingleDependency(rStrelba, [][]float64{uPRaznost, uPStrelba}, []string{"Прогонка u_p(r)", "Стрельба u_p(r)"}, "u_p(r)", "Радиус r, см", "u_p(r)", "up_compare.png")
	plot.PlotSingleDependency(rStrelba, [][]float64{divF}, []string{"Прогонка divF(r)"}, "divF(r)", "Радиус r, см", "divF(r)", "divF_progonka.png")

	absF := make([]float64, N+1)
	absU := make([]float64, N+1)
	for i := 0; i <= N; i++ {
		absF[i] = math.Abs(FDiff[i] - FStrelba[i])
		absU[i] = math.Abs(uRaznost[i] - uStrelba[i])
	}
	plot.PlotSingleDependency(rStrelba, [][]float64{absF}, []string{"Абс. разница F(r)"}, "F(r)", "Радиус r, см", "F(r)", "F_diff.png")
	plot.PlotSingleDependency(rStrelba, [][]float64{absU}, []string{"Абс. разница u(r)"}, "u(r)", "Радиус r, см", "u(r)", "u_diff.png")

	// 5. Двусторонняя стрельба (итерации по точке пересечения)
	rTwo, uTwo, FTwo, itTwo := shooting.SolveTwoSided(mBase, 1.0, 1.0, shootingTol, shootingMaxIter)
	fmt.Printf("\nДвусторонняя стрельба: итераций = %d\n", itTwo)
	plot.PlotSingleDependency(rTwo, [][]float64{uTwo}, []string{"Двусторонняя u(r)"}, "u(r)", "Радиус r, см", "u(r)", "u_twosided.png")
	plot.PlotSingleDependency(rTwo, [][]float64{FTwo}, []string{"Двусторонняя F(r)"}, "F(r)", "Радиус r, см", "F(r)", "F_twosided.png")
}
