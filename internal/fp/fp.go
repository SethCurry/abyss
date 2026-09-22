// Package fp provides functional programming helpers such as Map, Filter,
// and Apply for working with slices and callables in a generic, type-safe way.
package fp

import "fmt"

// Map applies callable to each element of vals and returns a new slice of
// the results.
func Map[T, P any](callable func(T) P, vals []T) []P {
	retVals := make([]P, len(vals))

	for k, v := range vals {
		retVals[k] = callable(v)
	}

	return retVals
}

// MapE applies callable to each element of vals and returns a new slice of the
// results, short-circuiting and returning an error if callable fails.
func MapE[T, P any](callable func(T) (P, error), vals []T) ([]P, error) {
	retVals := make([]P, len(vals))

	for k, v := range vals {
		newVal, err := callable(v)
		if err != nil {
			return retVals, fmt.Errorf("error applying callable at index %d: %w", k, err)
		}

		retVals[k] = newVal
	}

	return retVals, nil
}

// Filter returns the elements of vals for which predicate returns true.
func Filter[T any](predicate func(T) bool, vals []T) []T {
	retVals := make([]T, 0, len(vals))

	for _, v := range vals {
		if predicate(v) {
			retVals = append(retVals, v)
		}
	}

	return retVals
}

// FilterE returns the elements of vals for which predicate returns true,
// short-circuiting and returning an error if predicate fails.
func FilterE[T any](predicate func(T) (bool, error), vals []T) ([]T, error) {
	retVals := make([]T, 0, len(vals))

	for k, v := range vals {
		ok, err := predicate(v)
		if err != nil {
			return retVals, fmt.Errorf("error applying predicate at index %d: %w", k, err)
		}

		if ok {
			retVals = append(retVals, v)
		}
	}

	return retVals, nil
}

// Apply applies each callable to startVal in turn, returning the final value.
func Apply[T any](startVal T, callables ...func(T) T) T {
	for _, v := range callables {
		startVal = v(startVal)
	}

	return startVal
}

// ApplyE applies each callable to startVal in turn, short-circuiting and
// returning an error if a callable fails.
func ApplyE[T any](startVal T, callables ...func(T) (T, error)) (T, error) {
	var err error

	for k, v := range callables {
		startVal, err = v(startVal)

		if err != nil {
			return startVal, fmt.Errorf("error applying callable %d with value %v: %w", k, startVal, err)
		}
	}

	return startVal, nil
}
