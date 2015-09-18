package gpio

// Map of pins for the edison
var pinMap = map[int]PinInfo{
	0: PinInfo{
		pin:          130,
		resistor:     216,
		levelShifter: 248,
		pwm:          -1,
		mux:          []mux{},
	},
	1: PinInfo{
		pin:          131,
		resistor:     217,
		levelShifter: 249,
		pwm:          -1,
		mux:          []mux{},
	},
	2: PinInfo{
		pin:          128,
		resistor:     218,
		levelShifter: 250,
		pwm:          -1,
		mux:          []mux{},
	},
	3: PinInfo{
		pin:          12,
		resistor:     219,
		levelShifter: 251,
		pwm:          0,
		mux:          []mux{},
	},

	4: PinInfo{
		pin:          129,
		resistor:     220,
		levelShifter: 252,
		pwm:          -1,
		mux:          []mux{},
	},
	5: PinInfo{
		pin:          13,
		resistor:     221,
		levelShifter: 253,
		pwm:          1,
		mux:          []mux{},
	},
	6: PinInfo{
		pin:          182,
		resistor:     222,
		levelShifter: 254,
		pwm:          2,
		mux:          []mux{},
	},
	7: PinInfo{
		pin:          48,
		resistor:     223,
		levelShifter: 255,
		pwm:          -1,
		mux:          []mux{},
	},
	8: PinInfo{
		pin:          49,
		resistor:     224,
		levelShifter: 256,
		pwm:          -1,
		mux:          []mux{},
	},
	9: PinInfo{
		pin:          183,
		resistor:     225,
		levelShifter: 257,
		pwm:          3,
		mux:          []mux{},
	},
	10: PinInfo{
		pin:          41,
		resistor:     226,
		levelShifter: 258,
		pwm:          4,
		mux: []mux{
			mux{263, HIGH},
			mux{240, LOW},
		},
	},
	11: PinInfo{
		pin:          43,
		resistor:     227,
		levelShifter: 259,
		pwm:          5,
		mux: []mux{
			mux{262, HIGH},
			mux{241, LOW},
		},
	},
	12: PinInfo{
		pin:          42,
		resistor:     228,
		levelShifter: 260,
		pwm:          -1,
		mux: []mux{
			mux{242, LOW},
		},
	},
	13: PinInfo{
		pin:          40,
		resistor:     229,
		levelShifter: 261,
		pwm:          -1,
		mux: []mux{
			mux{243, LOW},
		},
	},
}
