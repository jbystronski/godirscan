package utils

import (
	"fmt"
	"strings"

	"github.com/jbystronski/godirscan/pkg/entry"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func printEmpty() {
	fmt.Printf("%s%s", "folder is empty", terminal.ResetFmt)
}

func fPrompt(prompt string) string {
	return fmt.Sprintf("%s%s %s %s", theme.BgPrompt, terminal.ColorsMap["white"]+terminal.BoldFmt, prompt, terminal.ResetFmt)
}

func highlightRow(sep string, en entry.Entry) {
	fmt.Printf("%s%s%s%s%s\n", theme.BgHighlight, terminal.ColorsMap["black"]+terminal.BoldFmt, sep, en.Name+terminal.Space+en.PrintSize(), terminal.ResetFmt)
}

func markRow(sep string, en entry.Entry) {
	fmt.Printf("%s%s%s%s%s\n", theme.BgPrompt, terminal.ColorsMap["black"]+terminal.BoldFmt, sep, en.Name+terminal.Space+en.PrintSize(), terminal.ResetFmt)
}

func r(s string, times int) string {
	return strings.Repeat(s, times)
}

func printBanner() {
	fmt.Print(theme.Main)

	fmt.Print("\u2554" + r(terminal.DbHoriz, 90) + "\u2557\n")
	fmt.Print(terminal.DbVert + r(terminal.Space, 90) + terminal.DbVert + "\n")

	fmt.Print(terminal.DbVert + r(" ", 6) + terminal.Fslash + r(terminal.Texture, 8) + terminal.Fslash + r(terminal.Space, 74) + terminal.DbVert + "\n")

	fmt.Print(terminal.DbVert + r(" ", 5) + terminal.Segment + r(terminal.Line, 6) + r(terminal.Space, 2) + terminal.Fslash + r(terminal.Texture, 7) + terminal.Fslash + r(terminal.Space, 1) + terminal.Fslash + r(terminal.Texture, 6) + terminal.Bslash + r(terminal.Space, 2) + terminal.Fslash + r(terminal.Texture, 1) + terminal.Fslash + r(terminal.Space, 1) + terminal.Fslash + r(terminal.Texture, 6) + terminal.Bslash + r(terminal.Space, 2) + terminal.Fslash + r(terminal.Texture, 6) + terminal.Fslash + r(terminal.Space, 1) + terminal.Fslash + r(terminal.Texture, 6) + terminal.Fslash + r(terminal.Space, 1) + terminal.Fslash + r(terminal.Texture, 7) + terminal.Fslash + r(terminal.Space, 1) + terminal.Fslash + r(terminal.Texture, 2) + terminal.Bslash + r(terminal.Space, 2) + terminal.Fslash + r(terminal.Texture, 1) + terminal.Fslash + r(terminal.Space, 3) + terminal.DbVert + "\n")
	fmt.Print(terminal.DbVert + r(" ", 4) + terminal.Segment + r(terminal.Space, 2) + terminal.Fslash + r(terminal.Texture, 3) + terminal.Fslash + r(terminal.Space, 1) + terminal.Segment + r(terminal.Line, 3) + terminal.Segment + terminal.Space + terminal.Segment + r(terminal.Line, 3) + terminal.Segment + r(terminal.Space, 1) + terminal.Segment + terminal.Space + terminal.Segment + r(terminal.Line, 3) + terminal.Segment + r(terminal.Space, 1) + terminal.Fslash + r(terminal.Texture, 5) + terminal.Bslash + r(terminal.Space, 2) + terminal.Segment + r(terminal.Line, 5) + r(terminal.Space, 1) + terminal.Segment + r(terminal.Line, 3) + terminal.Segment + r(terminal.Space, 1) + terminal.Segment + terminal.Bslash + terminal.Texture + terminal.Bslash + terminal.Segment + r(terminal.Space, 4) + terminal.DbVert + "\n")
	fmt.Print(terminal.DbVert + r(" ", 3) + terminal.Segment + r(terminal.Space, 3) + r(terminal.Line, 1) + terminal.Segment + terminal.Space + terminal.Segment + r(terminal.Space, 3) + terminal.Segment + terminal.Space + terminal.Segment + r(terminal.Space, 3) + terminal.Segment + r(terminal.Space, 1) + terminal.Segment + terminal.Space + terminal.Fslash + r(terminal.Texture, 7) + r(terminal.Line, 0) + r(terminal.Space, 0) + terminal.Fslash + r(terminal.Space, 2) + r(terminal.Line, 4) + terminal.Fslash + terminal.Texture + terminal.Fslash + terminal.Space + terminal.Segment + r(terminal.Space, 6) + terminal.Fslash + r(terminal.Texture, 7) + terminal.Fslash + terminal.Space + terminal.Segment + r(terminal.Space, 2) + terminal.Bslash + r(terminal.Texture, 2) + terminal.Fslash + r(terminal.Space, 5) + terminal.DbVert + "\n")
	fmt.Print(terminal.DbVert + r(" ", 2) + terminal.Fslash + r(terminal.Texture, 8) + terminal.Fslash + terminal.Space + terminal.Fslash + r(terminal.Texture, 7) + terminal.Fslash + terminal.Space + terminal.Fslash + r(terminal.Texture, 7) + terminal.Fslash + terminal.Space + terminal.Fslash + terminal.Texture + terminal.Fslash + r(terminal.Space, 1) + terminal.Segment + r(terminal.Line, 4) + terminal.Bslash + terminal.Texture + terminal.Bslash + r(terminal.Space, 1) + terminal.Fslash + r(terminal.Texture, 5) + terminal.Fslash + terminal.Space + terminal.Fslash + r(terminal.Texture, 6) + terminal.Fslash + terminal.Space + terminal.Segment + r(terminal.Line, 3) + terminal.Segment + terminal.Space + terminal.Segment + r(terminal.Space, 3) + terminal.Segment + r(terminal.Space, 6) + terminal.DbVert + "\n")

	fmt.Print(terminal.DbVert + r(" ", 2) + r(terminal.Line, 9) + r(terminal.Space, 2) + r(terminal.Line, 8) + r(terminal.Space, 2) + r(terminal.Line, 8) + r(terminal.Space, 2) + terminal.Line + terminal.Line + r(terminal.Space, 2) + r(terminal.Line, 2) + r(terminal.Space, 6) + r(terminal.Line, 2) + terminal.Space + r(terminal.Line, 6) + r(terminal.Space, 2) + r(terminal.Line, 7) + r(terminal.Space, 2) + r(terminal.Line, 2) + r(terminal.Space, 4) + r(terminal.Line, 2) + r(terminal.Space, 2) + r(terminal.Line, 2) + r(terminal.Space, 4) + r(terminal.Line, 2) + r(terminal.Space, 7) + terminal.DbVert + "\n")

	fmt.Print("\u255A" + r(terminal.DbHoriz, 90) + "\u255D")

	fmt.Print(terminal.ResetFmt)

	fmt.Print("\n\n")
}

func printRow(sep string, entry entry.Entry) {
	if entry.IsDir {
		fmt.Printf("%s%d%s\n", strings.Join([]string{terminal.ColorsMap["bright_white"], terminal.BoldFmt, sep, theme.Main, entry.Name, terminal.ColorsMap["bright_white"], terminal.Space}, ""), entry.Size, terminal.ResetFmt)
	} else {
		fmt.Printf("%s%d%s\n", strings.Join([]string{terminal.ColorsMap["bright_white"], sep, theme.Accent, entry.Name, terminal.BoldFmt, terminal.ColorsMap["bright_white"], terminal.Space}, ""), entry.Size, terminal.ResetFmt)
	}
}

func printHeader(header string) {
	fmt.Printf("%s%s%s%s\n", theme.BgPrompt+terminal.ColorsMap["black"]+" F2 - options "+terminal.ResetFmt, terminal.BackgroundsMap["yellow"], terminal.ColorsMap["black"]+terminal.BoldFmt+terminal.Space+header, terminal.ResetFmt)
}

// func printOptions(options []Option) {
// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

// 	for _, opt := range options {
// 		fmt.Fprintf(w, "%s\t%s\n", theme.Main+terminal.BoldFmt+opt.Command+terminal.ResetFmt, theme.Accent+strings.Join(opt.Description, " ")+terminal.ResetFmt)
// 	}

// 	w.Flush()
// }

func printHelp() {
	//for i, s := range selected {
	//	fmt.Printf("%d - %s", i, s)
	//}

	printBanner()
	// printOptions(optionsList)
}

func printBox(totalLines, totalCols, boxWidth, boxHeight int) {
}
