package fp

import (
	"errors"
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("ints to strings", func(t *testing.T) {
		in := []int{1, 2, 3}
		got := Map(func(v int) string { return fmt.Sprintf("n=%d", v) }, in)
		want := []string{"n=1", "n=2", "n=3"}
		assertSliceEqual(t, got, want)
	})

	t.Run("doubles ints", func(t *testing.T) {
		in := []int{1, 2, 3}
		got := Map(func(v int) int { return v * 2 }, in)
		want := []int{2, 4, 6}
		assertSliceEqual(t, got, want)
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		got := Map(func(v int) int { return v }, []int{})
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	t.Run("nil input returns empty slice", func(t *testing.T) {
		got := Map(func(v int) int { return v }, nil)
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	t.Run("preserves length and order", func(t *testing.T) {
		in := []int{5, 4, 3, 2, 1}
		got := Map(func(v int) int { return v * 10 }, in)
		if len(got) != len(in) {
			t.Fatalf("expected length %d, got %d", len(in), len(got))
		}
		for i, v := range in {
			if got[i] != v*10 {
				t.Fatalf("at index %d: expected %d, got %d", i, v*10, got[i])
			}
		}
	})
}

func TestMapE(t *testing.T) {
	t.Run("success path", func(t *testing.T) {
		in := []int{1, 2, 3}
		got, err := MapE(func(v int) (int, error) { return v * 2, nil }, in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSliceEqual(t, got, []int{2, 4, 6})
	})

	t.Run("returns error from callable", func(t *testing.T) {
	 boomErr := errors.New("boom")
		in := []int{1, 2, 3}
		_, err := MapE(func(v int) (int, error) {
			if v == 2 {
				return 0, boomErr
			}
			return v, nil
		}, in)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, boomErr) {
			t.Fatalf("expected error to wrap boomErr, got %v", err)
		}
	})

	t.Run("error message includes index", func(t *testing.T) {
		in := []int{1, 2, 3}
		_, err := MapE(func(v int) (int, error) {
			if v == 3 {
				return 0, errors.New("fail")
			}
			return v, nil
		}, in)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !contains(err.Error(), "index 2") {
			t.Fatalf("expected error message to mention index 2, got %q", err.Error())
		}
	})

	t.Run("empty input returns empty slice with no error", func(t *testing.T) {
		got, err := MapE(func(v int) (int, error) { return v, nil }, []int{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	t.Run("short-circuits on first error", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		calls := 0
		_, err := MapE(func(v int) (int, error) {
			calls++
			if v == 2 {
				return 0, errors.New("stop")
			}
			return v, nil
		}, in)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if calls != 2 {
			t.Fatalf("expected callable to be invoked 2 times, got %d", calls)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("filters even numbers", func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5, 6}
		got := Filter(func(v int) bool { return v%2 == 0 }, in)
		assertSliceEqual(t, got, []int{2, 4, 6})
	})

	t.Run("predicate always true returns all", func(t *testing.T) {
		in := []int{1, 2, 3}
		got := Filter(func(v int) bool { return true }, in)
		assertSliceEqual(t, got, in)
	})

	t.Run("predicate always false returns empty", func(t *testing.T) {
		in := []int{1, 2, 3}
		got := Filter(func(v int) bool { return false }, in)
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		got := Filter(func(v int) bool { return true }, []int{})
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	t.Run("filters strings by length", func(t *testing.T) {
		in := []string{"a", "ab", "abc", "abcd"}
		got := Filter(func(v string) bool { return len(v) >= 3 }, in)
		assertSliceEqual(t, got, []string{"abc", "abcd"})
	})

	t.Run("preserves order", func(t *testing.T) {
		in := []int{3, 1, 4, 1, 5, 9, 2, 6}
		got := Filter(func(v int) bool { return v > 2 }, in)
		assertSliceEqual(t, got, []int{3, 4, 5, 9, 6})
	})
}

func TestFilterE(t *testing.T) {
	t.Run("success path", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		got, err := FilterE(func(v int) (bool, error) { return v%2 == 0, nil }, in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSliceEqual(t, got, []int{2, 4})
	})

	t.Run("returns error from predicate", func(t *testing.T) {
		boomErr := errors.New("boom")
		in := []int{1, 2, 3}
		_, err := FilterE(func(v int) (bool, error) {
			if v == 2 {
				return false, boomErr
			}
			return true, nil
		}, in)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, boomErr) {
			t.Fatalf("expected error to wrap boomErr, got %v", err)
		}
	})

	t.Run("error message includes index", func(t *testing.T) {
		in := []int{1, 2, 3}
		_, err := FilterE(func(v int) (bool, error) {
			if v == 3 {
				return false, errors.New("fail")
			}
			return true, nil
		}, in)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !contains(err.Error(), "index 2") {
			t.Fatalf("expected error message to mention index 2, got %q", err.Error())
		}
	})

	t.Run("empty input returns empty slice with no error", func(t *testing.T) {
		got, err := FilterE(func(v int) (bool, error) { return true, nil }, []int{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	t.Run("short-circuits on first error", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		calls := 0
		_, err := FilterE(func(v int) (bool, error) {
			calls++
			if v == 2 {
				return false, errors.New("stop")
			}
			return true, nil
		}, in)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if calls != 2 {
			t.Fatalf("expected predicate to be invoked 2 times, got %d", calls)
		}
	})
}

func TestApply(t *testing.T) {
	t.Run("applies callables in order", func(t *testing.T) {
		got := Apply(0,
			func(v int) int { return v + 1 },
			func(v int) int { return v * 2 },
			func(v int) int { return v - 3 },
		)
		// ((0 + 1) * 2) - 3 = -1
		if got != -1 {
			t.Fatalf("expected -1, got %d", got)
		}
	})

	t.Run("no callables returns start value", func(t *testing.T) {
		got := Apply(42)
		if got != 42 {
			t.Fatalf("expected 42, got %d", got)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		got := Apply("hello",
			func(s string) string { return s + " world" },
			func(s string) string { return s + "!" },
		)
		if got != "hello world!" {
			t.Fatalf("expected %q, got %q", "hello world!", got)
		}
	})

	t.Run("single callable", func(t *testing.T) {
		got := Apply(10, func(v int) int { return v * v })
		if got != 100 {
			t.Fatalf("expected 100, got %d", got)
		}
	})
}

func TestApplyE(t *testing.T) {
	t.Run("applies callables in order", func(t *testing.T) {
		got, err := ApplyE(0,
			func(v int) (int, error) { return v + 1, nil },
			func(v int) (int, error) { return v * 2, nil },
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 2 {
			t.Fatalf("expected 2, got %d", got)
		}
	})

	t.Run("no callables returns start value", func(t *testing.T) {
		got, err := ApplyE(42)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 42 {
			t.Fatalf("expected 42, got %d", got)
		}
	})

	t.Run("returns error from callable", func(t *testing.T) {
		boomErr := errors.New("boom")
		got, err := ApplyE(0,
			func(v int) (int, error) { return v + 1, nil },
			func(v int) (int, error) { return 0, boomErr },
			func(v int) (int, error) { return v + 100, nil },
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, boomErr) {
			t.Fatalf("expected error to wrap boomErr, got %v", err)
		}
		if got != 0 {
			t.Fatalf("expected returned value to be 0, got %d", got)
		}
	})

	t.Run("short-circuits on first error", func(t *testing.T) {
		calls := 0
		_, err := ApplyE(0,
			func(v int) (int, error) {
				calls++
				return v + 1, nil
			},
			func(v int) (int, error) {
				calls++
				return 0, errors.New("stop")
			},
			func(v int) (int, error) {
				calls++
				return v + 1, nil
			},
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if calls != 2 {
			t.Fatalf("expected 2 callable invocations, got %d", calls)
		}
	})

	t.Run("error message includes index", func(t *testing.T) {
		_, err := ApplyE(0,
			func(v int) (int, error) { return v + 1, nil },
			func(v int) (int, error) { return 0, errors.New("fail") },
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !contains(err.Error(), "callable") {
			t.Fatalf("expected error message to mention callable, got %q", err.Error())
		}
	})
}

func assertSliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %v, want %v", got, want)
	}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("at index %d: got %v, want %v", i, v, want[i])
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}