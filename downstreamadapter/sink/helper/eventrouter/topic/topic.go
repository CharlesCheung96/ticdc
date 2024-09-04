package topic

import (
	"github.com/pingcap/tiflow/pkg/config"
)

type TopicGenerator interface {
	Substitute(schema, table string) string
}

type StaticTopicGenerator struct {
	topic string
}

// NewStaticTopicDispatcher returns a StaticTopicDispatcher.
func newStaticTopic(defaultTopic string) *StaticTopicGenerator {
	return &StaticTopicGenerator{
		topic: defaultTopic,
	}
}

// Substitute converts schema/table name in a topic expression to kafka topic name.
func (s *StaticTopicGenerator) Substitute(schema, table string) string {
	return s.topic
}

func (s *StaticTopicGenerator) String() string {
	return s.topic
}

// DynamicTopicGenerator is a topic generator which dispatches rows and DDLs
// dynamically to the target topics.
type DynamicTopicGenerator struct {
	expression Expression
}

// NewDynamicTopicDispatcher creates a DynamicTopicDispatcher.
func newDynamicTopicGenerator(topicExpr Expression) *DynamicTopicGenerator {
	return &DynamicTopicGenerator{
		expression: topicExpr,
	}
}

// Substitute converts schema/table name in a topic expression to kafka topic name.
func (d *DynamicTopicGenerator) Substitute(schema, table string) string {
	return d.expression.Substitute(schema, table)
}

func (d *DynamicTopicGenerator) String() string {
	return string(d.expression)
}

func GetTopicGenerator(
	rule string, defaultTopic string, protocol config.Protocol, scheme string,
) (TopicGenerator, error) {
	if rule == "" || isHardCode(rule) {
		return newStaticTopic(defaultTopic), nil
	}

	// check if this rule is a valid topic expression
	topicExpr := Expression(rule)
	err := validateTopicExpression(topicExpr, scheme, protocol)
	if err != nil {
		return nil, err
	}
	return newDynamicTopicGenerator(topicExpr), nil
}