
package gooa

import "fmt"

// Type Definitions
type ErrorLevel int
type ErrorRealm int
type Error interface {
	Print()
	Invalidate() bool
}
type ErrorHandler interface {
	Dump()
	Error(ErrorLevel, string, ...interface{}) Error
	SetErrorRealm(ErrorRealm) bool
	SetPosition(*Position)
	ShouldStop() bool
	SetStop(bool)
	ShouldImmediatelyStop() bool
}

// Error Level Definitions
const (
	ErrorWarning ErrorLevel = iota
	ErrorGeneral
	ErrorFatal
	ErrorUnsupported
)

var ErrorLevelMessages = map[ErrorLevel]string{
	ErrorWarning: 		"Warning",
	ErrorGeneral: 		"Error",
	ErrorFatal: 		"Fatal Error",
	ErrorUnsupported: 	"Unsupported Feature",
}

// Error Realm Definitions
const (
	ErrorRealmTokenizer ErrorRealm = iota
	ErrorRealmParser 
	ErrorRealmCompiler
)

var ErrorRealmMessages = map[ErrorRealm]string{
	ErrorRealmTokenizer: 		"during Tokenization",
	ErrorRealmParser: 			"while Parsing",
	ErrorRealmCompiler: 		"during Compilation",
}

// BaseError Definition
type BaseError struct {
	message		string
	level 		ErrorLevel
	realm 		ErrorRealm
	position 	Position
	invalid 	bool
}

func (self *BaseError) Print() {
	realm := ErrorRealmMessages[self.realm]
	level := ErrorLevelMessages[self.level]

	fmt.Printf("[Gooa %d:%d] [%s %s] %s\n", self.position.line, self.position.column, level, realm, self.message)
}

func (self *BaseError) Invalidate() bool {
	i := self.invalid
	self.invalid = true
	return i
}

// BaseErrorHandler Definition
type BaseErrorHandler struct {
	stop 		bool
	errors 		[]Error
	position 	*Position
	realm 		ErrorRealm
	fatal 		bool 
}

func (self *BaseErrorHandler) Dump() {
	for _, v := range self.errors {
		if v.Invalidate() {
			continue
		}
		v.Print()
	}
}

func (self *BaseErrorHandler) Error(lvl ErrorLevel, msg string, args ...interface{}) Error {
	if lvl != ErrorWarning {
		self.stop = true
	}

	if lvl == ErrorFatal {
		self.fatal = true
	}
	
	err 			:= &BaseError{}
	err.position 	= self.position.Copy()
	err.message 	= fmt.Sprintf(msg, args...)
	err.realm 		= self.realm
	err.level		= lvl

	self.errors = append(self.errors, err)

	return err
}

func (self *BaseErrorHandler) SetErrorRealm(rlm ErrorRealm) bool {
	self.realm = rlm
	
	return self.ShouldStop()
}

func (self *BaseErrorHandler) ShouldStop() bool {
	return self.stop
}

func (self *BaseErrorHandler) SetStop(b bool) {
	self.stop = b
}

func (self *BaseErrorHandler) SetPosition(pos *Position) {
	self.position = pos
}

func (self *BaseErrorHandler) ShouldImmediatelyStop() bool {
	return self.fatal
}