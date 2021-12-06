package main

import (
	"reflect"
	"testing"
)

func TestFindSubstringUsers(t *testing.T) {
	users := []string{"test", "testadmin", "123", "456", "bill", "billadmin"}
	expected := [][]string{{"test", "testadmin"}, {"bill", "billadmin"}}
	res := findSubstringUsers(users)
	if !reflect.DeepEqual(res, expected) {
		t.Fail()
	}

}
