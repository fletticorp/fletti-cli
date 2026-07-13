package command

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Device commands let you BE a drone, rover or kiosk from the terminal: they
// speak the device half of the protocol (ID-only login, heartbeat, command
// poll, ack) exactly as the on-board agents do. That is what makes it possible
// to drive and test the autonomous flows without powering any hardware.

func init() {
	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(deviceLoginCmd)
	deviceCmd.AddCommand(deviceAliveCmd)
	deviceCmd.AddCommand(deviceCommandsCmd)
	deviceCmd.AddCommand(deviceAckCmd)
	deviceCmd.AddCommand(devicePingCmd)
	deviceCmd.AddCommand(deviceAbortCmd)

	deviceAliveCmd.Flags().StringVar(&deviceStatus, "status", "standby", "device status to report (standby, online, off...)")
	deviceAliveCmd.Flags().Float64Var(&deviceLat, "lat", -34.7787778, "latitude to report")
	deviceAliveCmd.Flags().Float64Var(&deviceLng, "lng", -55.9625556, "longitude to report")
	deviceAckCmd.Flags().StringVar(&ackStatus, "status", "finished", "running, finished, failed (operations use completed)")
	deviceAckCmd.Flags().StringVar(&ackMessage, "message", "", "message to return with the ack")
}

var (
	deviceStatus string
	deviceLat    float64
	deviceLng    float64
	ackStatus    string
	ackMessage   string
)

var deviceCmd = &cobra.Command{
	Use:           "device",
	Short:         "Act as a drone, rover or operation (kiosk)",
	Long:          "Speaks the device half of the protocol: login, heartbeat, command poll and ack.\n\nType is one of: drones, rovers, operations.",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var deviceLoginCmd = &cobra.Command{
	Use:           "login <drones|rovers|operations> <deviceID>",
	Short:         "Get a device token (ID-only login, as the agents do)",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          deviceLogin,
}

var deviceAliveCmd = &cobra.Command{
	Use:           "alive <drones|rovers|operations> <deviceID>",
	Short:         "Send a heartbeat (a device must be standby to receive offers)",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          deviceAlive,
}

var deviceCommandsCmd = &cobra.Command{
	Use:           "commands <drones|rovers|operations> <deviceID>",
	Short:         "Poll the device's pending commands",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          deviceCommands,
}

var deviceAckCmd = &cobra.Command{
	Use:           "ack <drones|rovers|operations> <deviceID> <commandID>",
	Short:         "Mark a command running/finished/failed",
	Args:          cobra.ExactArgs(3),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          deviceAck,
}

var devicePingCmd = &cobra.Command{
	Use:               "ping <drones|rovers|operations> <deviceID>",
	Short:             "Queue a ping command for the device (proves a real agent is polling)",
	Args:              cobra.ExactArgs(2),
	SilenceErrors:     true,
	SilenceUsage:      true,
	PersistentPreRunE: ensureAuth,
	RunE:              devicePing,
}

var deviceAbortCmd = &cobra.Command{
	Use:               "abort <drones|rovers> <deviceID>",
	Short:             "Queue an abort command: the vehicle goes home",
	Args:              cobra.ExactArgs(2),
	SilenceErrors:     true,
	SilenceUsage:      true,
	PersistentPreRunE: ensureAuth,
	RunE:              deviceAbort,
}

func deviceType(kind string) (string, error) {
	switch kind {
	case "drones", "rovers", "operations":
		return kind, nil
	case "drone":
		return "drones", nil
	case "rover":
		return "rovers", nil
	case "operation", "kiosk":
		return "operations", nil
	}
	return "", fmt.Errorf("unknown device type: %s (use drones, rovers or operations)", kind)
}

// deviceToken logs the device in. Device tokens are not the user's token: the
// device authenticates with its ID alone, exactly as the agents do.
func deviceToken(kind, deviceID string) (string, error) {
	body, err := POST(fmt.Sprintf("%s/%s/%s/login", getUri(), kind, deviceID), nil, "device login")
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return "", err
	}

	token, ok := response["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access_token in login response: %s", body)
	}

	return token, nil
}

func deviceLogin(cmd *cobra.Command, args []string) error {
	kind, err := deviceType(args[0])
	if err != nil {
		return err
	}

	token, err := deviceToken(kind, args[1])
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", token)
	return nil
}

func deviceAlive(cmd *cobra.Command, args []string) error {
	kind, err := deviceType(args[0])
	if err != nil {
		return err
	}

	token, err := deviceToken(kind, args[1])
	if err != nil {
		return err
	}

	heartbeat := map[string]interface{}{
		"status":  deviceStatus,
		"message": "flesh",
		"location": map[string]interface{}{
			"latitude":  deviceLat,
			"longitude": deviceLng,
		},
	}

	body, err := POST(fmt.Sprintf("%s/%s/%s/alive?authorization=%s", getUri(), kind, args[1], token), heartbeat, "heartbeat")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func deviceCommands(cmd *cobra.Command, args []string) error {
	kind, err := deviceType(args[0])
	if err != nil {
		return err
	}

	token, err := deviceToken(kind, args[1])
	if err != nil {
		return err
	}

	body, err := GET(fmt.Sprintf("%s/%s/%s/commands?authorization=%s", getUri(), kind, args[1], token), "commands")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			fmt.Println("No pending commands.")
			return nil
		}
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func deviceAck(cmd *cobra.Command, args []string) error {
	kind, err := deviceType(args[0])
	if err != nil {
		return err
	}

	token, err := deviceToken(kind, args[1])
	if err != nil {
		return err
	}

	ack := map[string]interface{}{"status": ackStatus}
	if ackMessage != "" {
		ack["response"] = map[string]interface{}{"message": ackMessage}
	}

	body, err := PUT(fmt.Sprintf("%s/%s/%s/commands/%s?authorization=%s", getUri(), kind, args[1], args[2], token), ack, "command ack")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func devicePing(cmd *cobra.Command, args []string) error {
	kind, err := deviceType(args[0])
	if err != nil {
		return err
	}

	body, err := POST(fmt.Sprintf("%s/%s/%s/commands/ping?authorization=%s", getUri(), kind, args[1], getToken()), nil, "ping")
	if err != nil {
		return err
	}

	fmt.Printf("Ping queued at %s. The device acks it only if its agent is running.\n%s\n", time.Now().Format("15:04:05"), body)
	return nil
}

func deviceAbort(cmd *cobra.Command, args []string) error {
	kind, err := deviceType(args[0])
	if err != nil {
		return err
	}

	if kind == "operations" {
		return fmt.Errorf("operations cannot be aborted: abort is for drones and rovers")
	}

	body, err := POST(fmt.Sprintf("%s/%s/%s/commands/abort?authorization=%s", getUri(), kind, args[1], getToken()), nil, "abort")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
