package players

import (
	math2 "math"
	"sync"

	"github.com/irmine/gomine/entities"
	"github.com/irmine/gomine/entities/math"
	"github.com/irmine/gomine/interfaces"
	"github.com/irmine/gomine/net"
	"github.com/irmine/gomine/net/packets/data"
	"github.com/irmine/gomine/net/packets/types"
	"github.com/irmine/gomine/vectors"
	"github.com/irmine/goraklib/server"
	"github.com/irmine/gomine/permissions"
)

type Player struct {
	*entities.Human
	interfaces.IMinecraftSession

	playerName  string
	displayName string

	spawned bool

	permissions     map[string]*permissions.Permission
	permissionGroup *permissions.Group

	onGround     bool
	viewDistance int32

	skinId       string
	skinData     []byte
	capeData     []byte
	geometryName string
	geometryData string

	finalized bool

	server interfaces.IServer

	mux        sync.Mutex
	usedChunks map[int]interfaces.IChunk
}

/**
 * Returns a new player with the given name.
 */
func NewPlayer(server interfaces.IServer, name string) *Player {
	var player = &Player{}
	player.IMinecraftSession = &net.MinecraftSession{}

	player.playerName = name
	player.displayName = name

	player.usedChunks = make(map[int]interfaces.IChunk)

	player.permissions = make(map[string]*permissions.Permission)
	player.permissionGroup = server.GetPermissionManager().GetDefaultGroup()

	player.server = server

	return player
}

/**
 * Returns a new player with the given minecraft session.
 */
func (player *Player) New(server interfaces.IServer, session interfaces.IMinecraftSession, name string) interfaces.IPlayer {
	var pl = NewPlayer(server, name)
	pl.IMinecraftSession = session
	return pl
}

/**
 * Sets the player's minecraft session.
 */
func (player *Player) SetMinecraftSession(session interfaces.IMinecraftSession) {
	player.IMinecraftSession = session
}

/**
 * Returns a new minecraft session with the given server, session and login packet.
 */
func (player *Player) NewMinecraftSession(server interfaces.IServer, session *server.Session, data types.SessionData) interfaces.IMinecraftSession {
	return net.NewMinecraftSession(server, session, data)
}

/**
 * Returns if this player has been placed in a world.
 */
func (player *Player) IsInWorld() bool {
	return player.Dimension != nil && player.Level != nil
}

/**
 * Places this player inside of a level and dimension.
 */
func (player *Player) PlaceInWorld(position *vectors.TripleVector, rotation *math.Rotation, level interfaces.ILevel, dimension interfaces.IDimension) {
	player.Human = entities.NewHuman(player.GetDisplayName(), position, rotation, vectors.NewTripleVector(0, 0, 0), level, dimension)
}

/**
 * Checks if this player is finalized.
 */
func (player *Player) IsFinalized() bool {
	return player.finalized
}

/**
 * Sets this player finalized.
 */
func (player *Player) SetFinalized() {
	player.finalized = true
}

/**
 * Spawns this player to the given other player.
 */
func (player *Player) SpawnPlayerTo(player2 interfaces.IPlayer) {
	if !player2.HasSpawned() {
		return
	}
	player.SendAddPlayer(player2)
}

/**
 * Spawns this player to all other players.
 */
func (player *Player) SpawnPlayerToAll() {
	for _, p := range player.GetServer().GetPlayerFactory().GetPlayers() {
		if player == p {
			continue
		}
		player.SpawnPlayerTo(p)
	}
}

/**
 * Sets the view distance of this player.
 */
func (player *Player) SetViewDistance(distance int32) {
	player.viewDistance = distance
}

/**
 * Returns the view distance of this player.
 */
func (player *Player) GetViewDistance() int32 {
	return player.viewDistance
}

/**
 * Returns the main server.
 */
func (player *Player) GetServer() interfaces.IServer {
	return player.server
}

/**
 * Returns the username the player used to join the server.
 */
func (player *Player) GetName() string {
	return player.playerName
}

/**
 * Sets the player name of this player.
 * Note: This function is internal, and should not be used by plugins.
 */
func (player *Player) SetName(name string) {
	player.playerName = name
}

/**
 * Returns the name the player shows in-game.
 */
func (player *Player) GetDisplayName() string {
	return player.displayName
}

/**
 * Sets the name other players can see in-game.
 */
func (player *Player) SetDisplayName(name string) {
	player.displayName = name
}

/**
 * Returns the permission group this player is in.
 */
func (player *Player) GetPermissionGroup() *permissions.Group {
	return player.permissionGroup
}

/**
 * Sets the permission group of this player.
 */
func (player *Player) SetPermissionGroup(group *permissions.Group) {
	player.permissionGroup = group
}

/**
 * Checks if this player has a permission.
 */
func (player *Player) HasPermission(permission string) bool {
	if player.GetPermissionGroup().HasPermission(permission) {
		return true
	}
	var _, exists = player.permissions[permission]
	return exists
}

