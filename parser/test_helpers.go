package parser

import (
	"fmt"
	"testing"
)

func CompareRootChildren(t *testing.T, src string, got, expected []AstNode) {
	failed := false
	errs := []string{}
	gots := []string{}
	exps := []string{}
	if len(got) != len(expected) {
		errs = append(errs, fmt.Sprintf("Expected %d children, got %d", len(expected), len(got)))
		failed = true
	}
	// iterating over expected
	for i := 0; i < len(expected); i++ {
		iexp := expected[i]
		exps = append(exps, fmt.Sprintf("%v", iexp))

		var igot AstNode
		if i >= len(got) {
			igot = &ExprElision{}
		} else {
			igot = got[i]
		}
		if iexp.Type() != igot.Type() || iexp.Source() != igot.Source() {
			errs = append(errs, fmt.Sprintf("%d: Expected %s, got %s", i, iexp, igot))
			failed = true
		}
	}
	// iterating over got
	for i := 0; i < len(got); i++ {
		igot := got[i]
		gots = append(gots, fmt.Sprintf("%v", igot))

		var iexp AstNode
		if i >= len(expected) {
			iexp = &ExprElision{}
		} else {
			iexp = expected[i]
		}
		if iexp.Type() != igot.Type() || iexp.Source() != igot.Source() {
			errs = append(errs, fmt.Sprintf("%d: Expected %s, got %s", i, iexp, igot))
			failed = true
		}
	}
	if failed {
		t.Errorf("Error while parsing: %s", src)
		t.Errorf("Found differences (index: expected, got):")
		for _, err := range errs {
			t.Errorf(err)
		}
		t.Errorf("got:")
		for _, g := range gots {
			t.Errorf(g)
		}
		t.Errorf("expected:")
		for _, e := range exps {
			t.Errorf(e)
		}
		t.Fail()
	}
}

func CompareRootChildrenPointer(t *testing.T, src string, got, expected []*AstNode) {
	failed := false
	errs := []string{}
	gots := []string{}
	exps := []string{}
	if len(got) != len(expected) {
		errs = append(errs, fmt.Sprintf("Expected %d children, got %d", len(expected), len(got)))
		failed = true
	}
	// iterating over expected
	for i := 0; i < len(expected); i++ {
		iexp := expected[i]
		exps = append(exps, fmt.Sprintf("%v", iexp))

		var igot AstNode
		if i >= len(got) {
			igot = &ExprElision{}
		} else {
			igot = *got[i]
		}
		if (*iexp).Type() != (igot).Type() || (*iexp).Source() != (igot).Source() {
			errs = append(errs, fmt.Sprintf("%d: Expected %s, got %s", i, *iexp, igot))
			failed = true
		}
	}
	// iterating over got
	for i := 0; i < len(got); i++ {
		igot := got[i]
		gots = append(gots, fmt.Sprintf("%v", igot))

		var iexp AstNode
		if i >= len(expected) {
			iexp = &ExprElision{}
		} else {
			iexp = *expected[i]
		}
		if iexp.Type() != (*igot).Type() || iexp.Source() != (*igot).Source() {
			errs = append(errs, fmt.Sprintf("%d: Expected %s, got %s", i, iexp, *igot))
			failed = true
		}
	}
	if failed {
		t.Errorf("Error while parsing: %s", src)
		t.Errorf("Found differences (index: expected, got):")
		for _, err := range errs {
			t.Errorf(err)
		}
		t.Errorf("got:")
		for _, g := range gots {
			t.Errorf(g)
		}
		t.Errorf("expected:")
		for _, e := range exps {
			t.Errorf(e)
		}
		t.Fail()
	}
}
