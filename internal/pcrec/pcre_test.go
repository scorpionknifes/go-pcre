package pcrec

import (
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {
	var check = func(p string, groups int) {
		re, err := Compile(p, 0)
		if err != nil {
			t.Error(p, err)
		}
		if g := re.Groups(); g != groups {
			t.Error(p, g)
		}
	}
	check("", 0)
	check("^", 0)
	check("^$", 0)
	check("()", 1)
	check("(())", 2)
	check("((?:))", 1)
}

func TestCompileFail(t *testing.T) {
	var check = func(p, msg string, off int) {
		_, err := Compile(p, 0)
		if err == nil {
			t.Error(p)
		} else {
			cerr := err.(*CompileError)
			switch {
			case cerr.Message != msg:
				t.Error(p, "Message", cerr.Message)
			case cerr.Offset != off:
				t.Error(p, "Offset", cerr.Offset)
			}
		}
	}
	check("(", "missing )", 1)
	check("\\", "\\ at end of pattern", 1)
	check("abc\\", "\\ at end of pattern", 4)
	check("abc\000", "NUL byte in pattern", 3)
	check("a\000bc", "NUL byte in pattern", 1)
}

func strings(b [][]byte) (r []string) {
	r = make([]string, len(b))
	for i, v := range b {
		r[i] = string(v)
	}
	return
}

func equal(l, r []string) bool {
	if len(l) != len(r) {
		return false
	}
	for i, lv := range l {
		if lv != r[i] {
			return false
		}
	}
	return true
}

func checkmatch1(t *testing.T, dostring bool, m *Matcher,
	pattern, subject string, args ...interface{}) {
	re := MustCompile(pattern, 0)
	var prefix string
	if dostring {
		if m == nil {
			m = re.MatcherString(subject, 0)
		} else {
			m.ResetString(re, subject, 0)
		}
		prefix = "string"
	} else {
		if m == nil {
			m = re.Matcher([]byte(subject), 0)
		} else {
			m.Reset(re, []byte(subject), 0)
		}
		prefix = "[]byte"
	}
	if len(args) == 0 {
		if m.Matches() {
			t.Error(prefix, pattern, subject, "!Matches")
		}
	} else {
		if !m.Matches() {
			t.Error(prefix, pattern, subject, "Matches")
			return
		}
		if m.Groups() != len(args)-1 {
			t.Error(prefix, pattern, subject, "Groups", m.Groups())
			return
		}
		for i, arg := range args {
			if s, ok := arg.(string); ok {
				if !m.Present(i) {
					t.Error(prefix, pattern, subject,
						"Present", i)

				}
				if g := string(m.Group(i)); g != s {
					t.Error(prefix, pattern, subject,
						"Group", i, g, "!=", s)
				}
				if g := m.GroupString(i); g != s {
					t.Error(prefix, pattern, subject,
						"GroupString", i, g, "!=", s)
				}
			} else {
				if m.Present(i) {
					t.Error(prefix, pattern, subject,
						"!Present", i)
				}
			}
		}
	}
}

func TestMatcher(t *testing.T) {
	var m Matcher
	check := func(pattern, subject string, args ...interface{}) {
		checkmatch1(t, false, nil, pattern, subject, args...)
		checkmatch1(t, true, nil, pattern, subject, args...)
		checkmatch1(t, false, &m, pattern, subject, args...)
		checkmatch1(t, true, &m, pattern, subject, args...)
	}

	check(`^$`, "", "")
	check(`^abc$`, "abc", "abc")
	check(`^(X)*ab(c)$`, "abc", "abc", nil, "c")
	check(`^(X)*ab()c$`, "abc", "abc", nil, "")
	check(`^.*$`, "abc", "abc")
	check(`^.*$`, "a\000c", "a\000c")
	check(`^(.*)$`, "a\000c", "a\000c", "a\000c")
	check(`def`, "abcdefghi", "def")
}

func TestPartial(t *testing.T) {
	re := MustCompile(`^abc`, 0)
	defer re.FreeRegexp()
	// Check we get a partial match when we should
	m := re.MatcherString("ab", PARTIAL_SOFT)
	if !m.Matches() {
		t.Error("Failed to find any matches")
	} else if !m.Partial() {
		t.Error("The match was not partial")
	}

	// Check we get an exact match when we should
	m = re.MatcherString("abc", PARTIAL_SOFT)
	if !m.Matches() {
		t.Error("Failed to find any matches")
	} else if m.Partial() {
		t.Error("Match was partial but should have been exact")
	}

	m = re.Matcher([]byte("ab"), PARTIAL_SOFT)
	if !m.Matches() {
		t.Error("Failed to find any matches")
	} else if !m.Partial() {
		t.Error("The match was not partial")
	}

	m = re.Matcher([]byte("abc"), PARTIAL_SOFT)
	if !m.Matches() {
		t.Error("Failed to find any matches")
	} else if m.Partial() {
		t.Error("Match was partial but should have been exact")
	}
}

func TestCaseless(t *testing.T) {
	re := MustCompile("abc", CASELESS)
	defer re.FreeRegexp()
	m := re.MatcherString("...Abc...", 0)
	if !m.Matches() {
		t.Error("CASELESS")
	}
	re2 := MustCompile("abc", 0)
	defer re2.FreeRegexp()
	m = re2.MatcherString("Abc", 0)
	if m.Matches() {
		t.Error("!CASELESS")
	}
}

func TestNamed(t *testing.T) {
	pattern := "(?<L>a)(?<M>X)*bc(?<DIGITS>\\d*)"
	re := MustCompile(pattern, 0)
	defer re.FreeRegexp()
	m := re.MatcherString("abc12", 0)
	if !m.Matches() {
		t.Error("Matches")
	}
	if ok, err := m.NamedPresent("L"); !ok || err != nil {
		t.Errorf("NamedPresent(\"L\"): %v", err)
	}
	if ok, err := m.NamedPresent("M"); ok || err != nil {
		t.Errorf("NamedPresent(\"M\"): %v", err)
	}
	if ok, err := m.NamedPresent("DIGITS"); !ok || err != nil {
		t.Errorf("NamedPresent(\"DIGITS\"): %v", err)
	}
	if str, err := m.NamedString("DIGITS"); str != "12" || err != nil {
		t.Errorf("NamedString(\"DIGITS\"): %v", err)
	}
}

func TestMatcherIndex(t *testing.T) {
	re := MustCompile("bcd", 0)
	defer re.FreeRegexp()
	m := re.Matcher([]byte("abcdef"), 0)
	i := m.Index()
	if i[0] != 1 {
		t.Error("FindIndex start", i[0])
	}
	if i[1] != 4 {
		t.Error("FindIndex end", i[1])
	}
	re2 := MustCompile("xyz", 0)
	defer re2.FreeRegexp()
	m = re2.Matcher([]byte("abcdef"), 0)
	i = m.Index()
	if i != nil {
		t.Error("Index returned for non-match", i)
	}
}

func TestFindIndex(t *testing.T) {
	re := MustCompile("bcd", 0)
	defer re.FreeRegexp()
	i := re.FindIndex([]byte("abcdef"), 0)
	if i[0] != 1 {
		t.Error("FindIndex start", i[0])
	}
	if i[1] != 4 {
		t.Error("FindIndex end", i[1])
	}
}

func TestExtract(t *testing.T) {
	re := MustCompile("b(c)(d)", 0)
	defer re.FreeRegexp()
	m := re.MatcherString("abcdef", 0)
	i := m.ExtractString()
	if i[0] != "abcdef" {
		t.Error("Full line unavailable: ", i[0])
	}
	if i[1] != "c" {
		t.Error("First match group no as expected: ", i[1])
	}
	if i[2] != "d" {
		t.Error("Second match group no as expected: ", i[2])
	}
}

func TestReplaceAll(t *testing.T) {
	re := MustCompile("foo", 0)
	var result []byte
	var err error
	defer re.FreeRegexp()
	// Don't change at ends.
	if result, err = re.ReplaceAll(
		[]byte("I like foods."),
		[]byte("car"),
		0,
	); err != nil {
		t.Fatal(err)
	}
	if string(result) != "I like cards." {
		t.Error("ReplaceAll", result)
	}
	// Change at ends.
	if result, err = re.ReplaceAll(
		[]byte("food fight fools foo"),
		[]byte("car"),
		0,
	); err != nil {
		t.Fatal(err)
	}
	if string(result) != "card fight carls car" {
		t.Error("ReplaceAll2", result)
	}
}

func TestFreeRegexp(t *testing.T) {
	re := MustCompileJIT("\\d{3}", 0, STUDY_JIT_COMPILE)
	data := []string{"15asd213", "sadi32fjoi"}
	expected := []bool{true, false}
	for i := 0; i < len(data); i++ {
		m := re.MatcherString(data[i], 0)
		if m.Matches() != expected[i] {
			t.Error("Unexpected match for ", data[i])
		}
	}
	re.FreeRegexp()

	// Test double free.
	re.FreeRegexp()
}

func TestFindAll(t *testing.T) {
	re := MustCompile("\\d{2}x", 0)
	var matches []Match
	var err error
	defer re.FreeRegexp()
	data := "12x 12332xf 43bx62x"
	expected := []Match{
		Match{"12x", []int{0, 3}},
		Match{"32x", []int{7, 10}},
		Match{"62x", []int{16, 19}},
	}
	if matches, err = re.FindAll(data, 0); err != nil {
		t.Fatal(err)
	}
	verifyMatches(t, expected, matches)

	if matches, err = re.FindAll("", 0); err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Error("Expected no results, got: ", matches)
	}

	// Test zero-length matches.
	re2 := MustCompile("\\w*", 0)
	defer re2.FreeRegexp()
	data = "cat dog"
	expected = []Match{
		Match{"cat", []int{0, 3}},
		Match{"", []int{3, 3}},
		Match{"dog", []int{4, 7}},
	}
	matches, err = re2.FindAll(data, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyMatches(t, expected, matches)
}

func verifyMatches(t *testing.T, expected []Match, matches []Match) {
	if len(matches) != len(expected) {
		t.Errorf("Expected %d matches, got: %d", len(expected), len(matches))
	}
	for i := 0; i < len(expected); i++ {
		if !reflect.DeepEqual(matches[i], expected[i]) {
			t.Errorf("Expected match: %v, got: %v", expected[i], matches[i])
		}
	}
}
