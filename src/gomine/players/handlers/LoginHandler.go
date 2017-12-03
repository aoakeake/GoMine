package handlers

import (
	"gomine/players"
	"gomine/net/info"
	"gomine/interfaces"
	"gomine/net/packets"
	"goraklib/server"
)

type LoginHandler struct {
	*players.PacketHandler
}

func NewLoginHandler() LoginHandler {
	return LoginHandler{players.NewPacketHandler(info.LoginPacket)}
}

/**
 * Handles the main login process.
 */
func (handler LoginHandler) Handle(packet interfaces.IPacket, player interfaces.IPlayer, session *server.Session, server interfaces.IServer) bool {
	if loginPacket, ok := packet.(*packets.LoginPacket); ok {
		_, online := server.GetPlayerFactory().GetPlayerByName(loginPacket.Username)
		if online == nil {
			return false
		}

		var player = players.NewPlayer(server, session, loginPacket.Username, loginPacket.ClientUUID, loginPacket.ClientXUID, loginPacket.ClientId)
		player.SetLanguage(loginPacket.Language)
		player.SetSkinId(loginPacket.SkinId)
		player.SetSkinData(loginPacket.SkinData)
		player.SetCapeData(loginPacket.CapeData)
		player.SetGeometryName(loginPacket.GeometryName)
		player.SetGeometryData(string(loginPacket.GeometryData))

		pk := packets.NewPlayStatusPacket()
		pk.Status = 0
		server.GetRakLibAdapter().SendPacket(pk, session)

		pk3 := packets.NewResourcePackInfoPacket()
		server.GetRakLibAdapter().SendPacket(pk3, session)

		server.GetPlayerFactory().AddPlayer(player, session)
	}

	return true
}