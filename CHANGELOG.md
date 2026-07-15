# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Canonical `chat` addressing on all send tools.** Every send tool now accepts
  `chat` (a bare number, a user JID `@s.whatsapp.net`, a group JID `@g.us`, or a
  `@lid`) alongside the existing back-compat `to`. `chat` wins when both are set,
  and the address is passed through to the gateway as-is (no `FormatMSISDN`
  force-wrapping).
- **Reply and mention threading on all send tools** via `reply_to_id`,
  `reply_to_sender`, `reply_to_text`, and `mentions[]`.
- **Contact & group read tools**: `list_contacts` (paginated), `get_contact_info`,
  `get_avatar` (404 "not set" / 403 "hidden" surfaced as results, not errors),
  `list_groups`, and `get_group_info`.
- **Two-way conversation tools**: `mark_read` and `send_typing`
  (`composing`/`recording`/`paused`). These are conversation-affecting outbound
  actions governed by the gateway's outbound pacer (per-account pace +
  per-recipient cap); over-budget calls are paced or rejected with `429`.
- Guard test (`internal/server/manifest_test.go`) that asserts the exposed tool
  manifest contains no group/community mutation, admin, or metrics capability.

### Changed

- Bumped the gateway SDK dependency (`github.com/glennprays/whatsapp-gateway-sdk-go`)
  to `v0.7.0`.
- Aligned the MCP server's reported `Implementation.Version` to `0.7.0` (was the
  inconsistent hardcoded `1.0.0`).

### Excluded (by design — curated subset)

This release intentionally does **not** expose, and will not expose, the
following capabilities even though the SDK and gateway now support them, so an
autonomous agent cannot perform destructive or account-wide actions:

- Group/community **mutations**: create, leave, participants (add/remove/promote/
  demote), settings, name, topic, photo, invite, join, requests.
- **Community** operations (sub-group link/unlink).
- The operator **admin** plane (`/admin/sessions`) and **metrics** (`/metrics`).

Perform those operations by calling the gateway's REST API directly from a
trusted backend.
