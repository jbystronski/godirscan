package converter

import "math"

var (
	StorageUnits                            = [4]string{"  B", "KiB", "MiB", "GiB"}
	GbInBytes, MbInBytes, KbInBytes float64 = math.Pow(2, 30), 1024 * 1024, 1024
)

func convert(mult, pow float64) int64 {
	return int64(mult * math.Pow(1024, pow))
}

func ToBytes(unit int, min, max float64) (int64, int64) {
	var pow float64

	if unit >= 1 && unit <= 3 {
		pow = float64(unit)
		return convert(min, pow), convert(max, pow)
	} else {
		return int64(min), int64(max)
	}
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
