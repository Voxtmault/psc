package interfaces

import (
	"context"

	"github.com/voxtmault/psc/models"
)

type KasusInterface interface {
	Get(ctx context.Context, filter *models.KasusFilter) (*models.Response, error)
	Create(ctx context.Context, payload *models.KasusCreate) (*models.Response, error)
}
