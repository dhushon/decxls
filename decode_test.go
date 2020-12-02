package decxls

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


// TestNormalizeString serves to excercise the string normalizer used to compare key structures 
func TestNormalizeString(t *testing.T) {
	var ReplaceTests = []struct {
		in    string
		out	  string
	}{
		{"Lazy -_Dog", "LazyDog"},
	}
	for _,r := range(ReplaceTests) {
		assert.Equal(t,normalizeString(r.in),r.out, "strings don't match expected output")
	}
}