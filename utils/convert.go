package utils

import "fmt"

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

//ConvertBytes converts a bytes entry into appropriate unit (KB, MB,, GB, or TB) depending on how big it is
func ConvertBytes(b uint64) string {
	switch {
		case b < KILOBYTE:
			return fmt.Sprintf("%dB", b)
		case b > KILOBYTE && b < MEGABYTE:
			return fmt.Sprintf("%dKB", b/KILOBYTE)
		case b > MEGABYTE && b < GIGABYTE:
			return fmt.Sprintf("%dMB", b/MEGABYTE)
		case b > GIGABYTE && b < TERABYTE:
			return fmt.Sprintf("%dGB", b/GIGABYTE)
		case b > TERABYTE:
			return fmt.Sprintf("%dTB", b/TERABYTE)
	}
	return fmt.Sprintf("%dB", b)
}

