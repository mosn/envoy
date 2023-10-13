package plugins

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/http"
)

type config struct {
	AddReqHeaderName  string
	AddReqHeaderValue string
}

type parser struct {
}

func (p *parser) Parse(any *anypb.Any, callbacks api.ConfigCallbackHandler) (interface{}, error) {
	conf := &config{}
	return conf, nil
}

func (p *parser) Merge(parent interface{}, child interface{}) interface{} {
	return child
}

func addHeaderConfigFactory(c interface{}) api.StreamFilterFactory {
	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &addHeaderFilter{
			callbacks: callbacks,
			config:    c.(*config),
		}
	}
}

type addHeaderFilter struct {
	api.PassThroughStreamFilter
	callbacks api.FilterCallbackHandler
	config    *config
}

func (f *addHeaderFilter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	header.Add(f.config.AddReqHeaderName, f.config.AddReqHeaderValue)
	return api.LocalReply
}

func localReplyConfigFactory(interface{}) api.StreamFilterFactory {
	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &localReplyFilter{
			callbacks: callbacks,
		}
	}
}

type localReplyFilter struct {
	api.PassThroughStreamFilter
	callbacks api.FilterCallbackHandler
}

func (f *localReplyFilter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	f.callbacks.SendLocalReply(200, "", nil, 0, "")
	return api.LocalReply
}

func nonblockingConfigFactory(interface{}) api.StreamFilterFactory {
	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &nonblockingFilter{
			callbacks: callbacks,
		}
	}
}

type nonblockingFilter struct {
	api.PassThroughStreamFilter
	callbacks api.FilterCallbackHandler
}

func (f *nonblockingFilter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	// simulate time-consuming operation
	time.Sleep(1 * time.Millisecond)
	return api.Running
}

func panicConfigFactory(interface{}) api.StreamFilterFactory {
	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &panicFilter{
			callbacks: callbacks,
		}
	}
}

type panicFilter struct {
	api.PassThroughStreamFilter
	callbacks api.FilterCallbackHandler
}

func (f *panicFilter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	panic("ouch")
}

func logConfigFactory(interface{}) api.StreamFilterFactory {
	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &logFilter{
			callbacks: callbacks,
		}
	}
}

type logFilter struct {
	api.PassThroughStreamFilter
	callbacks api.FilterCallbackHandler
}

func (f *logFilter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	f.callbacks.Log(api.Critical, fmt.Sprintf("DecodeHeaders %v", endStream))
	return api.Continue
}

func (f *logFilter) OnLog() {
	f.callbacks.Log(api.Critical, "OnLog")
}

func (f *logFilter) OnLogDownstreamStart() {
	f.callbacks.Log(api.Critical, "OnLogDownstreamStart")
}

func (f *logFilter) OnLogDownstreamPeriodic() {
	f.callbacks.Log(api.Critical, "OnLogDownstreamPeriodic")
}

func (f *logFilter) OnDestroy(reason api.DestroyReason) {
	f.callbacks.Log(api.Critical, fmt.Sprintf("OnDestroy %v", reason))
}

func init() {
	http.RegisterHttpFilterConfigFactoryAndParser("addHeader", addHeaderConfigFactory, &parser{})
	http.RegisterHttpFilterConfigFactoryAndParser("localReply", localReplyConfigFactory, &parser{})
	http.RegisterHttpFilterConfigFactoryAndParser("nonblocking", nonblockingConfigFactory, &parser{})
	http.RegisterHttpFilterConfigFactoryAndParser("panic", panicConfigFactory, &parser{})
	http.RegisterHttpFilterConfigFactoryAndParser("log", logConfigFactory, &parser{})
}
