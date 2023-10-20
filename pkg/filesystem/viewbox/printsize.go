package viewbox

import "fmt"

func (v FsViewBox) PrintSizeAsString(size int) string {
	return fmt.Sprintf("%v", FormatSize(size))
}
