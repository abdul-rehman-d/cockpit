package usage

type Service struct {
	collectors []Collector
}

func NewService() *Service {
	return &Service{
		collectors: []Collector{
			NewMemoryCollector(),
			NewCPUCollector(),
			NewStorageCollector(),
		},
	}
}

func (s *Service) GetAllSamples() []Sample {
	var samples []Sample
	for _, collector := range s.collectors {
		sample, err := collector.Sample()
		if err != nil {
			panic("error should have been handled by the collector: " + err.Error())
		}
		samples = append(samples, sample)
	}
	return samples
}
