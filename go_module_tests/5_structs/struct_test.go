package structs

import "testing"

func TestPerimeter(t *testing.T) {
	got := Rectangle{3.0, 4.0}.Perimeter()
	expect := 14.0

	if got != expect {
		t.Errorf("The Perimeter should be %.2f but is %.2f", expect, got)
	}
}

func TestArea(t *testing.T) {

	checkArea := func(shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	}

	t.Run("Rectangle", func(t *testing.T) {
		expect := 12.0
		checkArea(Rectangle{3, 4}, expect)
	})

	t.Run("circles", func(t *testing.T) {
		expect := 314.1592653589793

		checkArea(Circle{100}, expect)
	})
}

func TestAreaWithTableDrivenTests(t *testing.T) {
	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{3, 4}, hasArea: 12.0},
		{name: "Circle", shape: Circle{100}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{6, 12}, hasArea: 36},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()

			if got != tt.hasArea {
				t.Errorf("The structure %#v got %g want %g", tt, got, tt.hasArea)
			}
		})
	}
}
