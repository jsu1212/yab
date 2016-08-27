// Code generated by thriftrw
// @generated

package fooserver

import (
	"github.com/thriftrw/thriftrw-go/protocol"
	"github.com/thriftrw/thriftrw-go/wire"
	"github.com/yarpc/yab/testdata/yarpc/integration/service/foo"
	yarpc "github.com/yarpc/yarpc-go"
	"github.com/yarpc/yarpc-go/encoding/thrift"
	"golang.org/x/net/context"
)

type Interface interface {
	Bar(ctx context.Context, reqMeta yarpc.ReqMeta, arg *int32) (int32, yarpc.ResMeta, error)
}

func New(impl Interface) thrift.Service {
	return service{handler{impl}}
}

type service struct{ h handler }

func (service) Name() string {
	return "Foo"
}

func (service) Protocol() protocol.Protocol {
	return protocol.Binary
}

func (s service) Handlers() map[string]thrift.Handler {
	return map[string]thrift.Handler{"bar": thrift.HandlerFunc(s.h.Bar)}
}

type handler struct{ impl Interface }

func (h handler) Bar(ctx context.Context, reqMeta yarpc.ReqMeta, body wire.Value) (thrift.Response, error) {
	var args foo.BarArgs
	if err := args.FromWire(body); err != nil {
		return thrift.Response{}, err
	}
	success, resMeta, err := h.impl.Bar(ctx, reqMeta, args.Arg)
	hadError := err != nil
	result, err := foo.BarHelper.WrapResponse(success, err)
	var response thrift.Response
	if err == nil {
		response.IsApplicationError = hadError
		response.Meta = resMeta
		response.Body = result
	}
	return response, err
}
