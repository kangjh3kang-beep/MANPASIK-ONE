// Manual gRPC service definitions for NLPService.
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
	NLPService_ParseHealthQuery_FullMethodName   = "/manpasik.v1.NLPService/ParseHealthQuery"
	NLPService_ExtractSymptoms_FullMethodName     = "/manpasik.v1.NLPService/ExtractSymptoms"
	NLPService_GetSuggestions_FullMethodName      = "/manpasik.v1.NLPService/GetSuggestions"
	NLPService_AnalyzeKoreanText_FullMethodName   = "/manpasik.v1.NLPService/AnalyzeKoreanText"
)

// ============================================================================
// Message types
// ============================================================================

type ParseHealthQueryRequest struct {
	UserId string `json:"user_id,omitempty"`
	Text   string `json:"text,omitempty"`
}

type ParseHealthQueryResponse struct {
	Id         string   `json:"id,omitempty"`
	UserId     string   `json:"user_id,omitempty"`
	RawText    string   `json:"raw_text,omitempty"`
	Intent     string   `json:"intent,omitempty"`
	Entities   []string `json:"entities,omitempty"`
	Confidence float64  `json:"confidence,omitempty"`
	Urgency    string   `json:"urgency,omitempty"`
}

type ExtractSymptomsRequest struct {
	Text string `json:"text,omitempty"`
}

type ExtractSymptomsResponse struct {
	Id          string             `json:"id,omitempty"`
	Text        string             `json:"text,omitempty"`
	Symptoms    []*NLPSymptom      `json:"symptoms,omitempty"`
	Urgency     string             `json:"urgency,omitempty"`
	ProcessedAt *timestamppb.Timestamp `json:"processed_at,omitempty"`
}

