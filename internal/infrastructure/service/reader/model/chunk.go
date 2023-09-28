package model

type Chunk struct {
	Len  int
	Data []byte
	Err  error
}

func NewChunk(size int) *Chunk {
	return &Chunk{Data: make([]byte, size)}
}

func (c *Chunk) GetLen() int {
	return c.Len
}

func (c *Chunk) SetLen(len int) {
	c.Len = len
}

func (c *Chunk) GetData() []byte {
	return c.Data
}

func (c *Chunk) SetData(data []byte) {
	c.Data = data
}

func (c *Chunk) GetError() error {
	return c.Err
}

func (c *Chunk) SetError(err error) {
	c.Err = err
}