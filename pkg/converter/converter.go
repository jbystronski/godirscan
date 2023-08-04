package converter

import "math"

var (
	StorageUnits                            = [4]string{"bytes", "Kb", "MB", "GB"}
	GbInBytes, MbInBytes, KbInBytes float64 = math.Pow(2, 30), 1048576, 1024
)

func ToBytes(from string, min, max float64) (int64, int64) {
	convert := func(multiplier, pow float64) int64 {
		return int64(multiplier * math.Pow(1024, pow))
	}

	var pow float64

	switch from {
	case StorageUnits[3]:
		{
			pow = 3
		}
	case StorageUnits[2]:
		{
			pow = 2
		}
	case StorageUnits[1]:
		{
			pow = 1
		}
	default:
		{
			return int64(min), int64(max)
		}
	}

	return convert(min, pow), convert(max, pow)
}

func BytesToFloat(bytes int) (float64, string) {
	size := float64(bytes)

	if size >= GbInBytes {
		return size / math.Pow(1024, 3), StorageUnits[3]
	}

	if size >= MbInBytes {
		return size / math.Pow(1024, 2), StorageUnits[2]
	}

	if size >= KbInBytes {
		return size / KbInBytes, StorageUnits[1]
	}

	return size, StorageUnits[0]
}
