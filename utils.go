package main

import (
	"strconv"
	"strings"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func StringToIntArray(input string) ([]int, error) {
	var array []int
	for _, value := range strings.Split(input, ",") {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		array = append(array, intValue)
	}
	return array, nil
}
