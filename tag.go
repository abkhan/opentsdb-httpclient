package tsdbh

import (
	"encoding/json"
	"errors"
	"strings"
)

// Tag struct is for one tsdb data point tag
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Tags is a type of Tag list
type Tags []Tag

// NewTag function does few checks and returns a tag from key/value pair
func NewTag(k, v string) (Tag, error) {
	if strings.Contains(k, " ") {
		return Tag{}, errors.New("Tag key cannot contain space")
	}
	if strings.Contains(v, " ") {
		return Tag{}, errors.New("Tag value cannot contain space")
	}
	return Tag{k, v}, nil
}

// MakeTags is a helper function to transform
// and amp[string]string into slice of tags
func MakeTags(tags map[string]string) ([]Tag, error) {
	t := []Tag{}
	for k, v := range tags {
		newTag, err := NewTag(k, v)
		if err != nil {
			return t, err
		}
		t = append(t, newTag)
	}
	return t, nil
}

// String function would return a string for a tag
func (t Tag) String() string {
	return t.Key + "=" + t.Value
}

// MarshalJSON a tag into []byte
func (t Tags) MarshalJSON() ([]byte, error) {
	m := map[string]string{}
	for _, tag := range t {
		m[tag.Key] = tag.Value
	}
	return json.Marshal(m)
}
