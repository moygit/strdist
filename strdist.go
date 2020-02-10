package strdist

import "math"

type DistanceType float64

const max_distance_val = DistanceType(math.MaxFloat64)
const min_distance_val = -DistanceType(math.MaxFloat64)

func New1dCostArray(val DistanceType) *[alphabet_size]DistanceType {
	var arr [alphabet_size]DistanceType
	for i := 0; i < alphabet_size; i++ {
		arr[i] = val
	}
	return &arr
}

func New2dCostArray(val DistanceType) *[alphabet_size][alphabet_size]DistanceType {
	var arr [alphabet_size][alphabet_size]DistanceType
	for i := 0; i < alphabet_size; i++ {
		for j := 0; j < alphabet_size; j++ {
			arr[i][j] = val
		}
	}
	return &arr
}

const alphabet_size int = 128

var unitArray [alphabet_size]DistanceType
var unit2dArray [alphabet_size][alphabet_size]DistanceType

func init() {
	unitArray = *New1dCostArray(1)
	unit2dArray = *New2dCostArray(1)
}

// golang 2-d slices are faster if they're backed by a 1-d slice
func makeFast2dDistanceSlice(numRows, numCols int) [][]DistanceType {
	array2d := make([][]DistanceType, numRows)
	flattenedArray := make([]DistanceType, numRows*numCols)
	for i := 0; i < numRows; i++ {
		array2d[i] = flattenedArray[(i * numCols):((i + 1) * numCols)]
	}
	return array2d
}

func max2(x1, x2 DistanceType) DistanceType {
	if x1 < x2 {
		return x2
	} else {
		return x1
	}
}

func min3(x1, x2, x3 DistanceType) DistanceType {
	minVal := x1
	if x2 < minVal {
		minVal = x2
	}
	if x3 < minVal {
		minVal = x3
	}
	return minVal
}

func min(x ...DistanceType) DistanceType {
	minVal := max_distance_val
	for _, i := range x {
		if i < minVal {
			minVal = i
		}
	}
	return minVal
}

func max(x ...DistanceType) DistanceType {
	maxVal := min_distance_val
	for _, i := range x {
		if i > maxVal {
			maxVal = i
		}
	}
	return maxVal
}
