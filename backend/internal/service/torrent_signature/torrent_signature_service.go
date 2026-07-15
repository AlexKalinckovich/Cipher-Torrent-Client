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
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent_signature/torrent_signature_repository_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/ports"
	userService "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security"
	storage_ports "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/storage/ports"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	torrentSignatureModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent_signature"
	"log"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

type decryptedKeys struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type Service struct {
	userRepo         userService.UserRepositoryPort
	signatureRepo    torrent_signature_repository_ports.TorrentSignatureRepositoryPort
	cryptoService    security.CryptoServicePort
	signatureBuilder ports.SignatureBuilder
	metaInfoSigner   ports.MetaInfoSigner
	minioRepo        storage_ports.TorrentStoragePort
	extractor        *meta_data_extractor.MetadataExtractor
}

func NewService(
	userRepo userService.UserRepositoryPort,
	signatureRepo torrent_signature_repository_ports.TorrentSignatureRepositoryPort,
	cryptoService security.CryptoServicePort,
	signatureBuilder ports.SignatureBuilder,
	minioRepo storage_ports.TorrentStoragePort,
	metaInfoSigner ports.MetaInfoSigner,
) *Service {
	return &Service{
		userRepo:         userRepo,
		signatureRepo:    signatureRepo,
		cryptoService:    cryptoService,
		signatureBuilder: signatureBuilder,
		metaInfoSigner:   metaInfoSigner,
		minioRepo:        minioRepo,
		extractor:        meta_data_extractor.NewMetadataExtractor(),
	}
}

func (s *Service) SignTorrent(ctx context.Context, req ports.SignTorrentServiceRequest) (torrentModel.TorrentDTO, error) {
	storageIdentityReq := storage_ports.StorageIdentityRequest{
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
	}
	fileBytes, err := s.minioRepo.DownloadBaseTorrent(ctx, storageIdentityReq)
	if err != nil {
		log.Println("Error downloading torrent file", err.Error())
		return torrentModel.TorrentDTO{}, err
	}

	return s.processSigning(ctx, fileBytes, req.UserID, req.InfoHash, req.CreatorPubKey)
}

func (s *Service) processSigning(ctx context.Context, fileBytes []byte, userID int64, infoHash []byte, torrentCreatorPubKey []byte) (torrentModel.TorrentDTO, error) {
	rootDict, mi, err := s.parseTorrentBytes(fileBytes)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	keys, err := s.getDecryptedKeyPair(ctx, userID)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	if err := s.ensureNotSigned(ctx, userID, infoHash, torrentCreatorPubKey); err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	return s.buildSignInjectAndSave(ctx, rootDict, mi, infoHash, userID, keys, torrentCreatorPubKey)
}

func (s *Service) parseTorrentBytes(fileBytes []byte) (map[string]interface{}, *metainfo.MetaInfo, error) {
	var rootDict map[string]interface{}
	if err := bencode.Unmarshal(fileBytes, &rootDict); err != nil {
		return nil, nil, err
	}
	mi, err := metainfo.Load(bytes.NewReader(fileBytes))
	return rootDict, mi, err
}

