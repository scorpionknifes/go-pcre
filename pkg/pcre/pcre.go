// Package used as a middleware for the internal implementation of pcre
package pcre

import (
	"github.com/scorpionknifes/go-pcre/internal/pcrec"
)

// Flags for Compile and Match functions.
const (
	ANCHORED          = pcrec.ANCHORED
	BSR_ANYCRLF       = pcrec.BSR_ANYCRLF
	BSR_UNICODE       = pcrec.BSR_UNICODE
	NEWLINE_ANY       = pcrec.NEWLINE_ANY
	NEWLINE_ANYCRLF   = pcrec.NEWLINE_ANYCRLF
	NEWLINE_CR        = pcrec.NEWLINE_CR
	NEWLINE_CRLF      = pcrec.NEWLINE_CRLF
	NEWLINE_LF        = pcrec.NEWLINE_LF
	NO_START_OPTIMIZE = pcrec.NO_START_OPTIMIZE
	NO_UTF8_CHECK     = pcrec.NO_UTF8_CHECK
)

// Flags for Compile functions
const (
	CASELESS          = pcrec.CASELESS
	DOLLAR_ENDONLY    = pcrec.DOLLAR_ENDONLY
	DOTALL            = pcrec.DOTALL
	DUPNAMES          = pcrec.DUPNAMES
	EXTENDED          = pcrec.EXTENDED
	EXTRA             = pcrec.EXTRA
	FIRSTLINE         = pcrec.FIRSTLINE
	JAVASCRIPT_COMPAT = pcrec.JAVASCRIPT_COMPAT
	MULTILINE         = pcrec.MULTILINE
	NEVER_UTF         = pcrec.NEVER_UTF
	NO_AUTO_CAPTURE   = pcrec.NO_AUTO_CAPTURE
	UNGREEDY          = pcrec.UNGREEDY
	UTF8              = pcrec.UTF8
	UCP               = pcrec.UCP
)

// Flags for Match functions
const (
	NOTBOL           = pcrec.NOTBOL
	NOTEOL           = pcrec.NOTEOL
	NOTEMPTY         = pcrec.NOTEMPTY
	NOTEMPTY_ATSTART = pcrec.NOTEMPTY_ATSTART
	PARTIAL_HARD     = pcrec.PARTIAL_HARD
	PARTIAL_SOFT     = pcrec.PARTIAL_SOFT
)

// Flags for Study function
const (
	STUDY_JIT_COMPILE              = pcrec.STUDY_JIT_COMPILE
	STUDY_JIT_PARTIAL_SOFT_COMPILE = pcrec.STUDY_JIT_PARTIAL_SOFT_COMPILE
	STUDY_JIT_PARTIAL_HARD_COMPILE = pcrec.STUDY_JIT_PARTIAL_HARD_COMPILE
)

// Exec-time and get/set-time error codes
const (
	ERROR_NOMATCH        = pcrec.ERROR_NOMATCH
	ERROR_NULL           = pcrec.ERROR_NULL
	ERROR_BADOPTION      = pcrec.ERROR_BADOPTION
	ERROR_BADMAGIC       = pcrec.ERROR_BADMAGIC
	ERROR_UNKNOWN_OPCODE = pcrec.ERROR_UNKNOWN_OPCODE
	ERROR_UNKNOWN_NODE   = pcrec.ERROR_UNKNOWN_NODE
	ERROR_NOMEMORY       = pcrec.ERROR_NOMEMORY
	ERROR_NOSUBSTRING    = pcrec.ERROR_NOSUBSTRING
	ERROR_MATCHLIMIT     = pcrec.ERROR_MATCHLIMIT
	ERROR_CALLOUT        = pcrec.ERROR_CALLOUT
	ERROR_BADUTF8        = pcrec.ERROR_BADUTF8
	ERROR_BADUTF8_OFFSET = pcrec.ERROR_BADUTF8_OFFSET
	ERROR_PARTIAL        = pcrec.ERROR_PARTIAL
	ERROR_BADPARTIAL     = pcrec.ERROR_BADPARTIAL
	ERROR_RECURSIONLIMIT = pcrec.ERROR_RECURSIONLIMIT
	ERROR_INTERNAL       = pcrec.ERROR_INTERNAL
	ERROR_BADCOUNT       = pcrec.ERROR_BADCOUNT
	ERROR_JIT_STACKLIMIT = pcrec.ERROR_JIT_STACKLIMIT
)

// Regexp holds a reference to a compiled regular expression.
// Use Compile or MustCompile to create such objects.
// Use FreeRegexp to free memory when done with the struct.
type Regexp struct {
	regexp pcrec.Regexp
}

