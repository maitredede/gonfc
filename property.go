package gonfc

type Property byte

const (
	/**
	 * Default command processing timeout
	 * Property value's (duration) unit is ms and 0 means no timeout (infinite).
	 * Default value is set by driver layer
	 */
	NP_TIMEOUT_COMMAND Property = iota
	/**
	* Timeout between ATR_REQ and ATR_RES
	* When the device is in initiator mode, a target is considered as mute if no
	* valid ATR_RES is received within this timeout value.
	* Default value for this property is 103 ms on PN53x based devices.
	 */
	NP_TIMEOUT_ATR
	/**
	* Timeout value to give up reception from the target in case of no answer.
	* Default value for this property is 52 ms).
	 */
	NP_TIMEOUT_COM
	/** Let the PN53X chip handle the CRC bytes. This means that the chip appends
	 * the CRC bytes to the frames that are transmitted. It will parse the last
	 * bytes from received frames as incoming CRC bytes. They will be verified
	 * against the used modulation and protocol. If an frame is expected with
	 * incorrect CRC bytes this option should be disabled. Example frames where
	 * this is useful are the ATQA and UID+BCC that are transmitted without CRC
	 * bytes during the anti-collision phase of the ISO14443-A protocol. */
	NP_HANDLE_CRC
	/** Parity bits in the network layer of ISO14443-A are by default generated and
	* validated in the PN53X chip. This is a very convenient feature. On certain
	* times though it is useful to get full control of the transmitted data. The
	* proprietary MIFARE Classic protocol uses for example custom (encrypted)
	* parity bits. For interoperability it is required to be completely
	* compatible, including the arbitrary parity bits. When this option is
	* disabled, the functions to communicating bits should be used. */
	NP_HANDLE_PARITY
	/** This option can be used to enable or disable the electronic field of the
	* NFC device. */
	NP_ACTIVATE_FIELD
	/** The internal CRYPTO1 co-processor can be used to transmit messages
	* encrypted. This option is automatically activated after a successful MIFARE
	* Classic authentication. */
	NP_ACTIVATE_CRYPTO1
	/** The default configuration defines that the PN53X chip will try indefinitely
	* to invite a tag in the field to respond. This could be desired when it is
	* certain a tag will enter the field. On the other hand, when this is
	* uncertain, it will block the application. This option could best be compared
	* to the (NON)BLOCKING option used by (socket)network programming. */
	NP_INFINITE_SELECT
	/** If this option is enabled, frames that carry less than 4 bits are allowed.
	* According to the standards these frames should normally be handles as
	* invalid frames. */
	NP_ACCEPT_INVALID_FRAMES
	/** If the NFC device should only listen to frames, it could be useful to let
	* it gather multiple frames in a sequence. They will be stored in the internal
	* FIFO of the PN53X chip. This could be retrieved by using the receive data
	* functions. Note that if the chip runs out of bytes (FIFO = 64 bytes long),
	* it will overwrite the first received frames, so quick retrieving of the
	* received data is desirable. */
	NP_ACCEPT_MULTIPLE_FRAMES
	/** This option can be used to enable or disable the auto-switching mode to
	* ISO14443-4 is device is compliant.
	* In initiator mode, it means that NFC chip will send RATS automatically when
	* select and it will automatically poll for ISO14443-4 card when ISO14443A is
	* requested.
	* In target mode, with a NFC chip compliant (ie. PN532), the chip will
	* emulate a 14443-4 PICC using hardware capability */
	NP_AUTO_ISO14443_4
	/** Use automatic frames encapsulation and chaining. */
	NP_EASY_FRAMING
	/** Force the chip to switch in ISO14443-A */
	NP_FORCE_ISO14443_A
	/** Force the chip to switch in ISO14443-B */
	NP_FORCE_ISO14443_B
	/** Force the chip to run at 106 kbps */
	NP_FORCE_SPEED_106
)

type PropertyInfo struct {
	b    byte
	Name string
}

var PropertyInfos map[Property]PropertyInfo = map[Property]PropertyInfo{
	NP_TIMEOUT_COMMAND:        mkPropInfo(NP_TIMEOUT_COMMAND, "NP_TIMEOUT_COMMAND "),
	NP_TIMEOUT_ATR:            mkPropInfo(NP_TIMEOUT_ATR, "NP_TIMEOUT_ATR"),
	NP_TIMEOUT_COM:            mkPropInfo(NP_TIMEOUT_COM, "NP_TIMEOUT_COM"),
	NP_HANDLE_CRC:             mkPropInfo(NP_HANDLE_CRC, "NP_HANDLE_CRC"),
	NP_HANDLE_PARITY:          mkPropInfo(NP_HANDLE_PARITY, "NP_HANDLE_PARITY"),
	NP_ACTIVATE_FIELD:         mkPropInfo(NP_ACTIVATE_FIELD, "NP_ACTIVATE_FIELD"),
	NP_ACTIVATE_CRYPTO1:       mkPropInfo(NP_ACTIVATE_CRYPTO1, "NP_ACTIVATE_CRYPTO1"),
	NP_INFINITE_SELECT:        mkPropInfo(NP_INFINITE_SELECT, "NP_INFINITE_SELECT"),
	NP_ACCEPT_INVALID_FRAMES:  mkPropInfo(NP_ACCEPT_INVALID_FRAMES, "NP_ACCEPT_INVALID_FRAMES"),
	NP_ACCEPT_MULTIPLE_FRAMES: mkPropInfo(NP_ACCEPT_MULTIPLE_FRAMES, "NP_ACCEPT_MULTIPLE_FRAMES"),
	NP_AUTO_ISO14443_4:        mkPropInfo(NP_AUTO_ISO14443_4, "NP_AUTO_ISO14443_4"),
	NP_EASY_FRAMING:           mkPropInfo(NP_EASY_FRAMING, "NP_EASY_FRAMING"),
	NP_FORCE_ISO14443_A:       mkPropInfo(NP_FORCE_ISO14443_A, "NP_FORCE_ISO14443_A"),
	NP_FORCE_ISO14443_B:       mkPropInfo(NP_FORCE_ISO14443_B, "NP_FORCE_ISO14443_B"),
	NP_FORCE_SPEED_106:        mkPropInfo(NP_FORCE_SPEED_106, "NP_FORCE_SPEED_106"),
}

func mkPropInfo(p Property, name string) PropertyInfo {
	return PropertyInfo{
		b:    byte(p),
		Name: name,
	}
}

const ( //friendly names
	TimeoutCommand       = NP_TIMEOUT_COMMAND
	TimeoutATR           = NP_TIMEOUT_ATR
	TimeoutCom           = NP_TIMEOUT_COM
	HandleCRC            = NP_HANDLE_CRC
	HandleParity         = NP_HANDLE_PARITY
	ActivateField        = NP_ACTIVATE_FIELD
	ActivateCrypto1      = NP_ACTIVATE_CRYPTO1
	InfiniteSelect       = NP_INFINITE_SELECT
	AcceptInvalidFrames  = NP_ACCEPT_INVALID_FRAMES
	AcceptMultipleFrames = NP_ACCEPT_MULTIPLE_FRAMES
	AutoISO14443_4       = NP_AUTO_ISO14443_4
	EasyFraming          = NP_EASY_FRAMING
	ForceISO14443a       = NP_FORCE_ISO14443_A
	ForceISO14443b       = NP_FORCE_ISO14443_B
	ForceSpeed106        = NP_FORCE_SPEED_106
)
