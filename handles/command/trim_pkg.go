package command

import "fmt"

const (
	StartFlow     = "START_FLOW"
	StopFlow      = "STOP_FLOW"
	StartOnce     = "START_ONCE"
	ErrorChecksum = "error with checksum"
	ErrorPacket   = "error with packet"
	ErrorCommand  = "error with unknown command"
)

func CommandTrim(commandMass []byte) (string, error) {
	n := len(commandMass)
	if n != 8 {
		return "", fmt.Errorf(ErrorPacket)
	}
	if (n + int(commandMass[n-1]) + 85) != 256 {
		return "", fmt.Errorf(ErrorChecksum)
	}
	switch commandMass[3] {
	case 16:
		return StartFlow, nil
	case 1:
		return StartOnce, nil
	case 32:
		return StopFlow, nil
	default:
		return "", fmt.Errorf(ErrorCommand)

	}

}
