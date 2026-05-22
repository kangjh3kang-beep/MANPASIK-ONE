// Manual gRPC service definitions for MarketplaceService.
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
	MarketplaceService_ListProducts_FullMethodName    = "/manpasik.v1.MarketplaceService/ListProducts"
	MarketplaceService_GetProduct_FullMethodName      = "/manpasik.v1.MarketplaceService/GetProduct"
	MarketplaceService_CreateProduct_FullMethodName   = "/manpasik.v1.MarketplaceService/CreateProduct"
	MarketplaceService_UpdateProduct_FullMethodName   = "/manpasik.v1.MarketplaceService/UpdateProduct"
	MarketplaceService_RegisterPartner_FullMethodName = "/manpasik.v1.MarketplaceService/RegisterPartner"
	MarketplaceService_GetPartnerStats_FullMethodName = "/manpasik.v1.MarketplaceService/GetPartnerStats"
)

// ============================================================================
// Message types
// ============================================================================

type ListMarketplaceProductsRequest struct {
	PartnerId string `json:"partner_id,omitempty"`
	Category  string `json:"category,omitempty"`
	Limit     int32  `json:"limit,omitempty"`
	Offset    int32  `json:"offset,omitempty"`
}

type ListMarketplaceProductsResponse struct {
	Products []*MarketplaceProduct `json:"products,omitempty"`
	Total    int32                 `json:"total,omitempty"`
}

type MarketplaceProduct struct {
	Id          string                 `json:"id,omitempty"`
	PartnerId   string                 `json:"partner_id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price,omitempty"`
	Category    string                 `json:"category,omitempty"`
	ImageUrl    string                 `json:"image_url,omitempty"`
	IsActive    bool                   `json:"is_active,omitempty"`
	CreatedAt   *timestamppb.Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `json:"updated_at,omitempty"`
}

type GetMarketplaceProductRequest struct {
	ProductId string `json:"product_id,omitempty"`
}

type CreateMarketplaceProductRequest struct {
	PartnerId   string  `json:"partner_id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Category    string  `json:"category,omitempty"`
	ImageUrl    string  `json:"image_url,omitempty"`
}

type UpdateMarketplaceProductRequest struct {
	ProductId   string  `json:"product_id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Category    string  `json:"category,omitempty"`
	ImageUrl    string  `json:"image_url,omitempty"`
	IsActive    bool    `json:"is_active,omitempty"`
}

type UpdateMarketplaceProductResponse struct {
	Success bool `json:"success,omitempty"`
}

type RegisterMarketplacePartnerRequest struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
}

type RegisterMarketplacePartnerResponse struct {
	PartnerId string `json:"partner_id,omitempty"`
	Status    string `json:"status,omitempty"`
}

type GetPartnerStatsRequest struct {
	PartnerId string `json:"partner_id,omitempty"`
}

type PartnerStatsResponse struct {
	PartnerId     string  `json:"partner_id,omitempty"`
	TotalProducts int32   `json:"total_products,omitempty"`
	TotalOrders   int32   `json:"total_orders,omitempty"`
	Revenue       float64 `json:"revenue,omitempty"`
}

// ============================================================================
// Server interface
// ============================================================================

type MarketplaceServiceServer interface {
	ListProducts(context.Context, *ListMarketplaceProductsRequest) (*ListMarketplaceProductsResponse, error)
	GetProduct(context.Context, *GetMarketplaceProductRequest) (*MarketplaceProduct, error)
	CreateProduct(context.Context, *CreateMarketplaceProductRequest) (*MarketplaceProduct, error)
	UpdateProduct(context.Context, *UpdateMarketplaceProductRequest) (*UpdateMarketplaceProductResponse, error)
	RegisterPartner(context.Context, *RegisterMarketplacePartnerRequest) (*RegisterMarketplacePartnerResponse, error)
	GetPartnerStats(context.Context, *GetPartnerStatsRequest) (*PartnerStatsResponse, error)
	mustEmbedUnimplementedMarketplaceServiceServer()
}

type UnimplementedMarketplaceServiceServer struct{}

func (UnimplementedMarketplaceServiceServer) ListProducts(context.Context, *ListMarketplaceProductsRequest) (*ListMarketplaceProductsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProducts not implemented")
}
func (UnimplementedMarketplaceServiceServer) GetProduct(context.Context, *GetMarketplaceProductRequest) (*MarketplaceProduct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}
func (UnimplementedMarketplaceServiceServer) CreateProduct(context.Context, *CreateMarketplaceProductRequest) (*MarketplaceProduct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}
func (UnimplementedMarketplaceServiceServer) UpdateProduct(context.Context, *UpdateMarketplaceProductRequest) (*UpdateMarketplaceProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProduct not implemented")
}
func (UnimplementedMarketplaceServiceServer) RegisterPartner(context.Context, *RegisterMarketplacePartnerRequest) (*RegisterMarketplacePartnerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterPartner not implemented")
}
func (UnimplementedMarketplaceServiceServer) GetPartnerStats(context.Context, *GetPartnerStatsRequest) (*PartnerStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPartnerStats not implemented")
}
func (UnimplementedMarketplaceServiceServer) mustEmbedUnimplementedMarketplaceServiceServer() {}

type UnsafeMarketplaceServiceServer interface {
	mustEmbedUnimplementedMarketplaceServiceServer()
}

// ============================================================================
// Registration
// ============================================================================