func (s *Service) ensureNotSigned(ctx context.Context, userID int64, infoHash []byte, torrentCreatorPubKey []byte) error {
	req := torrent_signature_repository_ports.UserTorrentIdentityRepositoryRequest{
		UserID:        userID,
		InfoHash:      infoHash,
		CreatorPubKey: torrentCreatorPubKey,
	}
	_, err := s.signatureRepo.GetByUserAndTorrentIdentity(ctx, req)
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

func (s *Service) buildSignInjectAndSave(ctx context.Context, rootDict map[string]interface{}, mi *metainfo.MetaInfo, infoHash []byte, userID int64, keys *decryptedKeys, torrentCreatorPubKey []byte) (torrentModel.TorrentDTO, error) {
	return s.signInjectAndPersist(ctx, rootDict, mi, infoHash, userID, keys, torrentCreatorPubKey)
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

func (s *Service) signInjectAndPersist(ctx context.Context, rootDict map[string]interface{}, mi *metainfo.MetaInfo, infoHash []byte, userID int64, keys *decryptedKeys, torrentCreatorPubKey []byte) (torrentModel.TorrentDTO, error) {
	sigResult, err := s.buildSignature(mi, infoHash, keys)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.injectAndPersist(ctx, rootDict, mi, infoHash, userID, keys, torrentCreatorPubKey, sigResult)
}

func (s *Service) buildSignature(mi *metainfo.MetaInfo, infoHash []byte, keys *decryptedKeys) (*ports.SignatureServiceResult, error) {
	req := ports.SignatureServiceRequest{
		InfoHash:     infoHash,
		AnnounceList: mi.AnnounceList,
		PrivateKey:   keys.PrivateKey,
	}
	return s.signatureBuilder.BuildPayloadAndSign(req)
}

func (s *Service) injectAndPersist(ctx context.Context, rootDict map[string]interface{}, mi *metainfo.MetaInfo, infoHash []byte, userID int64, keys *decryptedKeys, torrentCreatorPubKey []byte, sigResult *ports.SignatureServiceResult) (torrentModel.TorrentDTO, error) {
	signedBytes, err := s.injectSignature(rootDict, keys, sigResult)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.persistAndBuildDTO(ctx, signedBytes, infoHash, userID, torrentCreatorPubKey, sigResult)
}

func (s *Service) injectSignature(rootDict map[string]interface{}, keys *decryptedKeys, sigResult *ports.SignatureServiceResult) ([]byte, error) {
	req := ports.InjectionServiceRequest{
		RootDict:  rootDict,
		PubKey:    keys.PublicKey,
		Signature: sigResult.Signature,
		Timestamp: sigResult.Timestamp,
	}
	return s.metaInfoSigner.InjectSignature(req)
}

func (s *Service) persistAndBuildDTO(ctx context.Context, signedBytes []byte, infoHash []byte, userID int64, torrentCreatorPubKey []byte, sigResult *ports.SignatureServiceResult) (torrentModel.TorrentDTO, error) {
	if err := s.persistSignature(ctx, infoHash, userID, torrentCreatorPubKey, sigResult); err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	return s.buildDTO(signedBytes, infoHash)
}

func (s *Service) buildDTO(signedBytes []byte, originalInfoHash []byte) (torrentModel.TorrentDTO, error) {

	files, signatures := s.extractor.Extract(signedBytes)

	var info metainfo.Info
	if err := bencode.Unmarshal(signedBytes, &info); err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	return s.mapToDTO(&info, originalInfoHash, files, signatures), nil
}

func (s *Service) mapToDTO(info *metainfo.Info, originalInfoHash []byte, files []torrentModel.FileDTO, signatures []torrentModel.SignatureDTO) torrentModel.TorrentDTO {
	return torrentModel.TorrentDTO{
		InfoHash:   hex.EncodeToString(originalInfoHash),
		Name:       info.BestName(),
		SizeBytes:  info.TotalLength(),
		Status:     torrentModel.StatusIdle,
		Progress:   0.0,
		AddedAt:    time.Now(),
		Files:      files,
		Signatures: signatures,
	}
}
func (s *Service) persistSignature(ctx context.Context, infoHash []byte, userID int64, torrentCreatorPubKey []byte, sigResult *ports.SignatureServiceResult) error {
	sigEntity := s.buildSignatureEntity(infoHash, userID, sigResult)

	mapReq := torrent_signature_repository_ports.CreateSignatureMapRequest{
		TorrentHash:   infoHash,
		CreatorPubKey: torrentCreatorPubKey,
		SignerID:      userID,
		SignatureID:   sigEntity.ID,
		TrustLevel:    1,
	}

	return s.signatureRepo.CreateInTransaction(ctx, sigEntity, mapReq)
}

func (s *Service) buildSignatureEntity(infoHash []byte, userID int64, sigResult *ports.SignatureServiceResult) torrentSignatureModel.TorrentSignatureEntity {
	return torrentSignatureModel.TorrentSignatureEntity{
		TorrentHash:   infoHash,
		UserID:        userID,
		SignatureBlob: sigResult.Signature,
		PayloadHash:   sigResult.PayloadHash,
		CreatedAt:     time.Now(),
	}
}
