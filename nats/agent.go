package nats

// AgentInfo holds the metadata for an agent instance.
type AgentInfo struct {
	// InstanceID is the framework-assigned per-instance identifier (§3.4).
	InstanceID string
	// Agent is the metadata.agent value (e.g. "phero", "claude-code").
	Agent string
	// Owner is the metadata.owner value.
	Owner string
	// Session is the optional metadata.session value.
	Session string
	// Name is the instance name — the 5th token of the prompt endpoint subject.
	Name string
	// ProtocolVersion is the metadata.protocol_version value.
	ProtocolVersion string
	// PromptSubject is the subject to publish prompt requests to (§4.3).
	PromptSubject string
	// StatusSubject is the on-demand status request subject (§8.7).
	StatusSubject string
	// MaxPayloadBytes is the parsed max_payload endpoint metadata (§2.1).
	MaxPayloadBytes int64
	// AttachmentsOk mirrors the attachments_ok endpoint metadata flag (§2.1).
	AttachmentsOk bool
}
