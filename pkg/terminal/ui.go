package terminal

import (
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/entry"
)

func PrintEmpty() {
	fmt.Printf("%s%s", "folder is empty", ResetFmt)
	// moveToEndOfLine()
	// printBorder()
}

func Prompt(prompt string) string {
	return fmt.Sprintf("%s%s %s %s", CurrentTheme.BgPrompt, CurrentTheme.Prompt+BoldFmt, prompt, ResetFmt)
}

func highlightRow(sep string, en entry.Entry) {
	fmt.Printf("%s%s%s%s%s\n", CurrentTheme.BgHighlight, CurrentTheme.Highlight+BoldFmt, sep, en.Name+Space+en.PrintSize(), ResetFmt)
	// moveToEndOfLine()
	// printBorder()
}

func MarkRow(sep string, en entry.Entry) {
	fmt.Printf("%s%s%s%s%s\n", CurrentTheme.BgSelect, CurrentTheme.Select+BoldFmt, sep, en.Name+Space+en.PrintSize(), ResetFmt)
	// printRightBorder()
	// moveToEndOfLine()
	// printBorder()
}

func printRow(sep string, entry entry.Entry) {
	if entry.IsDir {
		fmt.Printf("%s%s%s\n", strings.Join([]string{ColorsMap["bright_white"], BoldFmt, sep, CurrentTheme.Main, entry.Name, ColorsMap["bright_white"], Space}, ""), entry.PrintSize(), ResetFmt)
	} else {
		fmt.Printf("%s%s%s\n", strings.Join([]string{ColorsMap["bright_white"], sep, CurrentTheme.Accent, entry.Name, BoldFmt, ColorsMap["bright_white"], Space}, ""), entry.PrintSize(), ResetFmt)
	}

	// printRightBorder()
}

func r(s string, times int) string {
	return strings.Repeat(s, times)
}

func PrintBanner() {
	fmt.Print(CurrentTheme.Main)

	fmt.Print("\u2554" + r(DbHoriz, 90) + "\u2557\n")
	fmt.Print(DbVert + r(Space, 90) + DbVert + "\n")

	fmt.Print(DbVert + r(" ", 6) + Fslash + r(Texture, 8) + Fslash + r(Space, 74) + DbVert + "\n")

	fmt.Print(DbVert + r(" ", 5) + Segment + r(Line, 6) + r(Space, 2) + Fslash + r(Texture, 7) + Fslash + r(Space, 1) + Fslash + r(Texture, 6) + Bslash + r(Space, 2) + Fslash + r(Texture, 1) + Fslash + r(Space, 1) + Fslash + r(Texture, 6) + Bslash + r(Space, 2) + Fslash + r(Texture, 6) + Fslash + r(Space, 1) + Fslash + r(Texture, 6) + Fslash + r(Space, 1) + Fslash + r(Texture, 7) + Fslash + r(Space, 1) + Fslash + r(Texture, 2) + Bslash + r(Space, 2) + Fslash + r(Texture, 1) + Fslash + r(Space, 3) + DbVert + "\n")
	fmt.Print(DbVert + r(" ", 4) + Segment + r(Space, 2) + Fslash + r(Texture, 3) + Fslash + r(Space, 1) + Segment + r(Line, 3) + Segment + Space + Segment + r(Line, 3) + Segment + r(Space, 1) + Segment + Space + Segment + r(Line, 3) + Segment + r(Space, 1) + Fslash + r(Texture, 5) + Bslash + r(Space, 2) + Segment + r(Line, 5) + r(Space, 1) + Segment + r(Line, 3) + Segment + r(Space, 1) + Segment + Bslash + Texture + Bslash + Segment + r(Space, 4) + DbVert + "\n")
	fmt.Print(DbVert + r(" ", 3) + Segment + r(Space, 3) + r(Line, 1) + Segment + Space + Segment + r(Space, 3) + Segment + Space + Segment + r(Space, 3) + Segment + r(Space, 1) + Segment + Space + Fslash + r(Texture, 7) + r(Line, 0) + r(Space, 0) + Fslash + r(Space, 2) + r(Line, 4) + Fslash + Texture + Fslash + Space + Segment + r(Space, 6) + Fslash + r(Texture, 7) + Fslash + Space + Segment + r(Space, 2) + Bslash + r(Texture, 2) + Fslash + r(Space, 5) + DbVert + "\n")
	fmt.Print(DbVert + r(" ", 2) + Fslash + r(Texture, 8) + Fslash + Space + Fslash + r(Texture, 7) + Fslash + Space + Fslash + r(Texture, 7) + Fslash + Space + Fslash + Texture + Fslash + r(Space, 1) + Segment + r(Line, 4) + Bslash + Texture + Bslash + r(Space, 1) + Fslash + r(Texture, 5) + Fslash + Space + Fslash + r(Texture, 6) + Fslash + Space + Segment + r(Line, 3) + Segment + Space + Segment + r(Space, 3) + Segment + r(Space, 6) + DbVert + "\n")

	fmt.Print(DbVert + r(" ", 2) + r(Line, 9) + r(Space, 2) + r(Line, 8) + r(Space, 2) + r(Line, 8) + r(Space, 2) + Line + Line + r(Space, 2) + r(Line, 2) + r(Space, 6) + r(Line, 2) + Space + r(Line, 6) + r(Space, 2) + r(Line, 7) + r(Space, 2) + r(Line, 2) + r(Space, 4) + r(Line, 2) + r(Space, 2) + r(Line, 2) + r(Space, 4) + r(Line, 2) + r(Space, 7) + DbVert + "\n")

	fmt.Print("\u255A" + r(DbHoriz, 90) + "\u255D")

	fmt.Print(ResetFmt)

	fmt.Print("\n\n")
}

func printHeader(header string) {
	fmt.Printf("%s%s", CurrentTheme.BgHeader+CurrentTheme.Header+" F2 - options "+Space+header, ResetFmt)
}

// func printOptions(options []Option) {
// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

// 	for _, opt := range options {
// 		fmt.Fprintf(w, "%s\t%s\n", theme.Main+BoldFmt+opt.Command+ResetFmt, theme.Accent+strings.Join(opt.Description, " ")+ResetFmt)
// 	}

// 	w.Flush()
// }

func PrintHelp() {
	//for i, s := range selected {
	//	fmt.Printf("%d - %s", i, s)
	//}

	PrintBanner()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println(err)
		}

		if key == keyboard.KeyEsc {
			return
		}

	}

	// printOptions(optionsList)
}

func printBox(totalLines, totalCols, boxWidth, boxHeight int) {
}
