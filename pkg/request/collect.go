package request

// Collects the result of a sign or auth order using orderRef as reference. RP should keep on calling collect every two seconds if status is pending. RP must abort if status indicates failed. The user identity is returned when complete.
type CollectRequest struct {
	// The orderRef returned from auth or sign.
	OrderRef string `json:"orderRef"`
}