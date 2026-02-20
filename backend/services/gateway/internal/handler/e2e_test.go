package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
)

// ============================================================================
// Mock gRPC 클라이언트
// ============================================================================

// mockAuthClient는 AuthServiceClient를 모킹합니다.
type mockAuthClient struct{}

func (m *mockAuthClient) Register(_ context.Context, in *v1.RegisterRequest, _ ...grpc.CallOption) (*v1.RegisterResponse, error) {
	return &v1.RegisterResponse{UserId: "user_001", Email: in.Email, DisplayName: in.DisplayName}, nil
}
func (m *mockAuthClient) Login(_ context.Context, _ *v1.LoginRequest, _ ...grpc.CallOption) (*v1.LoginResponse, error) {
	return &v1.LoginResponse{AccessToken: "at_test", RefreshToken: "rt_test"}, nil
}
func (m *mockAuthClient) SocialLogin(_ context.Context, _ *v1.SocialLoginRequest, _ ...grpc.CallOption) (*v1.LoginResponse, error) {
	return &v1.LoginResponse{AccessToken: "at_social", RefreshToken: "rt_social"}, nil
}
func (m *mockAuthClient) RefreshToken(_ context.Context, _ *v1.RefreshTokenRequest, _ ...grpc.CallOption) (*v1.LoginResponse, error) {
	return &v1.LoginResponse{AccessToken: "at_new", RefreshToken: "rt_new"}, nil
}
func (m *mockAuthClient) Logout(_ context.Context, _ *v1.LogoutRequest, _ ...grpc.CallOption) (*v1.LogoutResponse, error) {
	return &v1.LogoutResponse{}, nil
}
func (m *mockAuthClient) ValidateToken(_ context.Context, _ *v1.ValidateTokenRequest, _ ...grpc.CallOption) (*v1.ValidateTokenResponse, error) {
	return &v1.ValidateTokenResponse{UserId: "user_001", Email: "test@example.com"}, nil
}
func (m *mockAuthClient) ResetPassword(_ context.Context, _ *v1.ResetPasswordRequest, _ ...grpc.CallOption) (*v1.ResetPasswordResponse, error) {
	return &v1.ResetPasswordResponse{Success: true}, nil
}

// mockMeasurementClient는 MeasurementServiceClient를 모킹합니다.
type mockMeasurementClient struct{}

func (m *mockMeasurementClient) StartSession(_ context.Context, _ *v1.StartSessionRequest, _ ...grpc.CallOption) (*v1.StartSessionResponse, error) {
	return &v1.StartSessionResponse{SessionId: "sess_001"}, nil
}
func (m *mockMeasurementClient) StreamMeasurement(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[v1.MeasurementData, v1.MeasurementResult], error) {
	return nil, nil
}
func (m *mockMeasurementClient) EndSession(_ context.Context, _ *v1.EndSessionRequest, _ ...grpc.CallOption) (*v1.EndSessionResponse, error) {
	return &v1.EndSessionResponse{SessionId: "sess_001"}, nil
}
func (m *mockMeasurementClient) GetMeasurementHistory(_ context.Context, _ *v1.GetHistoryRequest, _ ...grpc.CallOption) (*v1.GetHistoryResponse, error) {
	return &v1.GetHistoryResponse{}, nil
}
func (m *mockMeasurementClient) ExportSingleMeasurement(_ context.Context, _ *v1.ExportSingleMeasurementRequest, _ ...grpc.CallOption) (*v1.ExportFHIRResponse, error) {
	return &v1.ExportFHIRResponse{}, nil
}
func (m *mockMeasurementClient) ExportToFHIRObservations(_ context.Context, _ *v1.ExportToFHIRObservationsRequest, _ ...grpc.CallOption) (*v1.ExportFHIRResponse, error) {
	return &v1.ExportFHIRResponse{}, nil
}

// mockDeviceClient는 DeviceServiceClient를 모킹합니다.
type mockDeviceClient struct{}

func (m *mockDeviceClient) RegisterDevice(_ context.Context, _ *v1.RegisterDeviceRequest, _ ...grpc.CallOption) (*v1.RegisterDeviceResponse, error) {
	return &v1.RegisterDeviceResponse{DeviceId: "dev_001"}, nil
}
func (m *mockDeviceClient) ListDevices(_ context.Context, _ *v1.ListDevicesRequest, _ ...grpc.CallOption) (*v1.ListDevicesResponse, error) {
	return &v1.ListDevicesResponse{}, nil
}
func (m *mockDeviceClient) StreamDeviceStatus(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[v1.DeviceStatusUpdate, v1.DeviceCommand], error) {
	return nil, nil
}
func (m *mockDeviceClient) RequestOtaUpdate(_ context.Context, _ *v1.OtaRequest, _ ...grpc.CallOption) (*v1.OtaResponse, error) {
	return &v1.OtaResponse{}, nil
}
func (m *mockDeviceClient) UpdateDeviceStatus(_ context.Context, _ *v1.UpdateDeviceStatusRequest, _ ...grpc.CallOption) (*v1.UpdateDeviceStatusResponse, error) {
	return &v1.UpdateDeviceStatusResponse{}, nil
}

// mockShopClient는 ShopServiceClient를 모킹합니다.
type mockShopClient struct{}

func (m *mockShopClient) ListProducts(_ context.Context, _ *v1.ListProductsRequest, _ ...grpc.CallOption) (*v1.ListProductsResponse, error) {
	return &v1.ListProductsResponse{}, nil
}
func (m *mockShopClient) GetProduct(_ context.Context, _ *v1.GetProductRequest, _ ...grpc.CallOption) (*v1.Product, error) {
	return &v1.Product{ProductId: "prod_001"}, nil
}
func (m *mockShopClient) AddToCart(_ context.Context, _ *v1.AddToCartRequest, _ ...grpc.CallOption) (*v1.Cart, error) {
	return &v1.Cart{}, nil
}
func (m *mockShopClient) GetCart(_ context.Context, _ *v1.GetCartRequest, _ ...grpc.CallOption) (*v1.Cart, error) {
	return &v1.Cart{}, nil
}
func (m *mockShopClient) RemoveFromCart(_ context.Context, _ *v1.RemoveFromCartRequest, _ ...grpc.CallOption) (*v1.Cart, error) {
	return &v1.Cart{}, nil
}
func (m *mockShopClient) CreateOrder(_ context.Context, _ *v1.CreateOrderRequest, _ ...grpc.CallOption) (*v1.Order, error) {
	return &v1.Order{OrderId: "order_001"}, nil
}
func (m *mockShopClient) GetOrder(_ context.Context, _ *v1.GetOrderRequest, _ ...grpc.CallOption) (*v1.Order, error) {
	return &v1.Order{OrderId: "order_001"}, nil
}
func (m *mockShopClient) ListOrders(_ context.Context, _ *v1.ListOrdersRequest, _ ...grpc.CallOption) (*v1.ListOrdersResponse, error) {
	return &v1.ListOrdersResponse{}, nil
}