/**
 * Adds a permission to the player.
 * Returns true if a permission with the same name was overwritten.
 */
func (player *Player) AddPermission(permission *permissions.Permission) bool {
	var hasPermission = player.HasPermission(permission.GetName())

	player.permissions[permission.GetName()] = permission

	return hasPermission
}

/**
 * Deletes a permission from the player.
 * This does not delete the permission from the group the player is in.
 */
func (player *Player) RemovePermission(permission string) bool {
	if !player.HasPermission(permission) {
		return false
	}

	delete(player.permissions, permission)

	return true
}

/**
 * Teleport player to a new position
 */
func (player *Player) Teleport(v *vectors.TripleVector, rot *math.Rotation) {
	player.SetPosition(v)
	player.SendMovePlayer(player, *v, *rot, data.MoveTeleport, player.onGround, 0)
}

/**
 * Sets the skin ID/name of the player.
 */
func (player *Player) SetSkinId(id string) {
	player.skinId = id
}

/**
 * Returns the skin ID/name of the player.
 */
func (player *Player) GetSkinId() string {
	return player.skinId
}

/**
 * Returns the skin data of the player. (RGBA byte array)
 */
func (player *Player) GetSkinData() []byte {
	return player.skinData
}

/**
 * Sets the skin data of the player. (RGBA byte array)
 */
func (player *Player) SetSkinData(data []byte) {
	player.skinData = data
}

/**
 * Returns the cape data of the player. (RGBA byte array)
 */
func (player *Player) GetCapeData() []byte {
	return player.capeData
}

/**
 * Sets the cape data of the player. (RGBA byte array)
 */
func (player *Player) SetCapeData(data []byte) {
	player.capeData = data
}

/**
 * Returns the geometry name of the player.
 */
func (player *Player) GetGeometryName() string {
	return player.geometryName
}

/**
 * Sets the geometry name of the player.
 */
func (player *Player) SetGeometryName(name string) {
	player.geometryName = name
}

/**
 * Returns the geometry data (json string) of the player.
 */
func (player *Player) GetGeometryData() string {
	return player.geometryData
}

/**
 * Sets the geometry data (json string) of the player.
 */
func (player *Player) SetGeometryData(data string) {
	player.geometryData = data
}

/**
 * Sends a chunk to the player.
 */
func (player *Player) SendChunk(chunk interfaces.IChunk, index int) {
	player.SendFullChunkData(chunk)
	player.mux.Lock()
	player.usedChunks[index] = chunk
	player.mux.Unlock()
}

/**
 * Synchronizes the server's player movement with the client movement and adjusts chunks.
 */
func (player *Player) SyncMove(x, y, z, pitch, yaw, headYaw float32, onGround bool) {
	player.SetPosition(vectors.NewTripleVector(x, y, z))
	player.Rotation.Pitch += pitch
	player.Rotation.Yaw += yaw
	player.Rotation.HeadYaw += headYaw
	player.onGround = onGround

	var chunkX = int32(math2.Floor(float64(x))) >> 4
	var chunkZ = int32(math2.Floor(float64(z))) >> 4

	var rs = player.GetViewDistance() * player.GetViewDistance()

	player.mux.Lock()
	for index, chunk := range player.usedChunks {
		xDist := chunkX - chunk.GetX()
		zDist := chunkZ - chunk.GetZ()

		if xDist*xDist+zDist*zDist > rs {
			chunk.RemoveViewer(player)
			delete(player.usedChunks, index)

			for _, entity := range chunk.GetEntities() {
				entity.DespawnFrom(player)
			}
		}
	}
	player.mux.Unlock()
}

/**
 * Checks if the player has a chunk with the given index in use.
 */
func (player *Player) HasChunkInUse(index int) bool {
	player.mux.Lock()
	_, ok := player.usedChunks[index]
	player.mux.Unlock()
	return ok
}

/**
 * Checks if the player has any used chunks.
 */
func (player *Player) HasAnyChunkInUse() bool {
	return len(player.usedChunks) > 0
}

func (player *Player) Tick() {
	if player.HasSpawned() {
		player.Entity.Tick()
	}
}

/**
 * Updates all entity attributes
 */
func (player *Player) UpdateAttributes() {
	player.SendUpdateAttributes(player, player.GetAttributeMap())
}

/**
 * Sends entity data
 */
func (player *Player) SendEntityData() {
}

/**
 * Sends a message to this player.
 */
func (player *Player) SendMessage(message string) {
	var text = types.Text{}
	text.Message = message
	player.SendText(text)
}

/**
 * Checks if this player has spawned.
 */
func (player *Player) HasSpawned() bool {
	return player.spawned
}

/**
 * Sets this player spawned.
 */
func (player *Player) SetSpawned(value bool) {
	player.spawned = value
}

/**
 * Checks if the player is initialized.
 */
func (player *Player) IsInitialized() bool {
	return player.IMinecraftSession != nil
}
