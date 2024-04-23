package commands

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/dorkowscy/lyslix/lifx/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Change the label of a LIFX bulb.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of a LIFX bulb",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		payload := &lifx.GetVersionMessage{}

		macaddr, err := net.ParseMAC(viper.GetString("mac"))
		if err != nil {
			return err
		}

		sourceaddr, err := hex.DecodeString(viper.GetString("source"))
		if err != nil {
			return err
		}
		packet := lifx.NewPacket(payload)
		packet.Header.FrameAddress.Target = lifx.MACAdressToFrameAddress(macaddr)
		packet.Header.Frame.Source = (uint32(sourceaddr[0]) << 8) | uint32(sourceaddr[1])
		packet.Header.FrameAddress.ResRequired = true

		buf := &bytes.Buffer{}
		err = packet.Write(buf)
		if err != nil {
			return err
		}

		fulladdr := viper.GetString("address") + ":56700"
		log.Debugf("Dialing UDP %s...", fulladdr)
		conn, err := net.Dial("udp", fulladdr)
		if err != nil {
			return err
		}
		defer conn.Close()

		log.Debugf("Sending version info request...")
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			return err
		}

		r := bufio.NewReader(conn)

		log.Debugf("Decoding response...")
		resp, err := lifx.DecodePacket(r)
		if err != nil {
			return err
		}

		state, ok := resp.Payload.(*lifx.StateVersionMessage)
		if !ok {
			return fmt.Errorf("response does not contain version info")
		}

		product, err := version.Lookup(int(state.Vendor), int(state.Product))

		log.Debugf("vendor=%v product=%v reserved6=%v", state.Vendor, state.Product, state.Reserved6)

		if err != nil {
			log.Warn(err)
		} else {
			log.Infof("Bulb reports itself as '%s'", product.Name)
			log.Infof("Features: %+v", product.Features)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	versionCmd.Flags().String("address", "255.255.255.255", "Hostname of bulb to send the rename command to")
	versionCmd.Flags().String("mac", "00:00:00:00:00:00", "MAC address of bulb")
	versionCmd.Flags().String("source", "beef", "16-bit source address")

	viper.BindPFlags(versionCmd.Flags())
}