func RegisterMarketplaceServiceServer(s grpc.ServiceRegistrar, srv MarketplaceServiceServer) {
	s.RegisterService(&MarketplaceService_ServiceDesc, srv)
}

func _MarketplaceService_ListProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMarketplaceProductsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketplaceServiceServer).ListProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: MarketplaceService_ListProducts_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketplaceServiceServer).ListProducts(ctx, req.(*ListMarketplaceProductsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketplaceService_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMarketplaceProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketplaceServiceServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: MarketplaceService_GetProduct_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketplaceServiceServer).GetProduct(ctx, req.(*GetMarketplaceProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketplaceService_CreateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMarketplaceProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketplaceServiceServer).CreateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: MarketplaceService_CreateProduct_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketplaceServiceServer).CreateProduct(ctx, req.(*CreateMarketplaceProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketplaceService_UpdateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMarketplaceProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketplaceServiceServer).UpdateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: MarketplaceService_UpdateProduct_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketplaceServiceServer).UpdateProduct(ctx, req.(*UpdateMarketplaceProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketplaceService_RegisterPartner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterMarketplacePartnerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketplaceServiceServer).RegisterPartner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: MarketplaceService_RegisterPartner_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketplaceServiceServer).RegisterPartner(ctx, req.(*RegisterMarketplacePartnerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketplaceService_GetPartnerStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPartnerStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketplaceServiceServer).GetPartnerStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: MarketplaceService_GetPartnerStats_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketplaceServiceServer).GetPartnerStats(ctx, req.(*GetPartnerStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ============================================================================
// Client
// ============================================================================

type MarketplaceServiceClient interface {
	ListProducts(ctx context.Context, in *ListMarketplaceProductsRequest, opts ...grpc.CallOption) (*ListMarketplaceProductsResponse, error)
	GetProduct(ctx context.Context, in *GetMarketplaceProductRequest, opts ...grpc.CallOption) (*MarketplaceProduct, error)
	CreateProduct(ctx context.Context, in *CreateMarketplaceProductRequest, opts ...grpc.CallOption) (*MarketplaceProduct, error)
	UpdateProduct(ctx context.Context, in *UpdateMarketplaceProductRequest, opts ...grpc.CallOption) (*UpdateMarketplaceProductResponse, error)
	RegisterPartner(ctx context.Context, in *RegisterMarketplacePartnerRequest, opts ...grpc.CallOption) (*RegisterMarketplacePartnerResponse, error)
	GetPartnerStats(ctx context.Context, in *GetPartnerStatsRequest, opts ...grpc.CallOption) (*PartnerStatsResponse, error)
}

type marketplaceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMarketplaceServiceClient(cc grpc.ClientConnInterface) MarketplaceServiceClient {
	return &marketplaceServiceClient{cc}
}

func (c *marketplaceServiceClient) ListProducts(ctx context.Context, in *ListMarketplaceProductsRequest, opts ...grpc.CallOption) (*ListMarketplaceProductsResponse, error) {
	out := new(ListMarketplaceProductsResponse)
	err := c.cc.Invoke(ctx, MarketplaceService_ListProducts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketplaceServiceClient) GetProduct(ctx context.Context, in *GetMarketplaceProductRequest, opts ...grpc.CallOption) (*MarketplaceProduct, error) {
	out := new(MarketplaceProduct)
	err := c.cc.Invoke(ctx, MarketplaceService_GetProduct_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketplaceServiceClient) CreateProduct(ctx context.Context, in *CreateMarketplaceProductRequest, opts ...grpc.CallOption) (*MarketplaceProduct, error) {
	out := new(MarketplaceProduct)
	err := c.cc.Invoke(ctx, MarketplaceService_CreateProduct_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketplaceServiceClient) UpdateProduct(ctx context.Context, in *UpdateMarketplaceProductRequest, opts ...grpc.CallOption) (*UpdateMarketplaceProductResponse, error) {
	out := new(UpdateMarketplaceProductResponse)
	err := c.cc.Invoke(ctx, MarketplaceService_UpdateProduct_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketplaceServiceClient) RegisterPartner(ctx context.Context, in *RegisterMarketplacePartnerRequest, opts ...grpc.CallOption) (*RegisterMarketplacePartnerResponse, error) {
	out := new(RegisterMarketplacePartnerResponse)
	err := c.cc.Invoke(ctx, MarketplaceService_RegisterPartner_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketplaceServiceClient) GetPartnerStats(ctx context.Context, in *GetPartnerStatsRequest, opts ...grpc.CallOption) (*PartnerStatsResponse, error) {
	out := new(PartnerStatsResponse)
	err := c.cc.Invoke(ctx, MarketplaceService_GetPartnerStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var MarketplaceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "manpasik.v1.MarketplaceService",
	HandlerType: (*MarketplaceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "ListProducts", Handler: _MarketplaceService_ListProducts_Handler},
		{MethodName: "GetProduct", Handler: _MarketplaceService_GetProduct_Handler},
		{MethodName: "CreateProduct", Handler: _MarketplaceService_CreateProduct_Handler},
		{MethodName: "UpdateProduct", Handler: _MarketplaceService_UpdateProduct_Handler},
		{MethodName: "RegisterPartner", Handler: _MarketplaceService_RegisterPartner_Handler},
		{MethodName: "GetPartnerStats", Handler: _MarketplaceService_GetPartnerStats_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "manpasik.proto",
}
