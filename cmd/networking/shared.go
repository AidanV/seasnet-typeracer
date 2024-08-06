package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type PlayerInfo struct {
	Name             string
	PercentCompleted uint
	Wpm              uint
}

func Serialize(p PlayerInfo) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		fmt.Println("here")
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeSerialize(data []byte) (PlayerInfo, error) {
	var playerInfo PlayerInfo
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&playerInfo)
	if err != nil {
		return playerInfo, err
	}
	return playerInfo, nil
}
