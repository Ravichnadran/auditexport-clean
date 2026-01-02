package run

type Standard string

const (
	ISO27001 Standard = "iso27001"
	SOC2     Standard = "soc2"
)

type Capabilities struct {
	AllowExtendedControls bool
	AllowCICD             bool
	AllowAuditLogs        bool
}

func CapabilitiesForStandard(s Standard) Capabilities {
	switch s {
	case SOC2:
		return Capabilities{
			AllowExtendedControls: true,
			AllowCICD:             true,
			AllowAuditLogs:        true,
		}
	default: // ISO 27001
		return Capabilities{
			AllowExtendedControls: false,
			AllowCICD:             false,
			AllowAuditLogs:        false,
		}
	}
}
