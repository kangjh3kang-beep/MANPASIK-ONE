// Manual gRPC service definitions for IoTGatewayService.
// These will be replaced by protoc-generated code once the proto definitions are regenerated.
package v1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// ============================================================================
// Method name constants
// ============================================================================

const (
	IoTGatewayService_RegisterDevice_FullMethodName    = "/manpasik.v1.IoTGatewayService/RegisterDevice"
	IoTGatewayService_GetDevice_FullMethodName          = "/manpasik.v1.IoTGatewayService/GetDevice"
	IoTGatewayService_SendDeviceCommand_FullMethodName  = "/manpasik.v1.IoTGatewayService/SendDeviceCommand"
	IoTGatewayService_GetDeviceTelemetry_FullMethodName = "/manpasik.v1.IoTGatewayService/GetDeviceTelemetry"
	IoTGatewayService_Heartbeat_FullMethodName          = "/manpasik.v1.IoTGatewayService/Heartbeat"
	IoTGatewayService_ListDeviceData_FullMethodName     = "/manpasik.v1.IoTGatewayService/ListDeviceData"
)

// ============================================================================
// Message types
// ============================================================================

type RegisterIoTDeviceRequest struct {
	DeviceId string            `json:"device_id,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type IoTDeviceResponse struct {
	Id         string                 `json:"id,omitempty"`
	DeviceId   string                 `json:"device_id,omitempty"`
	Protocol   string                 `json:"protocol,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Metadata   map[string]string      `json:"metadata,omitempty"`
	LastPingAt *timestamppb.Timestamp `json:"last_ping_at,omitempty"`
}

type GetIoTDeviceRequest struct {
	DeviceId string `json:"device_id,omitempty"`
}

