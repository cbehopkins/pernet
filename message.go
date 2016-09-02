package pernet

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Action string `json:"action"`
}

func UnmarshalMessage(input string) (bob Message, err error) {

	err = json.Unmarshal([]byte(input), &bob)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	return
}
func MarshalMessage(bob Message) (text string, err error) {
	tmp, err := json.Marshal(bob)
	text = string(tmp)
	return
}
