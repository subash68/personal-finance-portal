package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	token "github.com/subash68/ate/ate_token_service/pkg/api/token"
)

var jwtKey = []byte("JWTSECRETKEYFORTOKEN")

//TODO: get db instance here from service configuration
type tokenServiceServer struct {
	token.UnimplementedTokenServiceServer
}

func NewTokenServiceServer() token.TokenServiceServer {
	return &tokenServiceServer{}
}

type Claims struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	UserType int32  `json:"userType"`
	jwt.StandardClaims
}

func (s *tokenServiceServer) Validate(ctx context.Context, req *token.ValidateRequest) (*token.ValidateResponse, error) {

	tokenString := req.Token
	claims := &Claims{}

	tkn, _ := jwt.ParseWithClaims(tokenString, claims, func(jwtToken *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if !tkn.Valid {
		return &token.ValidateResponse{
			Message: "Invalid token",
			Status:  false,
		}, nil
	}

	// This claims should be modified based on the token string from authorization server!!!
	return &token.ValidateResponse{
		Id:       claims.Id,
		Email:    claims.Email,
		UserType: claims.UserType,
		Message:  "verified successfully", // this needs to be modified later
		Status:   true,
	}, nil
}

func (s *tokenServiceServer) Generate(ctx context.Context, req *token.GenerateRequest) (*token.GenerateResponse, error) {

	//Expiration time for token string
	expirationTime := time.Now().Add(160 * time.Minute)

	fmt.Printf("some string to print here...")

	//create claims object
	claims := &Claims{
		Id:       req.Id,
		Email:    req.Email,
		Fullname: req.Fullname,
		UserType: req.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Use HS256 as signing method
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := tokenClaims.SignedString(jwtKey)
	// if err != nil {
	// 	// Handle token generation exception here
	// }

	return &token.GenerateResponse{
		Token: tokenString,
		Id:    req.Id,
	}, nil
}

func (s *tokenServiceServer) Refresh(ctx context.Context, req *token.RefreshRequest) (*token.RefreshResponse, error) {

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(req.Token, claims, func(tokenStr *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return &token.RefreshResponse{}, err
		}

		return &token.RefreshResponse{}, err
	}

	if !tkn.Valid {
		return &token.RefreshResponse{}, err
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return &token.RefreshResponse{}, &time.ParseError{}
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	tokenMethod := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenMethod.SignedString(jwtKey)

	// refresh token should be
	return &token.RefreshResponse{
		Token: tokenString,
	}, nil
}
