package structs

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type Message[P []byte] struct {
	Hash            string       `json:"hash"`
	FromID          string       `json:"from_id"`
	TargetID        string       `json:"target_id"`
	TargetAction    TargetAction `json:"target_action"`
	Payload         P            `json:"payload"`
	WaitForResponse bool         `json:"wait_for_response"`
	ResponseHash    string       `json:"response_hash"`
}

func (m *Message[P]) UpdateHash() {
	hash := crypto.Keccak256([]byte(m.FromID), []byte(m.TargetID), []byte(m.TargetAction.Module), []byte(m.TargetAction.Action), m.Payload, []byte(time.Now().String()))
	m.Hash = crypto.Keccak256Hash(hash).String()
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
