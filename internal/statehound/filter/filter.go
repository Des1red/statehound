package filter

import (
	"time"
)

func BuildFilter() DiffFilter {
	return DiffFilter{
		ServiceFilters: []ServiceFilter{
			&IgnoreSafeServices{},
		},
		PortFilters: []PortFilter{
			&SuspiciousListenerFilter{},
		},
		FileFilters: []FileFilter{
			&PersistenceFileFilter{},
		},
		ConnectionFilters: []ConnectionFilter{
			&ConnectionRateFilter{Interval: 30 * time.Second},
		},
		ProcessFilters: []ProcessFilter{
			&SuspiciousProcessFilter{},
		},
	}
}
