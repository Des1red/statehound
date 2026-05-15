package statehound

import (
	"statehound/internal/statehound/filter"
	"time"
)

func buildFilter() filter.DiffFilter {
	return filter.DiffFilter{
		ServiceFilters: []filter.ServiceFilter{
			&filter.IgnoreSafeServices{},
		},
		PortFilters: []filter.PortFilter{
			&filter.SuspiciousListenerFilter{},
		},
		FileFilters: []filter.FileFilter{
			&filter.PersistenceFileFilter{},
		},
		ConnectionFilters: []filter.ConnectionFilter{
			&filter.ConnectionRateFilter{Interval: 30 * time.Second},
		},
		ProcessFilters: []filter.ProcessFilter{
			&filter.SuspiciousProcessFilter{},
		},
	}
}
