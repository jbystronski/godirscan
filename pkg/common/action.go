package common

type ControllerAction int

const (
	Delete ControllerAction = iota
	Move
	Copy
	Execute
	ExecuteCmd
	Rename
	CreateDirectory
	CreateFile
	Sort
	Help
	GoTo
	Find
	Scan
	Edit
	Select
	SelectAll
	MoveUp
	MoveDown
	MoveLeft
	MoveRight
	MoveToBottom
	MoveToTop
	PageUp
	PageDown
	SwitchController
	Resize
	Render
	Refresh
)
