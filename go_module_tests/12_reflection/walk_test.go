package reflection

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Profile Profile
}
type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		}, {
			"Struct with two string fields",
			struct {
				FName string
				Lname string
			}{"Morgan", "Freeman"},
			[]string{"Morgan", "Freeman"},
		}, {
			"Struct with two string fields and an int",
			struct {
				FName string
				Lname string
				Age   int
			}{"Morgan", "Freeman", 42},
			[]string{"Morgan", "Freeman"},
		}, {
			"Nested fields",
			Person{
				"Chris",
				Profile{42, "London"},
			},
			[]string{"Chris", "London"},
		}, {
			"Pointers to things",
			&Person{
				"Joe",
				Profile{42, "Burgas"},
			},
			[]string{"Joe", "Burgas"},
		}, {
			"Slice me",
			[]Profile{
				{42, "Sofia"},
				{82, "Varna"},
			},
			[]string{"Sofia", "Varna"},
		}, {
			"Arrays",
			[2]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		}}

	for _, test := range cases {
		t.Run("Walk with valid input", func(t *testing.T) {
			var got []string

			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			assertEquals(t, test.ExpectedCalls, got)
		})
	}

	t.Run("Walk with map", func(t *testing.T) {
		testMap := map[string]string{
			"Morgan":     "Freeman",
			"Stormlight": "Kaladin",
		}

		var got []string

		walk(testMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Freeman")
		assertContains(t, got, "Kaladin")
	})

	t.Run("Walk with channel", func(t *testing.T) {
		ch := make(chan Profile)

		go func() {
			ch <- Profile{42, "Adua"}
			ch <- Profile{68, "King's Landing"}
			close(ch)
		}()

		var got []string
		expected := []string{"Adua", "King's Landing"}

		walk(ch, func(input string) {
			got = append(got, input)
		})

		assertEquals(t, expected, got)
	})

	t.Run("Walk with function", func(t *testing.T) {
		fun := func() (Profile, Profile) {
			return Profile{42, "Tokyo"}, Profile{42, "Betelgeize"}
		}

		var got []string
		expected := []string{"Tokyo", "Betelgeize"}
		walk(fun, func(input string) {
			got = append(got, input)
		})

		assertEquals(t, expected, got)
	})
}

func assertEquals(t *testing.T, arr1 []string, arr2 []string) {
	t.Helper()
	if len(arr1) != len(arr2) {
		t.Fatalf("Wrong number of function calls want: %d but got %d", len(arr1), len(arr2))
	}

	if !reflect.DeepEqual(arr1, arr2) {
		t.Errorf("Expected the string: %v but got %v", arr1, arr2)
	}
}

func assertContains(t *testing.T, arr []string, key string) {
	t.Helper()
	var flag bool

	for _, v := range arr {
		if v == key {
			flag = true
		}
	}

	if !flag {
		t.Errorf("Word: %s was missing from result: %v", key, arr)
	}

}
