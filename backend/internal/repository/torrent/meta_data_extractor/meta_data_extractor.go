package meta_data_extractor

import (
	"encoding/hex"
	"path/filepath"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"

	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type MetadataExtractor struct{}

func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

func (e *MetadataExtractor) Extract(infoBytes []byte) ([]torrentModel.FileDTO, []torrentModel.SignatureDTO) {
	files := e.extractFiles(infoBytes)
	signatures := e.extractSignatures(infoBytes)
	return files, signatures
}

func (e *MetadataExtractor) extractFiles(infoBytes []byte) []torrentModel.FileDTO {
	var info metainfo.Info
	if err := bencode.Unmarshal(infoBytes, &info); err != nil {
		return nil
	}
	return e.buildFileList(&info)
}

func (e *MetadataExtractor) buildFileList(info *metainfo.Info) []torrentModel.FileDTO {
	upverted := info.UpvertedFiles()
	if len(upverted) == 0 {
		return e.buildSingleFileList(info)
	}
	return e.buildMultiFileList(info, upverted)
}

func (e *MetadataExtractor) buildSingleFileList(info *metainfo.Info) []torrentModel.FileDTO {
	return []torrentModel.FileDTO{
		{Path: info.BestName(), SizeBytes: info.Length},
	}
}

func (e *MetadataExtractor) buildMultiFileList(info *metainfo.Info, upverted []metainfo.FileInfo) []torrentModel.FileDTO {
	files := make([]torrentModel.FileDTO, len(upverted))
	rootDir := info.BestName()
	for i, f := range upverted {
		files[i] = e.buildFileDTO(rootDir, f)
	}
	return files
}

func (e *MetadataExtractor) buildFileDTO(rootDir string, f metainfo.FileInfo) torrentModel.FileDTO {
	return torrentModel.FileDTO{
		Path:      filepath.Join(append([]string{rootDir}, f.Path...)...),
		SizeBytes: f.Length,
	}
}

func (e *MetadataExtractor) extractSignatures(infoBytes []byte) []torrentModel.SignatureDTO {
	var rootDict map[string]interface{}
	if err := bencode.Unmarshal(infoBytes, &rootDict); err != nil {
		return nil
	}
	return e.parseSignaturesRoot(rootDict)
}

func (e *MetadataExtractor) parseSignaturesRoot(rootDict map[string]interface{}) []torrentModel.SignatureDTO {
	sigsInterface, ok := rootDict["signatures"].([]interface{})
	if !ok {
		return nil
	}
	return e.mapSignatures(sigsInterface)
}

func (e *MetadataExtractor) mapSignatures(sigsInterface []interface{}) []torrentModel.SignatureDTO {
	signatures := make([]torrentModel.SignatureDTO, 0, len(sigsInterface))
	for _, sigItem := range sigsInterface {
		e.appendSignatureIfValid(&signatures, sigItem)
	}
	return signatures
}

func (e *MetadataExtractor) appendSignatureIfValid(signatures *[]torrentModel.SignatureDTO, sigItem interface{}) {
	sigDict, ok := sigItem.(map[string]interface{})
	if !ok {
		return
	}
	dto := e.mapSignatureDict(sigDict)
	if dto != nil {
		*signatures = append(*signatures, *dto)
	}
}

func (e *MetadataExtractor) mapSignatureDict(sigDict map[string]interface{}) *torrentModel.SignatureDTO {
	pubKey, ok1 := sigDict["ed25519_pubkey"].([]byte)
	sigBytes, ok2 := sigDict["signature"].([]byte)
	timestamp, ok3 := sigDict["timestamp"].(int64)
	if !ok1 || !ok2 || !ok3 {
		return nil
	}
	return &torrentModel.SignatureDTO{
		SignerPublicKey: hex.EncodeToString(pubKey),
		SignatureBytes:  hex.EncodeToString(sigBytes),
		Timestamp:       timestamp,
	}
}
