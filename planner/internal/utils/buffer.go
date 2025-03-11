package utils

import "github.com/turbo-pioneer/planner/internal/item"

type Buffer struct {
	Data []*item.Item
}

func NewBuffer() *Buffer {
	return &Buffer{
		Data: make([]*item.Item, 0),
	}
}

func (b *Buffer) Push(item *item.Item) {
	b.Data = append(b.Data, item)
}

func (b *Buffer) Pop() *item.Item {
	if len(b.Data) == 0 {
		return nil
	}

	item := b.Data[len(b.Data)-1]
	b.Data = b.Data[:len(b.Data)-1]
	return item
}
