package commands

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Change the label of a LIFX bulb.
var setlabelCmd = &cobra.Command{
	Use:   "setlabel",
	Short: "Change the name of a LIFX bulb",
	Long:  `Sets the friendly name. Will be visible in the mobile app, and perhaps the hostname too.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true

		if len(viper.GetString("label")) == 0 || len(viper.GetString("address")) == 0 {
			return fmt.Errorf("you must specify --label and --address")
		}

		cmd.SilenceUsage = true

		payload := &lifx.SetLabelMessage{
			Label: viper.GetString("label"),
		}

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
		log.Infof("Dialing UDP %s...", fulladdr)
		conn, err := net.Dial("udp", fulladdr)
		if err != nil {
			return err
		}
		defer conn.Close()

		log.Infof("Sending label change command...")
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			return err
		}

		r := bufio.NewReader(conn)

		log.Infof("Decoding response...")
		resp, err := lifx.DecodePacket(r)
		if err != nil {
			return err
		}

		state, ok := resp.Payload.(*lifx.SetLabelMessage)
		if !ok {
			return fmt.Errorf("response does not contain bulb state")
		}

		log.Infof("Bulb label set to '%s'", state.Label)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setlabelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setlabelCmd.Flags().String("label", "", "New label for LIFX bulb")
	setlabelCmd.Flags().String("address", "", "Hostname of bulb to send the rename command to")
	setlabelCmd.Flags().String("mac", "00:00:00:00:00:00", "MAC address of bulb")
	setlabelCmd.Flags().String("source", "beef", "16-bit source address")

	viper.BindPFlags(setlabelCmd.Flags())
}
