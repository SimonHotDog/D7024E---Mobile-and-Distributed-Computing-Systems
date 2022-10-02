package util

import "github.com/stretchr/testify/mock"

func GetPointerOrNil[T any](args mock.Arguments, index int) *T {
	if args.Get(index) == nil {
		return nil
	} else {
		return args.Get(index).(*T)
	}
}

func GetArrayOrNil[T any](args mock.Arguments, index int) []T {
	if args.Get(index) == nil {
		return nil
	} else {
		return args.Get(index).([]T)
	}
}
