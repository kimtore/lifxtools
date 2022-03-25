package commands

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Change the label of a LIFX bulb.
var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Onboard the LIFX bulb on your Wifi network",
	Long: `Attempt to send an onboarding message with Wifi credentials.

Use this command only when connected to a LIFX bulb Wifi network.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		cmd.SilenceErrors = true

		ssid := viper.GetString("ssid")
		psk := viper.GetString("psk")
		address := viper.GetString("address")

		if len(ssid) == 0 || len(psk) == 0 {
			return fmt.Errorf("you must specify --ssid and --psk")
		}

		cmd.SilenceUsage = true

		payload := &lifx.SetAccessPointMessage{
			Reserved1: 0x02,
			SSID:      ssid,
			PSK:       psk,
			Security:  lifx.Security_WPA2_AES_PSK,
		}
		packet := lifx.NewPacket(payload)

		buf := &bytes.Buffer{}
		err := packet.Write(buf)
		if err != nil {
			return err
		}

		log.Infof("Connecting to TCP %s...", address)
		sock, err := tls.DialWithDialer(
			&net.Dialer{
				Timeout: time.Second * 5,
			},
			"tcp", address,
			&tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS10,
			},
		)
		if err != nil {
			return err
		}
		defer sock.Close()

		log.Infof("Sending onboarding payload...")
		_, err = sock.Write(buf.Bytes())
		if err != nil {
			return err
		}

		log.Infof("Done, happy hacking.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	onboardCmd.Flags().String("address", "172.16.0.1:56700", "Network address of LIFX bulb")
	onboardCmd.Flags().String("ssid", "", "Wifi SSID")
	onboardCmd.Flags().String("psk", "", "Wifi pre-shared key")
}
