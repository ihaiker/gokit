package maths

import (
	"github.com/shopspring/decimal"
)

type RoundingMode int

const (
	Down RoundingMode = iota
	Round
	Up
)

func Scale(num float32, mode RoundingMode, size int32) float32 {
	switch mode {
	case Down:
		return ScaleDown(num, size)
	case Round:
		return ScaleRound(num, size)
	case Up:
		return ScaleUp(num, size)
	}
	return num
}

//向下保留两位小数
func ScaleDown(num float32, size int32) float32 {
	f, _ := decimal.NewFromFloat32(num).Truncate(size).Float64()
	return float32(f)
}

//四舍五入保留两位小数
func ScaleRound(num float32, size int32) float32 {
	f, _ := decimal.NewFromFloat32(num).Round(size).Float64()
	return float32(f)
}

//向上保留两位小数
func ScaleUp(num float32, size int32) float32 {
	f, _ := decimal.NewFromFloat32(num).Shift(size).Ceil().Shift(-size).Truncate(size).Float64()
	return float32(f)
}

func Scale64(num float64, mode RoundingMode, size int32) float64 {
	switch mode {
	case Down:
		return Scale64Down(num, size)
	case Round:
		return Scale64Round(num, size)
	case Up:
		return Scale64Up(num, size)
	}
	return num
}

//向下保留两位小数
func Scale64Down(num float64, size int32) float64 {
	f, _ := decimal.NewFromFloat(num).Truncate(size).Float64()
	return f
}

//四舍五入保留两位小数
func Scale64Round(num float64, size int32) float64 {
	f, _ := decimal.NewFromFloat(num).Round(size).Float64()
	return f
}

//向上保留两位小数
func Scale64Up(num float64, size int32) float64 {
	f, _ := decimal.NewFromFloat(num).Shift(size).Ceil().Shift(-size).Truncate(size).Float64()
	return f
}
