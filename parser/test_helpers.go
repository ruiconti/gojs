package parser

import (
	"fmt"
	"testing"

	"github.com/ruiconti/gojs/internal"
)

func AssertExprEqual(t *testing.T, logger *internal.SimpleLogger, got, expected Node) {
	failure := false
	var errs []string
	if got.Type() != expected.Type() {
		failure = true
		errs = append(errs, "type differs")
		errs = append(errs, fmt.Sprintf("Expected %s, got %s", expected.Type(), got.Type()))
	}

	sgot := got.S()
	sexp := expected.S()
	if sexp != sgot {
		failure = true
		errs = append(errs, "PrettyPrint differs")
		errs = append(errs, "Expected:")
		errs = append(errs, fmt.Sprint(sexp))
		errs = append(errs, "Got:")
		errs = append(errs, fmt.Sprint(sgot))
	}

	if failure {
		logger.DumpLogs()
		for _, err := range errs {
			t.Errorf(err)
		}
		t.FailNow()
	}
}
