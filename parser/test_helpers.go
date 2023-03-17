package parser

import (
	"fmt"
	"testing"

	"github.com/ruiconti/gojs/internal"
)

func PrettyPrintExpr(t *testing.T, expr Node) {
	fmt.Printf("%s", expr.S())
}

func AssertExprEqual(t *testing.T, logger *internal.SimpleLogger, got, expected Node) {
	failure := false
	var errs []string
	if got.Type() != expected.Type() {
		failure = true
		errs = append(errs, fmt.Sprintf("Type differs"))
		errs = append(errs, fmt.Sprintf("Expected %s, got %s", expected.Type(), got.Type()))
	}

	sgot := got.S()
	sexp := expected.S()
	if sexp != sgot {
		failure = true
		errs = append(errs, fmt.Sprintf("PrettyPrint differs"))
		errs = append(errs, fmt.Sprintf("Expected:"))
		errs = append(errs, fmt.Sprint(sexp))
		errs = append(errs, fmt.Sprintf("Got:"))
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
