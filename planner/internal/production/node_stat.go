package production

type NodeStats struct {
	TotalItemsProduced int
	InputBuffers       map[string]int
	OutputBuffers      map[string]int
}

func NewNodeStats() *NodeStats {
	return &NodeStats{}
}

func (ns *NodeStats) SetItemsProduced(total int) {
	ns.TotalItemsProduced = total
}

func (ns *NodeStats) SetInputBuffer(itemName string, amount int) {
	ns.InputBuffers[itemName] = amount
}

func (ns *NodeStats) SetOutputBuffer(itemName string, amount int) {
	ns.OutputBuffers[itemName] = amount
}
