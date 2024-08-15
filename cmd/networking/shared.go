package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type PlayerInfo struct {
	Name             string
	PercentCompleted uint
	Wpm              uint
	ReadyToStart     bool
}

type Broadcast struct {
	Done        bool
	Started     bool
	StartTime   time.Time
	Paragraph   string
	PlayerInfos []PlayerInfo // ordered by position
}

func Serialize[T any](p T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		fmt.Println("here")
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeSerialize[T any](data []byte) (T, error) {
	var t T
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}
