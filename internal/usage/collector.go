package usage

type Collector interface {
	Key() string
	Sample() (Sample, error)
}
