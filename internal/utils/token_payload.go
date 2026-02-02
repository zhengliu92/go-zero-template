package utils

import "encoding/json"

// TokenPayload JWT token payload 结构
type TokenPayload struct {
	UserID int    `json:"userId"`
	Type   string `json:"type"` // "access" 或 "refresh"
}

func FromJSON(jsonStr string) (*TokenPayload, error) {
	var payload TokenPayload
	err := json.Unmarshal([]byte(jsonStr), &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