// mockSubscriptionClient는 SubscriptionServiceClient를 모킹합니다.
type mockSubscriptionClient struct{}

func (m *mockSubscriptionClient) CreateSubscription(_ context.Context, _ *v1.CreateSubscriptionRequest, _ ...grpc.CallOption) (*v1.SubscriptionDetail, error) {
	return &v1.SubscriptionDetail{}, nil
}
func (m *mockSubscriptionClient) GetSubscription(_ context.Context, _ *v1.GetSubscriptionDetailRequest, _ ...grpc.CallOption) (*v1.SubscriptionDetail, error) {
	return &v1.SubscriptionDetail{}, nil
}
func (m *mockSubscriptionClient) UpdateSubscription(_ context.Context, _ *v1.UpdateSubscriptionRequest, _ ...grpc.CallOption) (*v1.SubscriptionDetail, error) {
	return &v1.SubscriptionDetail{}, nil
}
func (m *mockSubscriptionClient) CancelSubscription(_ context.Context, _ *v1.CancelSubscriptionRequest, _ ...grpc.CallOption) (*v1.CancelSubscriptionResponse, error) {
	return &v1.CancelSubscriptionResponse{}, nil
}
func (m *mockSubscriptionClient) CheckFeatureAccess(_ context.Context, _ *v1.CheckFeatureAccessRequest, _ ...grpc.CallOption) (*v1.CheckFeatureAccessResponse, error) {
	return &v1.CheckFeatureAccessResponse{}, nil
}
func (m *mockSubscriptionClient) ListSubscriptionPlans(_ context.Context, _ *v1.ListSubscriptionPlansRequest, _ ...grpc.CallOption) (*v1.ListSubscriptionPlansResponse, error) {
	return &v1.ListSubscriptionPlansResponse{}, nil
}
func (m *mockSubscriptionClient) CheckCartridgeAccess(_ context.Context, _ *v1.CheckCartridgeAccessRequest, _ ...grpc.CallOption) (*v1.CheckCartridgeAccessResponse, error) {
	return &v1.CheckCartridgeAccessResponse{}, nil
}
func (m *mockSubscriptionClient) ListAccessibleCartridges(_ context.Context, _ *v1.ListAccessibleCartridgesRequest, _ ...grpc.CallOption) (*v1.ListAccessibleCartridgesResponse, error) {
	return &v1.ListAccessibleCartridgesResponse{}, nil
}

// mockCommunityClient는 CommunityServiceClient를 모킹합니다.
type mockCommunityClient struct{}

func (m *mockCommunityClient) CreatePost(_ context.Context, _ *v1.CreatePostRequest, _ ...grpc.CallOption) (*v1.Post, error) {
	return &v1.Post{PostId: "post_001"}, nil
}
func (m *mockCommunityClient) GetPost(_ context.Context, _ *v1.GetPostRequest, _ ...grpc.CallOption) (*v1.Post, error) {
	return &v1.Post{PostId: "post_001"}, nil
}
func (m *mockCommunityClient) ListPosts(_ context.Context, _ *v1.ListPostsRequest, _ ...grpc.CallOption) (*v1.ListPostsResponse, error) {
	return &v1.ListPostsResponse{}, nil
}
func (m *mockCommunityClient) LikePost(_ context.Context, _ *v1.LikePostRequest, _ ...grpc.CallOption) (*v1.LikePostResponse, error) {
	return &v1.LikePostResponse{}, nil
}
func (m *mockCommunityClient) CreateComment(_ context.Context, _ *v1.CreateCommentRequest, _ ...grpc.CallOption) (*v1.Comment, error) {
	return &v1.Comment{}, nil
}
func (m *mockCommunityClient) ListComments(_ context.Context, _ *v1.ListCommentsRequest, _ ...grpc.CallOption) (*v1.ListCommentsResponse, error) {
	return &v1.ListCommentsResponse{}, nil
}
func (m *mockCommunityClient) CreateChallenge(_ context.Context, _ *v1.CreateChallengeRequest, _ ...grpc.CallOption) (*v1.Challenge, error) {
	return &v1.Challenge{}, nil
}
func (m *mockCommunityClient) GetChallenge(_ context.Context, _ *v1.GetChallengeRequest, _ ...grpc.CallOption) (*v1.Challenge, error) {
	return &v1.Challenge{}, nil
}
func (m *mockCommunityClient) JoinChallenge(_ context.Context, _ *v1.JoinChallengeRequest, _ ...grpc.CallOption) (*v1.JoinChallengeResponse, error) {
	return &v1.JoinChallengeResponse{}, nil
}
func (m *mockCommunityClient) ListChallenges(_ context.Context, _ *v1.ListChallengesRequest, _ ...grpc.CallOption) (*v1.ListChallengesResponse, error) {
	return &v1.ListChallengesResponse{}, nil
}
func (m *mockCommunityClient) GetChallengeLeaderboard(_ context.Context, _ *v1.GetChallengeLeaderboardRequest, _ ...grpc.CallOption) (*v1.GetChallengeLeaderboardResponse, error) {
	return &v1.GetChallengeLeaderboardResponse{}, nil
}
func (m *mockCommunityClient) UpdateChallengeProgress(_ context.Context, _ *v1.UpdateChallengeProgressRequest, _ ...grpc.CallOption) (*v1.UpdateChallengeProgressResponse, error) {
	return &v1.UpdateChallengeProgressResponse{}, nil
}

