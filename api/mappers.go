package api

func MapAll[T any, R any](items []T, fn func(T) R) []R {
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, fn(item))
	}
	return result
}

func MapAllWithError[T any, R any](items []T, fn func(T) (R, error)) ([]R, error) {
	result := make([]R, 0, len(items))
	for _, item := range items {
		mapped, err := fn(item)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}
