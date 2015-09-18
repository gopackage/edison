// Package ble creates an easy to use, Go friendly interface to the bluetooth
// low energy (ble) functionality of the Edison.
package ble

// Bluetooth manages the Bluez bluetooth system.
type Bluetooth struct {
}

// NewBluetooth creates a new bluetooth manager.
func NewBluetooth() *Bluetooth {
	return &Bluetooth{}
}

// Start starts the bluetooth stack.
func (b *Bluetooth) Start(service *Service, ad *Advertisement) {

}

// Close closes the bluetooth stack.
func (b *Bluetooth) Close() {
	// notify service to disconnect - iterate through characteristics and stop notify
}

// Handler handles changes to GATT characteristics
type Handler struct {
	value []byte
}

func (h *Handler) Read() []byte {
	return h.value
}

func (h *Handler) Write(value []byte) {
	h.value = value
}

// Example clarifies where we are going by building out an example
func Example() {
	ble := NewBluetooth()
	defer ble.Close()

	// Configure service, avertisement
	handler := &Handler{[]byte("hello")}
	service := NewService(
		"helloService",
		NewUUID("88888888-1111-2222-3333-56789dddddd0"),
		true,
		NewCharacteristic(
			"helloChar",
			NewUUID("dddd5678-1234-5678-1234-56789dddddd1"),
			[]string{"read", "write", "notify"},
			NewDescriptor(
				"helloCharUserDescriptor",
				CharacteristicUserDescription,
				[]string{"read"},
				handler,
			),
		),
	)
	ad := NewAdvertisement(
		"advertisement", // ad name
		"perpheral",     // ad type
		NewUUIDs("88888888-1111-2222-3333-56789dddddd0"), // service UUIDs
		NewManufacturerData(0xffff, 0x00),                // manufacturer data
	)
	ble.Start(service, ad)

	// go channel watch for SIGINT | SIGTERM
}
