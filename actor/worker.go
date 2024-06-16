package actor

import (
	"errors"
	"fmt"
	"reflect"
)

func StartWorker(fn interface{}, argument interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return errors.New("first parameter must be function")
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered, need restarting:", r)
				//Handle error in called function
			}
		}()
		if argument != nil {
			reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(argument)})
		} else {
			reflect.ValueOf(fn).Call([]reflect.Value{})
		}
	}()
	return nil
}
