package timekit

import (
	"testing"
	"fmt"
)

func TestDuration(t *testing.T) {
	days := 7
	d := Duration(Day, days)
	fmt.Println(d)

	fmt.Println(Days(1))
}