// mockFamilyClient는 FamilyServiceClient를 모킹합니다.
type mockFamilyClient struct{}

func (m *mockFamilyClient) CreateFamilyGroup(_ context.Context, _ *v1.CreateFamilyGroupRequest, _ ...grpc.CallOption) (*v1.FamilyGroup, error) {
	return &v1.FamilyGroup{GroupId: "fam_001"}, nil
}
func (m *mockFamilyClient) GetFamilyGroup(_ context.Context, _ *v1.GetFamilyGroupRequest, _ ...grpc.CallOption) (*v1.FamilyGroup, error) {
	return &v1.FamilyGroup{GroupId: "fam_001"}, nil
}
func (m *mockFamilyClient) InviteMember(_ context.Context, _ *v1.InviteMemberRequest, _ ...grpc.CallOption) (*v1.FamilyInvitation, error) {
	return &v1.FamilyInvitation{}, nil
}
func (m *mockFamilyClient) RespondToInvitation(_ context.Context, _ *v1.RespondToInvitationRequest, _ ...grpc.CallOption) (*v1.RespondToInvitationResponse, error) {
	return &v1.RespondToInvitationResponse{}, nil
}
func (m *mockFamilyClient) RemoveMember(_ context.Context, _ *v1.RemoveMemberRequest, _ ...grpc.CallOption) (*v1.RemoveMemberResponse, error) {
	return &v1.RemoveMemberResponse{}, nil
}
func (m *mockFamilyClient) UpdateMemberRole(_ context.Context, _ *v1.UpdateMemberRoleRequest, _ ...grpc.CallOption) (*v1.FamilyMember, error) {
	return &v1.FamilyMember{}, nil
}
func (m *mockFamilyClient) ListFamilyMembers(_ context.Context, _ *v1.ListFamilyMembersRequest, _ ...grpc.CallOption) (*v1.ListFamilyMembersResponse, error) {
	return &v1.ListFamilyMembersResponse{}, nil
}
func (m *mockFamilyClient) SetSharingPreferences(_ context.Context, _ *v1.SetSharingPreferencesRequest, _ ...grpc.CallOption) (*v1.SharingPreferences, error) {
	return &v1.SharingPreferences{}, nil
}
func (m *mockFamilyClient) GetSharedHealthData(_ context.Context, _ *v1.GetSharedHealthDataRequest, _ ...grpc.CallOption) (*v1.GetSharedHealthDataResponse, error) {
	return &v1.GetSharedHealthDataResponse{}, nil
}
func (m *mockFamilyClient) ValidateSharingAccess(_ context.Context, _ *v1.ValidateSharingAccessRequest, _ ...grpc.CallOption) (*v1.ValidateSharingAccessResponse, error) {
	return &v1.ValidateSharingAccessResponse{}, nil
}

// mockHealthRecordClient는 HealthRecordServiceClient를 모킹합니다.
type mockHealthRecordClient struct{}

func (m *mockHealthRecordClient) CreateRecord(_ context.Context, _ *v1.CreateHealthRecordRequest, _ ...grpc.CallOption) (*v1.HealthRecord, error) {
	return &v1.HealthRecord{RecordId: "rec_001"}, nil
}
func (m *mockHealthRecordClient) GetRecord(_ context.Context, _ *v1.GetHealthRecordRequest, _ ...grpc.CallOption) (*v1.HealthRecord, error) {
	return &v1.HealthRecord{RecordId: "rec_001"}, nil
}
func (m *mockHealthRecordClient) ListRecords(_ context.Context, _ *v1.ListHealthRecordsRequest, _ ...grpc.CallOption) (*v1.ListHealthRecordsResponse, error) {
	return &v1.ListHealthRecordsResponse{}, nil
}
func (m *mockHealthRecordClient) UpdateRecord(_ context.Context, _ *v1.UpdateHealthRecordRequest, _ ...grpc.CallOption) (*v1.HealthRecord, error) {
	return &v1.HealthRecord{}, nil
}
func (m *mockHealthRecordClient) DeleteRecord(_ context.Context, _ *v1.DeleteHealthRecordRequest, _ ...grpc.CallOption) (*v1.DeleteHealthRecordResponse, error) {
	return &v1.DeleteHealthRecordResponse{}, nil
}
func (m *mockHealthRecordClient) ExportToFHIR(_ context.Context, _ *v1.ExportToFHIRRequest, _ ...grpc.CallOption) (*v1.ExportToFHIRResponse, error) {
	return &v1.ExportToFHIRResponse{FhirBundleJson: `{"resourceType":"Bundle","entry":[]}`}, nil
}
func (m *mockHealthRecordClient) ImportFromFHIR(_ context.Context, _ *v1.ImportFromFHIRRequest, _ ...grpc.CallOption) (*v1.ImportFromFHIRResponse, error) {
	return &v1.ImportFromFHIRResponse{}, nil
}
func (m *mockHealthRecordClient) GetHealthSummary(_ context.Context, _ *v1.GetHealthSummaryRequest, _ ...grpc.CallOption) (*v1.GetHealthSummaryResponse, error) {
	return &v1.GetHealthSummaryResponse{}, nil
}
func (m *mockHealthRecordClient) CreateDataSharingConsent(_ context.Context, _ *v1.CreateConsentRequest, _ ...grpc.CallOption) (*v1.DataSharingConsent, error) {
	return &v1.DataSharingConsent{}, nil
}
func (m *mockHealthRecordClient) RevokeDataSharingConsent(_ context.Context, _ *v1.RevokeConsentRequest, _ ...grpc.CallOption) (*v1.RevokeConsentResponse, error) {
	return &v1.RevokeConsentResponse{}, nil
}
func (m *mockHealthRecordClient) ListDataSharingConsents(_ context.Context, _ *v1.ListConsentsRequest, _ ...grpc.CallOption) (*v1.ListConsentsResponse, error) {
	return &v1.ListConsentsResponse{}, nil
}
func (m *mockHealthRecordClient) ShareWithProvider(_ context.Context, _ *v1.ShareWithProviderRequest, _ ...grpc.CallOption) (*v1.ShareWithProviderResponse, error) {
	return &v1.ShareWithProviderResponse{}, nil
}
func (m *mockHealthRecordClient) GetDataAccessLog(_ context.Context, _ *v1.GetDataAccessLogRequest, _ ...grpc.CallOption) (*v1.GetDataAccessLogResponse, error) {
	return &v1.GetDataAccessLogResponse{}, nil
}

