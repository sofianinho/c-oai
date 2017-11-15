package utils

import	"reflect"

//Contains returns true if a slice contains a value
func Contains(list interface{}, elem interface{}) bool {
	v := reflect.ValueOf(list)
	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == elem {
			return true
		}
	}
	return false
}

//Includes returns true if a slice's elements are all in another (different from "include")
func Includes(list interface{}, elem interface{}) bool {
	e := reflect.ValueOf(elem)
	for i := 0; i < e.Len(); i++ {
		if !Contains(list, e.Index(i).Interface()){
			return false
		}
	}
	return true
}