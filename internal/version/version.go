package version

// Version is the MCP server release version (see VERSION at repo root).
// Overridden at link time via -ldflags in CI releases.
var Version = "0.4.0"

// MulticaAPI is the Multica REST API version this client targets.
const MulticaAPI = "app-api@7ebe0bf58d99"
