package utils

import "github.com/spf13/viper"

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func ViperReturnFirstFound[T any](keys ...string) T {
	var val T

	for _, key := range keys {
		if viper.IsSet(key) {
			val = viper.Get(key).(T)
			return val
		}
	}

	return val
}

func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}

func FromT1ToT2[T1, T2 any](slice []T1, fn func(T1) T2) []T2 {
	var result []T2

	for _, v := range slice {
		result = append(result, fn(v))
	}

	return result
}

func ForEach[T any](slice []T, fn func(T)) {
	for _, v := range slice {
		fn(v)
	}
}
