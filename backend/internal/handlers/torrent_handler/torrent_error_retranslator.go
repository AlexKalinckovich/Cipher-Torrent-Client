package torrent_handler

import (
	"errors"
	"net/http"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/torrent_errors"
)

type TorrentHandlerErrorRetranslator struct{}

func NewTorrentHandlerErrorRetranslator() *TorrentHandlerErrorRetranslator {
	return &TorrentHandlerErrorRetranslator{}
}

func (r *TorrentHandlerErrorRetranslator) Retranslate(err error) error {
	if r.isMissingFileError(err) {
		return torrent_errors.NewTorrentFileMissingError()
	}
	return err
}

func (r *TorrentHandlerErrorRetranslator) isMissingFileError(err error) bool {
	return errors.Is(err, http.ErrMissingFile)
}
