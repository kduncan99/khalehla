package deviceMgr

import "khalehla/pkg"

type DiskChannelInfo struct {
	channelName    string
	nodeIdentifier NodeIdentifier
	channel        *DiskChannel
}

func NewDiskChannelInfo(channelName string) *DiskChannelInfo {
	return &DiskChannelInfo{
		channelName:    channelName,
		nodeIdentifier: NodeIdentifier(pkg.NewFromStringToFieldata(channelName, 1)[0]),
	}
}

func (dci *DiskChannelInfo) CreateNode() {
	dci.channel = NewDiskChannel()
}

func (dci *DiskChannelInfo) GetChannel() Channel {
	return dci.channel
}

func (dci *DiskChannelInfo) GetNodeIdentifier() NodeIdentifier {
	return dci.nodeIdentifier
}

func (dci *DiskChannelInfo) GetNodeName() string {
	return dci.channelName
}

func (dci *DiskChannelInfo) GetNodeStatus() NodeStatus {
	return NodeStatusUp
}

func (dci *DiskChannelInfo) GetNodeType() NodeType {
	return NodeTypeDisk
}
