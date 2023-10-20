package viewbox

func (v FsViewBox) Separator(isLast bool) string {
	if isLast {
		return "\u2514\u2500"
	}
	return "\u251c\u2500"
}
