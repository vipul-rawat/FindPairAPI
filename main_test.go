package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"gofr.dev/pkg/gofr"
	gofrHTTP "gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/http/response"
)

func Test_calculatePairs(t *testing.T) {
	testCases := []struct {
		inputArr    []int
		inputTarget int
		expOut      [][]int
	}{
		// standard input
		{
			inputArr:    []int{1, 2, 3},
			inputTarget: 4,
			expOut:      [][]int{{0, 2}},
		},

		// array with duplicate elements
		{
			inputArr:    []int{1, 3, 2, 4, 5, 10, 3},
			inputTarget: 6,
			expOut:      [][]int{{0, 4}, {1, 6}, {2, 3}},
		},

		// empty input array
		{
			inputArr:    []int{},
			inputTarget: 10,
			expOut:      [][]int{},
		},
	}

	for i, tc := range testCases {
		res := calculatePairs(tc.inputArr, tc.inputTarget)

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("TEST FAILED[%d]: expected %v, got %v", i, tc.expOut, res)
		}
	}
}

func Test_findPairs(t *testing.T) {
	testCases := []struct {
		in     input
		expErr error
		expOut res
	}{
		// standard success case
		{
			in: input{
				Numbers: []int{1, 2, 3},
				Target:  4,
			},
			expOut: res{Solutions: [][]int{{0, 2}}},
		},

		// error case
		{
			in: input{
				Numbers: []int{},
				Target:  4,
			},
			expErr: invalidInput{message: "no numbers in the input array"},
		},
	}

	ctx := &gofr.Context{
		Context: context.Background(),
	}

	for i, tc := range testCases {
		b, _ := json.Marshal(tc.in)
		req := httptest.NewRequest(http.MethodGet, "/find-pairs", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = gofrHTTP.NewRequest(req)

		out, err := findPairs(ctx)

		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("TEST FAILED[%d]: expected %v, got %v", i, tc.expErr, err)
		}

		if err == nil && !reflect.DeepEqual(out, response.Raw{Data: tc.expOut}) {
			t.Errorf("TEST FAILED[%d]: expected %v, got %v", i, tc.expOut, out)
		}
	}
}

func Test_validateInput(t *testing.T) {
	var (
		in  input
		err error
	)

	// success case
	in = input{Numbers: []int{1, 2, 3}, Target: 4}
	err = validateInput(&in)

	assert.NoError(t, err)

	// empty input array
	in = input{Numbers: []int{}, Target: 4}
	err = validateInput(&in)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "no numbers in the input array")

	// target not provided
	in = input{Numbers: []int{1, 2, 3}, Target: 0}
	err = validateInput(&in)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "no target in the input array")
}
