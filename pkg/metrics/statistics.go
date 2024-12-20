// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"time"

	"github.com/pingcap/ticdc/pkg/common"
	"github.com/pingcap/ticdc/pkg/config"

	commonEvent "github.com/pingcap/ticdc/pkg/common/event"
	"github.com/prometheus/client_golang/prometheus"
)

// NewStatistics creates a statistics
func NewStatistics(
	changefeed common.ChangeFeedID,
	sinkType string,
) *Statistics {
	statistics := &Statistics{
		sinkType:     sinkType,
		captureAddr:  config.GetGlobalServerConfig().AdvertiseAddr,
		changefeedID: changefeed,
	}

	namespcae := statistics.changefeedID.Namespace()
	changefeedID := statistics.changefeedID.Name()
	s := sinkType
	statistics.metricExecDDLHis = ExecDDLHistogram.WithLabelValues(namespcae, changefeedID, s)
	statistics.metricExecBatchHis = ExecBatchHistogram.WithLabelValues(namespcae, changefeedID, s)
	statistics.metricTotalWriteBytesCnt = TotalWriteBytesCounter.WithLabelValues(namespcae, changefeedID, s)
	statistics.metricEventSizeHis = EventSizeHistogram.WithLabelValues(namespcae, changefeedID)
	statistics.metricExecErrCnt = ExecutionErrorCounter.WithLabelValues(namespcae, changefeedID, s)
	statistics.metricExecDMLCnt = ExecDMLEventCounter.WithLabelValues(namespcae, changefeedID)
	return statistics
}

// Statistics maintains some status and metrics of the Sink
// Note: All methods of Statistics should be thread-safe.
type Statistics struct {
	sinkType     string
	captureAddr  string
	changefeedID common.ChangeFeedID

	// Histogram for DDL Executing duration.
	metricExecDDLHis prometheus.Observer
	// Histogram for DML batch size.
	metricExecBatchHis prometheus.Observer
	// Counter for total bytes of DML.
	metricTotalWriteBytesCnt prometheus.Counter
	// Histogram for Row size.
	metricEventSizeHis prometheus.Observer
	// Counter for sink error.
	metricExecErrCnt prometheus.Counter

	metricExecDMLCnt prometheus.Counter
}

// ObserveRows stats all received `RowChangedEvent`s.
func (b *Statistics) ObserveRows(events []*commonEvent.DMLEvent) {
	for _, event := range events {
		b.metricEventSizeHis.Observe(float64(event.Rows.MemoryUsage()))
		b.metricExecDMLCnt.Inc()
	}
}

// RecordBatchExecution stats batch executors which return (batchRowCount, error).
func (b *Statistics) RecordBatchExecution(executor func() (int, int64, error)) error {
	batchSize, batchWriteBytes, err := executor()
	if err != nil {
		b.metricExecErrCnt.Inc()
		return err
	}
	b.metricExecBatchHis.Observe(float64(batchSize))
	b.metricTotalWriteBytesCnt.Add(float64(batchWriteBytes))
	return nil
}

// RecordDDLExecution record the time cost of execute ddl
func (b *Statistics) RecordDDLExecution(executor func() error) error {
	start := time.Now()
	if err := executor(); err != nil {
		b.metricExecErrCnt.Inc()
		return err
	}
	b.metricExecDDLHis.Observe(time.Since(start).Seconds())
	return nil
}

// Close release some internal resources.
func (b *Statistics) Close() {
	namespace := b.changefeedID.Namespace()
	changefeedID := b.changefeedID.Name()
	ExecDDLHistogram.DeleteLabelValues(namespace, changefeedID)
	ExecBatchHistogram.DeleteLabelValues(namespace, changefeedID)
	EventSizeHistogram.DeleteLabelValues(namespace, changefeedID)
	ExecutionErrorCounter.DeleteLabelValues(namespace, changefeedID)
	TotalWriteBytesCounter.DeleteLabelValues(namespace, changefeedID)
	ExecDMLEventCounter.DeleteLabelValues(namespace, changefeedID)
}
