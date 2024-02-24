package facilitiesMgr

type FacNodeStatus uint

const (
	_ FacNodeStatus = iota
	FacNodeStatusUp
	FacNodeStatusReserved
	FacNodeStatusDown
	FacNodeStatusSuspended
)
