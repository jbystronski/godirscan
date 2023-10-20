package converter

import "math"

var (
	StorageUnits                            = [4]string{"  B", "KiB", "MiB", "GiB"}
	GbInBytes, MbInBytes, KbInBytes float64 = math.Pow(2, 30), 1024 * 1024, 1024
)

func ToBytes(unit int, min, max float64) (int64, int64) {
	convert := func(multiplier, pow float64) int64 {
		return int64(multiplier * math.Pow(1024, pow))
	}

	var pow float64

	switch unit {
	case 3:
		{
			pow = 3
		}
	case 2:
		{
			pow = 2
		}
	case 1:
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

	switch true {
	case size >= GbInBytes:
		return size / math.Pow(1024, 3), StorageUnits[3]

	case size >= MbInBytes:
		return size / math.Pow(1024, 2), StorageUnits[2]

	case size >= KbInBytes:
		return size / KbInBytes, StorageUnits[1]

	default:

		return size, StorageUnits[0]
	}
}

type Converter struct{}