// mockPrescriptionClient는 PrescriptionServiceClient를 모킹합니다.
type mockPrescriptionClient struct{}

func (m *mockPrescriptionClient) CreatePrescription(_ context.Context, _ *v1.CreatePrescriptionRequest, _ ...grpc.CallOption) (*v1.Prescription, error) {
	return &v1.Prescription{PrescriptionId: "presc_001"}, nil
}
func (m *mockPrescriptionClient) GetPrescription(_ context.Context, _ *v1.GetPrescriptionRequest, _ ...grpc.CallOption) (*v1.Prescription, error) {
	return &v1.Prescription{PrescriptionId: "presc_001"}, nil
}
func (m *mockPrescriptionClient) ListPrescriptions(_ context.Context, _ *v1.ListPrescriptionsRequest, _ ...grpc.CallOption) (*v1.ListPrescriptionsResponse, error) {
	return &v1.ListPrescriptionsResponse{}, nil
}
func (m *mockPrescriptionClient) UpdatePrescriptionStatus(_ context.Context, _ *v1.UpdatePrescriptionStatusRequest, _ ...grpc.CallOption) (*v1.Prescription, error) {
	return &v1.Prescription{}, nil
}
func (m *mockPrescriptionClient) AddMedication(_ context.Context, _ *v1.AddMedicationRequest, _ ...grpc.CallOption) (*v1.Prescription, error) {
	return &v1.Prescription{}, nil
}
func (m *mockPrescriptionClient) RemoveMedication(_ context.Context, _ *v1.RemoveMedicationRequest, _ ...grpc.CallOption) (*v1.Prescription, error) {
	return &v1.Prescription{}, nil
}
func (m *mockPrescriptionClient) CheckDrugInteraction(_ context.Context, _ *v1.CheckDrugInteractionRequest, _ ...grpc.CallOption) (*v1.CheckDrugInteractionResponse, error) {
	return &v1.CheckDrugInteractionResponse{}, nil
}
func (m *mockPrescriptionClient) GetMedicationReminders(_ context.Context, _ *v1.GetMedicationRemindersRequest, _ ...grpc.CallOption) (*v1.GetMedicationRemindersResponse, error) {
	return &v1.GetMedicationRemindersResponse{}, nil
}
func (m *mockPrescriptionClient) SelectPharmacyAndFulfillment(_ context.Context, _ *v1.SelectPharmacyRequest, _ ...grpc.CallOption) (*v1.SelectPharmacyResponse, error) {
	return &v1.SelectPharmacyResponse{Success: true, Message: "pharmacy_selected"}, nil
}
func (m *mockPrescriptionClient) SendPrescriptionToPharmacy(_ context.Context, _ *v1.SendToPharmacyRequest, _ ...grpc.CallOption) (*v1.SendToPharmacyResponse, error) {
	return &v1.SendToPharmacyResponse{}, nil
}
func (m *mockPrescriptionClient) GetPrescriptionByToken(_ context.Context, _ *v1.GetByTokenRequest, _ ...grpc.CallOption) (*v1.Prescription, error) {
	return &v1.Prescription{}, nil
}
func (m *mockPrescriptionClient) UpdateDispensaryStatus(_ context.Context, _ *v1.UpdateDispensaryStatusRequest, _ ...grpc.CallOption) (*v1.Prescription, error) {
	return &v1.Prescription{}, nil
}

// mockReservationClient는 ReservationServiceClient를 모킹합니다.
type mockReservationClient struct{}

func (m *mockReservationClient) SearchFacilities(_ context.Context, _ *v1.SearchFacilitiesRequest, _ ...grpc.CallOption) (*v1.SearchFacilitiesResponse, error) {
	return &v1.SearchFacilitiesResponse{}, nil
}
func (m *mockReservationClient) GetFacility(_ context.Context, _ *v1.GetFacilityRequest, _ ...grpc.CallOption) (*v1.Facility, error) {
	return &v1.Facility{}, nil
}
func (m *mockReservationClient) GetAvailableSlots(_ context.Context, _ *v1.GetAvailableSlotsRequest, _ ...grpc.CallOption) (*v1.GetAvailableSlotsResponse, error) {
	return &v1.GetAvailableSlotsResponse{}, nil
}
func (m *mockReservationClient) CreateReservation(_ context.Context, _ *v1.CreateReservationRequest, _ ...grpc.CallOption) (*v1.Reservation, error) {
	return &v1.Reservation{ReservationId: "rsv_001"}, nil
}
func (m *mockReservationClient) GetReservation(_ context.Context, _ *v1.GetReservationRequest, _ ...grpc.CallOption) (*v1.Reservation, error) {
	return &v1.Reservation{ReservationId: "rsv_001"}, nil
}
func (m *mockReservationClient) ListReservations(_ context.Context, _ *v1.ListReservationsRequest, _ ...grpc.CallOption) (*v1.ListReservationsResponse, error) {
	return &v1.ListReservationsResponse{}, nil
}
func (m *mockReservationClient) CancelReservation(_ context.Context, _ *v1.CancelReservationRequest, _ ...grpc.CallOption) (*v1.CancelReservationResponse, error) {
	return &v1.CancelReservationResponse{}, nil
}
func (m *mockReservationClient) ListDoctorsByFacility(_ context.Context, _ *v1.ListDoctorsByFacilityRequest, _ ...grpc.CallOption) (*v1.ListDoctorsByFacilityResponse, error) {
	return &v1.ListDoctorsByFacilityResponse{}, nil
}
func (m *mockReservationClient) GetDoctorAvailability(_ context.Context, _ *v1.GetDoctorAvailabilityRequest, _ ...grpc.CallOption) (*v1.GetDoctorAvailabilityResponse, error) {
	return &v1.GetDoctorAvailabilityResponse{}, nil
}
func (m *mockReservationClient) SelectDoctor(_ context.Context, _ *v1.SelectDoctorRequest, _ ...grpc.CallOption) (*v1.SelectDoctorResponse, error) {
	return &v1.SelectDoctorResponse{}, nil
}

