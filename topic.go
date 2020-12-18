// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package feed

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const ChunkAddressLength = 32

// TopicLength establishes the max length of a topic string
const TopicLength = ChunkAddressLength

// Topic represents what a feed is about
type Topic [TopicLength]byte

// ErrTopicTooLong is returned when creating a topic with a name/related content too long
var ErrTopicTooLong = fmt.Errorf("Topic is too long. Max length is %d", TopicLength)

// NewTopic creates a new topic from a provided name and "related content" byte array,
// merging the two together.
// If relatedContent or name are longer than TopicLength, they will be truncated and an error returned
// name can be an empty string
// relatedContent can be nil
func NewTopic(name string, relatedContent []byte) (topic Topic, err error) {
	if relatedContent != nil {
		contentLength := len(relatedContent)
		if contentLength > TopicLength {
			contentLength = TopicLength
			err = ErrTopicTooLong
		}
		copy(topic[:], relatedContent[:contentLength])
	}
	nameBytes := []byte(name)
	nameLength := len(nameBytes)
	if nameLength > TopicLength {
		nameLength = TopicLength
		err = ErrTopicTooLong
	}
	XORBytes(topic[:], topic[:], nameBytes[:nameLength])
	return topic, err
}

// Hex will return the topic encoded as an hex string
func (t *Topic) Hex() string {
	return hex.EncodeToString(t[:])
}

// FromHex will parse a hex string into this Topic instance
func (t *Topic) FromHex(h string) error {
	bytes, err := hex.DecodeString(h)
	if err != nil || len(bytes) != len(t) {
		return NewErrorf(ErrInvalidValue, "Cannot decode topic")
	}
	copy(t[:], bytes)
	return nil
}

// Name will try to extract the topic name out of the Topic
func (t *Topic) Name(relatedContent []byte) string {
	nameBytes := *t
	if relatedContent != nil {
		contentLength := len(relatedContent)
		if contentLength > TopicLength {
			contentLength = TopicLength
		}
		XORBytes(nameBytes[:], t[:], relatedContent[:contentLength])
	}
	z := bytes.IndexByte(nameBytes[:], 0)
	if z < 0 {
		z = TopicLength
	}
	return string(nameBytes[:z])

}

// UnmarshalJSON implements the json.Unmarshaller interface
func (t *Topic) UnmarshalJSON(data []byte) error {
	var hex string
	err := json.Unmarshal(data, &hex)
	if err != nil {
		return err
	}
	return t.FromHex(hex)
}

// MarshalJSON implements the json.Marshaller interface
func (t *Topic) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Hex())
}
