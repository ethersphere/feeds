package feed

import (
	"encoding/hex"
	"testing"
)

func TestTopic(t *testing.T) {
	related, _ := hex.DecodeString("abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789")
	topicName := "test-topic"
	topic, _ := NewTopic(topicName, related)
	hex := topic.Hex()
	expectedHex := "dfa89c750e3108f9c2aeef0123456789abcdef0123456789abcdef0123456789"
	if hex != expectedHex {
		t.Fatalf("Expected %s, got %s", expectedHex, hex)
	}

	var topic2 Topic
	topic2.FromHex(hex)
	if topic2 != topic {
		t.Fatal("Expected recovered topic to be equal to original one")
	}

	if topic2.Name(related) != topicName {
		t.Fatal("Retrieved name does not match")
	}

	bytes, err := topic2.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	expectedJSON := `"dfa89c750e3108f9c2aeef0123456789abcdef0123456789abcdef0123456789"`
	equal, err := areEqualJSON(expectedJSON, string(bytes))
	if err != nil {
		t.Fatal(err)
	}
	if !equal {
		t.Fatalf("Expected JSON to be %s, got %s", expectedJSON, string(bytes))
	}

	err = topic2.UnmarshalJSON(bytes)
	if err != nil {
		t.Fatal(err)
	}
	if topic2 != topic {
		t.Fatal("Expected recovered topic to be equal to original one")
	}

}
