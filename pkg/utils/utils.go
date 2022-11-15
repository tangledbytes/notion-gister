package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

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

func PrettyPrint(v interface{}) {
	byt, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}

	fmt.Println(string(byt))
}

func PrettyTime(t time.Time) string {
	return TimeInLocalZone(t, time.RFC1123)
}

func TimeInLocalZone(t time.Time, layout string) string {
	return TimeInZone(t, viper.GetString("timezone"), layout)
}

func TimeInZone(t time.Time, zone string, layout string) string {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		loc = time.UTC
	}

	return t.In(loc).Format(layout)
}
