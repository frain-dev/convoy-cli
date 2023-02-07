package convoy_cli

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListenRequest struct {
	HostName  string `json:"host_name"`
	ProjectID string `json:"project_id"`
	DeviceID  string `json:"device_id"`
	SourceID  string `json:"source_id"`

	Since     string `json:"-"`
	ForwardTo string `json:"-"`
	// EventTypes []string `json:"event_types"`
}

type LoginRequest struct {
	HostName string `json:"host_name"`
	DeviceID string `json:"device_id"`
}

type LoginResponse struct {
	Projects []ProjectDevice `json:"projects"`
	UserName string          `json:"user_name"`

	//Device   *Device   `json:"device"`
	//Project  *Project  `json:"project"`
	//Endpoint *Endpoint `json:"endpoint"`
}

type ProjectDevice struct {
	Project *Project `json:"project"`
	Device  *Device  `json:"device"`
}

type Device struct {
	UID        string             `json:"uid"`
	ProjectID  string             `json:"project_id,omitempty"`
	EndpointID string             `json:"endpoint_id,omitempty"`
	HostName   string             `json:"host_name,omitempty"`
	Status     string             `json:"status,omitempty"`
	LastSeenAt primitive.DateTime `json:"last_seen_at,omitempty"`
}

type Project struct {
	UID            string `json:"uid"`
	Name           string `json:"name"`
	OrganisationID string `json:"organisation_id"`
	Type           string `json:"type"`
}

type Endpoint struct {
	UID         string `json:"uid"`
	ProjectID   string `json:"project_id"`
	OwnerID     string `json:"owner_id,omitempty"`
	TargetURL   string `json:"target_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AckEventDelivery struct {
	UID string `json:"uid"`
}

type CLIEvent struct {
	UID     string              `json:"uid"`
	Headers map[string][]string `json:"headers"`
	Data    json.RawMessage     `json:"data"`
}
