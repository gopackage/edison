package ble

// UUID is a bluetooth compatible UUID value.
type UUID struct {
	Value string
}

// NewUUID creates a new UUID from a string.
func NewUUID(uuid string) *UUID {
	return &UUID{Value: uuid}
}

// NewUUIDs creates an array of UUIDs from a list of strings.
func NewUUIDs(uuid ...string) []*UUID {
	ids := make([]*UUID, len(uuid))
	for _, id := range uuid {
		ids = append(ids, NewUUID(id))
	}
	return ids
}
