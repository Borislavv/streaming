package reader

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"os"
)

const (
	ChunkSize    = 1024 * 1024 * 5 // 2.5MB
	ChunksBuffer = 10
)

type ResourceReader struct {
	logger logger.Logger
}

func NewReaderService(logger logger.Logger) *ResourceReader {
	return &ResourceReader{logger: logger}
}

// Read will read a resource and send file as butches of bytes
func (r *ResourceReader) Read(resource Resource) chan *Chunk {
	r.logger.Info(fmt.Sprintf("recourse '%v' reading started", resource.GetFilepath()))

	chunksCh := make(chan *Chunk, ChunksBuffer)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ResourceReader) handleRead(resource Resource, chunksCh chan *Chunk) {
	defer func() {
		close(chunksCh)
		r.logger.Info(fmt.Sprintf("recourse '%v' reading finished", resource.GetFilepath()))
	}()

	file, err := os.Open(resource.GetFilepath())
	if err != nil {
		r.logger.Error(err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			r.logger.Error(err)
			return
		}
	}()

	for {
		chunk := NewChunk(ChunkSize)

		chunk.Len, err = file.Read(chunk.Data)
		if err != nil {
			if err == io.EOF {
				break
			}
			r.logger.Error(err)
			return
		}

		r.sendChunk(chunk, chunksCh)
	}
}

func (r *ResourceReader) sendChunk(chunk *Chunk, chunksCh chan *Chunk) {
	if chunk.Len == 0 {
		return
	}

	if chunk.Len < ChunkSize {
		lastChunk := make([]byte, chunk.Len)
		lastChunk = chunk.Data[:chunk.Len]
		chunk.Data = lastChunk
	}

	if chunk.Len > 0 {
		r.logger.Info(fmt.Sprintf("%d bytes read and sent", chunk.Len))
		chunksCh <- chunk
	}
}
