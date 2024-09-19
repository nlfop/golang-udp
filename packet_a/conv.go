package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
)

type ReciveStructure struct {
	Prefix      byte                    `json:"prefix"`
	Size        int                     `json:"size"`
	Command     byte                    `json:"command"`
	ErrorCode   byte                    `json:"error_code"`
	info1       byte                    `json:"info_1"`
	info2       byte                    `json:"info_2"`
	ControlSumm byte                    `json:"control_summ"`
	Data        []ReceiveStructureBlock `json:"data"`
}

type ReceiveStructureBlock struct {
	Size                byte   `json:"size"`
	ErrorBlock          byte   `json:"error_block"`
	TypeStructure       byte   `json:"type_str"`
	NumberPlace         byte   `json:"num_block"`
	Info2               byte   `json:"info_2"`
	Info3               byte   `json:"info_3"`
	ControlSummSructure byte   `json:"control_str"`
	ControlSumm         byte   `json:"control_summ"`
	DataBlock           string `json:"data_block"`
}

func main() {
	file, err := os.Open("Crate.bin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	fileTXT, err := os.Create("convert.txt")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer fileTXT.Close()

	numBytes := 8
	for {
		data := make([]byte, numBytes)
		n, err := file.Read(data)
		if err != nil || n == 0 { // если конец файла
			fmt.Println(err)
			break // выходим из цикла
		}

		numBytes = EncodingPackage(data)

		// if numBytes != 8 {
		// 	fmt.Println(numBytes)
		// }
		// fileTXT.WriteString(fmt.Sprintf("%x\n", data))

	}

	// for key, val := range answers {
	// 	if key == 0 {
	// 		continue
	// 	}
	// 	fileTXT.WriteString(fmt.Sprintf("Новый пакет номер %d\n", key))

	// 	// fmt.Println("Новый пакет номер ", key)
	// 	for _, v := range val.Data {
	// 		fileTXT.WriteString(fmt.Sprintf("Блок номер: %d, данные:%x\n", int(v.NumberPlace), v.DataBlock))
	// 	}

	// }
	answers = answers[1:]
	dataAns, _ := json.Marshal(answers)
	fileTXT.WriteString(string(dataAns))
}

var answers []ReciveStructure
var answer = &ReciveStructure{}
var block = &ReceiveStructureBlock{}

func EncodingPackage(d []byte) int {
	switch {
	case len(d) != 8:
		block.DataBlock = fmt.Sprintf("%x", d)
		answer.Data = append(answer.Data, *block)
		return 8
	case d[0] == 83:
		answers = append(answers, *answer)
		sl := []byte{d[2], d[1]}
		dataSize := binary.BigEndian.Uint16(sl)
		answer = &ReciveStructure{
			Prefix:      d[0],
			Size:        int(dataSize),
			Command:     d[3],
			ErrorCode:   d[4],
			info1:       d[5],
			info2:       d[6],
			ControlSumm: d[7],
		}
		return 8

	default:
		block = &ReceiveStructureBlock{
			Size:                d[0],
			ErrorBlock:          d[1],
			TypeStructure:       d[2],
			NumberPlace:         d[3],
			Info2:               d[4],
			Info3:               d[5],
			ControlSummSructure: d[6],
			ControlSumm:         d[7],
		}

		if d[0] == 0 {
			answer.Data = append(answer.Data, *block)
			return 8
		}

		return int(d[0])
	}

}
