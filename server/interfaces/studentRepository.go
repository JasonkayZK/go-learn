package interfaces

import (
	"context"
	pb "grpc-sql-demo/server/proto"
)

type StudentRepo interface {
	Create(ctx context.Context, in *pb.Student) (*pb.StudentId, error)
	Read(ctx context.Context, in *pb.StudentId) (*pb.Student, error)
	Update(ctx context.Context, in *pb.Student) (*pb.StudentId, error)
	Delete(ctx context.Context, in *pb.StudentId) (*pb.StudentId, error)
}
