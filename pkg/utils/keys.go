package utils

import (
	k "github.com/eiannone/keyboard"
)

var (
	execKey   = k.KeyCtrlE
	deleteKey = k.KeyDelete

	deleteKey3 = k.KeyCtrlD
	acceptKey  = k.KeyEnter
	enterKey   = acceptKey
	rejectKey  = k.KeyEsc

	renameKey = k.KeyCtrlR
	quitKey   = rejectKey

	spaceKey   = k.KeySpace
	newFileKey = k.KeyCtrlW

	findKey     = k.KeyCtrlF
	findSizeKey = k.KeyCtrlL
	sortKey     = k.KeyCtrlS

	menuKey = k.KeyF2
	viewKey = k.KeyF3
	editKey = k.KeyF4
	copyKey = k.KeyCtrlV

	previewKey = k.KeyCtrlQ

	moveKey       = k.KeyF6
	newDirKey     = k.KeyF7
	deleteKey2    = k.KeyF8
	scanKey       = k.KeyF9
	quitKey2      = k.KeyF10
	goToKey       = k.KeyCtrlG
	homeKey       = k.KeyHome
	endKey        = k.KeyEnd
	pgDownKey     = k.KeyPgdn
	pgUpKey       = k.KeyPgup
	downKey       = k.KeyArrowDown
	leftKey       = k.KeyArrowLeft
	upKey         = k.KeyArrowUp
	rightKey      = k.KeyArrowRight
	selectKey     = k.KeyInsert
	selectKey2    = k.KeyCtrlI
	selectAllKey  = k.KeyCtrlA
	backSpaceKey  = k.KeyBackspace
	backSpaceKey2 = k.KeyBackspace2
	nextThemeKey  = k.KeyCtrlSlash

	rejectChar = 'n'
	acceptChar = 'y'
	quitChar   = 'q'
)
