package acr122usb

import (
	"fmt"
	"strings"

	"github.com/google/gousb"
)

func (pnd *Acr122UsbDevice) usbGetEndPoints() error {
	for adr, ep := range pnd.id.uif.Endpoints {
		if ep.TransferType != gousb.TransferTypeBulk {
			continue
		}
		if ep.Direction == gousb.EndpointDirectionIn {
			pnd.epInAdr = adr
			iep, err := pnd.intf.InEndpoint(int(ep.Address))
			if err != nil {
				return err
			}
			pnd.epIn = iep
			pnd.maxPacketSize = ep.MaxPacketSize
			ies, err := pnd.epIn.NewStream(pnd.maxPacketSize, 1)
			if err != nil {
				return err
			}
			pnd.epInStream = ies
		}
		if ep.Direction == gousb.EndpointDirectionOut {
			pnd.epOutAdr = adr
			oep, err := pnd.intf.OutEndpoint(int(ep.Address))
			if err != nil {
				return err
			}
			pnd.epOut = oep
			pnd.maxPacketSize = ep.MaxPacketSize
			oes, err := pnd.epOut.NewStream(pnd.maxPacketSize, 1)
			if err != nil {
				return err
			}
			pnd.epOutStream = oes
		}
	}
	return nil
}

func (pnd *Acr122UsbDevice) usbBulkRead(data []byte, timeout int) (int, error) {
	// n, err := pnd.epIn.Read(data)
	// if err != nil {
	// 	return 0, fmt.Errorf("usbBulkRead: %w", err)
	// }
	// return n, nil
	n, err := pnd.epInStream.Read(data)
	if err != nil {
		pnd.logger.Named("rx").Errorf("  usbBulkRead: %v", err)
	} else {
		read := data[:n]
		pnd.logger.Named("rx").Debugf("  usbBulkRead: n=%v d=%v", n, toHexString(read))
	}
	return n, err
}

func toHexString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for _, b := range data {
		sb.WriteString(fmt.Sprintf(" %02x", b))
	}
	return sb.String()[1:]
}

func (pnd *Acr122UsbDevice) usbBulkWrite(data []byte, timeout int) (int, error) {
	// n, err := pnd.epOut.Write(data)
	// if err != nil {
	// 	return 0, err
	// }
	// return n, nil
	pnd.logger.Named("tx").Debugf("  usbBulkWrite: n=%v d=%v", len(data), toHexString(data))
	n, err := pnd.epOutStream.Write(data)
	if err != nil {
		pnd.logger.Named("tx").Errorf("  usbBulkWrite: %v", err)
	} else {
		if n != len(data) {
			pnd.logger.Named("tx").Warnf("  usbBulkWrite: sent %v", n)
		}
	}
	return n, err
}
