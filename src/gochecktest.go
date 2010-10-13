package gocheck


import (
    "testing"
    "reflect"
    "fmt"
)


type T struct {
    *call
}


// -----------------------------------------------------------------------
// Basic succeeding/failing logic.

func (t *T) Failed() bool {
    return t.call.status == failedSt
}

func (t *T) Fail() {
    t.call.status = failedSt
}

func (t *T) FailNow() {
    t.Fail()
    t.stopNow()
}

func (t *T) Succeed() {
    t.call.status = succeededSt
}

func (t *T) SucceedNow() {
    t.Succeed()
    t.stopNow()
}

func (t *T) ExpectFailure(reason string) {
    t.expectedFailure = &reason
}


// -----------------------------------------------------------------------
// Basic logging.

func (t *T) GetLog() string {
    return t.call.logv
}

func (t *T) Log(args ...interface{}) {
    t.log(args...)
}

func (t *T) Logf(format string, args ...interface{}) {
    t.logf(format, args...)
}

func (t *T) Error(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    t.Fail()
}

func (t *T) Errorf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprintf("Error: " + format, args...))
    t.Fail()
}

func(t *T) Fatal(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    t.FailNow()
}

func(t *T) Fatalf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprintf(format, args...)))
    t.FailNow()
}


// -----------------------------------------------------------------------
// Equality testing.

func (t *T) CheckEqual(obtained interface{}, expected interface{},
                       issue ...interface{}) bool {
    summary := "CheckEqual(obtained, expected):"
    return t.internalCheckEqual(obtained, expected, true, summary, issue...)
}

func (t *T) CheckNotEqual(obtained interface{}, expected interface{},
                          issue ...interface{}) bool {
    summary := "CheckNotEqual(obtained, unexpected):"
    return t.internalCheckEqual(obtained, expected, false, summary, issue...)
}

func (t *T) AssertEqual(obtained interface{}, expected interface{},
                        issue ...interface{}) {
    summary := "AssertEqual(obtained, expected):"
    if !t.internalCheckEqual(obtained, expected, true, summary, issue...) {
        t.stopNow()
    }
}

func (t *T) AssertNotEqual(obtained interface{}, expected interface{},
                           issue ...interface{}) {
    summary := "AssertNotEqual(obtained, unexpected):"
    if !t.internalCheckEqual(obtained, expected, false, summary, issue...) {
        t.stopNow()
    }
}


func (t *T) internalCheckEqual(a interface{}, b interface{}, equal bool,
                               summary string, issue ...interface{}) bool {
    typeA := reflect.Typeof(a)
    typeB := reflect.Typeof(b)
    if (typeA == typeB && checkEqual(a, b)) != equal {
        t.logCaller(2, summary)
        if equal {
            t.logValue("Obtained", a)
            t.logValue("Expected", b)
        } else {
            t.logValue("Both", a)
        }
        if len(issue) != 0 {
            t.logString(fmt.Sprint(issue...))
        }
        t.logNewLine()
        t.Fail()
        return false
    }
    return true
}

// This will use a fast path to check for equality of normal types,
// and then fallback to reflect.DeepEqual if things go wrong.
func checkEqual(a interface{}, b interface{}) (result bool) {
    defer func() {
        if recover() != nil {
            result = reflect.DeepEqual(a, b)
        }
    }()
    return (a == b)
}


// -----------------------------------------------------------------------
// Error testing.

func (t *T) AssertErr(obtained interface{}, expected interface{},
                      issue ...interface{}) {
    var summary string
    if expected == nil {
        summary = "AssertErr(error, nil):"
    } else {
        summary = "AssertErr(error, expected):"
    }
    if !t.internalCheckErr(obtained, expected, true, summary, issue...) {
        t.stopNow()
    }
}

func (t *T) CheckErr(obtained interface{}, expected interface{},
                     issue ...interface{}) bool {
    var summary string
    if expected == nil {
        summary = "CheckErr(error, nil):"
    } else {
        summary = "CheckErr(error, expected):"
    }
    return t.internalCheckErr(obtained, expected, true, summary, issue...)
}

func (t *T) internalCheckErr(a interface{}, b interface{},
                             equal bool, summary string,
                             issue ...interface{}) bool {
    typeA := reflect.Typeof(a)
    typeB := reflect.Typeof(b)
    // If b is string, handle it as a match expression.
    if bStr, ok := b.(string); ok {
        var err string
        var matches bool
        aStr, aIsStr := a.(string)
        if !aIsStr {
            if aWithStr, aHasStr := a.(hasString); aHasStr {
                aStr, aIsStr = aWithStr.String(), true
            }
        }
        if aIsStr {
            matches, err = testing.MatchString("^" + bStr + "$", aStr)
        }
        if !matches || err != "" {
            t.logCaller(2, summary)
            var msg string
            if !matches {
                t.logValue("Error", a)
                msg = fmt.Sprintf("Expected to match expression: %#v", b)
            } else {
                msg = fmt.Sprintf("Can't compile match expression: %#v", b)
            }
            t.logString(msg)
            t.logNewLine()
            t.Fail()
            return false
        }
    } else if (typeA == typeB && checkEqual(a, b)) != equal {
        t.logCaller(2, summary)
        if b == nil {
            t.logValue("Error", a)
        } else {
            t.logValue("Error", a)
            t.logValue("Expected", b)
        }
        if len(issue) != 0 {
            t.logString(fmt.Sprint(issue...))
        }
        t.logNewLine()
        t.Fail()
        return false
    }
    return true
}

