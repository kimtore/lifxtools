package commands

import (
	"fmt"
	"net"
	"time"

	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const source = 0xbeef

// Change the label of a LIFX bulb.
var ztestCmd = &cobra.Command{
	Use:   "ztest",
	Short: "Run a test program on a LIFX Z strip",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true

		/*
			if len(viper.GetString("address")) == 0 {
				return fmt.Errorf("you must specify --address")
			}

		*/

		cmd.SilenceUsage = true

		//fulladdr := viper.GetString("address") + ":56700"
		fulladdr := "lifx66.iot.home.arpa:56700"
		log.Infof("Dialing UDP %s...", fulladdr)
		conn, err := net.Dial("udp", fulladdr)
		if err != nil {
			return err
		}
		defer conn.Close()

		const min = 0
		const max = 8 * 3

		send := func(start, end uint8, brightness uint16) {
			payload := &lifx.SetColorZonesMessage{
				StartIndex: start,
				EndIndex:   end,
				Color: lifx.HBSK{
					Hue:        62362,
					Saturation: 52000,
					Brightness: brightness,
					Kelvin:     0,
				},
				Duration: 0,
				Apply:    lifx.MultiZoneApply,
			}

			packet := lifx.NewPacket(payload)
			packet.Header.Frame.Source = source
			packet.Write(conn)
		}

		for cmd.Context().Err() == nil {
			var i uint8
			for i = min; i < max; i++ {
				if i > 0 {
					send(0, i-1, 0)
				}
				send(i, i, 52000)
				if i+1 < max {
					send(i+1, max, 0)
				}

				time.Sleep(200 * time.Millisecond)
				fmt.Print(".")
			}

			/*
				payload := &lifx.SetColorZonesMessage{
					StartIndex: min,
					EndIndex:   max,
					Apply:      lifx.MultiZoneApplyOnly,
				}

				packet := lifx.NewPacket(payload)
				packet.Header.Frame.Source = source
				packet.Write(conn)

			*/

			//	time.Sleep(1250 * time.Millisecond)
		}

		return cmd.Context().Err()
	},
}

func init() {
	rootCmd.AddCommand(ztestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ztestCmd.Flags().String("address", "", "Hostname of bulb to send the rename command to")

	viper.BindPFlags(setlabelCmd.Flags())
}
