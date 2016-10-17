package common

type ILogger interface {
	Criticalf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Tracef(message string, args ...interface{})
}
