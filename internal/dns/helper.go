package helper

import (
	"encoding/binary"
	"errors"
	"net/netip"
	"strings"
)

func ParseQuery(msg []byte) (name, qtype string, id uint16, err error) {
	// Headers not even set properly
	if len(msg) < 12 {
		return "", "", 0, errors.New("message too short")
	}

	id = binary.BigEndian.Uint16(msg[0:2])

	// skip 12 byte header
	offset := 12

	var labels []string
	for {
		if offset >= len(msg) {
			return "", "", 0, errors.New("question not found")
		}
		labelLength := int(msg[offset])
		offset += 1
		if labelLength == 0 {
			// Reached end of url
			break
		}

		if offset+labelLength >= len(msg) {
			return "", "", 0, errors.New("Lable out of bounds")
		}
		labels = append(labels, string(msg[offset:offset+labelLength]))
		offset += labelLength
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

func BuildResponse(id uint16, msg []byte, addrs []netip.Addr) []byte {
	// find end of question section
	offset := 12
	for msg[offset] != 0 {
		offset += int(msg[offset]) + 1
	}
	offset++    // skip null terminator
	offset += 4 // skip qtype and qclass
	questionSection := msg[12:offset]

	resp := make([]byte, 12)
	binary.BigEndian.PutUint16(resp[0:2], id)
	binary.BigEndian.PutUint16(resp[2:4], 0x8000)
	binary.BigEndian.PutUint16(resp[4:6], 1)                  // question count
	binary.BigEndian.PutUint16(resp[6:8], uint16(len(addrs))) // answer count

	resp = append(resp, questionSection...)

	for _, addr := range addrs {
		resp = append(resp, 0xC0, 0x0C)
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
