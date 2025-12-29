package utils

import "testing"

func TestGetAge(t *testing.T) {
	{
		age := GetAge("321023200505230706")
		t.Logf("age = %d", age)
	}
	{
		age := GetAge("321023200505240706")
		t.Logf("age = %d", age)
	}
}
