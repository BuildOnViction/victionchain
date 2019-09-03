package types

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
)

const (
	StatusUnread = "UNREAD"
	StatusRead   = "READ"

	TypeAnnounce = "ANNOUNCE"
	TypeAlert    = "ALERT"
	TypeLog      = "LOG"
)

//Message struct
type Message struct {
	MessageType string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
}

// Notification struct
type Notification struct {
	ID        bson.ObjectId  `json:"_id" bson:"_id"`
	Recipient common.Address `json:"recipient" bson:"recipient"`
	Message   Message        `json:"message" bson:"message"`
	Type      string         `json:"type" bson:"type"`
	Status    string         `json:"status" bson:"status"`
	CreatedAt time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" bson:"updatedAt"`
}

// NotificationRecord struct
type NotificationRecord struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Recipient string        `json:"recipient" bson:"recipient"`
	Message   Message       `json:"message" bson:"message"`
	Type      string        `json:"type" bson:"type"`
	Status    string        `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// NotificationBSONUpdate return BSON structure for NotificationSpec structure
type NotificationBSONUpdate struct {
	*Notification
}

// MarshalJSON returns the json encoded byte array representing the notification struct
func (n *Notification) MarshalJSON() ([]byte, error) {
	notification := map[string]interface{}{
		"id":        n.ID,
		"recipient": n.Recipient,
		"message":   n.Message,
		"type":      n.Type,
		"status":    n.Status,
		"createdAt": n.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt": n.UpdatedAt.Format(time.RFC3339Nano),
	}

	return json.Marshal(notification)
}

// UnmarshalJSON creates a notification object from a json byte string
func (n *Notification) UnmarshalJSON(b []byte) error {
	notification := map[string]interface{}{}

	err := json.Unmarshal(b, &notification)

	if err != nil {
		return err
	}

	if notification["_id"] != nil && bson.IsObjectIdHex(notification["_id"].(string)) {
		n.ID = bson.ObjectIdHex(notification["_id"].(string))
	}
	if notification["id"] != nil && bson.IsObjectIdHex(notification["id"].(string)) {
		n.ID = bson.ObjectIdHex(notification["id"].(string))
	}
	if notification["recipient"] == nil {
		// return errors.New("Order Hash is not set")
	} else {
		n.Recipient = common.HexToAddress(notification["recipient"].(string))
	}

	if notification["message"] != nil {
		n.Message = notification["message"].(Message)
	}

	if notification["type"] != nil {
		n.Type = notification["type"].(string)
	}

	if notification["status"] != nil {
		n.Status = notification["status"].(string)
	}

	if notification["createdAt"] != nil {
		nm, _ := time.Parse(time.RFC3339Nano, notification["createdAt"].(string))
		n.CreatedAt = nm
	}

	if notification["updatedAt"] != nil {
		nm, _ := time.Parse(time.RFC3339Nano, notification["updatedAt"].(string))
		n.UpdatedAt = nm
	}

	return nil
}

// GetBSON get Notification struct
func (n *Notification) GetBSON() (interface{}, error) {
	nr := NotificationRecord{
		ID:        n.ID,
		Recipient: n.Recipient.Hex(),
		Message:   n.Message,
		Status:    n.Status,
		Type:      n.Type,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
	return nr, nil
}

// SetBSON json to Notification
func (n *Notification) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID        bson.ObjectId `json:"_id" bson:"_id"`
		Recipient string        `json:"recipient" bson:"recipient"`
		Message   Message       `json:"message" bson:"message"`
		Type      string        `json:"type" bson:"type"`
		Status    string        `json:"status" bson:"status"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
		UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
	})

	err := raw.Unmarshal(decoded)

	if err != nil {
		return err
	}

	n.ID = decoded.ID
	n.Recipient = common.HexToAddress(decoded.Recipient)
	n.Message = decoded.Message
	n.Type = decoded.Type
	n.Status = decoded.Status
	n.CreatedAt = decoded.CreatedAt
	n.UpdatedAt = decoded.UpdatedAt

	return nil
}
