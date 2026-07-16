package main

import (
	"reflect"
	"sort"
)

type Rows []Row

func (r *Rows) Remove(i int) {
	if len((*r)) > i {
		(*r) = append((*r)[:i], (*r)[i+1:]...)
	}
}

func (r *Rows) SortBy(i int) {
	sort.Slice((*r), func(x, y int) bool {
		t := reflect.TypeOf((*r)[x][i].Value)
		v1 := reflect.ValueOf((*r)[x][i].Value)
		v2 := reflect.ValueOf((*r)[y][i].Value)

		switch t.Name() {
		case "int":
			return int(v1.Int()) < int(v2.Int())
		case "int64":
			return int64(v1.Int()) < int64(v2.Int())
		case "float64":
			return v1.Float() < v2.Float()
		case "string":
			return v1.String() < v2.String()
		case "bool":
			return !v1.Bool() // return small numbers first
		default:
			return false // return unmodified
		}
	})
}
