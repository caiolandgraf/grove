package controllers

import (
	"net/http"

	"github.com/caiolandgraf/grove/internal/dto"
	"github.com/caiolandgraf/grove/internal/models"
	"github.com/go-fuego/fuego"
)

func GetInvoice(c fuego.ContextNoBody) (*dto.InvoiceResponse, error) {
	id := c.PathParam("invoice_id")

	item, err := models.Invoices().Find(id)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return toInvoiceDTO(item), nil
}

func ListInvoices(c fuego.ContextNoBody) (*dto.InvoicesListResponse, error) {
	items, err := models.Invoices().All()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	result := make([]dto.InvoiceResponse, len(items))
	for i, item := range items {
		result[i] = *toInvoiceDTO(&item)
	}

	return &dto.InvoicesListResponse{
		Items: result,
		Total: len(result),
	}, nil
}

func CreateInvoice(c fuego.ContextWithBody[dto.CreateInvoiceRequest]) (*dto.InvoiceResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	item := &models.Invoice{
		// TODO: map fields from body
	}
	_ = body

	if err := models.Invoices().Create(item); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	return toInvoiceDTO(item), nil
}

func UpdateInvoice(c fuego.ContextWithBody[dto.UpdateInvoiceRequest]) (*dto.InvoiceResponse, error) {
	id := c.PathParam("invoice_id")

	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	repo := models.Invoices()

	item, err := repo.Find(id)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	// TODO: map updatable fields from body
	_ = body

	if err := repo.Update(item); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	return toInvoiceDTO(item), nil
}

func DeleteInvoice(c fuego.ContextNoBody) (map[string]string, error) {
	id := c.PathParam("invoice_id")

	if err := models.Invoices().Delete(id); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return map[string]string{"message": "Invoice deleted successfully"}, nil
}

func toInvoiceDTO(m *models.Invoice) *dto.InvoiceResponse {
	return &dto.InvoiceResponse{
		ID: m.ID,
	}
}
