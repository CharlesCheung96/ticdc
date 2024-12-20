// Copyright 2024 PingCAP, Inc.
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
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// SorterOutputEventCount is the metric that counts events output by the sorter.
	SorterOutputEventCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ticdc",
		Subsystem: "sorter",
		Name:      "output_event_count",
		Help:      "The number of events output by the sorter",
	}, []string{"namespace", "changefeed", "type"})

	// EventServiceSendEventCount is the metric that counts events sent by the event service.
	EventServiceSendEventCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ticdc",
		Subsystem: "event_service",
		Name:      "send_event_count",
		Help:      "The number of events sent by the event service",
	}, []string{"namespace", "changefeed", "type"})

	// EventServiceSendEventDuration is the metric that records the duration of sending events by the event service.
	EventServiceSendEventDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ticdc",
		Subsystem: "event_service",
		Name:      "send_event_duration",
		Help:      "The duration of sending events by the event service",
		Buckets:   prometheus.DefBuckets,
	}, []string{"type"})
	EventServiceResolvedTsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "resolved_ts",
			Help:      "resolved ts of eventService",
		})
	EventServiceResolvedTsLagGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "resolved_ts_lag",
			Help:      "resolved ts lag of eventService in seconds",
		}, []string{"type"})
	EventServiceScanDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "scan_duration",
			Help:      "The duration of scanning a data range from eventStore",
			Buckets:   prometheus.DefBuckets,
		})
	EventServiceDispatcherGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "dispatcher_count",
			Help:      "The number of dispatchers in event service",
		}, []string{"cluster"})
	EventServiceDropScanTaskCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "drop_scan_task_count",
			Help:      "The number of scan tasks dropped",
		})
	EventServiceDropResolvedTsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "drop_resolved_ts_count",
			Help:      "The number of resolved ts dropped",
		})
	EventServiceScanTaskQueueDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "scan_task_queue_duration",
			Help:      "The duration of a scan task being queued",
			Buckets:   prometheus.DefBuckets,
		})
	EventServiceHandleDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "handle_duration",
			Help:      "The duration of handling a scan task",
			Buckets:   prometheus.DefBuckets,
		})
	EventServiceDropNotificationCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "ticdc",
			Subsystem: "event_service",
			Name:      "drop_notification_count",
			Help:      "The number of notifications dropped",
		})
)

// InitMetrics registers all metrics in this file.
func InitEventServiceMetrics(registry *prometheus.Registry) {
	registry.MustRegister(SorterOutputEventCount)
	registry.MustRegister(EventServiceSendEventCount)
	registry.MustRegister(EventServiceSendEventDuration)
	registry.MustRegister(EventServiceResolvedTsGauge)
	registry.MustRegister(EventServiceResolvedTsLagGauge)
	registry.MustRegister(EventServiceScanDuration)
	registry.MustRegister(EventServiceDispatcherGauge)
	registry.MustRegister(EventServiceDropScanTaskCount)
	registry.MustRegister(EventServiceDropResolvedTsCount)
	registry.MustRegister(EventServiceScanTaskQueueDuration)
	registry.MustRegister(EventServiceHandleDuration)
	registry.MustRegister(EventServiceDropNotificationCount)
}
