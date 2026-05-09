package helper

import (
	"encoding/binary"
	"errors"
	"net/netip"
	"strings"
)

func ParseQuery(msg []byte) (name, qtype string, err error) {
	// Headers not even set properly
	if len(msg) < 12 {
		return "", "", errors.New("message too short")
	}

	// skip 12 byte header
	offset := 12

	var labels []string
	for {
		if offset >= len(msg) {
			return "", "", errors.New("question not found")
		}
		labelLength := int(msg[offset])
		offset += 1
		if labelLength == 0 {
			// Reached end of url
			break
		}

		if offset+labelLength >= len(msg) {
			return "", "", errors.New("Lable out of bounds")
		}
		labels = append(labels, string(msg[offset:offset+labelLength]))
		offset += labelLength
	}

	name = strings.Join(labels, ".")

	if offset+4 > len(msg) {
		return "", "", errors.New("missing qtype/qclass")
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

	return name, qtype, nil
}

func BuildNXDomain(msg []byte) []byte {
	resp := make([]byte, 12)
	id := msg[:2]
	resp = append(resp, id...)
	// flags: QR=1, RCODE=3 (NXDOMAIN)
	binary.BigEndian.PutUint16(resp[2:4], 0x8003)
	return resp
}

func BuildResponse(msg []byte, addrs []netip.Addr) []byte {
	var resp []byte

	resp = append(resp, msg...)
	// set QR bit to 1 in header flags
	resp[2] = resp[2] | 0x80
	// set ANCOUNT to count of answers
	binary.BigEndian.PutUint16(resp[6:8], uint16(len(addrs)))

	for _, addr := range addrs {
		resp = append(resp, 0xc0, 0x0c)

		if addr.Is4() {
			resp = binary.BigEndian.AppendUint16(resp, 1)
			resp = binary.BigEndian.AppendUint16(resp, 1)
			resp = binary.BigEndian.AppendUint32(resp, 60)
			resp = binary.BigEndian.AppendUint16(resp, 4)
			ip := addr.As4()
			resp = append(resp, ip[:]...)
		} else {
			resp = binary.BigEndian.AppendUint16(resp, 28)
			resp = binary.BigEndian.AppendUint16(resp, 1)
			resp = binary.BigEndian.AppendUint32(resp, 60)
			resp = binary.BigEndian.AppendUint16(resp, 16)
			ip := addr.As16()
			resp = append(resp, ip[:]...)
		}
	}

	return resp
}
