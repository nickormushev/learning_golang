package clockface

import (
	"fmt"
	"io"
	"math"
	"time"
)

//SVG represents the structure of an SVG file for the clock

const (
	secondHandleSize   float64 = 90
	minHandleSize      float64 = 80
	hourHandleSize     float64 = 60
	clockCentre        float64 = 150
	minutesInHalfClock float64 = 30
	secondsInHalfClock float64 = 30
	hoursInHalfClock   float64 = 6
	secondsInClock     float64 = 2 * secondsInHalfClock
	minutesInClock     float64 = 2 * minutesInHalfClock
	hoursInClock       float64 = 2 * hoursInHalfClock
)

//Point represents a point in Euclidian space
type Point struct {
	X float64
	Y float64
}

//Clock is a clock
type Clock struct {
}

//SecondHand takes in a time object and tells us where the second hand should be pointing at this time
func SecondHand(t time.Time) Point {
	return xHand(SecondsInRadians(t), secondHandleSize)
}

//MinuteHand takes in a time object and tells us where the second hand should be pointing at this time
func MinuteHand(t time.Time) Point {
	return xHand(MinutesInRadians(t), minHandleSize)
}

//HourHand takes in a time object and tells us where the second hand should be pointing at this time
func HourHand(t time.Time) Point {
	return xHand(HoursInRadians(t), hourHandleSize)
}

//xHand finds the coordinates of a handle in a unit circle and centers and resizes them
func xHand(rad float64, scale float64) Point {
	p := unitCoordinates(rad)
	p = centerPoint(p, scale)

	return p
}

//XInRadians takes in time and converts it to the angle of the clockhand
func XInRadians(hand int, scale float64) float64 {
	return (math.Pi / (scale / (float64(hand))))
}

//SecondsInRadians takes in time and returns the angle of the minute handle
func SecondsInRadians(t time.Time) float64 {
	return XInRadians(t.Second(), secondsInHalfClock)
}

//MinutesInRadians takes in time and returns the angle of the minute handle
func MinutesInRadians(t time.Time) float64 {
	return SecondsInRadians(t)/secondsInClock + XInRadians(t.Minute(), minutesInHalfClock)
}

//HoursInRadians takes in time and returns the angle of the hour handle
func HoursInRadians(t time.Time) float64 {
	return XInRadians(t.Hour(), hoursInHalfClock) + MinutesInRadians(t)/minutesInClock + SecondsInRadians(t)/(minutesInClock*secondsInClock)
}

//centerPoint  centers the clock handle points with a center of 150,150
func centerPoint(p Point, scale float64) Point {
	p = Point{p.X * scale, p.Y * scale} //scale
	p = Point{p.X, -p.Y}
	//Translates the point so it has a center of 150,150
	p = Point{p.X + clockCentre, p.Y + clockCentre}

	return p
}

//unitCoordinates takes in radians and tells us the coordinates of the handle in a unit circle
func unitCoordinates(rad float64) (p Point) {
	p.X = math.Sin(rad)
	p.Y = math.Cos(rad)

	return
}

func writeHours(w io.Writer, t time.Time) {
	p := HourHand(t)
	fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#000;stroke-width:3px;"/>`, p.X, p.Y)
}

func writeMins(w io.Writer, t time.Time) {
	p := MinuteHand(t)
	fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#000;stroke-width:3px;"/>`, p.X, p.Y)
}

func writeSecs(w io.Writer, t time.Time) {
	p := SecondHand(t)
	fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#f00;stroke-width:3px;"/>`, p.X, p.Y)
}

//SVGWriter writes svg to a file
func SVGWriter(w io.Writer, t time.Time) {
	io.WriteString(w, svgStart)
	io.WriteString(w, bezel)
	writeHours(w, t)
	writeMins(w, t)
	writeSecs(w, t)
	io.WriteString(w, svgEnd)
}

const svgStart = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="100%"
     height="100%"
     viewBox="0 0 300 300"
     version="2.0">`

const bezel = `<circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>`

const svgEnd = `</svg>`
