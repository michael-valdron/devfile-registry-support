package util

import (
	"reflect"
	"testing"

	"github.com/emirpasic/gods/lists/arraylist"
)

func TestConvertStringArrayToArrayList(t *testing.T) {
	tests := []struct {
		name       string
		array      []string
		wantResult *arraylist.List
	}{
		{
			name:       "Test singleton array",
			array:      []string{"abc"},
			wantResult: arraylist.New("abc"),
		},
		{
			name:       "Test Array with multiple values",
			array:      []string{"abc", "ab", "test"},
			wantResult: arraylist.New("abc", "ab", "test"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := ConvertStringArrayToArrayList(test.array)
			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Fatalf("Expected: %v, Got: %v\n", *test.wantResult, *gotResult)
			}
		})
	}
}

func TestConvertBoolArrayToArrayList(t *testing.T) {
	tests := []struct {
		name       string
		array      []bool
		wantResult *arraylist.List
	}{
		{
			name:       "Test singleton array",
			array:      []bool{true},
			wantResult: arraylist.New(true),
		},
		{
			name:       "Test Array with multiple values",
			array:      []bool{true, false, true},
			wantResult: arraylist.New(true, false, true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := ConvertBoolArrayToArrayList(test.array)
			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Fatalf("Expected: %v, Got: %v\n", *test.wantResult, *gotResult)
			}
		})
	}
}

func TestConvertIntArrayToArrayList(t *testing.T) {
	tests := []struct {
		name       string
		array      []int
		wantResult *arraylist.List
	}{
		{
			name:       "Test singleton array",
			array:      []int{0},
			wantResult: arraylist.New(0),
		},
		{
			name:       "Test Array with multiple values",
			array:      []int{0, 1, 2},
			wantResult: arraylist.New(0, 1, 2),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := ConvertIntArrayToArrayList(test.array)
			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Fatalf("Expected: %v, Got: %v\n", *test.wantResult, *gotResult)
			}
		})
	}
}

func TestConvertShortArrayToArrayList(t *testing.T) {
	tests := []struct {
		name       string
		array      []int16
		wantResult *arraylist.List
	}{
		{
			name:       "Test singleton array",
			array:      []int16{0},
			wantResult: arraylist.New(int16(0)),
		},
		{
			name:       "Test Array with multiple values",
			array:      []int16{0, 1, 2},
			wantResult: arraylist.New(int16(0), int16(1), int16(2)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := ConvertShortArrayToArrayList(test.array)
			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Fatalf("Expected: %v, Got: %v\n", *test.wantResult, *gotResult)
			}
		})
	}
}

func TestConvertLongArrayToArrayList(t *testing.T) {
	tests := []struct {
		name       string
		array      []int64
		wantResult *arraylist.List
	}{
		{
			name:       "Test singleton array",
			array:      []int64{0},
			wantResult: arraylist.New(int64(0)),
		},
		{
			name:       "Test Array with multiple values",
			array:      []int64{0, 1, 2},
			wantResult: arraylist.New(int64(0), int64(1), int64(2)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := ConvertLongArrayToArrayList(test.array)
			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Fatalf("Expected: %v, Got: %v\n", *test.wantResult, *gotResult)
			}
		})
	}
}

func TestConvertFloatArrayToArrayList(t *testing.T) {
	tests := []struct {
		name       string
		array      []float32
		wantResult *arraylist.List
	}{
		{
			name:       "Test singleton array",
			array:      []float32{0.23},
			wantResult: arraylist.New(float32(0.23)),
		},
		{
			name:       "Test Array with multiple values",
			array:      []float32{0.23, 1.1, 2.2},
			wantResult: arraylist.New(float32(0.23), float32(1.1), float32(2.2)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := ConvertFloatArrayToArrayList(test.array)
			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Fatalf("Expected: %v, Got: %v\n", *test.wantResult, *gotResult)
			}
		})
	}
}

func TestConvertDoubleArrayToArrayList(t *testing.T) {
	tests := []struct {
		name       string
		array      []float64
		wantResult *arraylist.List
	}{
		{
			name:       "Test singleton array",
			array:      []float64{0.23},
			wantResult: arraylist.New(0.23),
		},
		{
			name:       "Test Array with multiple values",
			array:      []float64{0.23, 1.1, 2.2},
			wantResult: arraylist.New(0.23, 1.1, 2.2),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := ConvertDoubleArrayToArrayList(test.array)
			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Fatalf("Expected: %v, Got: %v\n", *test.wantResult, *gotResult)
			}
		})
	}
}