// Free c allocated memory related to regexp.
func (re *Regexp) FreeRegexp() {
	re.regexp.FreeRegexp()
}

// Compile the pattern and return a compiled regexp.
// If compilation fails, the second return value holds a *CompileError.
func Compile(pattern string, flags int) (Regexp, error) {
	re, err := pcrec.Compile(pattern, flags)
	return Regexp{re}, err
}

// CompileJIT is a combination of Compile and Study. It first compiles
// the pattern and if this succeeds calls Study on the compiled pattern.
// comFlags are Compile flags, jitFlags are study flags.
// If compilation fails, the second return value holds a *CompileError.
func CompileJIT(pattern string, comFlags, jitFlags int) (Regexp, error) {
	re, err := pcrec.CompileJIT(pattern, comFlags, jitFlags)
	return Regexp{re}, err
}

// MustCompile compiles the pattern.  If compilation fails, panic.
func MustCompile(pattern string, flags int) Regexp {
	re := pcrec.MustCompile(pattern, flags)
	return Regexp{re}
}

// MustCompileJIT compiles and studies the pattern.  On failure it panics.
func MustCompileJIT(pattern string, comFlags, jitFlags int) Regexp {
	re := pcrec.MustCompileJIT(pattern, comFlags, jitFlags)
	return Regexp{re}
}

// Study adds Just-In-Time compilation to a Regexp. This may give a huge
// speed boost when matching. If an error occurs, return value is non-nil.
// Flags optionally specifies JIT compilation options for partial matches.
func (re *Regexp) Study(flags int) error {
	return re.regexp.Study(flags)
}

// Groups returns the number of capture groups in the compiled pattern.
func (re Regexp) Groups() int {
	return re.regexp.Groups()
}

// Matcher objects provide a place for storing match results.
// They can be created by the Matcher and MatcherString functions,
// or they can be initialized with Reset or ResetString.
type Matcher struct {
	matcher pcrec.Matcher
}

// NewMatcher creates a new matcher object for the given Regexp.
func (re Regexp) NewMatcher() (m *Matcher) {
	m = new(Matcher)
	m.Init(&re)
	return
}

// Matcher creates a new matcher object, with the byte slice as subject.
// It also starts a first match on subject. Test for success with Matches().
func (re Regexp) Matcher(subject []byte, flags int) (m *Matcher) {
	m = re.NewMatcher()
	m.Match(subject, flags)
	return
}

// MatcherString creates a new matcher, with the specified subject string.
// It also starts a first match on subject. Test for success with Matches().
func (re Regexp) MatcherString(subject string, flags int) (m *Matcher) {
	m = re.NewMatcher()
	m.MatchString(subject, flags)
	return

}

// Reset switches the matcher object to the specified regexp and subject.
// It also starts a first match on subject.
func (m *Matcher) Reset(re Regexp, subject []byte, flags int) bool {
	m.Init(&re)
	return m.Match(subject, flags)
}

// ResetString switches the matcher object to the given regexp and subject.
// It also starts a first match on subject.
func (m *Matcher) ResetString(re Regexp, subject string, flags int) bool {
	m.Init(&re)
	return m.MatchString(subject, flags)
}

// Init binds an existing Matcher object to the given Regexp.
func (m *Matcher) Init(re *Regexp) {
	m.matcher.Init(&re.regexp)
}

// Err returns first error encountered by Matcher.
func (m *Matcher) Err() error {
	return m.matcher.Err()
}

// Match tries to match the specified byte slice to
// the current pattern by calling Exec and collects the result.
// Returns true if the match succeeds.
// Match is a no-op if err is not nil.
func (m *Matcher) Match(subject []byte, flags int) bool {
	return m.matcher.Match(subject, flags)
}

// MatchString tries to match the specified subject string to
// the current pattern by calling ExecString and collects the result.
// Returns true if the match succeeds.
func (m *Matcher) MatchString(subject string, flags int) bool {
	return m.matcher.MatchString(subject, flags)
}

// Exec tries to match the specified byte slice to
// the current pattern. Returns the raw pcre_exec error code.
func (m *Matcher) Exec(subject []byte, flags int) int {
	return m.matcher.Exec(subject, flags)
}

// ExecString tries to match the specified subject string to
// the current pattern. It returns the raw pcre_exec error code.
func (m *Matcher) ExecString(subject string, flags int) int {
	return m.matcher.ExecString(subject, flags)
}

// Matches returns true if a previous call to Matcher, MatcherString, Reset,
// ResetString, Match or MatchString succeeded.
func (m *Matcher) Matches() bool {
	return m.matcher.Matches()
}

