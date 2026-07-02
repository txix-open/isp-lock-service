package routes

import (
	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/common_endpoints"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"

	"isp-lock-service/controller"
	"isp-lock-service/docs"
)

type Controllers struct {
	Locker       controller.Locker
	RateLimiter  controller.RateLimiter
	DailyLimiter controller.DailyLimiter
}

func EndpointDescriptors() []cluster.EndpointDescriptor {
	return concatDescriptors(endpointDescriptors(Controllers{}), commonEndpoints())
}

func Handler(wrapper endpoint.Wrapper, c Controllers) *grpc.Mux {
	muxer := grpc.NewMux()
	for _, descriptor := range endpointDescriptors(c) {
		muxer.Handle(descriptor.Path, wrapper.Endpoint(descriptor.Handler))
	}
	for _, descriptor := range commonEndpoints() {
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
		Path:             "isp-lock-service/try_lock",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.Locker.TryLock,
	}, {
		Path:             "isp-lock-service/unlock",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.Locker.UnLock,
	}, {
		Path:             "isp-lock-service/rate_limit",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.RateLimiter.Limit,
	}, {
		Path:             "isp-lock-service/daily_limit/increment",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.DailyLimiter.Increment,
	}, {
		Path:             "isp-lock-service/daily_limit/set",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.DailyLimiter.Set,
	}, {
		Path:             "isp-lock-service/daily_limit/get",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.DailyLimiter.Get,
	}, {
		Path:             "isp-lock-service/rate_limit/inmem",
		Inner:            true,
		UserAuthRequired: false,
		Handler:          c.RateLimiter.LimitInMem,
	}}
}

func commonEndpoints() []cluster.EndpointDescriptor {
	return common_endpoints.CommonEndpoints("isp-lock-service",
		common_endpoints.WithSwaggerEndpoint(docs.SwaggerJson),
	)
}

func concatDescriptors(descs ...[]cluster.EndpointDescriptor) []cluster.EndpointDescriptor {
	result := make([]cluster.EndpointDescriptor, 0)
	for _, desc := range descs {
		result = append(result, desc...)
	}
	return result
}
