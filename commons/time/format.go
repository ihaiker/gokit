package timeKit

import (
	"time"
	"strings"
	"strconv"
)

//format time like java, such as: yyyy-MM-dd HH:mm:ss
func JavaFormat(t time.Time, format string) string {

	//year
	if strings.ContainsAny(format, "y") {

		year := strconv.Itoa(t.Year())

		if strings.Count(format, "yy") == 1 && strings.Count(format, "y") == 2 {
			format = strings.Replace(format, "yy", year[2:], 1)
		} else if strings.Count(format, "yyyy") == 1 && strings.Count(format, "y") == 4 {
			format = strings.Replace(format, "yyyy", year, 1)
		} else {
			panic("format year error! please 'yyyy' or 'yy'")
		}
	}

	//month
	if strings.ContainsAny(format, "M") {

		var month string

		if int(t.Month()) < 10 {
			month = "0" + strconv.Itoa(int(t.Month()))
		} else {
			month = strconv.Itoa(int(t.Month()))
		}

		if strings.Count(format, "MM") == 1 && strings.Count(format, "M") == 2 {
			format = strings.Replace(format, "MM", month, 1)
		} else {
			panic("format month error! please 'MM'")
		}
	}

	//day
	if strings.ContainsAny(format, "d") {

		var day string

		if t.Day() < 10 {
			day = "0" + strconv.Itoa(t.Day())
		} else {
			day = strconv.Itoa(t.Day())
		}

		if strings.Count(format, "dd") == 1 && strings.Count(format, "d") == 2 {
			format = strings.Replace(format, "dd", day, 1)
		} else {
			panic("format day error! please 'dd'")
		}
	}

	//hour
	if strings.ContainsAny(format, "H") {

		var hour string

		if t.Hour() < 10 {
			hour = "0" + strconv.Itoa(t.Hour())
		} else {
			hour = strconv.Itoa(t.Hour())
		}

		if strings.Count(format, "HH") == 1 && strings.Count(format, "H") == 2 {
			format = strings.Replace(format, "HH", hour, 1)
		} else {
			panic("format hour error! please 'HH'")
		}
	}

	//minute
	if strings.ContainsAny(format, "m") {

		var minute string

		if t.Minute() < 10 {
			minute = "0" + strconv.Itoa(t.Minute())
		} else {
			minute = strconv.Itoa(t.Minute())
		}
		if strings.Count(format, "mm") == 1 && strings.Count(format, "m") == 2 {
			format = strings.Replace(format, "mm", minute, 1)
		} else {
			panic("format minute error! please 'mm'")
		}
	}

	//second
	if strings.ContainsAny(format, "s") {

		var second string

		if t.Second() < 10 {
			second = "0" + strconv.Itoa(t.Second())
		} else {
			second = strconv.Itoa(t.Second())
		}

		if strings.Count(format, "ss") == 1 && strings.Count(format, "s") == 2 {
			format = strings.Replace(format, "ss", second, 1)
		} else {
			panic("format second error! please 'ss'")
		}
	}

	return format
}
//2006-01-02 15:04:05.999999999 -0700 MST
//yyyy-MM-dd HH:mm:ss.SSS
func GoLayout(javaPattern string) string{
	layout := javaPattern
	//year
	if strings.ContainsAny(layout, "y") {
		if strings.Count(layout, "yy") == 1 && strings.Count(layout, "y") == 2 {
			layout = strings.Replace(layout, "yy", "06", 1)
		} else if strings.Count(layout, "yyyy") == 1 && strings.Count(layout, "y") == 4 {
			layout = strings.Replace(layout, "yyyy", "2006", 1)
		} else {
			panic("format year error! please 'yyyy' or 'yy'")
		}
	}
	//month
	if strings.ContainsAny(layout, "M") {
		if strings.Count(layout, "MM") == 1 && strings.Count(layout, "M") == 2 {
			layout = strings.Replace(layout, "MM", "01", 1)
		} else {
			panic("format month error! please 'MM'")
		}
	}
	//day
	if strings.ContainsAny(layout, "d") {
		if strings.Count(layout, "dd") == 1 && strings.Count(layout, "d") == 2 {
			layout = strings.Replace(layout, "dd", "02", 1)
		} else {
			panic("format day error! please 'dd'")
		}
	}
	//hour
	if strings.ContainsAny(layout, "H") {
		if strings.Count(layout, "HH") == 1 && strings.Count(layout, "H") == 2 {
			layout = strings.Replace(layout, "HH", "15", 1)
		} else {
			panic("format hour error! please 'HH'")
		}
	}
	
	//minute
	if strings.ContainsAny(layout, "m") {
		if strings.Count(layout, "mm") == 1 && strings.Count(layout, "m") == 2 {
			layout = strings.Replace(layout, "mm", "04", 1)
		} else {
			panic("format minute error! please 'mm'")
		}
	}
	
	//second
	if strings.ContainsAny(layout, "s") {
		if strings.Count(layout, "ss") == 1 && strings.Count(layout, "s") == 2 {
			layout = strings.Replace(layout, "ss", "05", 1)
		} else {
			panic("format second error! please 'ss'")
		}
	}
	if strings.ContainsAny(layout, "S") {
		layout = strings.Replace(layout, "S", "9", 10)
	}
	return layout
}