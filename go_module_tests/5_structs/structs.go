package structs

import "math"

//Shape represents all shapes
type Shape interface {
	Area() float64
}

//Triangle represents a triangle
type Triangle struct {
	Hypotenuse float64
	Height     float64
}

//Rectangle represents a rectangle
type Rectangle struct {
	Width  float64
	Height float64
}

//Circle represents a circle
type Circle struct {
	Radius float64
}

//Perimeter calculates the perimiter of a rectangle
func (rect Rectangle) Perimeter() float64 {
	return 2 * (rect.Width + rect.Height)
}

//Area calculates the area of a triangle
func (t Triangle) Area() float64 {
	return t.Height * t.Hypotenuse / 2
}

//Area calculates the area of a rectangle
func (rect Rectangle) Area() float64 {
	return rect.Width * rect.Height
}

//Area calculates the area of a circle
func (circle Circle) Area() float64 {
	return circle.Radius * math.Pi
}
