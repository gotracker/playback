package sampling

// Format is the format of the sample data
type Format uint8

const (
	// Format8BitUnsigned is for unsigned 8-bit data
	Format8BitUnsigned = Format(iota)
	// Format8BitSigned is for signed 8-bit data
	Format8BitSigned
	// Format16BitLEUnsigned is for unsigned, little-endian, 16-bit data
	Format16BitLEUnsigned
	// Format16BitLESigned is for signed, little-endian, 16-bit data
	Format16BitLESigned
	// Format16BitBEUnsigned is for unsigned, big-endian, 16-bit data
	Format16BitBEUnsigned
	// Format16BitBESigned is for signed, big-endian, 16-bit data
	Format16BitBESigned
	// Format32BitLEFloat is for little-endian, 32-bit floating-point data
	Format32BitLEFloat
	// Format32BitBEFloat is for big-endian, 32-bit floating-point data
	Format32BitBEFloat
	// Format64BitLEFloat is for little-endian, 64-bit floating-point data
	Format64BitLEFloat
	// Format64BitBEFloat is for big-endian, 64-bit floating-point data
	Format64BitBEFloat
)
