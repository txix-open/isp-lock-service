package routes

import (
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/grpc/endpoint"

	"isp-lock-service/controller"
)

type Controllers struct {
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
	return []cluster.EndpointDescriptor{{
		Path:             "isp-lock-service/lock",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.Locker.Lock,
	}, {
		Path:             "isp-lock-service/unlock",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.Locker.UnLock,
	}}
}
