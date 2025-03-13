package utils

import "github.com/turbo-pioneer/planner/internal/models"

type Buffer struct {
	Data []*models.Item
}

func NewBuffer() *Buffer {
	return &Buffer{
		Data: make([]*models.Item, 0),
	}
}

func (b *Buffer) Push(item *models.Item) {
	b.Data = append(b.Data, item)
}

func (b *Buffer) Pop() *models.Item {
	if len(b.Data) == 0 {
		return nil
	}

	item := b.Data[len(b.Data)-1]
	b.Data = b.Data[:len(b.Data)-1]
	return item
}
