package torrent_signing

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/meta_data_extractor"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent_signature/torrent_signature_repository_errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"

	torrentSignerPorts "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent_signature/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/ports"
	userService "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	torrentSignatureModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent_signature"
)

type decryptedKeys struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type Service struct {
	userRepo         userService.UserRepositoryPort
	signatureRepo    torrentSignerPorts.TorrentSignatureRepositoryPort
	cryptoService    security.CryptoServicePort
	signatureBuilder ports.SignatureBuilder
	metaInfoSigner   ports.MetaInfoSigner
	extractor        *meta_data_extractor.MetadataExtractor
}

func NewService(
	userRepo userService.UserRepositoryPort,
	signatureRepo torrentSignerPorts.TorrentSignatureRepositoryPort,
	cryptoService security.CryptoServicePort,
	signatureBuilder ports.SignatureBuilder,
	metaInfoSigner ports.MetaInfoSigner,
) *Service {
	return &Service{
		userRepo:         userRepo,
		signatureRepo:    signatureRepo,
		cryptoService:    cryptoService,
		signatureBuilder: signatureBuilder,
		metaInfoSigner:   metaInfoSigner,
		extractor:        meta_data_extractor.NewMetadataExtractor(),
	}
}

func (s *Service) SignTorrent(ctx context.Context, file multipart.File, userID int64) (torrentModel.TorrentDTO, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.processSigning(ctx, fileBytes, userID)
}

func (s *Service) processSigning(ctx context.Context, fileBytes []byte, userID int64) (torrentModel.TorrentDTO, error) {
	rootDict, mi, err := s.parseTorrentBytes(fileBytes)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.signAndPersist(ctx, rootDict, mi, userID)
}

func (s *Service) parseTorrentBytes(fileBytes []byte) (map[string]interface{}, *metainfo.MetaInfo, error) {
	var rootDict map[string]interface{}
	if err := bencode.Unmarshal(fileBytes, &rootDict); err != nil {
		return nil, nil, err
	}
	mi, err := metainfo.Load(bytes.NewReader(fileBytes))
	return rootDict, mi, err
}

func (s *Service) signAndPersist(ctx context.Context, rootDict map[string]interface{}, mi *metainfo.MetaInfo, userID int64) (torrentModel.TorrentDTO, error) {
	infoHash := mi.HashInfoBytes().Bytes()
	if err := s.ensureNotSigned(ctx, userID, infoHash); err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.buildSignInjectAndSave(ctx, rootDict, mi, infoHash, userID)
}

func (s *Service) ensureNotSigned(ctx context.Context, userID int64, infoHash []byte) error {
	_, err := s.signatureRepo.GetByUserAndTorrentHash(ctx, userID, infoHash)
	if err == nil {
		return repositoryErrors.NewDuplicateSignatureError()
	}
	return s.translateNotFoundError(err)
}

func (s *Service) translateNotFoundError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func (s *Service) buildSignInjectAndSave(ctx context.Context, rootDict map[string]interface{}, mi *metainfo.MetaInfo, infoHash []byte, userID int64) (torrentModel.TorrentDTO, error) {
	keys, err := s.getDecryptedKeyPair(ctx, userID)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.signInjectAndPersist(ctx, rootDict, mi, infoHash, userID, keys)
}

func (s *Service) getDecryptedKeyPair(ctx context.Context, userID int64) (*decryptedKeys, error) {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.decryptAndExtractKeys(userEntity.PrivateKeyEnc)
}

func (s *Service) decryptAndExtractKeys(ciphertext []byte) (*decryptedKeys, error) {
	privKeyBytes, err := s.cryptoService.DecryptPrivateKey(ciphertext)
	if err != nil {
		return nil, err
	}
	privKey := ed25519.PrivateKey(privKeyBytes)
	pubKey := privKey.Public().(ed25519.PublicKey)
	return &decryptedKeys{PrivateKey: privKey, PublicKey: pubKey}, nil
}

func (s *Service) signInjectAndPersist(ctx context.Context, rootDict map[string]interface{}, mi *metainfo.MetaInfo, infoHash []byte, userID int64, keys *decryptedKeys) (torrentModel.TorrentDTO, error) {
	sigResult, err := s.buildSignature(mi, infoHash, keys)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.injectAndPersist(ctx, rootDict, mi, infoHash, userID, keys, sigResult)
}

func (s *Service) buildSignature(mi *metainfo.MetaInfo, infoHash []byte, keys *decryptedKeys) (*ports.SignatureResult, error) {
	req := ports.SignatureRequest{
		InfoHash:     infoHash,
		AnnounceList: mi.AnnounceList,
		PrivateKey:   keys.PrivateKey,
	}
	return s.signatureBuilder.BuildPayloadAndSign(req)
}

func (s *Service) injectAndPersist(ctx context.Context, rootDict map[string]interface{}, mi *metainfo.MetaInfo, infoHash []byte, userID int64, keys *decryptedKeys, sigResult *ports.SignatureResult) (torrentModel.TorrentDTO, error) {
	signedBytes, err := s.injectSignature(rootDict, keys, sigResult)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.persistAndBuildDTO(ctx, signedBytes, mi, infoHash, userID, sigResult)
}

func (s *Service) injectSignature(rootDict map[string]interface{}, keys *decryptedKeys, sigResult *ports.SignatureResult) ([]byte, error) {
	req := ports.InjectionRequest{
		RootDict:  rootDict,
		PubKey:    keys.PublicKey,
		Signature: sigResult.Signature,
		Timestamp: sigResult.Timestamp,
	}
	return s.metaInfoSigner.InjectSignature(req)
}

func (s *Service) persistAndBuildDTO(ctx context.Context, signedBytes []byte, mi *metainfo.MetaInfo, infoHash []byte, userID int64, sigResult *ports.SignatureResult) (torrentModel.TorrentDTO, error) {
	if err := s.persistSignature(ctx, infoHash, userID, sigResult); err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.buildDTO(signedBytes, mi)
}

func (s *Service) persistSignature(ctx context.Context, infoHash []byte, userID int64, sigResult *ports.SignatureResult) error {
	entity := s.buildSignatureEntity(infoHash, userID, sigResult)
	return s.signatureRepo.Create(ctx, entity)
}

func (s *Service) buildSignatureEntity(infoHash []byte, userID int64, sigResult *ports.SignatureResult) torrentSignatureModel.TorrentSignatureEntity {
	return torrentSignatureModel.TorrentSignatureEntity{
		TorrentHash:   infoHash,
		UserID:        userID,
		SignatureBlob: sigResult.Signature,
		PayloadHash:   sigResult.PayloadHash,
		CreatedAt:     time.Now(),
	}
}

func (s *Service) buildDTO(signedBytes []byte, mi *metainfo.MetaInfo) (torrentModel.TorrentDTO, error) {
	files, signatures := s.extractor.Extract(signedBytes)
	info, err := mi.UnmarshalInfo()
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.mapToDTO(&info, mi, files, signatures), nil
}

func (s *Service) mapToDTO(info *metainfo.Info, mi *metainfo.MetaInfo, files []torrentModel.FileDTO, signatures []torrentModel.SignatureDTO) torrentModel.TorrentDTO {
	return torrentModel.TorrentDTO{
		InfoHash:   hex.EncodeToString(mi.HashInfoBytes().Bytes()),
		Name:       info.BestName(),
		SizeBytes:  info.TotalLength(),
		Status:     torrentModel.StatusIdle,
		Progress:   0.0,
		AddedAt:    time.Now(),
		Files:      files,
		Signatures: signatures,
	}
}
