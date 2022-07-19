package util

import "github.com/emirpasic/gods/lists/arraylist"

// ConvertStringArrayToArrayList converts string array to interface arraylist
func ConvertStringArrayToArrayList(sArray []string) *arraylist.List {
	arrayList := arraylist.New()

	for _, s := range sArray {
		arrayList.Add(s)
	}

	return arrayList
}

// ConvertBoolArrayToArrayList converts bool array to interface arraylist
func ConvertBoolArrayToArrayList(bArray []bool) *arraylist.List {
	arrayList := arraylist.New()

	for _, b := range bArray {
		arrayList.Add(b)
	}

	return arrayList
}

// ConvertIntArrayToArrayList converts int array to interface arraylist
func ConvertIntArrayToArrayList(iArray []int) *arraylist.List {
	arrayList := arraylist.New()

	for _, i := range iArray {
		arrayList.Add(i)
	}

	return arrayList
}

// ConvertShortArrayToArrayList converts int16 array to interface arraylist
func ConvertShortArrayToArrayList(iArray []int16) *arraylist.List {
	arrayList := arraylist.New()

	for _, i := range iArray {
		arrayList.Add(i)
	}

	return arrayList
}

// ConvertLongArrayToArrayList converts int64 array to interface arraylist
func ConvertLongArrayToArrayList(iArray []int64) *arraylist.List {
	arrayList := arraylist.New()

	for _, i := range iArray {
		arrayList.Add(i)
	}

	return arrayList
}

// ConvertFloatArrayToArrayList converts float32 array to interface arraylist
func ConvertFloatArrayToArrayList(fArray []float32) *arraylist.List {
	arrayList := arraylist.New()

	for _, f := range fArray {
		arrayList.Add(f)
	}

	return arrayList
}

// ConvertDoubleArrayToArrayList converts float64 array to interface arraylist
func ConvertDoubleArrayToArrayList(fArray []float64) *arraylist.List {
	arrayList := arraylist.New()

	for _, f := range fArray {
		arrayList.Add(f)
	}

	return arrayList
}
