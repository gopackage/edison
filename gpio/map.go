package gpio

var pinMap = map[int]fsPin{
	0: fsPin{
		pin:          130,
		resistor:     216,
		levelShifter: 248,
		pwmPin:       -1,
		mux:          []mux{},
	},
	1: fsPin{
		pin:          131,
		resistor:     217,
		levelShifter: 249,
		pwmPin:       -1,
		mux:          []mux{},
	},
	2: fsPin{
		pin:          128,
		resistor:     218,
		levelShifter: 250,
		pwmPin:       -1,
		mux:          []mux{},
	},
	3: fsPin{
		pin:          12,
		resistor:     219,
		levelShifter: 251,
		pwmPin:       0,
		mux:          []mux{},
	},

	4: fsPin{
		pin:          129,
		resistor:     220,
		levelShifter: 252,
		pwmPin:       -1,
		mux:          []mux{},
	},
	5: fsPin{
		pin:          13,
		resistor:     221,
		levelShifter: 253,
		pwmPin:       1,
		mux:          []mux{},
	},
	6: fsPin{
		pin:          182,
		resistor:     222,
		levelShifter: 254,
		pwmPin:       2,
		mux:          []mux{},
	},
	7: fsPin{
		pin:          48,
		resistor:     223,
		levelShifter: 255,
		pwmPin:       -1,
		mux:          []mux{},
	},
	8: fsPin{
		pin:          49,
		resistor:     224,
		levelShifter: 256,
		pwmPin:       -1,
		mux:          []mux{},
	},
	9: fsPin{
		pin:          183,
		resistor:     225,
		levelShifter: 257,
		pwmPin:       3,
		mux:          []mux{},
	},
	10: fsPin{
		pin:          41,
		resistor:     226,
		levelShifter: 258,
		pwmPin:       4,
		mux: []mux{
			mux{263, HIGH},
			mux{240, LOW},
		},
	},
	11: fsPin{
		pin:          43,
		resistor:     227,
		levelShifter: 259,
		pwmPin:       5,
		mux: []mux{
			mux{262, HIGH},
			mux{241, LOW},
		},
	},
	12: fsPin{
		pin:          42,
		resistor:     228,
		levelShifter: 260,
		pwmPin:       -1,
		mux: []mux{
			mux{242, LOW},
		},
	},
	13: fsPin{
		pin:          40,
		resistor:     229,
		levelShifter: 261,
		pwmPin:       -1,
		mux: []mux{
			mux{243, LOW},
		},
	},
}
