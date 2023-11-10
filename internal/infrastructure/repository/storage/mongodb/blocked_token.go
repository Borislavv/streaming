package mongodb

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const BlockedTokensCollection = "blockedTokens"

type BlockedTokenRepository struct {
	db      *mongo.Collection
	mu      *sync.Mutex
	logger  logger.Logger
	timeout time.Duration
}

func NewBlockedTokenRepository(db *mongo.Database, logger logger.Logger, timeout time.Duration) *BlockedTokenRepository {
	return &BlockedTokenRepository{
		db:      db.Collection(BlockedTokensCollection),
		logger:  logger,
		mu:      &sync.Mutex{},
		timeout: timeout,
	}
}

func (r *BlockedTokenRepository) Insert(ctx context.Context, token string) error {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	blockedToken := agg.BlockedToken{
		BlockedToken: entity.BlockedToken{
			Value:     token,
			BlockedAt: time.Now(),
		},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
		},
	}

	res, err := r.db.InsertOne(qCtx, blockedToken, options.InsertOne())
	if err != nil {
		return r.logger.ErrorPropagate(err)
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); !ok {
		return r.logger.LogPropagate(
			errors.NewInternalValidationError(
				fmt.Sprintf("error occurred while inserting a blocked token '%v'", token),
			),
		)
	}

	return nil
}

func (r *BlockedTokenRepository) Has(ctx context.Context, token string) (found bool, err error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err = r.db.FindOne(qCtx, bson.M{"value": bson.M{"$eq": token}}).Decode(&agg.BlockedToken{}); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, r.logger.ErrorPropagate(err)
	}

	return true, nil
}
