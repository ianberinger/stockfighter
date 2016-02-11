package api

import (
	"fmt"
	"sync"
)

type err struct {
	sync.Mutex
	v error
}

//apiError creates an error from an error message and a http status.
func apiError(str string, status string) error {
	if str == "" {
		str = "no error message"
	}
	return fmt.Errorf("API: %s; %s", status, str)
}

//setErr sets the error value on an instance only when the error isn't nil. Returns true if error was set.
//This is useful because we don't have to check if the error is nil before calling setErr() if we don't want a current error overwritten by a nil one.
func (i *Instance) setErr(err error) bool {
	if err != nil {
		if i.debug {
			fmt.Println("err: ", err)
		}
		i.err.Lock()
		i.err.v = err
		i.err.Unlock()
		return true
	}
	return false
}

//GetErr returns the last error of an API instance.
func (i *Instance) GetErr() error {
	i.err.Lock()
	defer i.err.Unlock()
	return i.err.v
}

//ResetErr resets the error of an API instance to nil.
func (i *Instance) ResetErr() {
	i.err.Lock()
	i.err.v = nil
	i.err.Unlock()
}

//Debug enables comprehensive logging for the instance.
func (i *Instance) Debug() {
	i.Lock()
	i.debug = true
	i.Unlock()
}
