package repository_interface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Audio interface {
	Insert(ctx context.Context, audio *agg.Audio) (string, error)
	Update(ctx context.Context, audio *agg.Audio) error
	Find(ctx context.Context, id vo.ID) (*agg.Audio, error)
	FindList(ctx context.Context, dto dto_interface.ListAudioRequest) ([]*agg.Audio, error)
}
