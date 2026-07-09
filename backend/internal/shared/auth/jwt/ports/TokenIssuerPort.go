package ports

type TokenIssuerPort interface {
	IssueAccessToken(userID int64, peerID string) (string, error)
	IssueRefreshToken() (string, error)
}
