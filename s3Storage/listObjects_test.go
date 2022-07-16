package s3Storage

import (
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {

	obj1 := newMockObj("b.txt")
	obj2 := newMockObj("a/b.txt")
	obj3 := newMockObj("a/a.txt")
	obj4 := newMockObj("a.txt")

	objSlice := []types.Object{*obj1, *obj2, *obj3, *obj4}
	correctOrder := []types.Object{*obj3, *obj2, *obj4, *obj1}

	sort.Sort(ByFileName(objSlice))

	for i, _ := range objSlice {
		objToTest := objSlice[i]
		correctObj := correctOrder[i]
		if objToTest.Key != correctObj.Key {
			t.Fatalf("Incorrect order! i: %v", i)
		}
	}

}

func newMockObj(s string) *types.Object {
	mockObj := new(types.Object)
	mockObj.Key = &s
	return mockObj
}
