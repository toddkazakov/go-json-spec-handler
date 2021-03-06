package jsh

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// List is just a wrapper around an object array that implements Sendable
type List []*Object

/*
Validate ensures that List is JSON API compatible.
*/
func (list List) Validate(r *http.Request, response bool) *Error {
	for _, object := range list {
		err := object.Validate(r, response)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
UnmarshalJSON allows us to manually decode a list via the json.Unmarshaler
interface.
*/
func (list *List) UnmarshalJSON(rawData []byte) error {
	// Create a sub-type here so when we call Unmarshal below, we don't recursively
	// call this function over and over
	type UnmarshalList List

	// if our "List" is a single object, modify the JSON to make it into a list
	// by wrapping with "[ ]"
	if rawData[0] == '{' {
		rawData = []byte(fmt.Sprintf("[%s]", rawData))
	}

	newList := UnmarshalList{}

	err := json.Unmarshal(rawData, &newList)
	if err != nil {
		return err
	}

	convertedList := List(newList)
	*list = convertedList

	return nil
}
