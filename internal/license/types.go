package license

type License struct {
	ID                   string   `json:"id"`
	CreatorID            string   `json:"creator_id"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	AITraingProhibited   bool     `json:"ai_training_prohibited"`
	BasePriceCents       int64    `json:"base_price_cents"`
	Currency             string   `json:"currency"`
	TerritoriesISO2Codes []string `json:"territories_iso2_codes"`
	Version              int64    `json:"version"`
}
