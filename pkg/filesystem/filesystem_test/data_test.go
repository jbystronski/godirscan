package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func initData() *filesystem.FsData {
	fsData := filesystem.NewFsData()

	fileNames := []string{"test_entry0", "test_entry1", "test_entry2"}
	sizes := []int{12, 202, 33}
	path, _ := os.Getwd()

	for k, v := range fileNames {

		en := &filesystem.FsEntry{}
		en.SetName(filepath.Join(path, v))
		en.SetSize(sizes[k])
		en.SetFsType(filesystem.File)
		fsData.Insert(en)

	}

	return fsData
}

func TestLen(t *testing.T) {
	data := initData()

	want := 3

	got := data.Len()

	if want != got {
		t.Errorf("Len = %v, want %v, got %v", got, want, got)
	}
}

func TestFind(t *testing.T) {
	data := initData()

	var nameTest func(int, string)

	nameTest = func(index int, want string) {
		en, _ := data.Find(index)
		got := en.Name()
		if got != want {
			t.Errorf("en.Name = %v, want %v", got, want)
		}
	}

	t.Run("test valid entry name", func(t *testing.T) { nameTest(1, "test_entry1") })
}

func TestSize(t *testing.T) {
	data := initData()

	got := data.Size()

	want := 12 + 202 + 33

	if want != got {
		t.Errorf("Size = %v, want %v, got %v", got, want, got)
	}
}

func TestFullPath(t *testing.T) {
	data := initData()

	entry, _ := data.Find(0)

	path, _ := os.Getwd()

	got := entry.FullPath()

	want := filepath.Join(path, entry.Name())

	if want != got {
		t.Errorf("FullPath test want %v, got %v", want, got)
	}
}

func TestPath(t *testing.T) {
	data := initData()

	entry, _ := data.Find(0)

	path, _ := os.Getwd()

	got := entry.Path()

	want := path

	if want != got {
		t.Errorf("Path test want %v, got %v", want, got)
	}
}

func TestReset(t *testing.T) {
	data := initData()

	data.Reset()
	got := data.Len()

	want := 0

	if want != got {
		t.Errorf("Reset test, Len after reset, want %v, got %v", want, got)
	}
}

func TestUpdateOne(t *testing.T) {
	data := initData()

	en, _ := data.Find(1)

	en.SetSize(100)

	want := 100

	updatedEn, _ := data.Find(1)

	got := updatedEn.Size()

	if want != got {
		t.Errorf("Test update entry, Size after update, want %v, got %v", want, got)
	}
}

func TestUpdateAll(t *testing.T) {
	data := initData()

	for _, v := range data.All() {
		v.SetSize(300)
	}

	en, _ := data.Find(1)

	want := 300

	got := en.Size()

	if want != got {
		t.Errorf("Test update entry, Size after update, want %v, got %v", want, got)
	}
}
