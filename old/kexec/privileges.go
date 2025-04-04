package kexec

type Privilege uint

const (
	_ Privilege = iota
	DLOCPrivilege
	SSConsolePrivilege
	UnconditionalFileSystemDelete
)
