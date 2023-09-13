package collector

import (
	"sync"

	"github.com/flashbots/mempool-dumpster/common"
	"go.uber.org/zap"
)

const (
	KeyStatsAll       = "all"
	KeyStatsFirst     = "first"
	KeyStatsUnique    = "unique"
	KeyStatsTxOnChain = "tx-onchain"
)

type SourceMetrics struct {
	lock   sync.Mutex
	counts map[string]map[string]map[string]uint64 // cntType -> source -> key -> count
}

func NewMetricsCounter() SourceMetrics {
	return SourceMetrics{ //nolint:exhaustruct
		counts: make(map[string]map[string]map[string]uint64),
	}
}

func (sc *SourceMetrics) Inc(cntType, source string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	if _, ok := sc.counts[cntType]; !ok {
		sc.counts[cntType] = make(map[string]map[string]uint64)
	}
	if _, ok := sc.counts[cntType][source]; !ok {
		sc.counts[cntType][source] = make(map[string]uint64)
	}

	sc.counts[cntType][source][cntType] += 1
}

func (sc *SourceMetrics) IncKey(cntType, source, key string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	if _, ok := sc.counts[cntType]; !ok {
		sc.counts[cntType] = make(map[string]map[string]uint64)
	}
	if _, ok := sc.counts[cntType][source]; !ok {
		sc.counts[cntType][source] = make(map[string]uint64)
	}

	sc.counts[cntType][source][key] += 1
}

func (sc *SourceMetrics) Get(cntType string) map[string]map[string]uint64 {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	return sc.counts[cntType]
}

func (sc *SourceMetrics) Reset() {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	sc.counts = make(map[string]map[string]map[string]uint64)
}

func (sc *SourceMetrics) Logger(log *zap.SugaredLogger, cntType string, useLen bool) *zap.SugaredLogger {
	for k, v := range sc.Get(cntType) {
		if useLen {
			log = log.With(k, common.Printer.Sprint(len(v)))
		} else {
			log = log.With(k, common.Printer.Sprint(v[cntType]))
		}
	}
	return log
}
