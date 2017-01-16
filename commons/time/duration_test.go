package timeKit

import (
	"testing"
	"fmt"
	"time"
)

func TestDuration(t *testing.T) {
	days := 7
	d := Duration(Day, days)
	fmt.Println(d)

	fmt.Println(Days(1))
}

func TestLayout(t *testing.T){
	t.Log(time.Now().Format(GoLayout("yyyy-MM-dd HH:mm")))
}