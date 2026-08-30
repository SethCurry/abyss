package fp

import "fmt"

func Map[T, P any](callable func(T) P, vals []T) []P {
	retVals := make([]P, len(vals))

	for k, v := range vals {
		retVals[k] = callable(v)
	}

	return retVals
}

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

func Filter[T any](predicate func(T) bool, vals []T) []T {
	retVals := make([]T, 0, len(vals))

	for _, v := range vals {
		if predicate(v) {
			retVals = append(retVals, v)
		}
	}

	return retVals
}

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

func Apply[T any](startVal T, callables ...func(T) T) T {
	for _, v := range callables {
		startVal = v(startVal)
	}

	return startVal
}

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
