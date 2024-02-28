package gojq_test_extention

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type compiled_regex struct {
	regex  map[string]*regexp.Regexp
	rwlock sync.RWMutex
}

var compiled_regex_a *compiled_regex = &compiled_regex{
	regex: make(map[string]*regexp.Regexp),
}

func compile_regexp(re string) (*regexp.Regexp, error) {
	re = strings.ReplaceAll(re, "(?<", "(?P<")
	r, err := regexp.Compile(re)
	if err != nil {
		return nil, fmt.Errorf("compile_regexp - invalid regular expression %q: %s", re, err)
	}
	return r, nil
}

func (cr *compiled_regex) get(re string) (*regexp.Regexp, error) {
	cr.rwlock.RLock()
	cre, ok := cr.regex[re]
	cr.rwlock.RUnlock()
	if !ok {
		cr.rwlock.Lock()
		defer cr.rwlock.Unlock()

		cre, ok = cr.regex[re]
		if !ok {
			ncre, err := compile_regexp(re)
			if err != nil {
				return nil, err
			}
			cr.regex[re] = ncre
			cre = ncre
		}
	}
	return cre, nil
}

func Compiled_test(in any, args []any) any {
	re := args[0]

	s, ok := in.(string)
	if !ok {
		return fmt.Errorf("compile_test - input is not a string %q", in)
	}
	restr, ok := re.(string)
	if !ok {
		return fmt.Errorf("compile_test - regex is not a string %q", re)
	}

	r, err := compiled_regex_a.get(restr)
	if err != nil {
		return err
	}

	got := r.FindStringSubmatchIndex(s)
	return got != nil
}
