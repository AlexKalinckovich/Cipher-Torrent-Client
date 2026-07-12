package meta_data_extractor

import (
	"bytes"
	"encoding/hex"
	"log"
	"path/filepath"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"

	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type MetadataExtractor struct{}

func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

func (e *MetadataExtractor) Extract(data []byte) ([]torrentModel.FileDTO, []torrentModel.SignatureDTO) {
	files := e.extractFiles(data)
	signatures := e.extractSignatures(data)
	return files, signatures
}

func (e *MetadataExtractor) extractFiles(data []byte) []torrentModel.FileDTO {
	reader := bytes.NewReader(data)
	mi, err := metainfo.Load(reader)
	if err != nil {
		return nil
	}
	return e.unmarshalAndBuild(mi)
}

func (e *MetadataExtractor) unmarshalAndBuild(mi *metainfo.MetaInfo) []torrentModel.FileDTO {
	info, err := mi.UnmarshalInfo()
	if err != nil {
		log.Println("Error unmarshalling metainfo:", err)
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
		{
			Path:      info.BestName(),
			SizeBytes: info.Length,
		},
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

func (e *MetadataExtractor) extractSignatures(data []byte) []torrentModel.SignatureDTO {
	var rootDict map[string]interface{}
	if err := bencode.Unmarshal(data, &rootDict); err != nil {
		log.Println("Error unmarshalling root dict:", err)
		return nil
	}
	return e.parseSignaturesRoot(rootDict)
}

func (e *MetadataExtractor) parseSignaturesRoot(rootDict map[string]interface{}) []torrentModel.SignatureDTO {
	logDict(rootDict)
	sigsInterface, ok := rootDict["signatures"].([]interface{})
	if !ok {
		return nil
	}
	return e.mapSignatures(sigsInterface)
}

func logDict(rootDict map[string]interface{}) {
	for k, _ := range rootDict {
		log.Println(k)
	}
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
