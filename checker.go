package golog

// Checkers is a slice of Checker
type Checkers []Checker

// Checker modifies an entry before it get written
type Checker interface {
	Check(Entry) bool
}

// CheckerFunc is a handy function which implements Checker
type CheckerFunc func(Entry) bool

// Check checks if an entry should proceed to be written when using a CheckLogger
func (fn CheckerFunc) Check(e Entry) bool {
	return fn(e)
}

// MinSeverity is the min log Level which can be written
type MinSeverity = Level

// LevelChecker is a Checker which ensure the log has an expected Level
type LevelChecker struct {
	MinSeverity MinSeverity
}

// NewLevelChecker returns a LevelChecker with the given MinSeverity
func NewLevelChecker(minSev MinSeverity) LevelChecker {
	return LevelChecker{MinSeverity: minSev}
}

// NewLevelCheckerOption returns an Option which applies a LevelChecker with the given MinSeverity
func NewLevelCheckerOption(minSev MinSeverity) Option {
	return OptionFunc(func(l StdLogger) StdLogger {
		return l.WithCheckers(LevelChecker{MinSeverity: minSev})
	})
}

// Check returns true if the Entry's Level is at least the one configured
func (lc LevelChecker) Check(e Entry) bool {
	return e.Level() >= lc.MinSeverity
}