type NLPSymptom struct {
	Name       string  `json:"name,omitempty"`
	Severity   string  `json:"severity,omitempty"`
	BodyPart   string  `json:"body_part,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

type GetSuggestionsRequest struct {
	QueryId string `json:"query_id,omitempty"`
}

type GetSuggestionsResponse struct {
	Suggestions []*NLPSuggestion `json:"suggestions,omitempty"`
}

type NLPSuggestion struct {
	Id       string `json:"id,omitempty"`
	QueryId  string `json:"query_id,omitempty"`
	Text     string `json:"text,omitempty"`
	Category string `json:"category,omitempty"`
	Priority int32  `json:"priority,omitempty"`
}

type AnalyzeKoreanTextRequest struct {
	Text string `json:"text,omitempty"`
}

type AnalyzeKoreanTextResponse struct {
	Text       string   `json:"text,omitempty"`
	Intent     string   `json:"intent,omitempty"`
	Confidence float64  `json:"confidence,omitempty"`
	Severity   int32    `json:"severity,omitempty"`
	Keywords   []string `json:"keywords,omitempty"`
	Urgency    string   `json:"urgency,omitempty"`
}

// ============================================================================
// Server interface
// ============================================================================

type NLPServiceServer interface {
	ParseHealthQuery(context.Context, *ParseHealthQueryRequest) (*ParseHealthQueryResponse, error)
	ExtractSymptoms(context.Context, *ExtractSymptomsRequest) (*ExtractSymptomsResponse, error)
	GetSuggestions(context.Context, *GetSuggestionsRequest) (*GetSuggestionsResponse, error)
	AnalyzeKoreanText(context.Context, *AnalyzeKoreanTextRequest) (*AnalyzeKoreanTextResponse, error)
	mustEmbedUnimplementedNLPServiceServer()
}

type UnimplementedNLPServiceServer struct{}

func (UnimplementedNLPServiceServer) ParseHealthQuery(context.Context, *ParseHealthQueryRequest) (*ParseHealthQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseHealthQuery not implemented")
}
func (UnimplementedNLPServiceServer) ExtractSymptoms(context.Context, *ExtractSymptomsRequest) (*ExtractSymptomsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtractSymptoms not implemented")
}
func (UnimplementedNLPServiceServer) GetSuggestions(context.Context, *GetSuggestionsRequest) (*GetSuggestionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSuggestions not implemented")
}
func (UnimplementedNLPServiceServer) AnalyzeKoreanText(context.Context, *AnalyzeKoreanTextRequest) (*AnalyzeKoreanTextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnalyzeKoreanText not implemented")
}
func (UnimplementedNLPServiceServer) mustEmbedUnimplementedNLPServiceServer() {}

type UnsafeNLPServiceServer interface {
	mustEmbedUnimplementedNLPServiceServer()
}

// ============================================================================
// Registration
// ============================================================================

func RegisterNLPServiceServer(s grpc.ServiceRegistrar, srv NLPServiceServer) {
	s.RegisterService(&NLPService_ServiceDesc, srv)
}

func _NLPService_ParseHealthQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseHealthQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NLPServiceServer).ParseHealthQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: NLPService_ParseHealthQuery_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NLPServiceServer).ParseHealthQuery(ctx, req.(*ParseHealthQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NLPService_ExtractSymptoms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtractSymptomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NLPServiceServer).ExtractSymptoms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: NLPService_ExtractSymptoms_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NLPServiceServer).ExtractSymptoms(ctx, req.(*ExtractSymptomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NLPService_GetSuggestions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSuggestionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NLPServiceServer).GetSuggestions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: NLPService_GetSuggestions_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NLPServiceServer).GetSuggestions(ctx, req.(*GetSuggestionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NLPService_AnalyzeKoreanText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeKoreanTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NLPServiceServer).AnalyzeKoreanText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: NLPService_AnalyzeKoreanText_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NLPServiceServer).AnalyzeKoreanText(ctx, req.(*AnalyzeKoreanTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ============================================================================
// Client
// ============================================================================

type NLPServiceClient interface {
	ParseHealthQuery(ctx context.Context, in *ParseHealthQueryRequest, opts ...grpc.CallOption) (*ParseHealthQueryResponse, error)
	ExtractSymptoms(ctx context.Context, in *ExtractSymptomsRequest, opts ...grpc.CallOption) (*ExtractSymptomsResponse, error)
	GetSuggestions(ctx context.Context, in *GetSuggestionsRequest, opts ...grpc.CallOption) (*GetSuggestionsResponse, error)
	AnalyzeKoreanText(ctx context.Context, in *AnalyzeKoreanTextRequest, opts ...grpc.CallOption) (*AnalyzeKoreanTextResponse, error)
}

type nlpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNLPServiceClient(cc grpc.ClientConnInterface) NLPServiceClient {
	return &nlpServiceClient{cc}
}

func (c *nlpServiceClient) ParseHealthQuery(ctx context.Context, in *ParseHealthQueryRequest, opts ...grpc.CallOption) (*ParseHealthQueryResponse, error) {
	out := new(ParseHealthQueryResponse)
	err := c.cc.Invoke(ctx, NLPService_ParseHealthQuery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nlpServiceClient) ExtractSymptoms(ctx context.Context, in *ExtractSymptomsRequest, opts ...grpc.CallOption) (*ExtractSymptomsResponse, error) {
	out := new(ExtractSymptomsResponse)
	err := c.cc.Invoke(ctx, NLPService_ExtractSymptoms_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nlpServiceClient) GetSuggestions(ctx context.Context, in *GetSuggestionsRequest, opts ...grpc.CallOption) (*GetSuggestionsResponse, error) {
	out := new(GetSuggestionsResponse)
	err := c.cc.Invoke(ctx, NLPService_GetSuggestions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nlpServiceClient) AnalyzeKoreanText(ctx context.Context, in *AnalyzeKoreanTextRequest, opts ...grpc.CallOption) (*AnalyzeKoreanTextResponse, error) {
	out := new(AnalyzeKoreanTextResponse)
	err := c.cc.Invoke(ctx, NLPService_AnalyzeKoreanText_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var NLPService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "manpasik.v1.NLPService",
	HandlerType: (*NLPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "ParseHealthQuery", Handler: _NLPService_ParseHealthQuery_Handler},
		{MethodName: "ExtractSymptoms", Handler: _NLPService_ExtractSymptoms_Handler},
		{MethodName: "GetSuggestions", Handler: _NLPService_GetSuggestions_Handler},
		{MethodName: "AnalyzeKoreanText", Handler: _NLPService_AnalyzeKoreanText_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "manpasik.proto",
}
