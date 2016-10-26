package mathKit

import (
	"fmt"
	"strconv"
)

//保留小数位数
func Scale(num float32, size int) float32 {
	str := fmt.Sprintf(("%." + strconv.Itoa(size) + "f"), num)
	n, _ := strconv.ParseFloat(str, 32)
	return float32(n)
}

func Scale64(num float64, size int) float64 {
	str := fmt.Sprintf(("%." + strconv.Itoa(size) + "f"), num)
	n, _ := strconv.ParseFloat(str, 64)
	return n
}
