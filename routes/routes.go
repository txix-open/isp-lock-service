package routes

import (
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/grpc/endpoint"

	"isp-lock-service/controller"
)

type Controllers struct {
	Object controller.Object
	Locker controller.Locker
}

func EndpointDescriptors() []cluster.EndpointDescriptor {
	return endpointDescriptors(Controllers{})
}

func Handler(wrapper endpoint.Wrapper, c Controllers) *grpc.Mux {
	muxer := grpc.NewMux()
	for _, descriptor := range endpointDescriptors(c) {
		muxer.Handle(descriptor.Path, wrapper.Endpoint(descriptor.Handler))
	}
	return muxer
}

func endpointDescriptors(c Controllers) []cluster.EndpointDescriptor {
	return nil
	return []cluster.EndpointDescriptor{{
		// 	Path:             "isp-lock-service/object/all",
		// 	Inner:            false,
		// 	UserAuthRequired: false,
		// 	Handler:          c.Object.All,
		// }, {
		// 	Path:             "isp-lock-service/object/get_by_id",
		// 	Inner:            false,
		// 	UserAuthRequired: false,
		// 	Handler:          c.Object.GetById,
	}}
}
