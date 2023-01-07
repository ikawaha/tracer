// Code generated by goa v3.10.2, DO NOT EDIT.
//
// calc HTTP client CLI support package
//
// Command:
// $ goa gen calc/design

package client

import (
	calc "calc/gen/calc"
	"fmt"
	"strconv"
)

// BuildAddPayload builds the payload for the calc add endpoint from CLI flags.
func BuildAddPayload(calcAddA string, calcAddB string) (*calc.AddPayload, error) {
	var err error
	var a int
	{
		var v int64
		v, err = strconv.ParseInt(calcAddA, 10, strconv.IntSize)
		a = int(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for a, must be INT")
		}
	}
	var b int
	{
		var v int64
		v, err = strconv.ParseInt(calcAddB, 10, strconv.IntSize)
		b = int(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for b, must be INT")
		}
	}
	v := &calc.AddPayload{}
	v.A = a
	v.B = b

	return v, nil
}
