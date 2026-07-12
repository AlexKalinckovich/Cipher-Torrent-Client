package validator

import (
	"net/url"
	"strings"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/validation"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type TorrentValidator struct{}

func NewTorrentValidator() *TorrentValidator {
	return &TorrentValidator{}
}

func (v *TorrentValidator) ValidateAddRequest(req torrentModel.AddTorrentRequest) error {
	agg := validation.NewAggregateError()
	v.checkSavePath(agg, req.SavePath)
	v.checkMagnetURI(agg, req.MagnetURI)
	return v.evaluateAggregate(agg)
}

func (v *TorrentValidator) evaluateAggregate(agg *validation.AggregateError) error {
	if agg.HasErrors() {
		return agg
	}
	return nil
}

func (v *TorrentValidator) checkSavePath(agg *validation.AggregateError, savePath string) {
	if strings.TrimSpace(savePath) == "" {
		agg.Add("save_path", savePath, "save_path is required")
	}
}

func (v *TorrentValidator) checkMagnetURI(agg *validation.AggregateError, uri string) {
	if uri == "" {
		return
	}
	v.validateMagnetFormat(agg, uri)
}

func (v *TorrentValidator) validateMagnetFormat(agg *validation.AggregateError, uri string) {
	if !strings.HasPrefix(uri, "magnet:?") {
		agg.Add("magnet_uri", uri, "invalid magnet URI prefix")
		return
	}
	v.parseURI(agg, uri)
}

func (v *TorrentValidator) parseURI(agg *validation.AggregateError, uri string) {
	_, err := url.Parse(uri)
	if err != nil {
		agg.Add("magnet_uri", uri, "invalid magnet URI format")
	}
}
