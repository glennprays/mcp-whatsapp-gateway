package tools

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

// ListContactsInput is the request payload for the list_contacts tool.
type ListContactsInput struct {
	// Limit caps the page size. Optional; the gateway defaults to 100 and caps at 500.
	Limit int `json:"limit,omitempty"`
	// Offset paginates the synced address book. Optional.
	Offset int `json:"offset,omitempty"`
}

// ListContactsDirect lists the account's locally-synced contacts. An empty
// address book is a normal result, never an error.
func ListContactsDirect(client gateway.GatewayClient, input ListContactsInput) (waga.ContactListResponse, error) {
	resp, err := client.ListContacts(context.Background(), input.Limit, input.Offset)
	if err != nil {
		return waga.ContactListResponse{}, fmt.Errorf("list_contacts: %w", err)
	}
	return *resp, nil
}

// GetContactInfoInput is the request payload for the get_contact_info tool.
type GetContactInfoInput struct {
	// Chat is the canonical recipient (a bare number, a user JID, or a "@lid").
	Chat string `json:"chat"`
}

// GetContactInfoDirect returns a server-side profile lookup for one user.
func GetContactInfoDirect(client gateway.GatewayClient, input GetContactInfoInput) (waga.ContactInfoResponse, error) {
	if input.Chat == "" {
		return waga.ContactInfoResponse{}, fmt.Errorf("chat is required")
	}
	resp, err := client.GetContactInfo(context.Background(), input.Chat)
	if err != nil {
		return waga.ContactInfoResponse{}, fmt.Errorf("get_contact_info: %w", err)
	}
	return *resp, nil
}

// GetAvatarInput is the request payload for the get_avatar tool.
type GetAvatarInput struct {
	// Chat is the canonical recipient (user or "@g.us" group).
	Chat string `json:"chat"`
	// Preview requests the low-res thumbnail instead of the full-resolution image.
	Preview bool `json:"preview,omitempty"`
}

// GetAvatarResult is the result of a get_avatar lookup. A missing picture (404)
// or a hidden one (403) is surfaced as Available=false with a Reason, not an
// error, so agents can distinguish "no avatar" from a real failure.
type GetAvatarResult struct {
	JID       string `json:"jid"`
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"` // "not_set" (404) or "hidden" (403) when unavailable
	URL       string `json:"url,omitempty"`    // time-limited WhatsApp CDN link
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"` // "image" (full res) or "preview" (thumbnail)
}

// GetAvatarDirect fetches a chat's profile picture, surfacing "not set" and
// "hidden" as results rather than errors.
func GetAvatarDirect(client gateway.GatewayClient, input GetAvatarInput) (GetAvatarResult, error) {
	if input.Chat == "" {
		return GetAvatarResult{}, fmt.Errorf("chat is required")
	}
	resp, err := client.GetAvatar(context.Background(), input.Chat, input.Preview)
	if err != nil {
		if waga.IsNotFound(err) {
			return GetAvatarResult{JID: input.Chat, Available: false, Reason: "not_set"}, nil
		}
		if waga.IsForbidden(err) {
			return GetAvatarResult{JID: input.Chat, Available: false, Reason: "hidden"}, nil
		}
		return GetAvatarResult{}, fmt.Errorf("get_avatar: %w", err)
	}
	return GetAvatarResult{
		JID:       resp.JID,
		Available: true,
		URL:       resp.URL,
		ID:        resp.ID,
		Type:      resp.Type,
	}, nil
}

// ListGroupsInput is the (empty) request payload for the list_groups tool.
type ListGroupsInput struct{}

// ListGroupsDirect lists the account's joined groups as lightweight summaries.
func ListGroupsDirect(client gateway.GatewayClient, _ ListGroupsInput) (waga.GroupListResponse, error) {
	resp, err := client.ListGroups(context.Background())
	if err != nil {
		return waga.GroupListResponse{}, fmt.Errorf("list_groups: %w", err)
	}
	return *resp, nil
}

// GetGroupInfoInput is the request payload for the get_group_info tool.
type GetGroupInfoInput struct {
	// Chat must be a group JID ("@g.us"). The account must be a member.
	Chat string `json:"chat"`
}

// GetGroupInfoDirect returns one group's full detail plus its participant roster.
func GetGroupInfoDirect(client gateway.GatewayClient, input GetGroupInfoInput) (waga.GroupInfoResponse, error) {
	if input.Chat == "" {
		return waga.GroupInfoResponse{}, fmt.Errorf("chat is required")
	}
	resp, err := client.GetGroupInfo(context.Background(), input.Chat)
	if err != nil {
		return waga.GroupInfoResponse{}, fmt.Errorf("get_group_info: %w", err)
	}
	return *resp, nil
}
