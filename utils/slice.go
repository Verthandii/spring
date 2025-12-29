package utils

import (
	"math/rand"
	"strings"
)

type reduceFn[T comparable] func(v T) T
type filterFn[T comparable] func(v T) bool

// InSlice 检查是否在切片中
func InSlice[T comparable](v T, slice []T) bool {
	for _, vv := range slice {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceImplode 将字符串切片转成字符
func SliceImplode(sep string, sl []string) (res string) {
	for _, v := range sl {
		res += v + sep
	}
	res = strings.TrimRight(res, sep)
	return
}

// SliceRandList generate an int slice from min to max.
func SliceRandList(min, max int) []int {
	if max < min {
		min, max = max, min
	}
	length := max - min + 1
	list := rand.Perm(length)
	for index := range list {
		list[index] += min
	}
	return list
}

// SliceMerge merges interface slices to one slice.
func SliceMerge[T comparable](slice1, slice2 []T) (res []T) {
	res = append(slice1, slice2...)
	return
}

// SliceReduce generates a new slice after parsing every value by reduce function
func SliceReduce[T comparable](slice []T, fn reduceFn[T]) (res []T) {
	for _, v := range slice {
		res = append(res, fn(v))
	}
	return
}

// SliceRand returns random one from slice.
func SliceRand[T comparable](a []T) (b T) {
	randnum := rand.Intn(len(a))
	b = a[randnum]
	return
}

// SliceSum sums all values in int64 slice.
func SliceSum(intslice []int64) (sum int64) {
	for _, v := range intslice {
		sum += v
	}
	return
}

// SliceFilter generates a new slice after filter function.
func SliceFilter[T comparable](slice []T, a filterFn[T]) (res []T) {
	for _, v := range slice {
		if a(v) {
			res = append(res, v)
		}
	}
	return
}

// SliceDiff returns diff slice of slice1 - slice2.
func SliceDiff[T comparable](slice1, slice2 []T) (res []T) {
	for _, v := range slice1 {
		if !InSlice(v, slice2) {
			res = append(res, v)
		}
	}
	return
}

// SliceIntersect returns slice that are present in all the slice1 and slice2.
func SliceIntersect[T comparable](slice1, slice2 []T) (res []T) {
	for _, v := range slice1 {
		if InSlice(v, slice2) {
			res = append(res, v)
		}
	}
	return
}

// SliceChunkSameSize separates one slice to same sized slice.
func SliceChunkSameSize[T comparable](slice []T, size int) (res [][]T) {
	if size >= len(slice) {
		res = append(res, slice)
		return
	}
	end := size
	for i := 0; i <= (len(slice) - size); i += size {
		res = append(res, slice[i:end])
		end += size
	}
	return
}

// SliceChunk separates one slice.
func SliceChunk[T comparable](slice []T, size int) (res [][]T) {
	if size >= len(slice) {
		res = append(res, slice)
		return
	}
	for i := 0; i < len(slice); i += size {
		if i+size >= len(slice) {
			res = append(res, slice[i:])
		} else {
			res = append(res, slice[i:i+size])
		}
	}
	return
}

// SliceRange generates a new slice from begin to end with step duration of int64 number.
func SliceRange(start, end, step int64) (res []int64) {
	for i := start; i <= end; i += step {
		res = append(res, i)
	}
	return
}

// SlicePad prepends size number of val into slice.
func SlicePad[T comparable](slice []T, size int, val T) []T {
	if size <= len(slice) {
		return slice
	}
	for i := 0; i < (size - len(slice)); i++ {
		slice = append(slice, val)
	}
	return slice
}

// SliceUnique cleans repeated values in slice.
func SliceUnique[T comparable](slice []T) (res []T) {
	for _, v := range slice {
		if !InSlice(v, res) {
			res = append(res, v)
		}
	}
	return
}

// SliceShuffle shuffles a slice.
func SliceShuffle[T comparable](slice []T) []T {
	for i := 0; i < len(slice); i++ {
		a := rand.Intn(len(slice))
		b := rand.Intn(len(slice))
		slice[a], slice[b] = slice[b], slice[a]
	}
	return slice
}

// SliceRemove 分片切割
func SliceRemove[T comparable](slice []T, i int) []T {
	slice = append(slice[:i], slice[i+1:]...)
	return slice
}

// SliceConvert []string转[]any
func SliceConvert(slice []string) (vals []any) {
	for _, v := range slice {
		vals = append(vals, v)
	}
	return
}

func LeftXOR[T comparable](arr1, arr2 []T) []T {
	m := make(map[T]struct{})
	for _, n := range arr2 {
		m[n] = struct{}{}
	}
	res := make([]T, 0)
	for _, n := range arr1 {
		if _, ok := m[n]; !ok {
			res = append(res, n)
		}
	}
	return res
}

func RightXOR[T comparable](arr1, arr2 []T) []T {
	return LeftXOR(arr2, arr1)
}