// mockTelemedicineClient는 TelemedicineServiceClient를 모킹합니다.
type mockTelemedicineClient struct{}

func (m *mockTelemedicineClient) CreateConsultation(_ context.Context, _ *v1.CreateConsultationRequest, _ ...grpc.CallOption) (*v1.Consultation, error) {
	return &v1.Consultation{ConsultationId: "consult_001"}, nil
}
func (m *mockTelemedicineClient) GetConsultation(_ context.Context, _ *v1.GetConsultationRequest, _ ...grpc.CallOption) (*v1.Consultation, error) {
	return &v1.Consultation{ConsultationId: "consult_001"}, nil
}
func (m *mockTelemedicineClient) ListConsultations(_ context.Context, _ *v1.ListConsultationsRequest, _ ...grpc.CallOption) (*v1.ListConsultationsResponse, error) {
	return &v1.ListConsultationsResponse{}, nil
}
func (m *mockTelemedicineClient) MatchDoctor(_ context.Context, _ *v1.MatchDoctorRequest, _ ...grpc.CallOption) (*v1.MatchDoctorResponse, error) {
	return &v1.MatchDoctorResponse{}, nil
}
func (m *mockTelemedicineClient) StartVideoSession(_ context.Context, _ *v1.StartVideoSessionRequest, _ ...grpc.CallOption) (*v1.VideoSession, error) {
	return &v1.VideoSession{}, nil
}
func (m *mockTelemedicineClient) EndVideoSession(_ context.Context, _ *v1.EndVideoSessionRequest, _ ...grpc.CallOption) (*v1.VideoSession, error) {
	return &v1.VideoSession{}, nil
}
func (m *mockTelemedicineClient) RateConsultation(_ context.Context, _ *v1.RateConsultationRequest, _ ...grpc.CallOption) (*v1.RateConsultationResponse, error) {
	return &v1.RateConsultationResponse{}, nil
}

// mockNotificationClient는 NotificationServiceClient를 모킹합니다.
type mockNotificationClient struct{}

func (m *mockNotificationClient) SendNotification(_ context.Context, _ *v1.SendNotificationRequest, _ ...grpc.CallOption) (*v1.Notification, error) {
	return &v1.Notification{NotificationId: "noti_001"}, nil
}
func (m *mockNotificationClient) ListNotifications(_ context.Context, _ *v1.ListNotificationsRequest, _ ...grpc.CallOption) (*v1.ListNotificationsResponse, error) {
	return &v1.ListNotificationsResponse{}, nil
}
func (m *mockNotificationClient) MarkAsRead(_ context.Context, _ *v1.MarkAsReadRequest, _ ...grpc.CallOption) (*v1.MarkAsReadResponse, error) {
	return &v1.MarkAsReadResponse{}, nil
}
func (m *mockNotificationClient) MarkAllAsRead(_ context.Context, _ *v1.MarkAllAsReadRequest, _ ...grpc.CallOption) (*v1.MarkAllAsReadResponse, error) {
	return &v1.MarkAllAsReadResponse{}, nil
}
func (m *mockNotificationClient) GetUnreadCount(_ context.Context, _ *v1.GetUnreadCountRequest, _ ...grpc.CallOption) (*v1.GetUnreadCountResponse, error) {
	return &v1.GetUnreadCountResponse{}, nil
}
func (m *mockNotificationClient) UpdateNotificationPreferences(_ context.Context, _ *v1.UpdateNotificationPreferencesRequest, _ ...grpc.CallOption) (*v1.NotificationPreferences, error) {
	return &v1.NotificationPreferences{}, nil
}
func (m *mockNotificationClient) GetNotificationPreferences(_ context.Context, _ *v1.GetNotificationPreferencesRequest, _ ...grpc.CallOption) (*v1.NotificationPreferences, error) {
	return &v1.NotificationPreferences{}, nil
}
func (m *mockNotificationClient) SendFromTemplate(_ context.Context, _ *v1.SendFromTemplateRequest, _ ...grpc.CallOption) (*v1.Notification, error) {
	return &v1.Notification{}, nil
}

// mockTranslationClient는 TranslationServiceClient를 모킹합니다.
type mockTranslationClient struct{}

func (m *mockTranslationClient) TranslateText(_ context.Context, _ *v1.TranslateTextRequest, _ ...grpc.CallOption) (*v1.TranslateTextResponse, error) {
	return &v1.TranslateTextResponse{TranslatedText: "translated"}, nil
}
func (m *mockTranslationClient) DetectLanguage(_ context.Context, _ *v1.DetectLanguageRequest, _ ...grpc.CallOption) (*v1.DetectLanguageResponse, error) {
	return &v1.DetectLanguageResponse{Languages: []*v1.DetectedLanguage{{LanguageCode: "ko", Confidence: 0.99}}}, nil
}
func (m *mockTranslationClient) ListSupportedLanguages(_ context.Context, _ *v1.ListSupportedLanguagesRequest, _ ...grpc.CallOption) (*v1.ListSupportedLanguagesResponse, error) {
	return &v1.ListSupportedLanguagesResponse{}, nil
}
func (m *mockTranslationClient) TranslateBatch(_ context.Context, _ *v1.TranslateBatchRequest, _ ...grpc.CallOption) (*v1.TranslateBatchResponse, error) {
	return &v1.TranslateBatchResponse{}, nil
}
func (m *mockTranslationClient) GetTranslationHistory(_ context.Context, _ *v1.GetTranslationHistoryRequest, _ ...grpc.CallOption) (*v1.GetTranslationHistoryResponse, error) {
	return &v1.GetTranslationHistoryResponse{}, nil
}
func (m *mockTranslationClient) GetTranslationUsage(_ context.Context, _ *v1.GetTranslationUsageRequest, _ ...grpc.CallOption) (*v1.GetTranslationUsageResponse, error) {
	return &v1.GetTranslationUsageResponse{}, nil
}
func (m *mockTranslationClient) TranslateRealtime(_ context.Context, _ *v1.TranslateRealtimeRequest, _ ...grpc.CallOption) (*v1.TranslateRealtimeResponse, error) {
	return &v1.TranslateRealtimeResponse{}, nil
}

