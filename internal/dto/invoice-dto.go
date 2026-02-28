package dto

// Request DTOs

type CreateInvoiceRequest struct {
	// TODO: add fields
}

type UpdateInvoiceRequest struct {
	// TODO: add fields
}

// Response DTOs

type InvoiceResponse struct {
	ID string `json:"id"`
}

type InvoicesListResponse struct {
	Items []InvoiceResponse `json:"items"`
	Total int                 `json:"total"`
}
