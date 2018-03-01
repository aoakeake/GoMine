package p200

import (
	"github.com/irmine/gomine/net/info"
	"github.com/irmine/gomine/net/packets"
	"github.com/irmine/worlds/entities/data"
)

type UpdateAttributesPacket struct {
	*packets.Packet
	RuntimeId  uint64
	Attributes data.AttributeMap
}

func NewUpdateAttributesPacket() *UpdateAttributesPacket {
	return &UpdateAttributesPacket{packets.NewPacket(info.PacketIds200[info.UpdateAttributesPacket]), 0, data.NewAttributeMap()}
}

func (pk *UpdateAttributesPacket) Encode() {
	pk.PutRuntimeId(pk.RuntimeId)
	pk.PutAttributeMap(pk.Attributes)
}

func (pk *UpdateAttributesPacket) Decode() {
	pk.RuntimeId = pk.GetRuntimeId()
	pk.Attributes = pk.GetAttributeMap()
}
