package convoy_cli

import (
	"encoding/json"

	"github.com/frain-dev/convoy/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListenRequest struct {
	HostName   string   `json:"host_name"`
	DeviceID   string   `json:"device_id"`
	SourceID   string   `json:"source_id"`
	EventTypes []string `json:"event_types"`
}

type LoginRequest struct {
	HostName string `json:"host_name"`
	DeviceID string `json:"device_id"`
}

type LoginResponse struct {
	Device   *Device             `json:"device"`
	Project  *Project            `json:"project"`
	Endpoint *datastore.Endpoint `json:"endpoint"`
}

type Device struct {
	UID        string             `json:"uid"`
	ProjectID  string             `json:"project_id,omitempty"`
	EndpointID string             `json:"endpoint_id,omitempty"`
	HostName   string             `json:"host_name,omitempty"`
	Status     string             `json:"status,omitempty"`
	LastSeenAt primitive.DateTime `json:"last_seen_at,omitempty"`
	CreatedAt  primitive.DateTime `json:"created_at,omitempty"`
}

type Project struct {
	UID            string              `json:"uid"`
	Name           string              `json:"name"`
	OrganisationID string              `json:"organisation_id"`
	Type           string              `json:"type"`
	CreatedAt      primitive.DateTime  `json:"created_at,omitempty"`
	UpdatedAt      primitive.DateTime  `json:"updated_at,omitempty"`
	DeletedAt      *primitive.DateTime `json:"deleted_at,omitempty"`
}

type Endpoint struct {
	UID         string              `json:"uid"`
	ProjectID   string              `json:"project_id"`
	OwnerID     string              `json:"owner_id,omitempty"`
	TargetURL   string              `json:"target_url"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	CreatedAt   primitive.DateTime  `json:"created_at,omitempty"`
	UpdatedAt   primitive.DateTime  `json:"updated_at,omitempty"`
	DeletedAt   *primitive.DateTime `json:"deleted_at,omitempty"`
}

type AckEventDelivery struct {
	UID string `json:"uid"`
}

type CLIEvent struct {
	UID     string              `json:"uid"`
	Headers map[string][]string `json:"headers"`
	Data    json.RawMessage     `json:"data"`

	// for filtering this event delivery
	EventType  string `json:"-"`
	DeviceID   string `json:"-"`
	EndpointID string `json:"-"`
	ProjectID  string `json:"-"`
}
