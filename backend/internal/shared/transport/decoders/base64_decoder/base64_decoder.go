package base64_decoder

import "encoding/base64"

func DecodeBase64Param(param string) ([]byte, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(param)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(param)
		if err != nil {
			return nil, err
		}
	}
	return decoded, nil
}
