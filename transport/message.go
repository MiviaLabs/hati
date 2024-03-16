package transport

import "encoding/json"

type Message[P []byte] struct {
	Hash            string       `json:"hash"`
	FromID          string       `json:"from_id"`
	TargetID        string       `json:"target_id"`
	TargetAction    TargetAction `json:"target_action"`
	Payload         P            `json:"payload"`
	WaitForResponse bool         `json:"wait_for_response"`
}

// func (m Message[P]) GetPayload()(P, error) {

// }

func (m Message[P]) Unmarshal(out *P) error {
	if err := json.Unmarshal(m.Payload, out); err != nil {
		return err
	}

	return nil
}

type TargetAction struct {
	Module string `json:"module"`
	Action string `json:"action"`
}
