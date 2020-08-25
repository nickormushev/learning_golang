package clockface_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	clockface "learning/16_clock"
	"math"
	"testing"
	"time"
)

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Version string   `xml:"version,attr"`
	Circle  Circle   `xml:"circle"`
	Line    []Line   `xml:"line"`
}

type Circle struct {
	Cx float64 `xml:"cx,attr"`
	Cy float64 `xml:"cy,attr"`
	R  float64 `xml:"r,attr"`
}

type Line struct {
	X1 float64 `xml:"x1,attr"`
	Y1 float64 `xml:"y1,attr"`
	X2 float64 `xml:"x2,attr"`
	Y2 float64 `xml:"y2,attr"`
}

func TestSecondHandAtMidnight(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := clockface.Point{X: 150, Y: 150 - 90}
	got := clockface.SecondHand(tm)

	if got != want {
		t.Errorf("Got: %v but wanted: %v", got, want)
	}
}

func TestSecondHandAt30seconds(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)

	want := clockface.Point{X: 150, Y: 150 + 90}
	got := clockface.SecondHand(tm)

	if got != want {
		t.Errorf("Got: %v but wanted: %v", got, want)
	}
}

func TestSVGWriterHourHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{simpleTime(12, 0, 0), Line{150, 150, 150, 90}},
		{simpleTime(0, 0, 0), Line{150, 150, 150, 90}},
		{simpleTime(6, 0, 0), Line{150, 150, 150, 210}},
	}

	for _, v := range cases {
		t.Run(testName(v.time), func(t *testing.T) {
			b := bytes.Buffer{}
			clockface.SVGWriter(&b, v.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(t, v.line, svg.Line) {
				t.Errorf("Expected to find the second hand line %+v, in the SVG lines %+v", v.line, svg.Line)
			}
		})
	}
}

func TestSVGWriterSecondHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{simpleTime(0, 30, 30), Line{150, 150, 150, 240}},
		{simpleTime(12, 0, 0), Line{150, 150, 150, 60}},
	}

	for _, v := range cases {
		t.Run(testName(v.time), func(t *testing.T) {
			b := bytes.Buffer{}
			clockface.SVGWriter(&b, v.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(t, v.line, svg.Line) {
				t.Errorf("Expected to find the second hand line %+v, in the SVG lines %+v", v.line, svg.Line)
			}
		})
	}
}

func TestSVGWriterMinuteHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{simpleTime(0, 30, 0), Line{150, 150, 150, 230}},
		{simpleTime(12, 0, 0), Line{150, 150, 150, 70}},
	}

	for _, v := range cases {
		t.Run(testName(v.time), func(t *testing.T) {
			b := bytes.Buffer{}
			clockface.SVGWriter(&b, v.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(t, v.line, svg.Line) {
				t.Errorf("Expected to find the hand line %+v, in the SVG lines %+v", v.line, svg.Line)
			}
		})
	}
}

func containsLine(t *testing.T, desiredLine Line, lines []Line) bool {
	t.Helper()

	for _, fileLine := range lines {
		if fileLine == desiredLine {
			return true
		}
	}

	return false
}

func simpleTime(hours, minutes, seconds int) time.Time {
	return time.Date(2022, time.September, 16, hours, minutes, seconds, 0, time.UTC)
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}

//Not sure if there is a point to this test and if xInRadians shouldn't be a private method
func TestXInRadians(t *testing.T) {
	cases := []struct {
		hand  int
		angle float64
	}{
		{30, math.Pi},
		{0, 0},
		{45, (math.Pi / 2) * 3},
		{15, math.Pi / 2},
		{7, (math.Pi / 30) * 7},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("Time %d", test.hand), func(t *testing.T) {
			got := clockface.XInRadians(test.hand, 30)

			if test.angle != got {
				t.Fatalf("Wanted %v radians, but got %v", test.angle, got)
			}
		})
	}
}

//func TestSecondsHandPoint(t *testing.T) {
//	cases := []struct {
//		time  time.Time
//		point clockface.Point
//	}{
//		{simpleTime(0, 0, 30), clockface.Point{0, -1}},
//		{simpleTime(12, 0, 0), clockface.Point{0, 1}},
//		{simpleTime(0, 0, 45), clockface.Point{-1, 0}},
//		{simpleTime(5, 0, 15), clockface.Point{1, 0}},
//	}
//
//	for _, test := range cases {
//		t.Run(testName(test.time), func(t *testing.T) {
//			got := clockface.SecondHand(test.time)
//
//			if !roughlyEqualPoint(got, test.point) {
//				t.Fatalf("Wanted %v radians, but got %v", test.point, got)
//			}
//		})
//	}
//}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b clockface.Point) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}
