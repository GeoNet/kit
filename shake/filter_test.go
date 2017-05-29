package shake

import (
	"math"
	"testing"
)

func TestFilter(t *testing.T) {

	var TestSlice = []struct {
		i int32
		a float64
		v float64
	}{
		{-532, 0.00000e+00, 0.00000e+00},
		{-530, 4.57240e-06, 4.46712e-08},
		{-535, -7.06915e-06, 1.82214e-08},
		{-530, 4.68737e-06, -5.88713e-09},
		{-534, -4.67327e-06, -5.47830e-09},
		{-539, -1.58891e-05, -2.06115e-07},
		{-538, -1.28712e-05, -4.77604e-07},
		{-531, 3.72493e-06, -5.44966e-07},
		{-530, 5.83960e-06, -4.26428e-07},
		{-537, -1.04327e-05, -4.51664e-07},
		{-535, -5.37988e-06, -5.85350e-07},
		{-537, -9.70453e-06, -7.05766e-07},
		{-529, 9.03195e-06, -6.79836e-07},
		{-533, -5.28768e-07, -5.65456e-07},
		{-532, 1.78178e-06, -5.27175e-07},
		{-527, 1.31307e-05, -3.57207e-07},
		{-541, -1.94807e-05, -4.02796e-07},
		{-531, 4.27835e-06, -5.32770e-07},
		{-529, 8.65372e-06, -3.81893e-07},
		{-541, -1.91792e-05, -4.67138e-07},
		{-531, 4.56602e-06, -5.88393e-07},
		{-530, 6.64196e-06, -4.51798e-07},
		{-540, -1.65259e-05, -5.27556e-07},
		{-532, 2.52472e-06, -6.40050e-07},
		{-530, 6.98085e-06, -5.17709e-07},
		{-536, -7.05781e-06, -4.94620e-07},
		{-527, 1.38430e-05, -4.05554e-07},
		{-530, 6.34692e-06, -1.89628e-07},
		{-539, -1.45211e-05, -2.60755e-07},
		{-537, -9.28004e-06, -4.81279e-07},
		{-528, 1.17231e-05, -4.35248e-07},
		{-532, 2.03845e-06, -2.80758e-07},
		{-540, -1.63450e-05, -4.07601e-07},
		{-528, 1.18421e-05, -4.32824e-07},
		{-538, -1.15653e-05, -4.10188e-07},
	}

	t.Run("highpass", func(t *testing.T) {

		q := 0.95395
		gain := 427336.1
		tol := 0.0001

		high := NewHighPass(gain, q)

		for i := range TestSlice {
			x := high.Sample((float64)(TestSlice[i].i))
			if math.Abs(TestSlice[i].a) == 0.0 && math.Abs(x) != 0.0 {
				t.Error("samples not within tolerance")
			} else if math.Abs(1.0-x/TestSlice[i].a) > tol {
				t.Error("samples outside tolerance")
			}
		}
	})

	t.Run("integrator", func(t *testing.T) {

		q := 0.95395
		gain := 427336.1
		tol := 0.0001
		dt := 1.0 / 50.0

		high := NewHighPass(gain, q)
		low := NewIntegrator(1.0, dt, q)

		for i := range TestSlice {
			x := low.Sample(high.Sample((float64)(TestSlice[i].i)))
			if math.Abs(TestSlice[i].v) == 0.0 && math.Abs(x) != 0.0 {
				t.Error("samples not within tolerance")
			} else if math.Abs(1.0-x/TestSlice[i].v) > tol {
				t.Error("samples outside tolerance")
			}
		}
	})

}