// Partial returns true if a previous call to Matcher, MatcherString, Reset,
// ResetString, Match or MatchString found a partial match.
func (m *Matcher) Partial() bool {
	return m.matcher.Partial()
}

// Groups returns the number of groups in the current pattern.
func (m *Matcher) Groups() int {
	return m.matcher.Groups()
}

// Present returns true if the numbered capture group is present in the last
// match (performed by Matcher, MatcherString, Reset, ResetString,
// Match, or MatchString).  Group numbers start at 1.  A capture group
// can be present and match the empty string.
func (m *Matcher) Present(group int) bool {
	return m.matcher.Present(group)
}

// Group returns the numbered capture group of the last match (performed by
// Matcher, MatcherString, Reset, ResetString, Match, or MatchString).
// Group 0 is the part of the subject which matches the whole pattern;
// the first actual capture group is numbered 1.  Capture groups which
// are not present return a nil slice.
func (m *Matcher) Group(group int) []byte {
	return m.matcher.Group(group)
}

// Extract returns a slice of byte slices for a single match.
// The first byte slice contains the complete match.
// Subsequent byte slices contain the captured groups.
// If there was no match then nil is returned.
func (m *Matcher) Extract() [][]byte {
	return m.matcher.Extract()
}

// ExtractString returns a slice of strings for a single match.
// The first string contains the complete match.
// Subsequent strings in the slice contain the captured groups.
// If there was no match then nil is returned.
func (m *Matcher) ExtractString() []string {
	return m.matcher.ExtractString()
}

// GroupIndices returns the numbered capture group positions of the last
// match (performed by Matcher, MatcherString, Reset, ResetString, Match,
// or MatchString). Group 0 is the part of the subject which matches
// the whole pattern; the first actual capture group is numbered 1.
// Capture groups which are not present return a nil slice.
func (m *Matcher) GroupIndices(group int) []int {
	return m.matcher.GroupIndices(group)
}

// GroupString returns the numbered capture group as a string.  Group 0
// is the part of the subject which matches the whole pattern; the first
// actual capture group is numbered 1.  Capture groups which are not
// present return an empty string.
func (m *Matcher) GroupString(group int) string {
	return m.matcher.GroupString(group)
}

// Index returns the start and end of the first match, if a previous
// call to Matcher, MatcherString, Reset, ResetString, Match or
// MatchString succeeded. loc[0] is the start and loc[1] is the end.
func (m *Matcher) Index() (loc []int) {
	return m.matcher.Index()
}

// Named returns the value of the named capture group.
// This is a nil slice if the capture group is not present.
// If the name does not refer to a group then error is non-nil.
func (m *Matcher) Named(group string) ([]byte, error) {
	return m.matcher.Named(group)
}

// NamedString returns the value of the named capture group,
// or an empty string if the capture group is not present.
// If the name does not refer to a group then error is non-nil.
func (m *Matcher) NamedString(group string) (string, error) {
	return m.matcher.NamedString(group)
}

// NamedPresent returns true if the named capture group is present.
// If the name does not refer to a group then error is non-nil.
func (m *Matcher) NamedPresent(group string) (bool, error) {
	return m.matcher.NamedPresent(group)
}

// FindIndex returns the start and end of the first match,
// or nil if no match.  loc[0] is the start and loc[1] is the end.
func (re *Regexp) FindIndex(bytes []byte, flags int) (loc []int) {
	return re.regexp.FindIndex(bytes, flags)
}

// ReplaceAll returns a copy of a byte slice
// where all pattern matches are replaced by repl.
func (re Regexp) ReplaceAll(bytes, repl []byte, flags int) ([]byte, error) {
	return re.regexp.ReplaceAll(bytes, repl, flags)
}

// ReplaceAllString is equivalent to ReplaceAll with string return type.
func (re Regexp) ReplaceAllString(in, repl string, flags int) (string, error) {
	return re.regexp.ReplaceAllString(in, repl, flags)
}

// Match holds details about a single successful regex match.
type Match struct {
	Finding string // Text that was found.
	Loc     []int  // Index bounds for location of finding.
}

// FindAll finds all instances that match the regex.
func (re Regexp) FindAll(subject string, flags int) ([]Match, error) {
	matches, err := re.regexp.FindAll(subject, flags)
	if err != nil {
		return nil, err
	}
	newMatches := make([]Match, 0, len(matches))
	for i, match := range matches {
		newMatches[i] = Match{match.Finding, match.Loc}
	}
	return newMatches, nil
}