// mockCoachingClient는 CoachingServiceClient를 모킹합니다.
type mockCoachingClient struct{}

func (m *mockCoachingClient) SetHealthGoal(_ context.Context, _ *v1.SetHealthGoalRequest, _ ...grpc.CallOption) (*v1.HealthGoal, error) {
	return &v1.HealthGoal{GoalId: "goal_001"}, nil
}
func (m *mockCoachingClient) GetHealthGoals(_ context.Context, _ *v1.GetHealthGoalsRequest, _ ...grpc.CallOption) (*v1.GetHealthGoalsResponse, error) {
	return &v1.GetHealthGoalsResponse{}, nil
}
func (m *mockCoachingClient) GenerateCoaching(_ context.Context, _ *v1.GenerateCoachingRequest, _ ...grpc.CallOption) (*v1.CoachingMessage, error) {
	return &v1.CoachingMessage{}, nil
}
func (m *mockCoachingClient) ListCoachingMessages(_ context.Context, _ *v1.ListCoachingMessagesRequest, _ ...grpc.CallOption) (*v1.ListCoachingMessagesResponse, error) {
	return &v1.ListCoachingMessagesResponse{}, nil
}
func (m *mockCoachingClient) GenerateDailyReport(_ context.Context, _ *v1.GenerateDailyReportRequest, _ ...grpc.CallOption) (*v1.DailyHealthReport, error) {
	return &v1.DailyHealthReport{}, nil
}
func (m *mockCoachingClient) GetWeeklyReport(_ context.Context, _ *v1.GetWeeklyReportRequest, _ ...grpc.CallOption) (*v1.WeeklyHealthReport, error) {
	return &v1.WeeklyHealthReport{}, nil
}
func (m *mockCoachingClient) GetRecommendations(_ context.Context, _ *v1.GetRecommendationsRequest, _ ...grpc.CallOption) (*v1.GetRecommendationsResponse, error) {
	return &v1.GetRecommendationsResponse{}, nil
}

// mockAdminClient는 AdminServiceClient를 모킹합니다.
type mockAdminClient struct{}

func (m *mockAdminClient) CreateAdmin(_ context.Context, _ *v1.CreateAdminRequest, _ ...grpc.CallOption) (*v1.AdminUser, error) {
	return &v1.AdminUser{AdminId: "admin_001"}, nil
}
func (m *mockAdminClient) GetAdmin(_ context.Context, _ *v1.GetAdminRequest, _ ...grpc.CallOption) (*v1.AdminUser, error) {
	return &v1.AdminUser{AdminId: "admin_001"}, nil
}
func (m *mockAdminClient) ListAdmins(_ context.Context, _ *v1.ListAdminsRequest, _ ...grpc.CallOption) (*v1.ListAdminsResponse, error) {
	return &v1.ListAdminsResponse{}, nil
}
func (m *mockAdminClient) UpdateAdminRole(_ context.Context, _ *v1.UpdateAdminRoleRequest, _ ...grpc.CallOption) (*v1.AdminUser, error) {
	return &v1.AdminUser{}, nil
}
func (m *mockAdminClient) DeactivateAdmin(_ context.Context, _ *v1.DeactivateAdminRequest, _ ...grpc.CallOption) (*v1.AdminUser, error) {
	return &v1.AdminUser{}, nil
}
func (m *mockAdminClient) ListUsers(_ context.Context, _ *v1.AdminListUsersRequest, _ ...grpc.CallOption) (*v1.AdminListUsersResponse, error) {
	return &v1.AdminListUsersResponse{}, nil
}
func (m *mockAdminClient) GetSystemStats(_ context.Context, _ *v1.GetSystemStatsRequest, _ ...grpc.CallOption) (*v1.GetSystemStatsResponse, error) {
	return &v1.GetSystemStatsResponse{}, nil
}
func (m *mockAdminClient) GetAuditLog(_ context.Context, _ *v1.GetAuditLogRequest, _ ...grpc.CallOption) (*v1.GetAuditLogResponse, error) {
	return &v1.GetAuditLogResponse{}, nil
}
func (m *mockAdminClient) SetSystemConfig(_ context.Context, _ *v1.SetSystemConfigRequest, _ ...grpc.CallOption) (*v1.SystemConfig, error) {
	return &v1.SystemConfig{}, nil
}
func (m *mockAdminClient) GetSystemConfig(_ context.Context, _ *v1.GetSystemConfigRequest, _ ...grpc.CallOption) (*v1.SystemConfig, error) {
	return &v1.SystemConfig{}, nil
}
func (m *mockAdminClient) ListAdminsByRegion(_ context.Context, _ *v1.ListAdminsByRegionRequest, _ ...grpc.CallOption) (*v1.ListAdminsResponse, error) {
	return &v1.ListAdminsResponse{}, nil
}
func (m *mockAdminClient) ListSystemConfigs(_ context.Context, _ *v1.ListSystemConfigsRequest, _ ...grpc.CallOption) (*v1.ListSystemConfigsResponse, error) {
	return &v1.ListSystemConfigsResponse{}, nil
}
func (m *mockAdminClient) GetConfigWithMeta(_ context.Context, _ *v1.GetConfigWithMetaRequest, _ ...grpc.CallOption) (*v1.ConfigWithMeta, error) {
	return &v1.ConfigWithMeta{}, nil
}
func (m *mockAdminClient) ValidateConfigValue(_ context.Context, _ *v1.ValidateConfigValueRequest, _ ...grpc.CallOption) (*v1.ValidateConfigValueResponse, error) {
	return &v1.ValidateConfigValueResponse{}, nil
}
func (m *mockAdminClient) BulkSetConfigs(_ context.Context, _ *v1.BulkSetConfigsRequest, _ ...grpc.CallOption) (*v1.BulkSetConfigsResponse, error) {
	return &v1.BulkSetConfigsResponse{}, nil
}
func (m *mockAdminClient) GetAuditLogDetails(_ context.Context, _ *v1.GetAuditLogDetailsRequest, _ ...grpc.CallOption) (*v1.GetAuditLogDetailsResponse, error) {
	return &v1.GetAuditLogDetailsResponse{}, nil
}
func (m *mockAdminClient) GetRevenueStats(_ context.Context, _ *v1.GetRevenueStatsRequest, _ ...grpc.CallOption) (*v1.GetRevenueStatsResponse, error) {
	return &v1.GetRevenueStatsResponse{}, nil
}
func (m *mockAdminClient) GetInventoryStats(_ context.Context, _ *v1.GetInventoryStatsRequest, _ ...grpc.CallOption) (*v1.GetInventoryStatsResponse, error) {
	return &v1.GetInventoryStatsResponse{}, nil
}

