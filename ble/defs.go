package ble

// Bluetooth IDs for characteristics types
const (
	CharacteristicExtendedProperties  = "2900"
	CharacteristicUserDescription     = "2901"
	ClientCharacteristicConfiguration = "2902"
	ServerCharacteristicConfiguration = "2903"
)

// DBus interface names for Bluez Bluetooth services
const (
	advertisementPathBase     = "/org/bluez/external/advertisements/"
	dbusProperties            = "org.freedesktop.DBus.Properties"
	bluezRoot                 = "/"
	bluezObjectManager        = "org.freedesktop.DBus.ObjectManager"
	bluezDBusServiceName      = "org.bluez"
	bluezAdapter              = bluezDBusServiceName + ".Adapter1"
	bluezDevice               = bluezDBusServiceName + ".Device1"
	bluezGATTManager          = bluezDBusServiceName + ".GattManager1"
	bluezGATTService          = "org.bluez.GattService1"
	bluezGATTCharacteristic   = "org.bluez.GattCharacteristic1"
	bluezGATTDescriptor       = "org.bluez.GattDescriptor1"
	bluezLEAdvertisingManager = "org.bluez.LEAdvertisingManager1"
	bluezLEAdvertisement      = "org.bluez.LEAdvertisement1"
)
