package response

import "github.com/BAN1ce/Tree/proto"

func NewClientOnlineResponse() *proto.ClientOnlineResponse {
	return &proto.ClientOnlineResponse{
		Base: &proto.BaseResponse{
			Success: false,
			Code:    0,
			Message: "",
		},
		LatestVersion: 0,
	}
}
