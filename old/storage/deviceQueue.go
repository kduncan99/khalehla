package storage

type deviceQueue interface {
	Close() AggregatorResult
	Open(bool, bool) AggregatorResult
}
