package methods

// Point — точка (x, u)
type Point struct {
	X, U float64
}

// Euler решает u' = f(x, u) методом Эйлера.
// Возвращает массив точек от x0 до x0+n*h.
func Euler(f func(x, u float64) float64, x0, u0, h float64, n int) []Point {
	pts := make([]Point, n+1)
	pts[0] = Point{x0, u0}
	x, u := x0, u0
	for i := 1; i <= n; i++ {
		u = u + h*f(x, u)
		x = x0 + float64(i)*h
		pts[i] = Point{x, u}
	}
	return pts
}
