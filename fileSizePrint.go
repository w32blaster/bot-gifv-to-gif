package main

// Copied file from GitHub :)

import (
	"fmt"
)

// ByteSize is a datatype to have human readable byte sizes.
type ByteSize float64

const (
	_ = iota // ignore first value by assigning to blank identifier
	// KB defines a KiloByte
	KB ByteSize = 1 << (10 * iota)
	// MB defines a MegaByte
	MB
	// GB defines a GigaByte
	GB
	// TB defines a TerraByte
	TB
	// PB defines a PetaByte
	PB
	// EB defines an ExoByte
	EB
	// ZB defines a ZetaByte
	ZB
	// YB defines a YotaByte
	YB
)

// (b ByteSize) String() is the Stringer function for ByteSize.
// It returns a human readable string for the given size of bytes
func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}
