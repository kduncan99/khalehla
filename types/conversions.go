package types

var PackedBytesPerBlockFromWords = map[BlockSize]BlockSize{
	28:   128,  // slop 2 bytes
	56:   256,  // slop 4 bytes
	112:  512,  // slop 8 bytes
	224:  1024, // slop 16 bytes
	448:  2048, // slop 32 bytes
	896:  4096, // slop 64 bytes
	1792: 8192, // slop 128 bytes
}

var RawBytesPerBlockFromWords = map[BlockSize]BlockSize{
	28:   28 * 8,
	56:   56 * 8,
	112:  112 * 8,
	224:  224 * 8,
	448:  448 * 8,
	896:  896 * 8,
	1792: 1792 * 8,
}

var AsciiFromFieldata = []byte{
	'@', '[', ']', '#', '^', ' ', 'A', 'B',
	'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
	'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	')', '-', '+', '<', '=', '>', '&', '$',
	'*', '(', '%', ':', '?', '!', ',', '\\',
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', '\'', ';', '/', '.', '"', '_',
}

var FieldataFromAscii = []int{
	005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005,
	005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005,
	005, 055, 076, 003, 047, 052, 046, 072, 051, 040, 050, 042, 056, 041, 075, 074,
	060, 061, 062, 063, 064, 065, 066, 067, 070, 071, 053, 073, 043, 044, 045, 054,
	000, 006, 007, 010, 011, 012, 013, 014, 015, 016, 017, 020, 021, 022, 023, 024,
	025, 026, 027, 030, 031, 032, 033, 034, 035, 036, 037, 001, 057, 060, 004, 077,
	000, 006, 007, 010, 011, 012, 013, 014, 015, 016, 017, 020, 021, 022, 023, 024,
	025, 026, 027, 030, 031, 032, 033, 034, 035, 036, 037, 054, 057, 055, 004, 077,
}
