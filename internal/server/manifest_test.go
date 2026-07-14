package server

import (
	"context"
	"strings"
	"testing"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registeredToolNames boots the real stdio server, connects an in-memory client,
// and returns the exact set of tool names the server exposes over tools/list.
func registeredToolNames(t *testing.T) []string {
	t.Helper()

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
		AppEnv:       config.Dev,
		Transport:    "stdio",
	}
	srv, err := NewStdioServer(cfg, &mockGatewayClient{})
	if err != nil {
		t.Fatalf("NewStdioServer() failed: %v", err)
	}

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	ss, err := srv.server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect failed: %v", err)
	}
	defer ss.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "manifest-test", Version: "0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer cs.Close()

	res, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	names := make([]string, 0, len(res.Tools))
	for _, tool := range res.Tools {
		names = append(names, tool.Name)
	}
	return names
}

// TestToolManifest_ExcludesMutationsAndAdmin locks the curated-subset guarantee:
// the MCP must expose reads and messaging, but NEVER group/community mutations,
// the admin plane, or metrics — even though the SDK now has those methods.
func TestToolManifest_ExcludesMutationsAndAdmin(t *testing.T) {
	names := registeredToolNames(t)

	// Guard against a vacuous pass: the manifest must actually be populated with
	// the curated reads/two-way actions this MCP is supposed to expose.
	mustHave := []string{
		"list_groups", "get_group_info", "list_contacts",
		"get_contact_info", "get_avatar", "mark_read", "send_typing",
	}
	have := make(map[string]bool, len(names))
	for _, n := range names {
		have[n] = true
	}
	for _, n := range mustHave {
		if !have[n] {
			t.Errorf("expected curated tool %q to be registered; manifest=%v", n, names)
		}
	}

	// Forbidden capability tokens. None of these may appear as a substring of any
	// exposed tool name. Chosen to be unambiguous — they do not collide with the
	// allowed reads (list_groups / get_group_info carry no mutation verb).
	forbidden := []string{
		"create_group", "leave_group", "join_group", "group_join",
		"participant", "promote", "demote",
		"group_settings", "group_name", "group_topic", "group_photo", "group_requests",
		"invite", "community", "subgroup",
		"admin", "sessions", "metrics",
	}
	for _, name := range names {
		lower := strings.ToLower(name)
		for _, bad := range forbidden {
			if strings.Contains(lower, bad) {
				t.Errorf("tool %q exposes a forbidden capability token %q (curated subset excludes group mutations, community, admin, metrics)", name, bad)
			}
		}
	}
}
