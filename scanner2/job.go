package scanner2

import (
	"context"
	"io/fs"
	"sync"
	"sync/atomic"
	"time"

	"github.com/navidrome/navidrome/model"
)

type scanJob struct {
	lib         model.Library
	fs          fs.FS
	ds          model.DataStore
	startTime   time.Time
	lastUpdates map[string]time.Time
	lock        sync.RWMutex
	fullRescan  bool
	numFolders  atomic.Int64
}

func newScanJob(ctx context.Context, ds model.DataStore, lib model.Library, fullRescan bool) (*scanJob, error) {
	//lastUpdates, err := ds.Folder(ctx).GetLastUpdates(lib) FIXME
	//if err != nil {
	//	return nil, fmt.Errorf("error getting last updates: %w", err)
	//}
	lastUpdates := map[string]time.Time{}
	return &scanJob{
		lib:         lib,
		fs:          lib.FS(),
		ds:          ds,
		startTime:   time.Now(),
		lastUpdates: lastUpdates,
		fullRescan:  fullRescan,
	}, nil
}

func (s *scanJob) getLastUpdatedInDB(folderId string) time.Time {
	s.lock.RLock()
	defer s.lock.RUnlock()

	t, ok := s.lastUpdates[folderId]
	if !ok {
		return time.Time{}
	}
	return t
}
