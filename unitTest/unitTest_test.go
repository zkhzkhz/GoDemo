package main

import "testing"

func TestAdd(t *testing.T) {
	sum := Add1(1, 2)
	if sum == 3 {
		t.Log("the result is ok")
	} else {
		t.Fatal("the result is wrong")
	}
}
func Add1(a, b int) int {
	return a + b
}

// Unkeyed can be good to remind future you that you've changed the signature of a struct and
// are now not populating all the fields.  That cuts both ways, since it means you *have* to go to every place
// you've created a value of that struct and update it with a new value... but it also means that
// if you don't have a good default for that new field, that you won't miss places it's used.
//
// However, there's a negative here, too.  If you change the *order* of fields in a struct,
// all your code using unkeyed fields is now (probably) wrong.
//
// In my experience, it's almost always better to use keyed fields.
