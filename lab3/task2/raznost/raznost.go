package raznost

import "modeling/lab3/task2/core"

func Solve(m *core.Model) (u []float64, r, T, k, up []float64, h float64) {
	N := m.N
	h = m.R / float64(N)

	r = make([]float64, N+1)
	T = make([]float64, N+1)
	k = make([]float64, N+1)
	up = make([]float64, N+1)

	for i := 0; i <= N; i++ {
		r[i] = float64(i) * h
		T[i] = m.TOfR(r[i])
		k[i] = m.KOfT(T[i])
		up[i] = m.UPlanck(T[i])
	}

	kappa := make([]float64, N)
	for i := 0; i < N; i++ {
		kappa[i] = 2.0 / (3.0 * (k[i] + k[i+1]))
	}

	rHalf := make([]float64, N)
	for i := 0; i < N; i++ {
		rHalf[i] = r[i] + h/2.0
	}

	V := make([]float64, N+1)
	for i := 1; i < N; i++ {
		V[i] = (rHalf[i]*rHalf[i] - rHalf[i-1]*rHalf[i-1]) / 2.0
	}
	V[0] = rHalf[0] * rHalf[0] / 2.0
	V[N] = (m.R*m.R - rHalf[N-1]*rHalf[N-1]) / 2.0

	A, B, C, F := make([]float64, N+1), make([]float64, N+1), make([]float64, N+1), make([]float64, N+1)

	for i := 1; i < N; i++ {
		A[i] = rHalf[i-1] * kappa[i-1] / h
		C[i] = rHalf[i] * kappa[i] / h
		B[i] = A[i] + C[i] + k[i]*V[i]
		F[i] = k[i] * up[i] * V[i]
	}

	A[0] = 0.0
	C[0] = rHalf[0] * kappa[0] / h
	B[0] = C[0] + k[0]*V[0]
	F[0] = k[0] * up[0] * V[0]

	A[N] = rHalf[N-1] * kappa[N-1] / h
	C[N] = 0.0
	B[N] = A[N] + 0.39*m.R + k[N]*V[N]
	F[N] = k[N] * up[N] * V[N]

	xi, eta := make([]float64, N+1), make([]float64, N+1)

	xi[0] = C[0] / B[0]
	eta[0] = F[0] / B[0]
	for i := 1; i <= N; i++ {
		denom := B[i] - A[i]*xi[i-1]
		xi[i] = C[i] / denom
		eta[i] = (F[i] + A[i]*eta[i-1]) / denom
	}

	u = make([]float64, N+1)
	u[N] = eta[N]
	for i := N - 1; i >= 0; i-- {
		u[i] = xi[i]*u[i+1] + eta[i]
	}
	return u, r, T, k, up, h
}
