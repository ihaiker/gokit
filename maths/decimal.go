package maths

import (
	"github.com/shopspring/decimal"
)

func Add(a, b float32) float32 {
	f, _ := decimal.NewFromFloat32(a).Add(decimal.NewFromFloat32(b)).Float64()
	return float32(f)
}

func Subtract(a, b float32) float32 {
	f, _ := decimal.NewFromFloat32(a).Sub(decimal.NewFromFloat32(b)).Float64()
	return float32(f)
}

func Multiply(a, b float32) float32 {
	f, _ := decimal.NewFromFloat32(a).Mul(decimal.NewFromFloat32(b)).Float64()
	return float32(f)
}

func Divide(a, b float32) float32 {
	f, _ := decimal.NewFromFloat32(a).Div(decimal.NewFromFloat32(b)).Float64()
	return float32(f)
}

func Add64(a, b float64) float64 {
	f, _ := decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)).Float64()
	return f
}

func Subtract64(a, b float64) float64 {
	f, _ := decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).Float64()
	return f
}

func Multiply64(a, b float64) float64 {
	f, _ := decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(b)).Float64()
	return f
}

func Divide64(a, b float64) float64 {
	f, _ := decimal.NewFromFloat(a).Div(decimal.NewFromFloat(b)).Float64()
	return f
}
