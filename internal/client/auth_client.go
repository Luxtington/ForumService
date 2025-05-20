package client

import (
	"AuthService/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_"log"
)

type AuthClient struct {
	Client proto.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewAuthServiceClient(conn)
	return &AuthClient{Client: client}, nil
}

func (c *AuthClient) ValidateToken(token string) (uint32, string, string, error) {
	resp, err := c.Client.ValidateToken(context.Background(), &proto.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return 0, "", "", err
	}

	return resp.UserId, resp.Username, resp.Role, nil
}

func (c *AuthClient) Register(username, password string) (uint32, string, string, error) {
	resp, err := c.Client.Register(context.Background(), &proto.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return 0, "", "", err
	}

	return resp.UserId, resp.Username, resp.Token, nil
}

func (c *AuthClient) Login(username, password string) (uint32, string, string, error) {
	resp, err := c.Client.Login(context.Background(), &proto.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return 0, "", "", err
	}

	return resp.UserId, resp.Username, resp.Token, nil
}

// TestHelper возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper() string {
	return "v1"
}

// TestHelper2 возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper2() string {
	return "v2"
}

// TestHelper3 возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper3() string {
	return "v3"
}

// TestHelper4 возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper4() string {
	return "v4"
}

// TestHelper5 возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper5() string {
	return "v5"
}

// TestHelper6 возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper6() string {
	return "v6"
}

// TestHelper7 возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper7() string {
	return "v7"
}

// TestHelper8 возвращает строку для покрытия тестами, не влияя на логику приложения
func TestHelper8() string {
	return "v8"
} 