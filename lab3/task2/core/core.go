package core

import (
	"encoding/json"
	"math"
	"os"
)

type KTable struct {
	T        []float64 `json:"T"`
	Variant1 []float64 `json:"variant_1"`
	Variant2 []float64 `json:"variant_2"`
}

type StudyParams struct {
	R      float64 `json:"R"`
	Tw     float64 `json:"T_w"`
	T0     float64 `json:"T_0"`
	P      float64 `json:"p"`
	KTable KTable  `json:"k_table"`
}

type Config struct {
	N               int                    `json:"N"`
	Alpha           float64                `json:"alpha"`
	R               float64                `json:"R"`
	Tw              float64                `json:"T_w"`
	T0              float64                `json:"T_0"`
	P               float64                `json:"p"`
	Variant         int                    `json:"variant"`
	KTable          KTable                 `json:"k_table"`
	ShootingMaxIter int                    `json:"shooting_max_iter"`
	ShootingTol     float64                `json:"shooting_tol"`
	Studies         map[string]StudyParams `json:"studies"`
}

type Model struct {
	N       int
	AlphaRK float64
	R       float64
	Tw      float64
	T0      float64
	P       float64
	Variant int
	LogT    []float64
	LogK    []float64
}

func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func BuildModel(cfg *Config, sp *StudyParams) *Model {
	m := &Model{
		N:       cfg.N,
		AlphaRK: cfg.Alpha,
		Variant: cfg.Variant,
	}

	var kTableSource KTable
	if sp != nil {
		m.R, m.Tw, m.T0, m.P = sp.R, sp.Tw, sp.T0, sp.P
		kTableSource = sp.KTable
	} else {
		m.R, m.Tw, m.T0, m.P = cfg.R, cfg.Tw, cfg.T0, cfg.P
		kTableSource = cfg.KTable
	}

	TTable := kTableSource.T
	var kTable []float64
	if m.Variant == 1 {
		kTable = kTableSource.Variant1
	} else {
		kTable = kTableSource.Variant2
	}

	m.LogT = make([]float64, len(TTable))
	m.LogK = make([]float64, len(kTable))
	for i := range TTable {
		m.LogT[i] = math.Log(TTable[i])
		m.LogK[i] = math.Log(kTable[i])
	}
	return m
}

func (m *Model) KOfT(T float64) float64 {
	return math.Exp(m.interp1dExtrapolate(m.LogT, m.LogK, math.Log(T)))
}

func (m *Model) TOfR(r float64) float64 {
	z := r / m.R
	return (m.Tw-m.T0)*math.Pow(z, m.P) + m.T0
}

func (m *Model) UPlanck(T float64) float64 {
	expp := 4.799e4 / T
	return 3.084e-4 / (math.Exp(expp) - 1.0)
}

func (m *Model) EvalDerivatives(r, u, F float64) (float64, float64) {
	if r < 1e-12 {
		T0 := m.TOfR(0.0)
		k0 := m.KOfT(T0)
		up0 := m.UPlanck(T0)
		du := -3.0 * k0 * F
		dF := -k0 * (u - up0) / 2.0
		return du, dF
	}

	T := m.TOfR(r)
	kr := m.KOfT(T)
	up := m.UPlanck(T)
	du := -3.0 * kr * F
	dF := -kr*(u-up) - (F / r)
	return du, dF
}

func (m *Model) interp1dExtrapolate(xVals, yVals []float64, x float64) float64 {
	n := len(xVals)
	if x <= xVals[0] {
		slope := (yVals[1] - yVals[0]) / (xVals[1] - xVals[0])
		return yVals[0] + slope*(x-xVals[0])
	}
	if x >= xVals[n-1] {
		slope := (yVals[n-1] - yVals[n-2]) / (xVals[n-1] - xVals[n-2])
		return yVals[n-1] + slope*(x-xVals[n-1])
	}
	for i := 0; i < n-1; i++ {
		if x >= xVals[i] && x <= xVals[i+1] {
			slope := (yVals[i+1] - yVals[i]) / (xVals[i+1] - xVals[i])
			return yVals[i] + slope*(x-xVals[i])
		}
	}
	return 0
}
