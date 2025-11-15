package model

type Rule struct {
	Action    string // prefix, suffix, replace, extension, lowercase, uppercase
	Parameter string
}

type RenameResult struct {
	OldName string
	NewName string
	Success bool
	Error   error
}
