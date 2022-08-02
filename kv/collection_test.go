package kv

import (
	"fmt"
	"testing"
)

func TestCollect(t *testing.T) {
	items := []string{"foo", "bar", "baz"}

	collection := Collect(items...)

	for k, v := range items {
		foundV, ok := collection[k]

		if !ok {
			t.Errorf("key %d should've been found", k)
		}

		if v != foundV {
			t.Errorf("found value should be %s. got %s", v, foundV)
		}
	}
}

func TestCollectSlice(t *testing.T) {
	items := []string{"foo", "bar", "baz"}

	collection := CollectSlice(items)

	for k, v := range items {
		foundV, ok := collection[k]

		if !ok {
			t.Errorf("key %d should've been found", k)
		}

		if v != foundV {
			t.Errorf("found value should be %s. got %s", v, foundV)
		}
	}
}

func TestCollectMap(t *testing.T) {
	items := map[string]string{"a": "foo", "b": "bar", "c": "baz"}

	collection := CollectMap(items)

	for k, v := range items {
		foundV, ok := collection[k]

		if !ok {
			t.Errorf("key %s should've been found", k)
		}

		if v != foundV {
			t.Errorf("found value should be %s. got %s", v, foundV)
		}
	}
}

func TesetCombineE(t *testing.T) {
	testCases := []struct {
		description string
		keys        []string
		values      []string
		expected    map[string]string
		err         error
	}{
		{
			"mismatching keys and values lenghts",
			[]string{"a", "b"},
			[]string{"foo", "bar", "baz"},
			map[string]string{},
			fmt.Errorf(""),
		},
	}

	for _, tc := range testCases {
		combined, err := CombineE(tc.keys, tc.values)

		if tc.err != nil {
			if err == nil {
				t.Errorf("expected error '%s'. got nil", tc.err.Error())
			} else if tc.err.Error() != err.Error() {
				t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
			}
		}

		for k, v := range tc.expected {
			foundV, ok := combined[k]

			if !ok {
				t.Errorf("key %s should've been found", k)
			}

			if v != foundV {
				t.Errorf("found value should be %s. got %s", v, foundV)
			}
		}
	}
}
