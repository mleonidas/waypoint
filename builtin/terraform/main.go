package terraform

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

//go:generate protoc -I ../../.. -I ../../thirdparty/proto/opaqueany -I ../../thirdparty/proto --go_out=../../.. waypoint/builtin/terraform/plugin.proto

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Infra{}),
}
