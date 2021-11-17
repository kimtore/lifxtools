package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"net"

	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
)

var (
	label  = flag.String("label", "", "New label for LIFX bulb")
	addr   = flag.String("addr", "", "Hostname of bulb to send the rename command to")
	mac    = flag.String("mac", "00:00:00:00:00:00", "MAC address of bulb")
	source = flag.String("source", "beef", "16-bit source address")
)

// Change the label of a LIFX bulb.
func main() {
	flag.Parse()

	payload := &lifx.SetLabelMessage{
		Label: *label,
	}

	macaddr, err := net.ParseMAC(*mac)
	if err != nil {
		panic(err)
	}

	sourceaddr, err := hex.DecodeString(*source)
	if err != nil {
		panic(err)
	}
	packet := lifx.NewPacket(payload)
	packet.Header.FrameAddress.Target = lifx.MACAdressToFrameAddress(macaddr)
	packet.Header.Frame.Source = (uint32(sourceaddr[0]) << 8) | uint32(sourceaddr[1])
	packet.Header.FrameAddress.ResRequired = true

	buf := &bytes.Buffer{}
	err = packet.Write(buf)
	if err != nil {
		panic(err)
	}

	fulladdr := *addr+":56700"
	log.Infof("Dialing UDP %s...",fulladdr)
	conn, err := net.Dial("udp", fulladdr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Infof("Sending label change command...")
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(conn)

	log.Infof("Decoding response...")
	resp, err := lifx.DecodePacket(r)
	if err != nil {
		panic(err)
	}

	state, ok := resp.Payload.(*lifx.SetLabelMessage)
	if !ok {
		panic("response does not contain bulb state")
	}

	log.Infof("Bulb label set to '%s'", state.Label)
}
