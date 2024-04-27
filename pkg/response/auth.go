package response

type AuthResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`

	// Used to compile the start url according to launching.
	// See https://www.bankid.com/utvecklare/guider/teknisk-integrationsguide/programstart
	AutoStartToken string `json:"autoStartToken"`

	// Used to compute the animated QR code.
	QRStartToken string `json:"qrStartToken"`

	// Used to compute the animated QR code.
	QRStartSecret string `json:"qrStartSecret"`
}