package eventservice

import (
	"github.com/flowbehappy/tigate/pkg/common"
	"github.com/flowbehappy/tigate/utils/dynstream"
	"github.com/pingcap/log"
)

type dispatcherEventsHandler struct {
}

func (h *dispatcherEventsHandler) Path(task scanTask) common.DispatcherID {
	return task.dispatcherStat.info.GetID()
}

// Handle implements the dynstream.Handler interface.
// If the event is processed successfully, it should return false.
// If the event is processed asynchronously, it should return true. The later events of the path are blocked
// until a wake signal is sent to DynamicStream's Wake channel.
func (h *dispatcherEventsHandler) Handle(broker *eventBroker, tasks ...scanTask) bool {
	if len(tasks) != 1 {
		log.Panic("only one task is allowed")
	}
	task := tasks[0]
	needScan, _ := broker.checkNeedScan(task)
	if !needScan {
		task.handle()
		return false
	}
	// The dispatcher has new events. We need to push the task to the task pool.
	return broker.taskPool.pushTask(task)
}

func (h *dispatcherEventsHandler) GetType(event scanTask) dynstream.EventType {
	// scanTask is only a signal to trigger the scan.
	// We make it a RepeatedSignal to make the new scan task squeeze out the old one.
	return dynstream.EventType{DataGroup: 0, Property: dynstream.RepeatedSignal}
}

func (h *dispatcherEventsHandler) GetSize(event scanTask) int                              { return 0 }
func (h *dispatcherEventsHandler) GetArea(path common.DispatcherID, dest *eventBroker) int { return 0 }
func (h *dispatcherEventsHandler) GetTimestamp(event scanTask) dynstream.Timestamp         { return 0 }
func (h *dispatcherEventsHandler) IsPaused(event scanTask) bool                            { return false }
func (h *dispatcherEventsHandler) OnDrop(event scanTask)                                   {}