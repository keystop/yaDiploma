package luhn

import (
	"strconv"

	"github.com/keystop/yaDiploma/pkg/logger"
)

func CheckInteger(i int) bool {
	str := strconv.Itoa(i)
	return CheckString(str)
}

func CheckString(s string) bool {
	var rArr []int
	for _, v := range s {
		// Десятичные значения
		res, err := strconv.Atoi(string(v))
		if err != nil {
			logger.Info("Luhn", "Ошибка конвертации строки", err)
		}
		rArr = append(rArr, res)
	}
	return Check(rArr)
}

func Check(arr []int) bool {
	l := len(arr)
	var v int
	var cd int
	for i, j := l-1, 0; i >= 0; i, j = i-1, j+1 {
		v = arr[i]
		if j%2 != 0 {
			v *= 2
			if v > 9 {
				v -= 9
			}
		}
		cd += v
	}

	return cd%10 == 0
}
