package roma

import (
	"strings"
)

//ConvertToRomanV1 converts integers to roman numerals version1
func ConvertToRomanV1(num int) string {
	switch num {
	case 1:
		return "I"
	case 2:
		return "II"
	case 3:
		return "III"
	default:
		return ""
	}
}

//ConvertToRomanV2 converts integers to roman numerals version2
func ConvertToRomanV2(num int) string {
	var res strings.Builder

	for num > 0 {
		switch {
		case num > 9:
			res.WriteString("X")
			num -= 10
		case num > 8:
			res.WriteString("IX")
			num -= 9
		case num > 4:
			res.WriteString("V")
			num -= 5
		case num > 3:
			res.WriteString("IV")
			num -= 4
		default:
			res.WriteString("I")
			num--

		}
	}

	return res.String()
}

//RomanNumeral is a structure that has an arabic numeral and it's counterpart roman numeral
type RomanNumeral struct {
	arabicNum uint16
	romanNum  string
}

var romanNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func valueOf(symbol string) uint16 {
	for _, v := range romanNumerals {
		if v.romanNum == symbol {
			return v.arabicNum
		}
	}

	return 0
}

//ConvertToRoman final version
func ConvertToRoman(num uint16) string {
	var res strings.Builder

	for _, v := range romanNumerals {
		for num >= v.arabicNum {
			res.WriteString(v.romanNum)
			num -= v.arabicNum
		}
	}

	return res.String()
}

func arabicNumberExists(symbol string) bool {
	return valueOf(symbol) != 0
}

//ConvertToArabic takes in a roman numeral and converts it to arabic
func ConvertToArabic(num string) uint16 {
	var res uint16

	for i := 0; i < len(num); i++ {
		if i+1 < len(num) && couldBeSubtrackted(num[i]) && arabicNumberExists(num[i:i+2]) {
			res += valueOf(num[i : i+2])
			i++
		} else {
			res += valueOf(string(num[i]))
		}
	}

	return res
}

func couldBeSubtrackted(symbol byte) bool {
	return symbol == 'X' || symbol == 'C' || symbol == 'I'
}
