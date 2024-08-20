package main

import (
	"net/http"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
)

type input struct {
	Numbers []int `json:"numbers"`
	Target  int   `json:"target"`
}

type res struct {
	Solutions interface{} `json:"solutions"`
}

func main() {
	app := gofr.New()

	app.POST("/find-pairs", findPairs)

	app.Run()
}

// findPairs is the handler that takes the input, validates it and returns the calculated
// pairs as the output
func findPairs(ctx *gofr.Context) (interface{}, error) {
	var in input

	// bind the incoming request body into the input structure
	err := ctx.Bind(&in)
	if err != nil {
		return nil, invalidInput{message: "invalid input in body, numbers are not integer"}
	}

	err = validateInput(&in)
	if err != nil {
		return nil, err
	}

	pairs := calculatePairs(in.Numbers, in.Target)

	return response.Raw{
		Data: res{Solutions: pairs},
	}, nil
}

func calculatePairs(arr []int, t int) [][]int {
	ans := make([][]int, 0)
	m := make(map[int]int)

	for i := 0; i < len(arr); i++ {
		m[arr[i]] = i
	}

	// calculate the pairs
	for i := 0; i < len(arr); i++ {
		v, ok := m[t-arr[i]]
		if ok && v != i {
			ans = append(ans, []int{i, v})
			delete(m, arr[i])
		}
	}

	return ans
}

func validateInput(in *input) error {
	if len(in.Numbers) == 0 {
		return invalidInput{message: "no numbers in the input array"}
	}

	if in.Target == 0 {
		return invalidInput{message: "no target in the input array"}
	}

	return nil
}

// custom error struct to send bad request status code with the appropriate message
type invalidInput struct {
	message string
}

func (i invalidInput) Error() string {
	return i.message
}

func (invalidInput) StatusCode() int {
	return http.StatusBadRequest
}
