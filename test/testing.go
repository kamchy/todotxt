package testing

import (
	"strings"
	"testing"
)

func AssertEquals(t *testing.T, desc string, exp string, act string) {
	if act != exp {
		t.Errorf("%v:\nexpected:\n%v\nactual:\n%v\n", desc, exp, act)
	}
}

func AssertStartsWith(t *testing.T, desc string, exp string, act string) {
	if !strings.Contains(act, exp) {
		t.Errorf("%v:\nexpected:\n%v\nactual:\n%v\n", desc, exp, act)
	}
}

func AssertEqualArrays(t *testing.T, desc string, act []string, exp []string) {
	if actlen, explen := len(act), len(exp); actlen != explen {
		t.Errorf("%s: expected lenght %d, given lenght %d", desc, actlen, explen)
	}
	for i, a := range act {
		if a != exp[i] {
			t.Errorf("%s: values at pos %d differ: expected %s, actual %s.", 
			desc, i, a, exp[i])
		}
	}
}

