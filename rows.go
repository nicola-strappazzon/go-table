package table

import (
	"reflect"
	"sort"
)

type Rows []Row

func (r *Rows) Remove(i int) {
	if r.Len() > i {
		(*r) = append((*r)[:i], (*r)[i+1:]...)
	}
}

func (r *Rows) SortBy(i int) {
	sort.Slice((*r), func(x, y int) bool {
		t := reflect.TypeOf((*r)[x].Fields[i].Value)
		v1 := reflect.ValueOf((*r)[x].Fields[i].Value)
		v2 := reflect.ValueOf((*r)[y].Fields[i].Value)

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

func (r Rows) Len() int {
	return len(r)
}

func (r Rows) First() Row {
	if r.IsEmpty() {
		return Row{}
	}
	return r[0]
}

func (r Rows) Last() Row {
	if r.IsEmpty() {
		return Row{}
	}
	return r[r.Len()-1]
}

func (r Rows) IsEmpty() bool {
	return r.Len() == 0
}
