package client

import (
	"AuthService/proto"
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	proto.RegisterAuthServiceServer(s, &mockAuthServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

type mockAuthServer struct {
	proto.UnimplementedAuthServiceServer
}

func (s *mockAuthServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	if req.Token == "valid_token" {
		return &proto.ValidateTokenResponse{
			UserId:   1,
			Username: "test_user",
			Role:     "user",
		}, nil
	}
	return nil, assert.AnError
}

func (s *mockAuthServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req.Username == "test_user" && req.Password == "test_password" {
		return &proto.RegisterResponse{
			UserId:   1,
			Username: "test_user",
			Token:    "valid_token",
		}, nil
	}
	return nil, assert.AnError
}

func (s *mockAuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req.Username == "test_user" && req.Password == "test_password" {
		return &proto.LoginResponse{
			UserId:   1,
			Username: "test_user",
			Token:    "valid_token",
		}, nil
	}
	return nil, assert.AnError
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestNewAuthClient(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := proto.NewAuthServiceClient(conn)
	authClient := &AuthClient{Client: client}

	// Тест успешной валидации токена
	userId, username, role, err := authClient.ValidateToken("valid_token")
	require.NoError(t, err)
	assert.Equal(t, uint32(1), userId)
	assert.Equal(t, "test_user", username)
	assert.Equal(t, "user", role)

	// Тест невалидного токена
	userId, username, role, err = authClient.ValidateToken("invalid_token")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), userId)
	assert.Empty(t, username)
	assert.Empty(t, role)

	// Тест успешной регистрации
	userId, username, token, err := authClient.Register("test_user", "test_password")
	require.NoError(t, err)
	assert.Equal(t, uint32(1), userId)
	assert.Equal(t, "test_user", username)
	assert.Equal(t, "valid_token", token)

	// Тест неуспешной регистрации
	userId, username, token, err = authClient.Register("invalid_user", "invalid_password")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), userId)
	assert.Empty(t, username)
	assert.Empty(t, token)

	// Тест успешного входа
	userId, username, token, err = authClient.Login("test_user", "test_password")
	require.NoError(t, err)
	assert.Equal(t, uint32(1), userId)
	assert.Equal(t, "test_user", username)
	assert.Equal(t, "valid_token", token)

	// Тест неуспешного входа
	userId, username, token, err = authClient.Login("invalid_user", "invalid_password")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), userId)
	assert.Empty(t, username)
	assert.Empty(t, token)
}

func TestTestHelper(t *testing.T) {
	assert.Equal(t, "v1", TestHelper())
}

func TestTestHelper2(t *testing.T) {
	assert.Equal(t, "v2", TestHelper2())
}

func TestTestHelper3(t *testing.T) {
	assert.Equal(t, "v3", TestHelper3())
}

func TestTestHelper4(t *testing.T) {
	assert.Equal(t, "v4", TestHelper4())
}

func TestTestHelper5(t *testing.T) {
	assert.Equal(t, "v5", TestHelper5())
}

func TestTestHelper6(t *testing.T) {
	assert.Equal(t, "v6", TestHelper6())
}

func TestTestHelper7(t *testing.T) {
	assert.Equal(t, "v7", TestHelper7())
}

func TestTestHelper8(t *testing.T) {
	assert.Equal(t, "v8", TestHelper8())
} 