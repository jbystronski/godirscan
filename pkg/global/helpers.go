package global

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/app/config"

	"github.com/jbystronski/godirscan/pkg/lib/converter"
	t "github.com/jbystronski/godirscan/pkg/lib/termui"
)

const (
	Bold  = t.Bold
	Reset = t.Reset
	Space = t.Space
)

func Rows() int {
	return t.NewTerminal().Rows()
}

func Cols() int {
	return t.NewTerminal().Cols()
}

func Cell(x, y int) {
	t.Cell(x, y)
}

func Clear(x, y, len int) {
	t.ClearRow(x, y, len)
}

func Truncate(path string) error {
	err := os.Truncate(path, 0)
	if err != nil {
		return err
	}
	return nil
}

func ClearScreen() {
	t.ClearScreen()
}

func TrimEnd(input string, maxLen, trimLen, swaps int, swapChar rune) string {
	return t.TrimEnd(input, maxLen, trimLen, swaps, swapChar)
}

func BuildString(substrings ...string) string {
	return t.BuildString(substrings...)
}

func FmtBold(strings ...string) string {
	st := BuildString(strings...)

	return fmt.Sprint(Bold, st, Reset)
}

func ThemeMain() string {
	return config.CurrentTheme.Main
}

func ThemeAccent() string {
	return config.CurrentTheme.Accent
}

func ThemeBgHighlight() string {
	return config.CurrentTheme.BgHighlight
}

func ThemeHighlight() string {
	return config.CurrentTheme.Highlight
}

func ThemeBgHeader() string {
	return config.CurrentTheme.BgHeader
}

func ThemeHeader() string {
	return config.CurrentTheme.Header
}

func ThemeBgSelect() string {
	return config.CurrentTheme.BgSelect
}

func ThemeBgPrompt() string {
	return config.CurrentTheme.BgPrompt
}

func ThemePrompt() string {
	return config.CurrentTheme.Prompt
}

func ThemeSelect() string {
	return config.CurrentTheme.Select
}

func StrLen(st string) int {
	return t.StrLen(st)
}

func HideCurson() {
	t.HideCursor()
}

func ShowCursor() {
	t.ShowCursor()
}

func AlignCenter(maxWidth int, st, padding string) string {
	return t.AlignCenter(maxWidth, st, padding)
}

func AlignLeft(maxWidth int, st, padding string) string {
	return t.AlignLeft(maxWidth, st, padding)
}

func AlignRight(maxWidth int, st, padding string) string {
	return t.AlignRight(maxWidth, st, padding)
}

func UpdateDimensions() error {
	return t.NewTerminal().UpdateDimensions()
}

func FmtPrompt(prompt string) string {
	return FmtBold(ThemeBgPrompt(), ThemePrompt(), prompt)
}

func OpenCommandLine() {
	t.NewTerminal().OpenCommandLine()
}

func CloseCommandLine() {
	t.NewTerminal().CloseCommandLine()
}

func ReceiveCommandLine() <-chan (struct {
	Key  keyboard.Key
	Char rune
}) {
	return t.NewTerminal().ReceiveCommandLine()
}

func FormatSize(bytes int) string {
	const sizeMaxLen = 11

	if bytes < int(converter.KbInBytes) {
		st := fmt.Sprintf("%d %s", bytes, converter.StorageUnits[0]+" ")
		if len(st) < sizeMaxLen {
			st = strings.Repeat(" ", sizeMaxLen-len(st)) + st
		}
		return FmtBold(Space, st)

	}

	floatSize, unit := converter.BytesToFloat(bytes)

	st := fmt.Sprintf("%.1f %s", floatSize, unit+" ")

	if len(st) < sizeMaxLen {
		st = strings.Repeat(" ", sizeMaxLen-len(st)) + st
	}
	return FmtBold(Space, st)
}

func IsSymlink(info fs.FileInfo) bool {
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}

	return false
}

var virtualFsMap = map[string]struct{}{
	"/proc": {},
	"/dev":  {},
	"/sys":  {},
}

func IsVirtualFs(dirName string) bool {
	//	name := strings.Split(dirName, string(os.PathSeparator))[0]

	if _, ok := virtualFsMap[dirName]; ok {
		return true
	}

	return false
}

type DirReader struct{}

func (d *DirReader) Read(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func GetParentDirectory(dir string) (string, bool) {
	if GetRootDirectory() == dir {
		return "", false
	}

	parent, _ := filepath.Split(dir)
	parent = strings.TrimSuffix(parent, string(filepath.Separator))

	if parent == "" {
		parent = GetRootDirectory()
	}

	return parent, true
}

func GetRootDirectory() string {
	wd, _ := os.Getwd()

	return filepath.VolumeName(wd) + string(filepath.Separator)
}

func ReadIgnorePermission(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsPermission(err) {
			return []fs.DirEntry{}, nil
		} else {
			return nil, err
		}
	}

	return entries, nil
}

func ResolveUserDirectory(path *string) {
	switch runtime.GOOS {
	case "darwin":
		{
			break
		}
	case "windows":
		{
			break
		}
	default:
		{
			if strings.HasPrefix(*path, "~") {
				currentUser, err := user.Current()
				if err != nil {
					panic(err)
				}

				*path = strings.Replace(*path, "~", currentUser.HomeDir, 1)

			}
		}
	}
}

func Search(path, pattern string, entryChan chan<- struct {
	Path string
	Info fs.FileInfo
}, test func(fs.FileInfo) bool,
) error {
	if dirEntries, err := ReadIgnorePermission(path); err != nil {
		return err
	} else {
		for _, en := range dirEntries {
			info, err := en.Info()
			if err != nil {
				return err
			}

			if info.IsDir() {
				err := Search(filepath.Join(path, en.Name()), pattern, entryChan, test)
				if err != nil {
					return err
				}

			} else {
				if strings.Contains(info.Name(), pattern) {
					if ok := test(info); ok {
						entryChan <- struct {
							Path string
							Info fs.FileInfo
						}{filepath.Join(path, info.Name()), info}
					}
				}
			}

		}
	}
	return nil
}

func Copy(sourcePath, targetPath string) error {
	sourceFile, err := os.Open(sourcePath)
	defer sourceFile.Close()

	if err != nil {
		return err
	}

	targetFile, err := os.Create(targetPath)

	defer targetFile.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return err
	}

	info, _ := os.Stat(sourcePath)

	err = os.Chmod(targetPath, info.Mode())
	if err != nil {
		return err
	}
	return nil
}

type CancelContext struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

func NewCancelContext() *CancelContext {
	ctx := &CancelContext{}
	ctx.Create()

	return ctx
}

func (c *CancelContext) Create() {
	c.Ctx, c.CancelFunc = context.WithCancel(context.Background())
}

func (c *CancelContext) Cancel() {
	c.CancelFunc()
}

func (c *CancelContext) Observe(task func()) {
	var wg sync.WaitGroup
	done := make(chan struct{}, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {

			case <-done:

				return

			case <-c.Ctx.Done():

				return
			//	done <- struct{}{}
			default:

				task()
				done <- struct{}{}

			}
		}
	}()

	wg.Wait()
}

// type Coords struct {
// 	Y int
// 	X int
// }
