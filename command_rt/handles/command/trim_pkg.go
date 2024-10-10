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
	check := 0
	for i := 0; i < n-1; i++ {
		check += int(commandMass[i])
	}
	check = check % 256
	if check != int(commandMass[n-1]) {
		return "", fmt.Errorf(ErrorChecksum)
	}
	switch commandMass[2] {
	case 7:
		switch commandMass[3] {
		case 1:
			return StartFlow, nil
		case 2:
			return StopFlow, nil
		default:
			return "", fmt.Errorf(ErrorCommand)

		}
	// case 32:
	// 	return StopFlow, nil
	default:
		return "", fmt.Errorf(ErrorCommand)
	}

}
