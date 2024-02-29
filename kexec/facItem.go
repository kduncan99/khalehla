// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

// FacilitiesItem structs are store in the RCE for all assigned facilities

type FacilitiesItem interface {
	GetInternalFileName() string
	GetFileName() string
	GetQualifier() string
	GetEquipmentCode() uint
}

/*
All facility items:
+00,W   internal file Name - Fieldata LJSF
+01,W   (internal file Name cont)
+02,W   file Name - Fieldata LJSF
+03,W   (file Name cont)
+04,W   qualifier - Fieldata LJSF
+05,W   (qualifier cont)
+06,S1  equipment code
         000 file has not been assigned (@USE exists, but @ASG has not been done)
         015 9-track tape
         016 virtual tape handler
         017 cartridge tape, DVD tape
         024 word-addressable mass storage
         036 sector-addressable mass storage
         077 arbitrary device
+07,S1  attributes
+07,b10:b35 @ASG options

Unit record and non-standard peripherals
+07,S1  attributes
         040 tape labeling is supported
         020 file is temporary
         010 internal Name is a use Name

Sector-formatted mass storage
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 file is write inhibited
         002 file is read inhibited
         001 word-addressable (always clear)
+06,S3  granularity
         zero -> track, nonzero -> position
+06,S4  relative file-cycle
+06,T3  absolute file-cycle
+07,S1  attributes
         020 file is temporary
         010 internal Name is a use Name
         004 shared file
         002 large file
+010,H1 initial granule count (initial reserve)
+010,H2 max granule count
+011,H1 highest track referenced
+011,H2 highest granule assigned
+012,S4 total pack count if removable (63 -> 63 or greater)
+012,S5 equipment code - same as +06,S1
+012,S6 subcode - zero

Magnetic tape peripherals
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 file is write inhibited
         002 file is read inhibited
+06,S3  unit count (I presume, the docs are not helpful)
		number of units assigned (0?, 1, 2)
+07,S1  attributes
         040 tape labeling is supported
         020 file is temporary
         010 internal Name is a use Name
         004 file is a shared file
+010,S1 total reel count
+010,S2 logical channel
+010,S3 noise constant
+012,T1 expiration period
+012,S3 reel index
+012,S4 files extended
+012,T3 blocks extended
+013,W  current reel number
+014,W  next reel number

Word addressable
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 write inhibited
         002 read inhibited
         001 word-addressable (always set)
+06,S3  granularity
         zero -> track, nonzero -> position
+06,S4  relative file-cycle
+06,T3  absolute file-cycle
+07,S1  attributes
         020 file is temporary
         010 internal Name is a use Name
         004 shared file
+010,W  length of file in words
+011,W  maximum file length in words
+012,S4 total pack count if removable (63 -> 63 or greater)
+012,S5 equipment code - same as +06,S1
+012,S6 subcode - zero
*/
