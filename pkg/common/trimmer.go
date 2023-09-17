package common

type Trimmer struct {
	output string
}

// TODO:fix trimmer
func (t *Trimmer) Trim(input string, maxLen, swaps int) *Trimmer {
	toTrim := []rune(input)

	if len(toTrim) > maxLen {

		toTrim = toTrim[0:maxLen]
		toTrim = t.SwapEnd(toTrim, swaps)

	}

	t.output = string(toTrim)
	return t
}

func (t *Trimmer) TrimEnd(input string, maxLen, swaps int) *Trimmer {
	toTrim := []rune(input)

	if len(toTrim) > maxLen {

		toTrim = toTrim[0:maxLen]
		toTrim = t.SwapEnd(toTrim, swaps)

	}

	t.output = string(toTrim)
	return t
}

func (t *Trimmer) TrimMiddle(input string, maxLen, swaps int) *Trimmer {
	toTrim := []rune(input)

	if len(toTrim) > maxLen {

		toTrim = toTrim[0:maxLen]
		toTrim = t.SwapEnd(toTrim, swaps)

	}

	t.output = string(toTrim)
	return t
}

func (t *Trimmer) TrimStart(input string, maxLen, swaps int) *Trimmer {
	toTrim := []rune(input)

	if len(toTrim) > maxLen {

		toTrim = toTrim[0:maxLen]
		toTrim = t.SwapEnd(toTrim, swaps)

	}

	t.output = string(toTrim)
	return t
}

func (t *Trimmer) SwapEnd(chars []rune, swaps int) []rune {
	swapFormat := '.'

	for swaps > 0 {

		chars[len(chars)-swaps] = swapFormat

		swaps--
	}

	return chars
}

func (t *Trimmer) String() string {
	return t.output
}

func (t *Trimmer) Swap(chars []rune, swaps int) []rune {
	swapFormat := '.'

	for swaps > 0 {

		chars[len(chars)-swaps] = swapFormat

		swaps--
	}

	return chars
}
