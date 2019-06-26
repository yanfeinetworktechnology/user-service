package main

import (
	"log"
	"net"

	"github.com/spf13/viper"

	base_common "base_service/common"
	user_common "user_service/common"

	"user_service/common"
	"user_service/model"
	pb "user_service/service/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) Check(ctx context.Context, in *pb.TokenRequest) (*pb.CheckResponse, error) {
	token := in.Token
	j := common.NewJWT()
	claims, err := j.ParseToken(token)

	if err == common.ErrTokenExpired {
		return &pb.CheckResponse{Status: -1}, nil
	}
	if err == common.ErrTokenInvalid {
		return &pb.CheckResponse{Status: -2}, nil
	}

	db := base_common.GetMySQL()
	userID := claims.UserID
	var existUser model.User
	err = db.Where("id = ?", userID).First(&existUser).Error

	if err != nil {
		return &pb.CheckResponse{Status: -3}, nil
	}
	if existUser.InfoID == 0 {
		return &pb.CheckResponse{UserId: userID, Status: -4}, nil
	}

	return &pb.CheckResponse{UserId: userID, Status: 1}, nil
}

// init 在 main 之前执行
func init() {
	// init config
	user_common.DefaultConfig()
	user_common.SetConfig()
	user_common.WatchConfig()

	// init Database
	db := base_common.InitMySQL()
	db.SingularTable(true)
}

func main() {
	lis, err := net.Listen("tcp", viper.GetString("service.address"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTokenServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
