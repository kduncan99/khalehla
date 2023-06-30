package fileSystem

/*
Assigning Mass Storage Files
@ASG[,options] filename
               [,type/reserve/granule/maximum/placement]
               [,pack-id-1/.../pack-id-n,,,ACR-name]

Changing Assignment Fields for Mass Storage Files
You can change the type, reserve, maximum, or placement fields for a file with an
@ASG statement when the following applies:
 You have already assigned the file with the A option. Change by adding the X and
I options. You can assign the file again with the same options and specify a new
type, reserve, maximum, or placement.
 You have already assigned the file with a C or U option. You can assign the file
again with the same option and specify a new type, reserve, maximum, or
placement.
 You have already cataloged the file but not assigned it. You can assign the file
with an @ASG,A statement and give new specifications for the type, reserve,
maximum, or placement.
To change the subfields that can be changed, the file must be write enabled.
You cannot change the granule subfield.
The format for changing assignment fields is
@ASG[,options] filename[,type/reserve//maximum/placement]

Assigning Tape Files to Your Run
@ASG,options filename,
             type[/units/log/noise/processor/tape/format/data-converter/block-numbering/data-compression/buffered-write/expanded-buffer/,
             reel-1/reel-2/.../reel-n,
             expiration/mmspec,
             ring-indicator,
             ACR-name,
             CTL-pool ]



*/

// Not sure if the following are, or will be, useful

type MFDClientIdentifier struct {
	hostIdentifier uint64 // what should this *really* be?
	runId          string
	userId         string
}

type MFDFileAddressability uint
type MFDFileIdentifier uint64
type MFDGranularity uint
type TrackId uint64
type TrackCount uint64

const MFDTrackGranularity = 0
const MFDPositionGranularity = 1

const MFDSectorAddressable MFDFileAddressability = 0
const MFDWordAddressable MFDFileAddressability = 1

type MFDResult struct {
}

type MasterFileDirectory interface {
	//	TODO what to do about ER MSCON$ ?
	//	TODO what to do for retrieving file / cycle info?

	AllocateTracks(
		clientIdentifier MFDClientIdentifier,
		fileIdentifier MFDFileIdentifier,
		firstTrackId TrackId,
		trackCount TrackCount) MFDResult

	// TODO AssignExistingFileCycle()
	// TODO AssignNewDiskFileCycle()
	// TODO CatalogDiskFileCycle()

	// TODO Lots of Change* stuffs

	CreateTapeFileCycle(
		clientIdentifier MFDClientIdentifier,
		qualifier string,
		filename string,
		cycleSpecifier int,
		readKey string,
		writeKey string,
		maxFCycleCount uint,
		// lots of other stuff
	) (MFDFileIdentifier, MFDResult)

	CreateTemporaryDiskFile(
		clientIdentifier MFDClientIdentifier,
		qualifier string,
		filename string,
		cycleSpecifier int,
		granularity MFDGranularity,
		initialReserve uint64,
		maxGranules uint64) (MFDFileIdentifier, MFDResult)

	DeleteFileCycle(
		clientIdentifier MFDClientIdentifier,
		qualifier string,
		filename string,
		cycleSpecifier int,
		writeKey string) MFDResult

	ReleaseAllFiles(
		clientIdentifier MFDClientIdentifier) MFDResult

	ReleaseFile(
		clientIdentifier MFDClientIdentifier,
		fileIdentifier MFDFileIdentifier) MFDResult

	//	TODO ReleaseFileAndDelete (or should we have delete flag in ReleaseFile() ?
}

/*
type LeadItem struct {
	qualifier                string
	filename                 string
	projectId                string
	readKey                  string
	writeKey                 string
	fileType                 uint //	00 MS, 01 Tape, 040 REM
	fCycleCount              uint //	number of f-cycles which actually exist, not including to-be-cataloged or to-be-dropped
	maxRange                 uint
	currentRange             uint
	highestAbsoluteFCycle    uint //	whether it actually exists or not
	guardedFile              bool
	plusOneFCycleExists      bool
	fileNameChangeInProgress bool
	directoryIndex           uint // 00 local, 01 shared (we may do something different)
	accessType               uint //	the access types specified in the ACR (meaning what?)
	// security words
	// links to main items
	//		main item links contain status bits: to-be-cataloged, to-be-dropped
	//		we might rather have these flags in the various MainItem structs...
}

type MainItem struct {
	//	DAD table info
	qualifier        string //	some of these are repeated from lead item, to support removable disk
	filename         string
	projectId        string
	accountNumber    string
	timeOfFirstWrite time.Time // or time of unlock

	//	disabled flags...
	disabledDirectoryError            bool
	disabledAssignedWrittenBeforeStop bool
	disabledInaccessibleBackup        bool
	disabledCacheDrainFailure         bool

	//	link to lead item only for fixed devices. do we need it at all?

	//	descriptor flags...
	unloaded          bool
	backedUp          bool
	saveOnCheckpoint  bool
	toBeCataloged     bool
	tapeFile          bool
	removableDiskFile bool
	fileToBeWriteOnly bool
	fileToBeReadOnly  bool
	fileToBeDropped   bool

	//	file flags...
	largeFile bool
	writtenTo bool

	//	PCHAR flags
	positionGranularity bool
	wordAddressable     bool

	assignMnemonic string
	//	link to initial SMOQUE entry
	cumulativeAssignCount uint
	//	link to shared file extension item if shared
	sharedFileIndex uint

	//	inhibit flags
	guardedFile   bool
	inhibitUnload bool
	privatePublic bool // definition varies according to other states (owned, unowned, ACR, etc)
	exclusiveUse  bool
	writeOnly     bool
	readOnly      bool

	assignedIndicator      uint // for this f-cycle
	absoluteFCycle         uint
	timeOfLastReference    time.Time
	timeOfCatalog          time.Time
	initialGranules        uint
	maxGranules            uint
	highestGranuleAssigned uint
	highestTrackWritten    uint
	//	readkey again? (yes, for removable)
	//	writekey again? (ditto)

	//	user unit selection (applies only to Fixed)
	fileCreatedWithDevicePlacement      bool
	fileCreatedWithControlUnitPlacement bool
	fileCreatedWithLogicalPlacement     bool
	fileSpreadAcrossDevices             bool

	//	number of granules for quota groups

	//	backupwords
	backupCreationTime                  time.Time
	maxBackupLevels                     uint
	currentBackupLevels                 uint
	numberOfTextBlocks                  uint
	startingFilePositionFirstBackupReel uint64
	reelNumberFirstBackupReel           string
	reelNumberSecondBackupReel          string
}
*/
