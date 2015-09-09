// Package mraa provides simplified Go development for Intel Edison.
package mraa

import ()

func Init() error {
	_ = InitEdison()
	return nil
}

// Init initializes the Edison gpio pins. This global initialization will
// be removed and individual pins will initialize as they are created.
func Init2() error {
	var err error
	/*tristate := newDigitalPin(214)
	if err = tristate.Export(); err != nil {
		return err
	}
	if err = tristate.Direction(OUT); err != nil {
		return err
	}
	if err = tristate.DigitalWrite(LOW); err != nil {
		return err
	}

	for _, i := range []int{263, 262} {
		io := newDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(OUT); err != nil {
			return err
		}
		if err = io.DigitalWrite(HIGH); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		io := newDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(OUT); err != nil {
			return err
		}
		if err = io.DigitalWrite(LOW); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}

	}

	for _, i := range []int{111, 115, 114, 109} {
		if err = changePinMode(i, "1"); err != nil {
			return err
		}
	}

	for _, i := range []int{131, 129, 40} {
		if err = changePinMode(i, "0"); err != nil {
			return err
		}
	}

	err = tristate.DigitalWrite(HIGH)*/
	return err
}