type SendDeviceCommandRequest struct {
	DeviceId    string `json:"device_id,omitempty"`
	CommandType string `json:"command_type,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

type DeviceCommandResponse struct {
	Id          string                 `json:"id,omitempty"`
	DeviceId    string                 `json:"device_id,omitempty"`
	CommandType string                 `json:"command_type,omitempty"`
	Payload     string                 `json:"payload,omitempty"`
	Status      string                 `json:"status,omitempty"`
	CreatedAt   *timestamppb.Timestamp `json:"created_at,omitempty"`
}

type GetDeviceTelemetryRequest struct {
	DeviceId string `json:"device_id,omitempty"`
	Limit    int32  `json:"limit,omitempty"`
	Offset   int32  `json:"offset,omitempty"`
}

type DeviceTelemetryResponse struct {
	Data []*IoTDataPoint `json:"data,omitempty"`
}

type IoTDataPoint struct {
	Id         string                 `json:"id,omitempty"`
	DeviceId   string                 `json:"device_id,omitempty"`
	DataType   string                 `json:"data_type,omitempty"`
	Value      float64                `json:"value,omitempty"`
	Unit       string                 `json:"unit,omitempty"`
	ReceivedAt *timestamppb.Timestamp `json:"received_at,omitempty"`
}

type HeartbeatRequest struct {
	DeviceId string `json:"device_id,omitempty"`
}

type ListDeviceDataRequest struct {
	DeviceId string `json:"device_id,omitempty"`
	Limit    int32  `json:"limit,omitempty"`
	Offset   int32  `json:"offset,omitempty"`
}

type ListDeviceDataResponse struct {
	Data []*IoTDataPoint `json:"data,omitempty"`
}

// ============================================================================
// Server interface
// ============================================================================

type IoTGatewayServiceServer interface {
	RegisterDevice(context.Context, *RegisterIoTDeviceRequest) (*IoTDeviceResponse, error)
	GetDevice(context.Context, *GetIoTDeviceRequest) (*IoTDeviceResponse, error)
	SendDeviceCommand(context.Context, *SendDeviceCommandRequest) (*DeviceCommandResponse, error)
	GetDeviceTelemetry(context.Context, *GetDeviceTelemetryRequest) (*DeviceTelemetryResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*IoTDeviceResponse, error)
	ListDeviceData(context.Context, *ListDeviceDataRequest) (*ListDeviceDataResponse, error)
	mustEmbedUnimplementedIoTGatewayServiceServer()
}

type UnimplementedIoTGatewayServiceServer struct{}

func (UnimplementedIoTGatewayServiceServer) RegisterDevice(context.Context, *RegisterIoTDeviceRequest) (*IoTDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDevice not implemented")
}
func (UnimplementedIoTGatewayServiceServer) GetDevice(context.Context, *GetIoTDeviceRequest) (*IoTDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDevice not implemented")
}
func (UnimplementedIoTGatewayServiceServer) SendDeviceCommand(context.Context, *SendDeviceCommandRequest) (*DeviceCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDeviceCommand not implemented")
}
func (UnimplementedIoTGatewayServiceServer) GetDeviceTelemetry(context.Context, *GetDeviceTelemetryRequest) (*DeviceTelemetryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceTelemetry not implemented")
}
func (UnimplementedIoTGatewayServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*IoTDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedIoTGatewayServiceServer) ListDeviceData(context.Context, *ListDeviceDataRequest) (*ListDeviceDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDeviceData not implemented")
}
func (UnimplementedIoTGatewayServiceServer) mustEmbedUnimplementedIoTGatewayServiceServer() {}

type UnsafeIoTGatewayServiceServer interface {
	mustEmbedUnimplementedIoTGatewayServiceServer()
}

// ============================================================================
// Registration
// ============================================================================

func RegisterIoTGatewayServiceServer(s grpc.ServiceRegistrar, srv IoTGatewayServiceServer) {
	s.RegisterService(&IoTGatewayService_ServiceDesc, srv)
}

func _IoTGatewayService_RegisterDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterIoTDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IoTGatewayServiceServer).RegisterDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: IoTGatewayService_RegisterDevice_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IoTGatewayServiceServer).RegisterDevice(ctx, req.(*RegisterIoTDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IoTGatewayService_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIoTDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IoTGatewayServiceServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: IoTGatewayService_GetDevice_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IoTGatewayServiceServer).GetDevice(ctx, req.(*GetIoTDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IoTGatewayService_SendDeviceCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendDeviceCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IoTGatewayServiceServer).SendDeviceCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: IoTGatewayService_SendDeviceCommand_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IoTGatewayServiceServer).SendDeviceCommand(ctx, req.(*SendDeviceCommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IoTGatewayService_GetDeviceTelemetry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceTelemetryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IoTGatewayServiceServer).GetDeviceTelemetry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: IoTGatewayService_GetDeviceTelemetry_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IoTGatewayServiceServer).GetDeviceTelemetry(ctx, req.(*GetDeviceTelemetryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IoTGatewayService_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IoTGatewayServiceServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: IoTGatewayService_Heartbeat_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IoTGatewayServiceServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IoTGatewayService_ListDeviceData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDeviceDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IoTGatewayServiceServer).ListDeviceData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: IoTGatewayService_ListDeviceData_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IoTGatewayServiceServer).ListDeviceData(ctx, req.(*ListDeviceDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ============================================================================
// Client
// ============================================================================

type IoTGatewayServiceClient interface {
	RegisterDevice(ctx context.Context, in *RegisterIoTDeviceRequest, opts ...grpc.CallOption) (*IoTDeviceResponse, error)
	GetDevice(ctx context.Context, in *GetIoTDeviceRequest, opts ...grpc.CallOption) (*IoTDeviceResponse, error)
	SendDeviceCommand(ctx context.Context, in *SendDeviceCommandRequest, opts ...grpc.CallOption) (*DeviceCommandResponse, error)
	GetDeviceTelemetry(ctx context.Context, in *GetDeviceTelemetryRequest, opts ...grpc.CallOption) (*DeviceTelemetryResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*IoTDeviceResponse, error)
	ListDeviceData(ctx context.Context, in *ListDeviceDataRequest, opts ...grpc.CallOption) (*ListDeviceDataResponse, error)
}

type ioTGatewayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIoTGatewayServiceClient(cc grpc.ClientConnInterface) IoTGatewayServiceClient {
	return &ioTGatewayServiceClient{cc}
}

func (c *ioTGatewayServiceClient) RegisterDevice(ctx context.Context, in *RegisterIoTDeviceRequest, opts ...grpc.CallOption) (*IoTDeviceResponse, error) {
	out := new(IoTDeviceResponse)
	err := c.cc.Invoke(ctx, IoTGatewayService_RegisterDevice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ioTGatewayServiceClient) GetDevice(ctx context.Context, in *GetIoTDeviceRequest, opts ...grpc.CallOption) (*IoTDeviceResponse, error) {
	out := new(IoTDeviceResponse)
	err := c.cc.Invoke(ctx, IoTGatewayService_GetDevice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ioTGatewayServiceClient) SendDeviceCommand(ctx context.Context, in *SendDeviceCommandRequest, opts ...grpc.CallOption) (*DeviceCommandResponse, error) {
	out := new(DeviceCommandResponse)
	err := c.cc.Invoke(ctx, IoTGatewayService_SendDeviceCommand_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ioTGatewayServiceClient) GetDeviceTelemetry(ctx context.Context, in *GetDeviceTelemetryRequest, opts ...grpc.CallOption) (*DeviceTelemetryResponse, error) {
	out := new(DeviceTelemetryResponse)
	err := c.cc.Invoke(ctx, IoTGatewayService_GetDeviceTelemetry_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ioTGatewayServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*IoTDeviceResponse, error) {
	out := new(IoTDeviceResponse)
	err := c.cc.Invoke(ctx, IoTGatewayService_Heartbeat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ioTGatewayServiceClient) ListDeviceData(ctx context.Context, in *ListDeviceDataRequest, opts ...grpc.CallOption) (*ListDeviceDataResponse, error) {
	out := new(ListDeviceDataResponse)
	err := c.cc.Invoke(ctx, IoTGatewayService_ListDeviceData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var IoTGatewayService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "manpasik.v1.IoTGatewayService",
	HandlerType: (*IoTGatewayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "RegisterDevice", Handler: _IoTGatewayService_RegisterDevice_Handler},
		{MethodName: "GetDevice", Handler: _IoTGatewayService_GetDevice_Handler},
		{MethodName: "SendDeviceCommand", Handler: _IoTGatewayService_SendDeviceCommand_Handler},
		{MethodName: "GetDeviceTelemetry", Handler: _IoTGatewayService_GetDeviceTelemetry_Handler},
		{MethodName: "Heartbeat", Handler: _IoTGatewayService_Heartbeat_Handler},
		{MethodName: "ListDeviceData", Handler: _IoTGatewayService_ListDeviceData_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "manpasik.proto",
}
