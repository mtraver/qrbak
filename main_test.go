package main

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name string
		s    string
		n    int
		want []string
	}{
		{
			"negative_parts",
			"foo",
			-1,
			[]string{},
		},
		{
			"empty_0",
			"",
			0,
			[]string{},
		},
		{
			"nonempty_0",
			"foo",
			0,
			[]string{},
		},
		{
			"one_part",
			"foo",
			1,
			[]string{"foo"},
		},
		{
			"len_equal_to_parts",
			"foo",
			3,
			[]string{"f", "o", "o"},
		},
		{
			"len_greater_than_parts",
			"foo",
			5,
			[]string{"f", "o", "o"},
		},
		{
			"equal_parts",
			"foobarbaz",
			3,
			[]string{"foo", "bar", "baz"},
		},
		{
			"remainder",
			"foobarbazz",
			3,
			[]string{"foob", "arba", "zz"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := split(c.s, c.n)
			if diff := pretty.Compare(got, c.want); diff != "" {
				t.Errorf("Unexpected result (-got +want):\n%s", diff)
			}
		})
	}
}
