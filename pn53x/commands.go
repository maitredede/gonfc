package pn53x

type Command byte

const ( // Miscellaneous
	Diagnose           Command = 0x00
	GetFirmwareVersion Command = 0x02
	GetGeneralStatus   Command = 0x04
	ReadRegister       Command = 0x06
	WriteRegister      Command = 0x08
	ReadGPIO           Command = 0x0C
	WriteGPIO          Command = 0x0E
	SetSerialBaudRate  Command = 0x10
	SetParameters      Command = 0x12
	SAMConfiguration   Command = 0x14
	PowerDown          Command = 0x16
	AlparCommandForTDA Command = 0x18

	// RC-S360 has another command 0x18 for reset &..?
)

const ( // RF communication
	RFConfiguration  Command = 0x32
	RFRegulationTest Command = 0x58
)

const ( // Initiator
	InJumpForDEP                Command = 0x56
	InJumpForPSL                Command = 0x46
	InListPassiveTarget         Command = 0x4A
	InATR                       Command = 0x50
	InPSL                       Command = 0x4E
	InDataExchange              Command = 0x40
	InCommunicateThru           Command = 0x42
	InQuartetByteExchange       Command = 0x38
	InDeselect                  Command = 0x44
	InRelease                   Command = 0x52
	InSelect                    Command = 0x54
	InActivateDeactivatePaypass Command = 0x48
	InAutoPoll                  Command = 0x60
)

const ( // Target
	TgInitAsTarget        Command = 0x8C
	TgSetGeneralBytes     Command = 0x92
	TgGetData             Command = 0x86
	TgSetData             Command = 0x8E
	TgSetDataSecure       Command = 0x96
	TgSetMetaData         Command = 0x94
	TgSetMetaDataSecure   Command = 0x98
	TgGetInitiatorCommand Command = 0x88
	TgResponseToInitiator Command = 0x90
	TgGetTargetStatus     Command = 0x8A
)

type commandInfo struct {
	b     byte
	chips ChipType
	name  string
}

func mkCommand(b Command, chips ChipType, name string) commandInfo {
	return commandInfo{
		b:     byte(b),
		chips: chips,
		name:  name,
	}
}

var pn53xCommands map[Command]commandInfo = map[Command]commandInfo{
	// Miscellaneous
	Diagnose:           mkCommand(Diagnose, PN531|PN532|PN533|RCS360, "Diagnose"),
	GetFirmwareVersion: mkCommand(GetFirmwareVersion, PN531|PN532|PN533|RCS360, "GetFirmwareVersion"),
	GetGeneralStatus:   mkCommand(GetGeneralStatus, PN531|PN532|PN533|RCS360, "GetGeneralStatus"),
	ReadRegister:       mkCommand(ReadRegister, PN531|PN532|PN533|RCS360, "ReadRegister"),
	WriteRegister:      mkCommand(WriteRegister, PN531|PN532|PN533|RCS360, "WriteRegister"),
	ReadGPIO:           mkCommand(ReadGPIO, PN531|PN532|PN533, "ReadGPIO"),
	WriteGPIO:          mkCommand(WriteGPIO, PN531|PN532|PN533, "WriteGPIO"),
	SetSerialBaudRate:  mkCommand(SetSerialBaudRate, PN531|PN532|PN533, "SetSerialBaudRate"),
	SetParameters:      mkCommand(SetParameters, PN531|PN532|PN533|RCS360, "SetParameters"),
	SAMConfiguration:   mkCommand(SAMConfiguration, PN531|PN532, "SAMConfiguration"),
	PowerDown:          mkCommand(PowerDown, PN531|PN532, "PowerDown"),
	AlparCommandForTDA: mkCommand(AlparCommandForTDA, PN533|RCS360, "AlparCommandForTDA"), // Has another usage on RC-S360...

	// RF communication
	RFConfiguration:  mkCommand(RFConfiguration, PN531|PN532|PN533|RCS360, "RFConfiguration"),
	RFRegulationTest: mkCommand(RFRegulationTest, PN531|PN532|PN533, "RFRegulationTest"),

	// Initiator
	InJumpForDEP:                mkCommand(InJumpForDEP, PN531|PN532|PN533|RCS360, "InJumpForDEP"),
	InJumpForPSL:                mkCommand(InJumpForPSL, PN531|PN532|PN533, "InJumpForPSL"),
	InListPassiveTarget:         mkCommand(InListPassiveTarget, PN531|PN532|PN533|RCS360, "InListPassiveTarget"),
	InATR:                       mkCommand(InATR, PN531|PN532|PN533, "InATR"),
	InPSL:                       mkCommand(InPSL, PN531|PN532|PN533, "InPSL"),
	InDataExchange:              mkCommand(InDataExchange, PN531|PN532|PN533, "InDataExchange"),
	InCommunicateThru:           mkCommand(InCommunicateThru, PN531|PN532|PN533|RCS360, "InCommunicateThru"),
	InQuartetByteExchange:       mkCommand(InQuartetByteExchange, PN533, "InQuartetByteExchange"),
	InDeselect:                  mkCommand(InDeselect, PN531|PN532|PN533|RCS360, "InDeselect"),
	InRelease:                   mkCommand(InRelease, PN531|PN532|PN533|RCS360, "InRelease"),
	InSelect:                    mkCommand(InSelect, PN531|PN532|PN533, "InSelect"),
	InAutoPoll:                  mkCommand(InAutoPoll, PN532, "InAutoPoll"),
	InActivateDeactivatePaypass: mkCommand(InActivateDeactivatePaypass, PN533, "InActivateDeactivatePaypass"),

	// Target
	TgInitAsTarget:        mkCommand(TgInitAsTarget, PN531|PN532|PN533, "TgInitAsTarget"),
	TgSetGeneralBytes:     mkCommand(TgSetGeneralBytes, PN531|PN532|PN533, "TgSetGeneralBytes"),
	TgGetData:             mkCommand(TgGetData, PN531|PN532|PN533, "TgGetData"),
	TgSetData:             mkCommand(TgSetData, PN531|PN532|PN533, "TgSetData"),
	TgSetDataSecure:       mkCommand(TgSetDataSecure, PN533, "TgSetDataSecure"),
	TgSetMetaData:         mkCommand(TgSetMetaData, PN531|PN532|PN533, "TgSetMetaData"),
	TgSetMetaDataSecure:   mkCommand(TgSetMetaDataSecure, PN533, "TgSetMetaDataSecure"),
	TgGetInitiatorCommand: mkCommand(TgGetInitiatorCommand, PN531|PN532|PN533, "TgGetInitiatorCommand"),
	TgResponseToInitiator: mkCommand(TgResponseToInitiator, PN531|PN532|PN533, "TgResponseToInitiator"),
	TgGetTargetStatus:     mkCommand(TgGetTargetStatus, PN531|PN532|PN533, "TgGetTargetStatus"),
}

func (pnd *Chip) cmdTrace(cmd byte) {
	info := pn53xCommands[Command(cmd)]
	pnd.logger.Debugf("  cmd:%s", info.name)
}
