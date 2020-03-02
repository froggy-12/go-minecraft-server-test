package prot

import (
	"minecraft-server/apis/logs"
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
	"minecraft-server/impl/game/mode"
	"minecraft-server/impl/prot/states"
)

type packets struct {
	util.Watcher

	logger  *logs.Logging
	packetI map[base.PacketState]map[int32]func() base.PacketI // UUID to I server_data

	join chan base.PlayerAndConnection
	quit chan base.PlayerAndConnection
}

func NewPackets(join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) base.Packets {
	packets := &packets{
		Watcher: util.NewWatcher(),

		logger:  logs.NewLogging("protocol", logs.BasicLevel...),
		packetI: createPacketI(),
	}

	mode.HandleState0(packets)
	mode.HandleState1(packets)
	mode.HandleState2(packets, join)
	mode.HandleState3(packets, packets.logger, join)

	return packets
}

func (p *packets) GetPacketI(uuid int32, state base.PacketState) base.PacketI {
	creator := p.packetI[state][uuid]
	if creator == nil {
		return nil
	}

	return creator()
}

func createPacketI() map[base.PacketState]map[int32]func() base.PacketI {
	return map[base.PacketState]map[int32]func() base.PacketI{
		base.Shake: {
			0x00: func() base.PacketI {
				return &states.PacketIHandshake{}
			},
		},
		base.Status: {
			0x00: func() base.PacketI {
				return &states.PacketIRequest{}
			},
			0x01: func() base.PacketI {
				return &states.PacketIPing{}
			},
		},
		base.Login: {
			0x00: func() base.PacketI {
				return &states.PacketILoginStart{}
			},
			0x01: func() base.PacketI {
				return &states.PacketIEncryptionResponse{}
			},
			0x02: func() base.PacketI {
				return &states.PacketILoginPluginResponse{}
			},
		},
	}
}
