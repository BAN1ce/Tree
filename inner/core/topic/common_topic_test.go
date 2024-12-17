package topic

import (
	"fmt"
	"github.com/BAN1ce/Tree/pkg/util"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"github.com/google/uuid"
	"testing"
)

func TestCommonTopic_matchTopicForWildcard1(t1 *testing.T) {

	var (
		commonTopic          = NewCommonTopic()
		topics               = util.GenerateTopics(3, []string{"a", "b", "c"}, []string{"+", "#"})
		noWildCardTopicCount int
	)

	for _, topic := range topics {
		fmt.Println(topic)
		_ = commonTopic.createSub(uuid.NewString(), &proto.SubOption{
			Topic: topic,
		})
		if !utils.HasWildcard(topic) {
			noWildCardTopicCount++
		}
	}

	result := commonTopic.matchTopicForWildcard("/a/+/c")
	if len(result) != 3 {
		t1.Fatalf("match wildcard topic failed, client count not match got = %d, expect = %d", len(result), 3)
	} else {
		t1.Logf("match wildcard topic success %d", len(result))

	}
}

func TestState_matchWildcardTopic(t1 *testing.T) {
	var (
		commonTopic = NewCommonTopic()
		topics      = map[string]string{
			"/a/b/c": uuid.NewString(), // Exact match
			"/a/+/c": uuid.NewString(), // Single-level wildcard +
			"/a/+/+": uuid.NewString(), // Single-level wildcard +, twice
			"/a/+/#": uuid.NewString(), // Single-level wildcard + and multi-level wildcard #
			"/a/#":   uuid.NewString(), // Multi-level wildcard #

			"/a/b/+": uuid.NewString(), // Single-level wildcard +
			"/a/b/#": uuid.NewString(), // Multi-level wildcard #

			"/+/#":   uuid.NewString(), // Single-level wildcard + and multi-level wildcard #
			"/+/b/c": uuid.NewString(), // Single-level wildcard +
			"/+/b/+": uuid.NewString(), // Single-level wildcard +, twice

			"/+/+/c": uuid.NewString(), // Single-level wildcard +, twice
			"/+/b/#": uuid.NewString(), // Single-level wildcard + and multi-level wildcard #
			"/+/+/#": uuid.NewString(), // Single-level wildcard +, twice and multi-level wildcard #

			"/+/+/+": uuid.NewString(), // Single-level wildcard +, three times

			"/#":       uuid.NewString(), // Multi-level wildcard, matches any topic
			"/a/b/c/+": uuid.NewString(), // Three-level topic followed by single-level wildcard
			"/a/b/c/#": uuid.NewString(), // Three-level topic followed by multi-level wildcard

			// Additional test cases
			"":         uuid.NewString(), // Empty topic
			"+/a/b":    uuid.NewString(), // Root level single-level wildcard
			"#/a/b":    uuid.NewString(), // Root level multi-level wildcard
			"/a/+/b/c": uuid.NewString(), // Single-level wildcard in the middle
			"/a/b/c/d": uuid.NewString(), // Non-matching topic
		}
	)

	clientID := map[string]struct{}{}
	for topic, clientID := range topics {
		if err := commonTopic.createSub(clientID, &proto.SubOption{
			Topic: topic,
		}); err != nil {
			t1.Log("create wildcard sub failed", err, topic)
		}
	}
	for _, v := range topics {
		clientID[v] = struct{}{}
	}

	result := commonTopic.matchTopic("/a/b/c")

	if len(result) != 15 {
		for t := range topics {
			if _, ok := result[t]; !ok {
				t1.Log("not match topic", t)
			}
		}

		for t := range result {
			t1.Log("match topic", t)
		}
		t1.Fatalf("match wildcard topic failed, client count not match got = %d, expect = %d", len(result), 14)
	} else {
		t1.Logf("match wildcard topic success %d", len(result))
		t1.Log("client", result)
	}

}

func TestCreateSub(t *testing.T) {
	var (
		commonTopic = NewCommonTopic()
	)
	commonTopic.createSub("client1", &proto.SubOption{
		Topic: "/a/b/c",
		QoS:   1,
	})

	commonTopic.createSub("client2", &proto.SubOption{
		Topic: "/a/b/c",
		QoS:   2,
	})
	commonTopic.createSub("client2", &proto.SubOption{
		Topic: "/#",
		QoS:   1,
	})

	result := commonTopic.matchTopic("/a/b/c")
	if result["/a/b/c"] != 2 {
		t.Fatalf("match topic failed, client count not match got = %d, expect = %d", result["/a/b/c"], 2)
	} else {
		t.Logf("match topic success %d", result["/a/b/c"])
	}

	if result["/#"] != 1 {
		t.Errorf("match topic failed, client count not match got = %d, expect = %d", result["/#"], 1)
	}
}

func TestState_matchTopic(t1 *testing.T) {
	var (
		commonTopic = NewCommonTopic()
		topics      = map[string]string{
			"/a":      uuid.NewString(),
			"/retain": uuid.NewString(),
		}
	)

	clientID := map[string]struct{}{}
	for topic, clientID := range topics {
		if err := commonTopic.createSub(clientID, &proto.SubOption{
			Topic: topic,
		}); err != nil {
			t1.Log("create wildcard sub failed", err, topic)
		}
	}
	for _, v := range topics {
		clientID[v] = struct{}{}
	}

	result := commonTopic.matchTopic("/a")
	if len(result) != 1 {
		t1.Fatalf("match topic failed, client count not match got = %d, expect = %d", len(result), 1)
	} else {
		t1.Logf("match topic success %d", len(result))
	}

	result = commonTopic.matchTopic("/retain")
	if len(result) != 1 {
		t1.Fatalf("match topic failed, client count not match got = %d, expect = %d", len(result), 1)
	} else {
		t1.Logf("match topic success %d", len(result))
	}

}
