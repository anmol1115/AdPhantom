package helper

import (
	"encoding/binary"
	"errors"
	"net/netip"
	"strings"
)

func ParseQuery(msg []byte) (name, qtype string, id uint16, err error) {
	if len(msg) < 12 {
		return "", "", 0, errors.New("message too short")
	}

	id = binary.BigEndian.Uint16(msg[0:2])

	// skip 12 byte header
	offset := 12

	// walk labels
	var labels []string
	for {
		if offset >= len(msg) {
			return "", "", 0, errors.New("malformed question")
		}
		length := int(msg[offset])
		offset++
		if length == 0 {
			break
		}
		if offset+length > len(msg) {
			return "", "", 0, errors.New("label out of bounds")
		}
		labels = append(labels, string(msg[offset:offset+length]))
		offset += length
	}

	name = strings.Join(labels, ".")

	if offset+4 > len(msg) {
		return "", "", 0, errors.New("missing qtype/qclass")
	}

	qtypeInt := binary.BigEndian.Uint16(msg[offset : offset+2])
	switch qtypeInt {
	case 1:
		qtype = "A"
	case 28:
		qtype = "AAAA"
	default:
		qtype = "UNSUPPORTED"
	}

	return name, qtype, id, nil
}

func BuildNXDomain(id uint16) []byte {
	resp := make([]byte, 12)
	binary.BigEndian.PutUint16(resp[0:2], id)
	// flags: QR=1, RCODE=3 (NXDOMAIN)
	binary.BigEndian.PutUint16(resp[2:4], 0x8003)
	return resp
}

func BuildResponse(id uint16, addrs []netip.Addr) []byte {
	resp := make([]byte, 12)
	binary.BigEndian.PutUint16(resp[0:2], id)
	// flags: QR=1, AA=0, RCODE=0
	binary.BigEndian.PutUint16(resp[2:4], 0x8000)
	// answer count
	binary.BigEndian.PutUint16(resp[6:8], uint16(len(addrs)))

	for _, addr := range addrs {
		// name pointer back to offset 12 (0xC00C)
		resp = append(resp, 0xC0, 0x0C)

		if addr.Is4() {
			binary.BigEndian.AppendUint16(resp, 1)  // type A
			binary.BigEndian.AppendUint16(resp, 1)  // class IN
			binary.BigEndian.AppendUint32(resp, 60) // TTL
			binary.BigEndian.AppendUint16(resp, 4)  // rdlength
			ip := addr.As4()
			resp = append(resp, ip[:]...)
		} else {
			binary.BigEndian.AppendUint16(resp, 28) // type AAAA
			binary.BigEndian.AppendUint16(resp, 1)  // class IN
			binary.BigEndian.AppendUint32(resp, 60) // TTL
			binary.BigEndian.AppendUint16(resp, 16) // rdlength
			ip := addr.As16()
			resp = append(resp, ip[:]...)
		}
	}

	return resp
}
