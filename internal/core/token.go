package core

type TokenMetadata struct {
	UserId  int
	Role    string
	Expires int64
}

type TokenWithClaims struct {
	Token  string
	UserId int
	Role   string
}

type RefreshTokenMetadata struct {
	UserId  int
	Expires int64
}

type RefreshTokenWithClaims struct {
	Token  string
	UserId int
}
