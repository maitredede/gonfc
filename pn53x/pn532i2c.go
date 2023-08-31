package pn53x

import (
	"time"
)

const (
	PN532_BUFFER_LEN int = PN53x_EXTENDED_FRAME__DATA_MAX_LEN + PN53x_EXTENDED_FRAME__OVERHEAD

	/*
	 * When sending lots of data, the pn532 occasionally fails to respond in time.
	 * Since it happens so rarely, lets try to fix it by re-sending the data. This
	 * define allows for fine tuning the number of retries.
	 */
	PN532_SEND_RETRIES int = 3
	/*
	 * Bus free time (in ms) between a STOP condition and START condition. See
	 * tBuf in the PN532 data sheet, section 12.25: Timing for the I2C interface,
	 * table 320. I2C timing specification, page 211, rev. 3.2 - 2007-12-07.
	 */
	PN532_BUS_FREE_TIME time.Duration = 5 * time.Millisecond
)

var (
	/* preamble and start bytes, see pn532-internal.h for details */
	PN53X_PREAMBLE_AND_START     []byte = []byte{0x00, 0x00, 0xff}
	PN53X_PREAMBLE_AND_START_LEN int    = len(PN53X_PREAMBLE_AND_START)
)
