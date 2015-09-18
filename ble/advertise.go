package ble

/*
var AdvertisementIface = {
    name: Defs.Dbus.LE_ADVERTISEMENT_IFACE,
    methods: {
      Release: ['', '']
    },
    signals: {
    },
    properties: {
        Type: 's',
        ServiceUUIDs: 'as',
        ManufacturerData: 'a{qay}',
        ServiceData: 'a{say}',
        IncludeTxPower: 'b'
    }
};
*/

// ManufacturerData captures the data used in bluetooth advertisements for manufacturers.
type ManufacturerData struct {
	Code uint16
	Data []byte
}

// NewManufacturerData creates a new data representation for the provided code and data.
func NewManufacturerData(code uint16, data ...byte) *ManufacturerData {
	m := &ManufacturerData{}
	m.Code = code
	for _, d := range data {
		m.Data = append(m.Data, d)
	}
	return m
}

// Advertisement captures the data used in a bluetooth advertisement.
type Advertisement struct {
	Name           string
	Type           string
	UUIDs          []*UUID
	Manufacturer   *ManufacturerData
	IncludeTxPower bool
}

// NewAdvertisement creates a new avertisement object.
//
//Difference from Bluez example is that dbus-native lib does not call GetAll
//method, instead it gets the properties that are defined in the interface.properties
//that is exported via exportInterface method.
func NewAdvertisement(name, adType string, serviceUUIDs []*UUID, manuf *ManufacturerData) *Advertisement {
	ad := &Advertisement{}
	ad.Name = name
	ad.Type = adType
	ad.UUIDs = serviceUUIDs
	ad.Manufacturer = manuf
	ad.IncludeTxPower = true

	/*
	   this.include_tx_power = true
	   if (impl.include_tx_power !== undefined) {
	       this.include_tx_power = impl.include_tx_power
	   }

	   ad.setDbusProperties()
	*/
	return ad
}

//Set properties as defined in the interface so that
//org.freedesktop.DBus.Properties.GetAll function can access them.
//Also they need to be formatted as dbus-native expects them.
func (a *Advertisement) setDbusProperties() {
	/*
	   a.Type = a.ad_type
	   a.ServiceUUIDs = [a.service_uuids]

	   a.ManufacturerData = a.manufacturer_data ? [[a.manufacturer_data]] : [[]]
	   a.ServiceData = a.service_data ? [[a.service_data]] : [[]]

	   a.IncludeTxPower = a.include_tx_power
	*/
}

func (a *Advertisement) Path() string {
	return advertisementPathBase + a.Name
}

/*
func (a *Advertisement) Release() {
    log.Println("Advertising Released", a.Path())
};
*/
