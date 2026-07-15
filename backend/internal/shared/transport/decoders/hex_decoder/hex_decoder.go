package hex_decoder

import "encoding/hex"

func DecodeHexParam(param string) ([]byte, error) {
	decoded, err := hex.DecodeString(param)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