// ============================================================================
// 테스트 헬퍼
// ============================================================================

// newMockHandler는 모든 gRPC mock 클라이언트를 주입한 RestHandler를 생성합니다.
func newMockHandler() *RestHandler {
	return &RestHandler{
		auth:         &mockAuthClient{},
		user:         nil, // 의도적 nil — 시나리오 #5 (인증 없이 접근) 테스트용
		measurement:  &mockMeasurementClient{},
		device:       &mockDeviceClient{},
		subscription: &mockSubscriptionClient{},
		shop:         &mockShopClient{},
		payment:      nil,
		aiInference:  nil,
		cartridge:    nil,
		calibration:  nil,
		coaching:     &mockCoachingClient{},
		reservation:  &mockReservationClient{},
		admin:        &mockAdminClient{},
		family:       &mockFamilyClient{},
		healthRecord: &mockHealthRecordClient{},
		prescription: &mockPrescriptionClient{},
		community:    &mockCommunityClient{},
		video:        nil,
		notification: &mockNotificationClient{},
		translation:  &mockTranslationClient{},
		telemedicine: &mockTelemedicineClient{},
	}
}

// e2eSetup은 mock 핸들러와 mux를 한번에 만들어 반환합니다.
func e2eSetup() *http.ServeMux {
	h := newMockHandler()
	return h.SetupRoutes()
}

// doRequest는 HTTP 요청을 수행하고 상태 코드와 응답 바디를 반환합니다.
func doRequest(t *testing.T, mux *http.ServeMux, method, path, body string) (int, []byte) {
	t.Helper()
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reqBody)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w.Code, w.Body.Bytes()
}

// assertStatusCode는 기대하는 상태 코드와 비교합니다.
func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("HTTP 상태 코드 = %d, 기대값 = %d", got, want)
	}
}

// assertJSONKey는 JSON 응답에 특정 키가 존재하는지 확인합니다.
func assertJSONKey(t *testing.T, body []byte, key string) {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("JSON 파싱 실패: %v (body: %s)", err, string(body))
	}
	if _, ok := m[key]; !ok {
		t.Errorf("응답 JSON에 키 %q가 없습니다. 응답: %s", key, string(body))
	}
}

// assertJSONHasError는 JSON 응답에 "error" 키가 있는지 확인합니다.
func assertJSONHasError(t *testing.T, body []byte) {
	t.Helper()
	assertJSONKey(t, body, "error")
}

// ============================================================================
// E2E 테스트 시나리오 30건
// ============================================================================

// ---------------------------------------------------------------------------
// Auth (5건)
// ---------------------------------------------------------------------------

// E2E #1: POST /api/v1/auth/register → 201 + userId/email/displayName
func TestE2E_AuthRegister(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/auth/register",
		`{"email":"test@example.com","password":"pass123","display_name":"Tester"}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "userId")
	assertJSONKey(t, body, "email")
	assertJSONKey(t, body, "displayName")
}

// E2E #2: POST /api/v1/auth/login → 200 + tokens
func TestE2E_AuthLogin(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/auth/login",
		`{"email":"test@example.com","password":"pass123"}`)
	assertStatusCode(t, code, http.StatusOK)
	assertJSONKey(t, body, "accessToken")
	assertJSONKey(t, body, "refreshToken")
}

// E2E #3: POST /api/v1/auth/refresh → 200 + newAccessToken
func TestE2E_AuthRefresh(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/auth/refresh",
		`{"refresh_token":"rt_old_token"}`)
	assertStatusCode(t, code, http.StatusOK)
	assertJSONKey(t, body, "accessToken")
}

// E2E #4: POST /api/v1/auth/social-login → 200 + tokens
func TestE2E_AuthSocialLogin(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/auth/social-login",
		`{"provider":"GOOGLE","id_token":"google_id_token_123"}`)
	assertStatusCode(t, code, http.StatusOK)
	assertJSONKey(t, body, "accessToken")
	assertJSONKey(t, body, "refreshToken")
}

// E2E #5: GET /api/v1/users/{userId}/profile (user 클라이언트 nil → 인증없이) → 503
func TestE2E_AuthMe_Unauthenticated(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "GET", "/api/v1/users/user_001/profile", "")
	assertStatusCode(t, code, http.StatusServiceUnavailable)
	assertJSONHasError(t, body)
}

// ---------------------------------------------------------------------------
// 측정/디바이스 (4건)
// ---------------------------------------------------------------------------

// E2E #6: GET /api/v1/devices → 200 + array
func TestE2E_ListDevices(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/devices?user_id=user_001", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #7: POST /api/v1/measurements/sessions → 201 + sessionId
func TestE2E_StartMeasurementSession(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/measurements/sessions",
		`{"device_id":"dev_001","user_id":"user_001","cartridge_id":"cart_001"}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "sessionId")
}

// E2E #8: POST /api/v1/measurements/sessions/{id}/end → 200 + summary
func TestE2E_EndMeasurementSession(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/measurements/sessions/sess_001/end", "")
	assertStatusCode(t, code, http.StatusOK)
	assertJSONKey(t, body, "sessionId")
}

// E2E #9: GET /api/v1/measurements/history → 200 + items
func TestE2E_GetMeasurementHistory(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/measurements/history?user_id=user_001", "")
	assertStatusCode(t, code, http.StatusOK)
}

// ---------------------------------------------------------------------------
// 마켓/결제 (3건)
// ---------------------------------------------------------------------------

