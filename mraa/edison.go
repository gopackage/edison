// Package mraa provides simplified Go development for Intel Edison.
package mraa

import ()

func Init() error {
	_ = InitEdison()
	return nil
}
