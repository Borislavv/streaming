package repository_interface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query/interface"
)

type Resource interface {
	FindOneByID(context.Context, query_interface.FindOneResourceByID) (*agg.Resource, error)
	Insert(context.Context, *agg.Resource) (*agg.Resource, error)
	Remove(context.Context, *agg.Resource) error
}
