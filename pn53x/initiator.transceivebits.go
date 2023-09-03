package pn53x

import (
	"errors"

	"github.com/maitredede/gonfc"
)

// InitiatorTransceiveBits
// chips/pn53x.c pn53x_initiator_transceive_bits
func (pnd *Chip) InitiatorTransceiveBits(pbtTx []byte, szTxBits int, pbtTxPar []byte, pbtRx []byte, pbtRxPar []byte) (int, error) {

	var szFrameBits int

	abtCmd := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	abtCmd[0] = byte(InCommunicateThru)

	// Check if we should prepare the parity bits ourself
	if (!pnd.bPar.Get()) && (szTxBits > 0) {
		// Convert data with parity to a frame
		frame, err := WrapFrame(pbtTx, szTxBits, pbtTxPar)
		if err != nil {
			return 0, err
		}
		for i := 0; i < min(len(frame), len(abtCmd)-1); i++ {
			abtCmd[i+1] = frame[i]
		}
		szFrameBits = len(frame)
	} else {
		szFrameBits = szTxBits
	}

	// Retrieve the leading bits
	ui8Bits := byte(szFrameBits) % 8

	// Get the amount of frame bytes + optional (1 byte if there are leading bits)
	var lead int
	if ui8Bits == 0 {
		lead = 0
	} else {
		lead = 1
	}
	szFrameBytes := (szFrameBits / 8) + lead

	// When the parity is handled before us, we just copy the data
	if pnd.bPar.Get() {
		for i := 0; i < szFrameBytes; i++ {
			abtCmd[i+1] = pbtTx[i]
		}
	}

	// Set the amount of transmission bits in the PN53X chip register
	if err := pnd.setTxBits(ui8Bits); err != nil {
		return 0, err
	}

	// Send the frame to the PN53X chip and get the answer
	// We have to give the amount of bytes + (the command byte 0x42)
	abtRx := make([]byte, PN53x_EXTENDED_FRAME__DATA_MAX_LEN)
	szRx, err := pnd.transceive(abtCmd[:szFrameBytes+1], abtRx, -1)
	if err != nil {
		return 0, err
	}
	// Get the last bit-count that is stored in the received byte
	ui8rcc, err := pnd.readRegister(PN53X_REG_CIU_Control)
	if err != nil {
		return 0, err
	}
	ui8Bits = ui8rcc & SYMBOL_RX_LAST_BITS

	// Recover the real frame length in bits
	var tern int
	if ui8Bits == 0 {
		tern = 0
	} else {
		tern = 1
	}
	szFrameBits = ((szRx - 1 - tern) * 8) + int(ui8Bits)

	var szRxBits int
	if len(pbtRx) > 0 {
		// Ignore the status byte from the PN53X here, it was checked earlier in pn53x_transceive()
		// Check if we should recover the parity bits ourself
		if !pnd.bPar.Get() {
			// Unwrap the response frame
			frame, framePar, err := UnwrapFrame(abtRx[1:], szFrameBits)
			if err != nil {
				return 0, err
			}
			szRxBits = len(frame)
			for i := 0; i < szRxBits; i++ {
				pbtRx[i] = frame[i]
				pbtRxPar[i] = framePar[i]
			}
			// if ((res = pn53x_unwrap_frame(abtRx + 1, szFrameBits, pbtRx, pbtRxPar)) < 0)
			//   return res;
		} else {
			// Save the received bits
			szRxBits = szFrameBits
			// Copy the received bytes
			for i := 0; i < szRx-1; i++ {
				pbtRx[i] = abtRx[i+1]
			}
		}
	}
	// Everything went successful
	return szRxBits, nil
}

// WrapFrame
// chips/pn53x.c pn53x_wrap_frame
func WrapFrame(pbtTx []byte, szTxBits int, pbtTxPar []byte) ([]byte, error) {
	szBitsLeft := szTxBits

	// Make sure we should frame at least something
	if szBitsLeft == 0 {
		return nil, gonfc.NFC_ECHIP
	}

	panic(errors.New("TODO : pn53x.WrapFrame"))
	// var szFrameBits int
	// // Handle a short response (1byte) as a special case
	// if szBitsLeft < 9 {
	// 	frame := pbtTx[:szTxBits]
	// 	return frame, nil
	// }

	// // We start by calculating the frame length in bits
	// szFrameBits = szTxBits + (szTxBits / 8)

	// // Parse the data bytes and add the parity bits
	// // This is really a sensitive process, mirror the frame bytes and append parity bits
	// // buffer = mirror(frame-byte) + parity + mirror(frame-byte) + parity + ...
	// // split "buffer" up in segments of 8 bits again and mirror them
	// // air-bytes = mirror(buffer-byte) + mirror(buffer-byte) + mirror(buffer-byte) + ..
	// for {
	// 	// Reset the temporary frame byte;
	// 	var btFrame byte = 0
	// 	for uiBitPos := byte(0); uiBitPos < 8; uiBitPos++ {

	// 	}
	// 	// Every 8 data bytes we lose one frame byte to the parities
	// 	pbtFrame++
	// }
}

// WrapFrame
// chips/pn53x.c pn53x_unwrap_frame
func UnwrapFrame(pbtFrame []byte, szFrameBits int) ([]byte, []byte, error) {
	var pbtRx []byte
	var pbtRxPar []byte
	szBitsLeft := szFrameBits
	// Make sure we should frame at least something
	if szBitsLeft == 0 {
		return pbtRx, pbtRxPar, gonfc.NFC_ECHIP
	}

	panic(errors.New("TODO : pn53x.UnwrapFrame"))
}
