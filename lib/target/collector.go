package target

import (
	"github.com/maxkulish/dnscrypt-list/lib/config"
	"github.com/maxkulish/dnscrypt-list/lib/files"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/validator"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	// ErrorValidationTarget shows that target is not valid
	ErrorValidationTarget = errors.New("no valid targets found")
)

// TargetsStore keeps the list of Targets
type TargetsStore struct {
	Targets []*Target
}

// NewTargetsStore creates an instance of TargetsStore
func NewTargetsStore() *TargetsStore {
	return &TargetsStore{
		Targets: []*Target{},
	}
}

// Length count the number of elements in the TargetsStore
func (t *TargetsStore) Length() int {
	return len(t.Targets)
}

// CollectTargets iterates targets in BlackList and WhiteList
func CollectTargets(conf *config.Config) (*TargetsStore, error) {

	targetsStore := NewTargetsStore()

	for _, rawTarget := range conf.SourceList.Targets {

		var target *Target

		var validURL bool
		var err error
		if rawTarget.URL != "" {
			validURL, err = validator.IsURL(rawTarget.URL)
			if err != nil {
				logger.Debug("raw url validation error", zap.Error(err))
				continue
			}
		}

		var validPath bool
		if rawTarget.File != "" {
			validPath = files.IsPathExist(rawTarget.File)
		}

		if validURL || validPath {
			target = NewTargetFromRaw(rawTarget)
		} else {
			logger.Error("no valid targets found", zap.Error(ErrorValidationTarget))
			continue
		}

		targetsStore.Targets = append(targetsStore.Targets, target)
	}

	return targetsStore, nil
}
