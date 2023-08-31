package pn53x

// TargetMode PN53x target mode enumeration
// chips/pn53x.h pn53x_target_type
type TargetMode byte

const (
	/** Configure the PN53x to accept all initiator mode */
	PTM_NORMAL TargetMode = 0x00
	/** Configure the PN53x to accept to be initialized only in passive mode */
	PTM_PASSIVE_ONLY TargetMode = 0x01
	/** Configure the PN53x to accept to be initialized only as DEP target */
	PTM_DEP_ONLY TargetMode = 0x02
	/** Configure the PN532 to accept to be initialized only as ISO/IEC14443-4 PICC */
	PTM_ISO14443_4_PICC_ONLY TargetMode = 0x04
)
