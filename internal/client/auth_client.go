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