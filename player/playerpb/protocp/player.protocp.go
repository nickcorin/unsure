package protocp

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/luno/jettison/errors"
	"github.com/nickcorin/unsure/player"
	pb "github.com/nickcorin/unsure/player/playerpb"
)

// PartFromProto converts a pb.Part to a player.Part.
func PartFromProto(in *pb.Part) (*player.Part, error) {
	createdAt, err := ptypes.Timestamp(in.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	updatedAt, err := ptypes.Timestamp(in.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	return &player.Part{
		ID:        in.Id,
		RoundID:   in.RoundId,
		Player:    in.Player,
		Rank:      in.Rank,
		Value:     in.Value,
		Submitted: in.Submitted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// PartToProto converts a player.Part to a pb.Part.
func PartToProto(in *player.Part) (*pb.Part, error) {
	createdAt, err := ptypes.TimestampProto(in.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	updatedAt, err := ptypes.TimestampProto(in.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	return &pb.Part{
		Id:        in.ID,
		RoundId:   in.RoundID,
		Player:    in.Player,
		Rank:      in.Rank,
		Value:     in.Value,
		Submitted: in.Submitted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func RoundFromProto(in *pb.Round) (*player.Round, error) {
	createdAt, err := ptypes.Timestamp(in.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	updatedAt, err := ptypes.Timestamp(in.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	return &player.Round{
		ID:         in.Id,
		ExternalID: in.ExternalId,
		Player:     in.Player,
		Status:     player.RoundStatus(in.Status),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func RoundToProto(in *player.Round) (*pb.Round, error) {
	createdAt, err := ptypes.TimestampProto(in.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	updatedAt, err := ptypes.TimestampProto(in.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	return &pb.Round{
		Id:         in.ID,
		ExternalId: in.ExternalID,
		Player:     in.Player,
		Status:     int32(in.Status),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
