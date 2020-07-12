package scalers

import (
	"reflect"
	"testing"
)

type parseKafkaMetadataTestData struct {
	metadata            map[string]string
	isError             bool
	numBrokers          int
	brokers             []string
	group               string
	topic               string
	consumerOffsetReset offsetResetPolicy
}

// A complete valid metadata example for reference
var validMetadata = map[string]string{
	"bootstrapServers": "broker1:9092,broker2:9092",
	"consumerGroup":    "my-group",
	"topic":            "my-topic",
}

// A complete valid authParams example for sasl, with username and passwd
var validWithAuthParams = map[string]string{
	"authMode": "sasl_plaintext",
	"username": "admin",
	"password": "admin",
}

// A complete valid authParams example for sasl, without username and passwd
var validWithoutAuthParams = map[string]string{}

var parseKafkaMetadataTestDataset = []parseKafkaMetadataTestData{
	// failure, no bootstrapServers
	{map[string]string{}, true, 0, nil, "", "", ""},
	// failure, no consumer group
	{map[string]string{"bootstrapServers": "foobar:9092"}, true, 1, []string{"foobar:9092"}, "", "", "earliest"},
	// failure, no topic
	{map[string]string{"bootstrapServers": "foobar:9092", "consumerGroup": "my-group"}, true, 1, []string{"foobar:9092"}, "my-group", "", offsetResetPolicy("earliest")},
	// success
	{map[string]string{"bootstrapServers": "foobar:9092", "consumerGroup": "my-group", "topic": "my-topic"}, false, 1, []string{"foobar:9092"}, "my-group", "my-topic", offsetResetPolicy("earliest")},
	// success, more brokers
	{map[string]string{"bootstrapServers": "foo:9092,bar:9092", "consumerGroup": "my-group", "topic": "my-topic"}, false, 2, []string{"foo:9092", "bar:9092"}, "my-group", "my-topic", offsetResetPolicy("earliest")},
	// success, consumerOffsetReset policy earliest
	{map[string]string{"bootstrapServers": "foo:9092,bar:9092", "consumerGroup": "my-group", "topic": "my-topic", "consumerOffsetReset": "earliest"}, false, 2, []string{"foo:9092", "bar:9092"}, "my-group", "my-topic", offsetResetPolicy("earliest")},
	// failure, consumerOffsetReset policy wrong
	{map[string]string{"bootstrapServers": "foo:9092,bar:9092", "consumerGroup": "my-group", "topic": "my-topic", "consumerOffsetReset": "foo"}, true, 2, []string{"foo:9092", "bar:9092"}, "my-group", "my-topic", ""},
	// success, consumerOffsetReset policy latest
	{map[string]string{"bootstrapServers": "foo:9092,bar:9092", "consumerGroup": "my-group", "topic": "my-topic", "consumerOffsetReset": "latest"}, false, 2, []string{"foo:9092", "bar:9092"}, "my-group", "my-topic", offsetResetPolicy("latest")},
}

func TestGetBrokers(t *testing.T) {
	for _, testData := range parseKafkaMetadataTestDataset {
		meta, err := parseKafkaMetadata(nil, testData.metadata, validWithAuthParams)

		if err != nil && !testData.isError {
			t.Error("Expected success but got error", err)
		}
		if testData.isError && err == nil {
			t.Error("Expected error but got success")
		}
		if len(meta.bootstrapServers) != testData.numBrokers {
			t.Errorf("Expected %d bootstrap servers but got %d\n", testData.numBrokers, len(meta.bootstrapServers))
		}
		if !reflect.DeepEqual(testData.brokers, meta.bootstrapServers) {
			t.Errorf("Expected %v but got %v\n", testData.brokers, meta.bootstrapServers)
		}
		if meta.group != testData.group {
			t.Errorf("Expected group %s but got %s\n", testData.group, meta.group)
		}
		if meta.topic != testData.topic {
			t.Errorf("Expected topic %s but got %s\n", testData.topic, meta.topic)
		}
		if err == nil && meta.consumerOffsetReset != testData.consumerOffsetReset {
			t.Errorf("Expected consumerOffsetReset %s but got %s\n", testData.consumerOffsetReset, meta.consumerOffsetReset)
		}

		meta, err = parseKafkaMetadata(nil, testData.metadata, validWithoutAuthParams)

		if err != nil && !testData.isError {
			t.Error("Expected success but got error", err)
		}
		if testData.isError && err == nil {
			t.Error("Expected error but got success")
		}
		if len(meta.bootstrapServers) != testData.numBrokers {
			t.Errorf("Expected %d bootstrap servers but got %d\n", testData.numBrokers, len(meta.bootstrapServers))
		}
		if !reflect.DeepEqual(testData.brokers, meta.bootstrapServers) {
			t.Errorf("Expected %v but got %v\n", testData.brokers, meta.bootstrapServers)
		}
		if meta.group != testData.group {
			t.Errorf("Expected group %s but got %s\n", testData.group, meta.group)
		}
		if meta.topic != testData.topic {
			t.Errorf("Expected topic %s but got %s\n", testData.topic, meta.topic)
		}
		if err == nil && meta.consumerOffsetReset != testData.consumerOffsetReset {
			t.Errorf("Expected consumerOffsetReset %s but got %s\n", testData.consumerOffsetReset, meta.consumerOffsetReset)
		}
	}
}
