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
		name     string
		arg      any
		expected []string
	}{
		{
			"struct with one string field",
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},
		{
			"struct with two string fields",
			struct {
				Name string
				City string
			}{"Chris", "London"},
			[]string{"Chris", "London"},
		},
		{
			"struct with non string field",
			struct {
				Age int
			}{33},
			nil,
		},
		{
			"nested fields",
			Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"pointers",
			&Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"slices",
			[]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		},
		{
			"arrays",
			[2]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		},
		// {
		// 	"Maps",
		// 	map[string]string{
		// 		"Foo": "Bar",
		// 		"Baz": "Boz",
		// 	},
		// 	[]string{"Bar", "Boz"},
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			Walk(tt.arg, func(x string) {
				got = append(got, x)
			})
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}

	t.Run("maps", func(t *testing.T) {
		m := map[string]string{
			"abc": "Bar",
			"xyz": "Boz",
		}

		var got []string
		Walk(m, func(x string) {
			got = append(got, x)
		})

		assertContains(t, got, "Bar")
		assertContains(t, got, "Boz")
	})

	t.Run("channels", func(t *testing.T) {
		ch := make(chan Profile)

		go func() {
			ch <- Profile{33, "Berlin"}
			ch <- Profile{34, "Katowice"}
			close(ch)
		}()

		var got []string
		expected := []string{"Berlin", "Katowice"}

		Walk(ch, func(x string) {
			got = append(got, x)
		})

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("functions", func(t *testing.T) {
		f := func() (Profile, Profile) {
			return Profile{33, "Berlin"}, Profile{34, "Katowice"}
		}

		var got []string
		expected := []string{"Berlin", "Katowice"}

		Walk(f, func(x string) {
			got = append(got, x)
		})

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})
}

func assertContains(t *testing.T, got []string, v string) {
	t.Helper()
	contains := false
	for _, x := range got {
		if x == v {
			contains = true
		}
	}
	if !contains {
		t.Errorf("expected %+v to contain %q but it didn't", got, v)
	}
}
