package dto


// payload dto
type CreditMemoPayloadDTO struct {
	Reference        string  `json:"reference"`
	Date             string  `json:"date"`
	DepositTo        int     `json:"deposit_to"`
	LiabilityAccount int     `json:"liability_account"`
	PeopleID         int64     `json:"people_id"`
	BuildingID       int64     `json:"building_id"`
	UnitID           int64     `json:"unit_id"`
	Amount           float64 `json:"amount"`
	Description      string  `json:"description"`
}
type CreateCreditMemoRequest struct {
	CreditMemoPayloadDTO
}

type UpdateCreditMemoRequest struct {
	ID               int     `json:"id"`
	CreditMemoPayloadDTO
}