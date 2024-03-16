package structs

import "encoding/json"

type Message[P []byte] struct {
	Hash            string       `json:"hash"`
	FromID          string       `json:"from_id"`
	TargetID        string       `json:"target_id"`
	TargetAction    TargetAction `json:"target_action"`
	Payload         P            `json:"payload"`
	WaitForResponse bool         `json:"wait_for_response"`
}

func (m *Message[P]) MarshalMessage() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (m *Message[P]) UnmarshalPayload(out any) error {
	if err := json.Unmarshal(m.Payload, out); err != nil {
		return err
	}

	return nil
}