// E2E #10: GET /api/v1/products → 200 + products
func TestE2E_ListProducts(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/products", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #11: POST /api/v1/orders → 201 + orderId
func TestE2E_CreateOrder(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/orders",
		`{"user_id":"user_001","shipping_address":"Seoul","payment_method":"card"}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "orderId")
}

// E2E #12: GET /api/v1/subscriptions/plans → 200 + plans
func TestE2E_ListSubscriptionPlans(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/subscriptions/plans", "")
	assertStatusCode(t, code, http.StatusOK)
}

// ---------------------------------------------------------------------------
// 건강기록 (3건)
// ---------------------------------------------------------------------------

// E2E #13: GET /api/v1/health-records → 200 + records
func TestE2E_ListHealthRecords(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/health-records?user_id=user_001", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #14: POST /api/v1/health-records → 201 + record
func TestE2E_CreateHealthRecord(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/health-records",
		`{"user_id":"user_001","record_type":1,"title":"혈압 기록","description":"정상 범위"}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "recordId")
}

// E2E #15: POST /api/v1/health-records/export/fhir → 200 + FHIR bundle
func TestE2E_ExportHealthRecordFHIR(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/health-records/export/fhir",
		`{"user_id":"user_001","record_ids":["rec_001"],"target_type":1}`)
	assertStatusCode(t, code, http.StatusOK)
	assertJSONKey(t, body, "fhirBundleJson")
}

// ---------------------------------------------------------------------------
// 처방 (2건)
// ---------------------------------------------------------------------------

// E2E #16: GET /api/v1/prescriptions → 200 + prescriptions
func TestE2E_ListPrescriptions(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/prescriptions?user_id=user_001", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #17: POST /api/v1/prescriptions/{id}/pharmacy → 200 + success/message
func TestE2E_SelectPharmacy(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/prescriptions/presc_001/pharmacy",
		`{"pharmacy_id":"pharm_001","pharmacy_name":"건강약국","fulfillment_type":"pickup"}`)
	assertStatusCode(t, code, http.StatusOK)
	assertJSONKey(t, body, "success")
	assertJSONKey(t, body, "message")
}

// ---------------------------------------------------------------------------
// 진료/예약 (3건)
// ---------------------------------------------------------------------------

// E2E #18: GET /api/v1/facilities → 200 + facilities
func TestE2E_SearchFacilities(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/facilities?query=병원&latitude=37.5665&longitude=126.978", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #19: POST /api/v1/reservations → 201 + reservation
func TestE2E_CreateReservation(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/reservations",
		`{"user_id":"user_001","facility_id":"fac_001","slot_id":"slot_001","doctor_id":"doc_001","reason":"진료"}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "reservationId")
}

// E2E #20: GET /api/v1/telemedicine/doctors → 200 + doctors
func TestE2E_SearchDoctors(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/telemedicine/doctors?specialty=1", "")
	assertStatusCode(t, code, http.StatusOK)
}

// ---------------------------------------------------------------------------
// 커뮤니티 (2건)
// ---------------------------------------------------------------------------

// E2E #21: GET /api/v1/posts → 200 + posts
func TestE2E_ListPosts(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/posts", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #22: POST /api/v1/posts → 201 + postId
func TestE2E_CreatePost(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/posts",
		`{"author_id":"user_001","title":"건강 꿀팁","content":"매일 30분 산책하세요!","category":1}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "postId")
}

// ---------------------------------------------------------------------------
// 가족 (1건)
// ---------------------------------------------------------------------------

// E2E #23: POST /api/v1/family/groups → 201 + groupId
func TestE2E_CreateFamilyGroup(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/family/groups",
		`{"owner_user_id":"user_001","name":"우리 가족","description":"가족 건강 관리 그룹"}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "groupId")
}

// ---------------------------------------------------------------------------
// 알림/번역 (2건)
// ---------------------------------------------------------------------------

// E2E #24: GET /api/v1/notifications → 200 + notifications
func TestE2E_ListNotifications(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/notifications?user_id=user_001", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #25: POST /api/v1/translations/detect → 200 + languages
func TestE2E_DetectLanguage(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/translations/detect",
		`{"text":"안녕하세요 건강 관리 앱입니다"}`)
	assertStatusCode(t, code, http.StatusOK)
	assertJSONKey(t, body, "languages")
}

// ---------------------------------------------------------------------------
// AI/코칭 (2건)
// ---------------------------------------------------------------------------

// E2E #26: GET /api/v1/coaching/goals → 200 + reports(goals)
func TestE2E_GetCoachingGoals(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/coaching/goals?user_id=user_001", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #27: POST /api/v1/coaching/goals → 201 + goal
func TestE2E_SetHealthGoal(t *testing.T) {
	mux := e2eSetup()
	code, body := doRequest(t, mux, "POST", "/api/v1/coaching/goals",
		`{"user_id":"user_001","metric_name":"steps","target_value":10000,"unit":"steps","description":"만보 걷기"}`)
	assertStatusCode(t, code, http.StatusCreated)
	assertJSONKey(t, body, "goalId")
}

// ---------------------------------------------------------------------------
// 관리자 (2건)
// ---------------------------------------------------------------------------

// E2E #28: GET /api/v1/admin/users → 200 + users
func TestE2E_AdminListUsers(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/admin/users?limit=10", "")
	assertStatusCode(t, code, http.StatusOK)
}

// E2E #29: GET /api/v1/admin/audit-log → 200 + logs
func TestE2E_AdminAuditLog(t *testing.T) {
	mux := e2eSetup()
	code, _ := doRequest(t, mux, "GET", "/api/v1/admin/audit-log?limit=20", "")
	assertStatusCode(t, code, http.StatusOK)
}

// ---------------------------------------------------------------------------
// 부하 테스트 (1건)
// ---------------------------------------------------------------------------

// E2E #30: 동시 10건 요청 전송 → 모두 200 (goroutine)
func TestE2E_ConcurrentRequests(t *testing.T) {
	mux := e2eSetup()

	const concurrency = 10
	var wg sync.WaitGroup
	wg.Add(concurrency)

	results := make([]int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			results[idx] = w.Code
		}(i)
	}

	wg.Wait()

	for i, code := range results {
		if code != http.StatusOK {
			t.Errorf("동시 요청 #%d: 상태 코드 = %d, 기대값 = 200", i, code)
		}
	}
}
