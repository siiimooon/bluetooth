package bluetooth

// This file implements 16-bit and 128-bit UUIDs as defined in the Bluetooth
// specification.

// UUID is a single UUID as used in the Bluetooth stack. It is represented as a
// [4]uint32 instead of a [16]byte for efficiency.
type UUID [4]uint32

// New16BitUUID returns a new 128-bit UUID based on a 16-bit UUID.
//
// Note: only use registered UUIDs. See
// https://www.bluetooth.com/specifications/gatt/services/ for a list.
func New16BitUUID(shortUUID uint16) UUID {
	// https://stackoverflow.com/questions/36212020/how-can-i-convert-a-bluetooth-16-bit-service-uuid-into-a-128-bit-uuid
	var uuid UUID
	uuid[0] = 0x5F9B34FB
	uuid[1] = 0x80000080
	uuid[2] = 0x00001000
	uuid[3] = uint32(shortUUID)
	return uuid
}

// NewUUID returns a new UUID based on the 128-bit (or 16-byte) input.
func NewUUID(uuid [16]byte) UUID {
	u := UUID{}
	u[0] = uint32(uuid[15]) | uint32(uuid[14])<<8 | uint32(uuid[13])<<16 | uint32(uuid[12])<<24
	u[1] = uint32(uuid[11]) | uint32(uuid[10])<<8 | uint32(uuid[9])<<16 | uint32(uuid[8])<<24
	u[2] = uint32(uuid[7]) | uint32(uuid[6])<<8 | uint32(uuid[5])<<16 | uint32(uuid[4])<<24
	u[3] = uint32(uuid[3]) | uint32(uuid[2])<<8 | uint32(uuid[1])<<16 | uint32(uuid[0])<<24
	return u
}

// Replace16BitComponent returns a new UUID where bits 16..32 have been replaced
// with the bits given in the argument. These bits are the same bits that vary
// in the 16-bit compressed UUID form.
//
// This is especially useful for the Nordic SoftDevice, because it is able to
// store custom UUIDs more efficiently when only these bits vary between them.
func (uuid UUID) Replace16BitComponent(component uint16) UUID {
	uuid[3] &^= 0x0000ffff       // clear the new component bits
	uuid[3] |= uint32(component) // set the component bits
	return uuid
}

// Is16Bit returns whether this UUID is a 16-bit BLE UUID.
func (uuid UUID) Is16Bit() bool {
	return uuid.Is32Bit() && uuid[3] == uint32(uint16(uuid[3]))
}

// Is32Bit returns whether this UUID is a 32-bit or 16-bit BLE UUID.
func (uuid UUID) Is32Bit() bool {
	return uuid[0] == 0x5F9B34FB && uuid[1] == 0x80000080 && uuid[2] == 0x00001000
}

// Get16Bit returns the 16-bit version of this UUID. This is only valid if it
// actually is a 16-bit UUID, see Is16Bit.
func (uuid UUID) Get16Bit() uint16 {
	// Note: using a Get* function as a getter because method names can't start
	// with a number.
	return uint16(uuid[3])
}

// Bytes returns a 16-byte array containing the raw UUID.
func (uuid UUID) Bytes() [16]byte {
	buf := [16]byte{}
	buf[0] = byte(uuid[0])
	buf[1] = byte(uuid[0] >> 8)
	buf[2] = byte(uuid[0] >> 16)
	buf[3] = byte(uuid[0] >> 24)
	buf[4] = byte(uuid[1])
	buf[5] = byte(uuid[1] >> 8)
	buf[6] = byte(uuid[1] >> 16)
	buf[7] = byte(uuid[1] >> 24)
	buf[8] = byte(uuid[2])
	buf[9] = byte(uuid[2] >> 8)
	buf[10] = byte(uuid[2] >> 16)
	buf[11] = byte(uuid[2] >> 24)
	buf[12] = byte(uuid[3])
	buf[13] = byte(uuid[3] >> 8)
	buf[14] = byte(uuid[3] >> 16)
	buf[15] = byte(uuid[3] >> 24)
	return buf
}

// String returns a human-readable version of this UUID, such as
// 00001234-0000-1000-8000-00805F9B34FB.
func (uuid UUID) String() string {
	// TODO: make this more efficient.
	s := ""
	raw := uuid.Bytes()
	for i := range raw {
		// Insert a hyphen at the correct locations.
		if i == 4 || i == 6 || i == 8 || i == 10 {
			s += "-"
		}

		// The character to convert to hex.
		c := raw[15-i]

		// First nibble.
		nibble := c >> 4
		if nibble <= 9 {
			s += string(nibble + '0')
		} else {
			s += string(nibble + 'A' - 10)
		}

		// Second nibble.
		nibble = c & 0x0f
		if nibble <= 9 {
			s += string(nibble + '0')
		} else {
			s += string(nibble + 'A' - 10)
		}
	}

	return s
}
