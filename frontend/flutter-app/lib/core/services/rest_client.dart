import 'package:dio/dio.dart';

/// ManPaSik REST Gateway Client
///
/// Provides typed methods for all REST API endpoints exposed by the gateway.
/// Uses Dio for HTTP with configurable base URL, timeouts, and auth tokens.
///
/// Usage:
/// ```dart
/// final client = ManPaSikRestClient(baseUrl: 'http://10.0.2.2:8080/api/v1');
/// final loginResp = await client.login('user@test.com', 'password');
/// client.setAuthToken(loginResp['access_token']);
/// ```
class ManPaSikRestClient {
  ManPaSikRestClient({String? baseUrl})
      : _dio = Dio(BaseOptions(
          baseUrl: baseUrl ?? 'http://localhost:8080/api/v1',
          connectTimeout: const Duration(seconds: 10),
          receiveTimeout: const Duration(seconds: 30),
          headers: {'Content-Type': 'application/json'},
        ));

  final Dio _dio;

  /// Expose the Dio instance for advanced configuration (interceptors, etc.)
  Dio get dio => _dio;

  /// Set the Authorization bearer token for authenticated requests.
  void setAuthToken(String token) {
    _dio.options.headers['Authorization'] = 'Bearer $token';
  }

  /// Clear the auth token (e.g. on logout).
  void clearAuthToken() {
    _dio.options.headers.remove('Authorization');
  }

  // ==========================================================================
  // Auth Service
  // ==========================================================================

  /// Register a new user account.
  Future<Map<String, dynamic>> register(
    String email,
    String password,
    String displayName,
  ) async {
    final resp = await _dio.post('/auth/register', data: {
      'email': email,
      'password': password,
      'display_name': displayName,
    });
    return _asMap(resp);
  }

  /// Login with email and password. Returns access/refresh tokens.
  Future<Map<String, dynamic>> login(String email, String password) async {
    final resp = await _dio.post('/auth/login', data: {
      'email': email,
      'password': password,
    });
    return _asMap(resp);
  }

  /// Refresh an access token using a refresh token.
  Future<Map<String, dynamic>> refreshToken(String refreshToken) async {
    final resp = await _dio.post('/auth/refresh', data: {
      'refresh_token': refreshToken,
    });
    return _asMap(resp);
  }

  /// Logout a user.
  Future<Map<String, dynamic>> logout(String userId) async {
    final resp = await _dio.post('/auth/logout', data: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// Request a password reset. Sends a reset code to the given email.
  Future<Map<String, dynamic>> resetPassword(String email) async {
    final resp = await _dio.post('/auth/reset-password', data: {
      'email': email,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // User Service
  // ==========================================================================

  /// Get a user's profile.
  Future<Map<String, dynamic>> getProfile(String userId) async {
    final resp = await _dio.get('/users/$userId/profile');
    return _asMap(resp);
  }

  /// Update a user's profile.
  Future<Map<String, dynamic>> updateProfile(
    String userId, {
    String? displayName,
    String? avatarUrl,
    String? language,
    String? timezone,
  }) async {
    final resp = await _dio.put('/users/$userId/profile', data: {
      if (displayName != null) 'display_name': displayName,
      if (avatarUrl != null) 'avatar_url': avatarUrl,
      if (language != null) 'language': language,
      if (timezone != null) 'timezone': timezone,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Measurement Service
  // ==========================================================================

  /// Start a new measurement session.
  Future<Map<String, dynamic>> startSession({
    required String deviceId,
    required String userId,
    String? cartridgeId,
    int? cartridgeCategory,
    int? cartridgeTypeIndex,
  }) async {
    final resp = await _dio.post('/measurements/sessions', data: {
      'device_id': deviceId,
      'user_id': userId,
      if (cartridgeId != null) 'cartridge_id': cartridgeId,
      if (cartridgeCategory != null) 'cartridge_category': cartridgeCategory,
      if (cartridgeTypeIndex != null)
        'cartridge_type_index': cartridgeTypeIndex,
    });
    return _asMap(resp);
  }

  /// End a measurement session.
  Future<Map<String, dynamic>> endSession(String sessionId) async {
    final resp =
        await _dio.post('/measurements/sessions/$sessionId/end');
    return _asMap(resp);
  }

  /// Get measurement history for a user.
  Future<Map<String, dynamic>> getMeasurementHistory(
    String userId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/measurements/history', queryParameters: {
      'user_id': userId,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Device Service
  // ==========================================================================

  /// Register a new device.
  Future<Map<String, dynamic>> registerDevice({
    required String deviceId,
    required String userId,
    String? serialNumber,
    String? firmwareVersion,
  }) async {
    final resp = await _dio.post('/devices', data: {
      'device_id': deviceId,
      'user_id': userId,
      if (serialNumber != null) 'serial_number': serialNumber,
      if (firmwareVersion != null) 'firmware_version': firmwareVersion,
    });
    return _asMap(resp);
  }

  /// List devices for a user.
  Future<Map<String, dynamic>> listDevices(String userId) async {
    final resp = await _dio.get('/devices', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Reservation Service (Facilities + Reservations)
  // ==========================================================================

  /// Search facilities by location and query.
  Future<Map<String, dynamic>> searchFacilities({
    String? query,
    String? type,
    double? latitude,
    double? longitude,
    double? radiusKm,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/facilities', queryParameters: {
      if (query != null) 'query': query,
      if (type != null) 'type': type,
      if (latitude != null) 'latitude': latitude,
      if (longitude != null) 'longitude': longitude,
      if (radiusKm != null) 'radius_km': radiusKm,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get a single facility by ID.
  Future<Map<String, dynamic>> getFacility(String facilityId) async {
    final resp = await _dio.get('/facilities/$facilityId');
    return _asMap(resp);
  }

  /// Create a reservation.
  Future<Map<String, dynamic>> createReservation({
    required String userId,
    required String facilityId,
    String? slotId,
    String? doctorId,
    int? specialty,
    String? reason,
    String? notes,
  }) async {
    final resp = await _dio.post('/reservations', data: {
      'user_id': userId,
      'facility_id': facilityId,
      if (slotId != null) 'slot_id': slotId,
      if (doctorId != null) 'doctor_id': doctorId,
      if (specialty != null) 'specialty': specialty,
      if (reason != null) 'reason': reason,
      if (notes != null) 'notes': notes,
    });
    return _asMap(resp);
  }

  /// List reservations for a user.
  Future<Map<String, dynamic>> listReservations(
    String userId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/reservations', queryParameters: {
      'user_id': userId,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get a single reservation.
  Future<Map<String, dynamic>> getReservation(String reservationId) async {
    final resp = await _dio.get('/reservations/$reservationId');
    return _asMap(resp);
  }

  // ==========================================================================
  // Prescription Service
  // ==========================================================================

  /// Select a pharmacy for a prescription.
  Future<Map<String, dynamic>> selectPharmacy(
    String prescriptionId, {
    required String pharmacyId,
    String? pharmacyName,
    String? fulfillmentType,
    String? shippingAddress,
  }) async {
    final resp =
        await _dio.post('/prescriptions/$prescriptionId/pharmacy', data: {
      'pharmacy_id': pharmacyId,
      if (pharmacyName != null) 'pharmacy_name': pharmacyName,
      if (fulfillmentType != null) 'fulfillment_type': fulfillmentType,
      if (shippingAddress != null) 'shipping_address': shippingAddress,
    });
    return _asMap(resp);
  }

  /// Send prescription to pharmacy.
  Future<Map<String, dynamic>> sendToPharmacy(String prescriptionId) async {
    final resp =
        await _dio.post('/prescriptions/$prescriptionId/send');
    return _asMap(resp);
  }

  /// Get prescription by fulfillment token.
  Future<Map<String, dynamic>> getPrescriptionByToken(String token) async {
    final resp = await _dio.get('/prescriptions/token/$token');
    return _asMap(resp);
  }

  // ==========================================================================
  // Subscription Service
  // ==========================================================================

  /// List available subscription plans.
  Future<Map<String, dynamic>> listSubscriptionPlans() async {
    final resp = await _dio.get('/subscriptions/plans');
    return _asMap(resp);
  }

  /// Get subscription details for a user.
  Future<Map<String, dynamic>> getSubscription(String userId) async {
    final resp = await _dio.get('/subscriptions/$userId');
    return _asMap(resp);
  }

  /// Create a subscription.
  Future<Map<String, dynamic>> createSubscription({
    required String userId,
    required int tier,
  }) async {
    final resp = await _dio.post('/subscriptions', data: {
      'user_id': userId,
      'tier': tier,
    });
    return _asMap(resp);
  }

  /// Cancel a subscription.
  Future<Map<String, dynamic>> cancelSubscription(
    String subscriptionId, {
    String? userId,
    String? reason,
  }) async {
    final resp = await _dio.delete('/subscriptions/$subscriptionId', data: {
      if (userId != null) 'user_id': userId,
      if (reason != null) 'reason': reason,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Shop Service
  // ==========================================================================

  /// List products with optional category filter.
  Future<Map<String, dynamic>> listProducts({
    int? category,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/products', queryParameters: {
      if (category != null) 'category': category,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get a single product.
  Future<Map<String, dynamic>> getProduct(String productId) async {
    final resp = await _dio.get('/products/$productId');
    return _asMap(resp);
  }

  /// Add an item to the cart.
  Future<Map<String, dynamic>> addToCart({
    required String userId,
    required String productId,
    int quantity = 1,
  }) async {
    final resp = await _dio.post('/cart', data: {
      'user_id': userId,
      'product_id': productId,
      'quantity': quantity,
    });
    return _asMap(resp);
  }

  /// Get a user's cart.
  Future<Map<String, dynamic>> getCart(String userId) async {
    final resp = await _dio.get('/cart/$userId');
    return _asMap(resp);
  }

  /// Create an order.
  Future<Map<String, dynamic>> createOrder({
    required String userId,
    String? shippingAddress,
    String? paymentMethod,
  }) async {
    final resp = await _dio.post('/orders', data: {
      'user_id': userId,
      if (shippingAddress != null) 'shipping_address': shippingAddress,
      if (paymentMethod != null) 'payment_method': paymentMethod,
    });
    return _asMap(resp);
  }

  /// List orders for a user.
  Future<Map<String, dynamic>> listOrders(
    String userId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/orders', queryParameters: {
      'user_id': userId,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Product Reviews
  // ==========================================================================

  /// Get reviews for a product.
  Future<Map<String, dynamic>> getProductReviews(
    String productId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/products/$productId/reviews', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Create a review for a product.
  Future<Map<String, dynamic>> createProductReview({
    required String productId,
    required String userId,
    required int rating,
    required String content,
  }) async {
    final resp = await _dio.post('/products/$productId/reviews', data: {
      'user_id': userId,
      'rating': rating,
      'content': content,
    });
    return _asMap(resp);
  }

  /// Upload avatar image for a user.
  Future<Map<String, dynamic>> uploadAvatar(String userId, List<int> imageBytes, String filename) async {
    final formData = FormData.fromMap({
      'file': MultipartFile.fromBytes(imageBytes, filename: filename),
    });
    final resp = await _dio.post('/users/$userId/avatar', data: formData);
    return _asMap(resp);
  }

  /// Change admin role for a user.
  Future<Map<String, dynamic>> adminChangeRole(String userId, String newRole) async {
    final resp = await _dio.put('/admin/users/$userId/role', data: {
      'role': newRole,
    });
    return _asMap(resp);
  }

  /// Bulk admin action on users.
  Future<Map<String, dynamic>> adminBulkAction({
    required List<String> userIds,
    required String action,
  }) async {
    final resp = await _dio.post('/admin/users/bulk', data: {
      'user_ids': userIds,
      'action': action,
    });
    return _asMap(resp);
  }

  /// Translate text.
  Future<Map<String, dynamic>> translateText({
    required String text,
    required String sourceLanguage,
    required String targetLanguage,
  }) async {
    final resp = await _dio.post('/translations/translate', data: {
      'text': text,
      'source_language': sourceLanguage,
      'target_language': targetLanguage,
    });
    return _asMap(resp);
  }

  /// Save emergency settings.
  Future<Map<String, dynamic>> saveEmergencySettings({
    required String userId,
    required bool autoReport119,
    required List<String> emergencyContacts,
  }) async {
    final resp = await _dio.put('/users/$userId/emergency-settings', data: {
      'auto_report_119': autoReport119,
      'emergency_contacts': emergencyContacts,
    });
    return _asMap(resp);
  }

  /// 긴급 위치 공유 — 보호자에게 GPS 위치를 전송합니다.
  Future<Map<String, dynamic>> shareEmergencyLocation({
    required String userId,
    required double latitude,
    required double longitude,
    required List<String> contactPhones,
  }) async {
    final resp = await _dio.post('/users/$userId/emergency-location', data: {
      'latitude': latitude,
      'longitude': longitude,
      'contact_phones': contactPhones,
      'timestamp': DateTime.now().toIso8601String(),
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Telemedicine Service
  // ==========================================================================

  /// Search doctors by specialty.
  Future<Map<String, dynamic>> searchDoctors({
    String? specialty,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/telemedicine/doctors', queryParameters: {
      if (specialty != null) 'specialty': specialty,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Create a telemedicine consultation.
  Future<Map<String, dynamic>> createConsultation({
    required String userId,
    required String doctorId,
    required String specialty,
    String? reason,
  }) async {
    final resp = await _dio.post('/telemedicine/consultations', data: {
      'user_id': userId,
      'doctor_id': doctorId,
      'specialty': specialty,
      if (reason != null) 'reason': reason,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Family Service
  // ==========================================================================

  /// List family groups for a user.
  Future<Map<String, dynamic>> listFamilyGroups(String userId) async {
    final resp = await _dio.get('/family/groups', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// Get family group health report.
  Future<Map<String, dynamic>> getFamilyGroupReport(String groupId) async {
    final resp = await _dio.get('/family/groups/$groupId/report');
    return _asMap(resp);
  }

  // ==========================================================================
  // Social Auth
  // ==========================================================================

  /// Social login (Google/Apple OAuth).
  Future<Map<String, dynamic>> socialLogin({
    required String provider,
    required String idToken,
  }) async {
    final resp = await _dio.post('/auth/social-login', data: {
      'provider': provider,
      'id_token': idToken,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Payment Service
  // ==========================================================================

  /// Create a payment.
  Future<Map<String, dynamic>> createPayment({
    required String userId,
    String? orderId,
    String? subscriptionId,
    int? paymentType,
    int? amountKrw,
    String? paymentMethod,
  }) async {
    final resp = await _dio.post('/payments', data: {
      'user_id': userId,
      if (orderId != null) 'order_id': orderId,
      if (subscriptionId != null) 'subscription_id': subscriptionId,
      if (paymentType != null) 'payment_type': paymentType,
      if (amountKrw != null) 'amount_krw': amountKrw,
      if (paymentMethod != null) 'payment_method': paymentMethod,
    });
    return _asMap(resp);
  }

  /// Confirm a payment after PG processing.
  Future<Map<String, dynamic>> confirmPayment(
    String paymentId, {
    required String pgTransactionId,
    String? pgProvider,
  }) async {
    final resp = await _dio.post('/payments/$paymentId/confirm', data: {
      'pg_transaction_id': pgTransactionId,
      if (pgProvider != null) 'pg_provider': pgProvider,
    });
    return _asMap(resp);
  }

  /// Get payment details.
  Future<Map<String, dynamic>> getPayment(String paymentId) async {
    final resp = await _dio.get('/payments/$paymentId');
    return _asMap(resp);
  }

  /// Refund a payment.
  Future<Map<String, dynamic>> refundPayment(
    String paymentId, {
    String? reason,
    int? refundAmountKrw,
  }) async {
    final resp = await _dio.post('/payments/$paymentId/refund', data: {
      if (reason != null) 'reason': reason,
      if (refundAmountKrw != null) 'refund_amount_krw': refundAmountKrw,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Health Record Service
  // ==========================================================================

  /// Create a health record.
  Future<Map<String, dynamic>> createHealthRecord({
    required String userId,
    int? recordType,
    String? title,
    String? description,
    String? provider,
    Map<String, String>? metadata,
    String? measurementId,
  }) async {
    final resp = await _dio.post('/health-records', data: {
      'user_id': userId,
      if (recordType != null) 'record_type': recordType,
      if (title != null) 'title': title,
      if (description != null) 'description': description,
      if (provider != null) 'provider': provider,
      if (metadata != null) 'metadata': metadata,
      if (measurementId != null) 'measurement_id': measurementId,
    });
    return _asMap(resp);
  }

  /// List health records for a user.
  Future<Map<String, dynamic>> listHealthRecords(
    String userId, {
    int? typeFilter,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/health-records', queryParameters: {
      'user_id': userId,
      if (typeFilter != null) 'type_filter': typeFilter,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get a single health record.
  Future<Map<String, dynamic>> getHealthRecord(String recordId) async {
    final resp = await _dio.get('/health-records/$recordId');
    return _asMap(resp);
  }

  /// Export health records to FHIR format.
  Future<Map<String, dynamic>> exportToFHIR({
    required String userId,
    List<String>? recordIds,
    int? targetType,
  }) async {
    final resp = await _dio.post('/health-records/export/fhir', data: {
      'user_id': userId,
      if (recordIds != null) 'record_ids': recordIds,
      if (targetType != null) 'target_type': targetType,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Notification Service
  // ==========================================================================

  /// Get unread notification count.
  Future<Map<String, dynamic>> getUnreadCount(String userId) async {
    final resp =
        await _dio.get('/notifications/unread-count', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// List notifications for a user.
  Future<Map<String, dynamic>> listNotifications(
    String userId, {
    bool? unreadOnly,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/notifications', queryParameters: {
      'user_id': userId,
      if (unreadOnly != null) 'unread_only': unreadOnly,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Mark a notification as read.
  Future<Map<String, dynamic>> markNotificationAsRead(
      String notificationId) async {
    final resp =
        await _dio.post('/notifications/$notificationId/read');
    return _asMap(resp);
  }

  // ==========================================================================
  // Community Service
  // ==========================================================================

  /// List community posts.
  Future<Map<String, dynamic>> listPosts({
    int? category,
    String? authorId,
    String? query,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/posts', queryParameters: {
      if (category != null) 'category': category,
      if (authorId != null) 'author_id': authorId,
      if (query != null) 'query': query,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Create a new post.
  Future<Map<String, dynamic>> createPost({
    required String authorId,
    required String title,
    required String content,
    int? category,
    List<String>? tags,
  }) async {
    final resp = await _dio.post('/posts', data: {
      'author_id': authorId,
      'title': title,
      'content': content,
      if (category != null) 'category': category,
      if (tags != null) 'tags': tags,
    });
    return _asMap(resp);
  }

  /// Get a single post.
  Future<Map<String, dynamic>> getPost(String postId) async {
    final resp = await _dio.get('/posts/$postId');
    return _asMap(resp);
  }

  /// Like a post.
  Future<Map<String, dynamic>> likePost(
      String postId, String userId) async {
    final resp = await _dio.post('/posts/$postId/like', data: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Admin Service
  // ==========================================================================

  /// Get system statistics.
  Future<Map<String, dynamic>> getSystemStats() async {
    final resp = await _dio.get('/admin/stats');
    return _asMap(resp);
  }

  /// List users (admin).
  Future<Map<String, dynamic>> adminListUsers({
    String? query,
    int? tierFilter,
    bool? activeOnly,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/admin/users', queryParameters: {
      if (query != null) 'query': query,
      if (tierFilter != null) 'tier_filter': tierFilter,
      if (activeOnly != null) 'active_only': activeOnly,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get audit log (admin).
  Future<Map<String, dynamic>> getAuditLog({
    String? adminId,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/admin/audit-log', queryParameters: {
      if (adminId != null) 'admin_id': adminId,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // AI Inference Service
  // ==========================================================================

  /// Analyze a measurement using AI models.
  Future<Map<String, dynamic>> analyzeMeasurement({
    required String userId,
    required String measurementId,
    List<int>? models,
  }) async {
    final resp = await _dio.post('/ai/analyze', data: {
      'user_id': userId,
      'measurement_id': measurementId,
      if (models != null) 'models': models,
    });
    return _asMap(resp);
  }

  /// Get AI-computed health score for a user.
  Future<Map<String, dynamic>> getHealthScore(
    String userId, {
    int? days,
  }) async {
    final resp = await _dio.get('/ai/health-score/$userId', queryParameters: {
      if (days != null) 'days': days,
    });
    return _asMap(resp);
  }

  /// Predict a trend for a health metric.
  Future<Map<String, dynamic>> predictTrend({
    required String userId,
    required String metricName,
    int? historyDays,
    int? predictionDays,
  }) async {
    final resp = await _dio.post('/ai/predict-trend', data: {
      'user_id': userId,
      'metric_name': metricName,
      if (historyDays != null) 'history_days': historyDays,
      if (predictionDays != null) 'prediction_days': predictionDays,
    });
    return _asMap(resp);
  }

  /// List available AI models.
  Future<Map<String, dynamic>> listAiModels() async {
    final resp = await _dio.get('/ai/models');
    return _asMap(resp);
  }

  // ==========================================================================
  // Cartridge Service
  // ==========================================================================

  /// Read cartridge data from NFC tag.
  Future<Map<String, dynamic>> readCartridge({
    required List<int> nfcTagData,
    int tagVersion = 2,
  }) async {
    final resp = await _dio.post('/cartridges/read', data: {
      'nfc_tag_data': nfcTagData,
      'tag_version': tagVersion,
    });
    return _asMap(resp);
  }

  /// Record cartridge usage.
  Future<Map<String, dynamic>> recordCartridgeUsage({
    required String userId,
    required String sessionId,
    required String cartridgeUid,
    int? categoryCode,
    int? typeIndex,
  }) async {
    final resp = await _dio.post('/cartridges/usage', data: {
      'user_id': userId,
      'session_id': sessionId,
      'cartridge_uid': cartridgeUid,
      if (categoryCode != null) 'category_code': categoryCode,
      if (typeIndex != null) 'type_index': typeIndex,
    });
    return _asMap(resp);
  }

  /// List cartridge categories/types.
  Future<Map<String, dynamic>> listCartridgeTypes() async {
    final resp = await _dio.get('/cartridges/types');
    return _asMap(resp);
  }

  /// Get remaining uses for a cartridge.
  Future<Map<String, dynamic>> getRemainingUses(String cartridgeId) async {
    final resp = await _dio.get('/cartridges/$cartridgeId/remaining');
    return _asMap(resp);
  }

  /// Validate a cartridge (NFC UID + expiry + remaining uses).
  Future<Map<String, dynamic>> validateCartridge({
    required String cartridgeUid,
    int? categoryCode,
    int? typeIndex,
    String? userId,
  }) async {
    final resp = await _dio.post('/cartridges/validate', data: {
      'cartridge_uid': cartridgeUid,
      if (categoryCode != null) 'category_code': categoryCode,
      if (typeIndex != null) 'type_index': typeIndex,
      if (userId != null) 'user_id': userId,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Calibration Service
  // ==========================================================================

  /// Register factory calibration data.
  Future<Map<String, dynamic>> registerFactoryCalibration({
    required String deviceId,
    required int cartridgeCategory,
    required int cartridgeTypeIndex,
    double? alpha,
    List<double>? channelOffsets,
    List<double>? channelGains,
    double? tempCoefficient,
    double? humidityCoefficient,
    String? referenceStandard,
    String? calibratedBy,
  }) async {
    final resp = await _dio.post('/calibration/factory', data: {
      'device_id': deviceId,
      'cartridge_category': cartridgeCategory,
      'cartridge_type_index': cartridgeTypeIndex,
      if (alpha != null) 'alpha': alpha,
      if (channelOffsets != null) 'channel_offsets': channelOffsets,
      if (channelGains != null) 'channel_gains': channelGains,
      if (tempCoefficient != null) 'temp_coefficient': tempCoefficient,
      if (humidityCoefficient != null)
        'humidity_coefficient': humidityCoefficient,
      if (referenceStandard != null) 'reference_standard': referenceStandard,
      if (calibratedBy != null) 'calibrated_by': calibratedBy,
    });
    return _asMap(resp);
  }

  /// Perform field calibration (user calibration).
  Future<Map<String, dynamic>> performFieldCalibration({
    required String deviceId,
    required String userId,
    required int cartridgeCategory,
    required int cartridgeTypeIndex,
    List<double>? referenceValues,
    List<double>? measuredValues,
    double? temperatureC,
    double? humidityPct,
  }) async {
    final resp = await _dio.post('/calibration/field', data: {
      'device_id': deviceId,
      'user_id': userId,
      'cartridge_category': cartridgeCategory,
      'cartridge_type_index': cartridgeTypeIndex,
      if (referenceValues != null) 'reference_values': referenceValues,
      if (measuredValues != null) 'measured_values': measuredValues,
      if (temperatureC != null) 'temperature_c': temperatureC,
      if (humidityPct != null) 'humidity_pct': humidityPct,
    });
    return _asMap(resp);
  }

  /// Check calibration status for a device.
  Future<Map<String, dynamic>> checkCalibrationStatus(
    String deviceId, {
    int? cartridgeCategory,
    int? cartridgeTypeIndex,
  }) async {
    final resp = await _dio
        .get('/calibration/$deviceId/status', queryParameters: {
      if (cartridgeCategory != null) 'cartridge_category': cartridgeCategory,
      if (cartridgeTypeIndex != null)
        'cartridge_type_index': cartridgeTypeIndex,
    });
    return _asMap(resp);
  }

  /// List calibration models.
  Future<Map<String, dynamic>> listCalibrationModels() async {
    final resp = await _dio.get('/calibration/models');
    return _asMap(resp);
  }

  // ==========================================================================
  // Coaching Service
  // ==========================================================================

  /// Set a health goal.
  Future<Map<String, dynamic>> setHealthGoal({
    required String userId,
    int? category,
    String? metricName,
    double? targetValue,
    String? unit,
    String? description,
    String? targetDate,
  }) async {
    final resp = await _dio.post('/coaching/goals', data: {
      'user_id': userId,
      if (category != null) 'category': category,
      if (metricName != null) 'metric_name': metricName,
      if (targetValue != null) 'target_value': targetValue,
      if (unit != null) 'unit': unit,
      if (description != null) 'description': description,
      if (targetDate != null) 'target_date': targetDate,
    });
    return _asMap(resp);
  }

  /// Get health goals for a user.
  Future<Map<String, dynamic>> getHealthGoals(
    String userId, {
    int? statusFilter,
  }) async {
    final resp =
        await _dio.get('/coaching/goals/$userId', queryParameters: {
      if (statusFilter != null) 'status_filter': statusFilter,
    });
    return _asMap(resp);
  }

  /// Generate AI coaching message.
  Future<Map<String, dynamic>> generateCoaching({
    required String userId,
    String? measurementId,
    int? coachingType,
  }) async {
    final resp = await _dio.post('/coaching/generate', data: {
      'user_id': userId,
      if (measurementId != null) 'measurement_id': measurementId,
      if (coachingType != null) 'coaching_type': coachingType,
    });
    return _asMap(resp);
  }

  /// Generate daily health report.
  Future<Map<String, dynamic>> generateDailyReport(String userId) async {
    final resp = await _dio.get('/coaching/daily-report/$userId');
    return _asMap(resp);
  }

  /// Get personalized recommendations.
  Future<Map<String, dynamic>> getRecommendations(
    String userId, {
    int? typeFilter,
    int? limit,
  }) async {
    final resp = await _dio
        .get('/coaching/recommendations/$userId', queryParameters: {
      if (typeFilter != null) 'type_filter': typeFilter,
      if (limit != null) 'limit': limit,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Video Service (WebRTC)
  // ==========================================================================

  /// Join a video room and get WebRTC token + ICE servers.
  Future<Map<String, dynamic>> joinVideoRoom({
    required String roomId,
    required String userId,
    String? displayName,
  }) async {
    final resp = await _dio.post('/video/rooms/$roomId/join', data: {
      'user_id': userId,
      if (displayName != null) 'display_name': displayName,
    });
    return _asMap(resp);
  }

  /// Leave a video room.
  Future<Map<String, dynamic>> leaveVideoRoom({
    required String roomId,
    required String userId,
  }) async {
    final resp = await _dio.post('/video/rooms/$roomId/leave', data: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// Analyze food image with AI.
  Future<Map<String, dynamic>> analyzeFoodImage({
    required String userId,
    required String imagePath,
  }) async {
    final formData = FormData.fromMap({
      'user_id': userId,
      'image': await MultipartFile.fromFile(imagePath),
    });
    final resp = await _dio.post('/ai/food-analyze', data: formData);
    return _asMap(resp);
  }

  /// Import health data from external sources (HealthKit/Google Health Connect).
  Future<Map<String, dynamic>> importExternalHealthData({
    required String userId,
    required String source,
    required List<Map<String, dynamic>> records,
  }) async {
    final resp = await _dio.post('/health-records/import', data: {
      'user_id': userId,
      'source': source,
      'records': records,
    });
    return _asMap(resp);
  }

  /// Analyze exercise video for calorie estimation.
  Future<Map<String, dynamic>> analyzeExerciseVideo({
    required String userId,
    required String videoPath,
  }) async {
    final formData = FormData.fromMap({
      'user_id': userId,
      'video': await MultipartFile.fromFile(videoPath),
    });
    final resp = await _dio.post('/ai/exercise-analyze', data: formData);
    return _asMap(resp);
  }

  // ==========================================================================
  // Community Extended (Challenges, Q&A)
  // ==========================================================================

  /// List health challenges.
  Future<Map<String, dynamic>> getChallenges({
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/community/challenges', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Join a health challenge.
  Future<Map<String, dynamic>> joinChallenge({
    required String challengeId,
    required String userId,
  }) async {
    final resp = await _dio.post('/community/challenges/$challengeId/join', data: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// List Q&A questions.
  Future<Map<String, dynamic>> getQnaQuestions({
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/community/qna', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Create a post with image attachment (multipart).
  Future<Map<String, dynamic>> createPostWithImage({
    required String authorId,
    required String title,
    required String content,
    int? category,
    List<String>? imagePaths,
  }) async {
    final map = <String, dynamic>{
      'author_id': authorId,
      'title': title,
      'content': content,
      if (category != null) 'category': category,
    };
    if (imagePaths != null && imagePaths.isNotEmpty) {
      map['images'] = await Future.wait(
        imagePaths.map((p) => MultipartFile.fromFile(p)),
      );
    }
    final resp = await _dio.post('/posts', data: FormData.fromMap(map));
    return _asMap(resp);
  }

  // ==========================================================================
  // Family Extended
  // ==========================================================================

  /// Create a family group.
  Future<Map<String, dynamic>> createFamilyGroup({
    required String userId,
    required String name,
    String? inviteMethod,
  }) async {
    final resp = await _dio.post('/family/groups', data: {
      'user_id': userId,
      'name': name,
      if (inviteMethod != null) 'invite_method': inviteMethod,
    });
    return _asMap(resp);
  }

  /// Update a family member's role/mode.
  Future<Map<String, dynamic>> updateFamilyMember({
    required String groupId,
    required String memberId,
    String? role,
    String? mode,
    Map<String, bool>? permissions,
  }) async {
    final resp = await _dio.put('/family/groups/$groupId/members/$memberId', data: {
      if (role != null) 'role': role,
      if (mode != null) 'mode': mode,
      if (permissions != null) 'permissions': permissions,
    });
    return _asMap(resp);
  }

  /// Get guardian dashboard data.
  Future<Map<String, dynamic>> getGuardianDashboard({
    required String groupId,
  }) async {
    final resp = await _dio.get('/family/groups/$groupId/guardian-dashboard');
    return _asMap(resp);
  }

  /// Get alert detail.
  Future<Map<String, dynamic>> getAlertDetail(String alertId) async {
    final resp = await _dio.get('/notifications/alerts/$alertId');
    return _asMap(resp);
  }

  // ==========================================================================
  // Medical Extended (Consultation Result)
  // ==========================================================================

  /// Get consultation result.
  Future<Map<String, dynamic>> getConsultationResult(String consultationId) async {
    final resp = await _dio.get('/telemedicine/consultations/$consultationId/result');
    return _asMap(resp);
  }

  // ==========================================================================
  // Market Extended (Order Detail, Plans)
  // ==========================================================================

  /// Get order detail.
  Future<Map<String, dynamic>> getOrderDetail(String orderId) async {
    final resp = await _dio.get('/orders/$orderId');
    return _asMap(resp);
  }

  /// Get subscription plans comparison.
  Future<Map<String, dynamic>> getSubscriptionPlans() async {
    final resp = await _dio.get('/subscriptions/plans/compare');
    return _asMap(resp);
  }

  // ==========================================================================
  // Admin Extended (Monitor, Hierarchy, Compliance)
  // ==========================================================================

  /// Get system metrics (CPU, memory, network, etc.).
  Future<Map<String, dynamic>> getSystemMetrics() async {
    final resp = await _dio.get('/admin/metrics');
    return _asMap(resp);
  }

  /// Get organization hierarchy.
  Future<Map<String, dynamic>> getHierarchy() async {
    final resp = await _dio.get('/admin/hierarchy');
    return _asMap(resp);
  }

  /// Get compliance checklist.
  Future<Map<String, dynamic>> getComplianceChecklist() async {
    final resp = await _dio.get('/admin/compliance');
    return _asMap(resp);
  }

  // ==========================================================================
  // Settings Extended (Inquiry)
  // ==========================================================================

  /// Create a 1:1 inquiry.
  Future<Map<String, dynamic>> createInquiry({
    required String userId,
    required String type,
    required String title,
    required String content,
    bool? notifyByPush,
    bool? notifyByEmail,
  }) async {
    final resp = await _dio.post('/support/inquiries', data: {
      'user_id': userId,
      'type': type,
      'title': title,
      'content': content,
      if (notifyByPush != null) 'notify_by_push': notifyByPush,
      if (notifyByEmail != null) 'notify_by_email': notifyByEmail,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // AI Chat Streaming (C1)
  // ==========================================================================

  /// Stream chat with AI (SSE-like, returns full response).
  Future<Map<String, dynamic>> streamChat({
    required String userId,
    required String message,
    List<Map<String, String>>? history,
  }) async {
    final resp = await _dio.post('/ai/chat/stream', data: {
      'user_id': userId,
      'message': message,
      if (history != null) 'history': history,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Challenge Leaderboard (C8)
  // ==========================================================================

  /// Get challenge leaderboard.
  Future<Map<String, dynamic>> getChallengeLeaderboard(
    String challengeId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get(
      '/community/challenges/$challengeId/leaderboard',
      queryParameters: {'limit': limit, 'offset': offset},
    );
    return _asMap(resp);
  }

  /// Update challenge progress.
  Future<Map<String, dynamic>> updateChallengeProgress({
    required String challengeId,
    required String userId,
    required int progressValue,
  }) async {
    final resp = await _dio.post(
      '/community/challenges/$challengeId/progress',
      data: {
        'user_id': userId,
        'progress_value': progressValue,
      },
    );
    return _asMap(resp);
  }

  // ==========================================================================
  // Admin Revenue & Inventory (C12)
  // ==========================================================================

  /// Get revenue statistics.
  Future<Map<String, dynamic>> getRevenueStats({
    String? period,
    int? months,
  }) async {
    final resp = await _dio.get('/admin/revenue', queryParameters: {
      if (period != null) 'period': period,
      if (months != null) 'months': months,
    });
    return _asMap(resp);
  }

  /// Get inventory statistics.
  Future<Map<String, dynamic>> getInventoryStats() async {
    final resp = await _dio.get('/admin/inventory');
    return _asMap(resp);
  }

  // ==========================================================================
  // Realtime Translation (C6)
  // ==========================================================================

  /// Translate text in realtime with medical term support.
  Future<Map<String, dynamic>> translateRealtime({
    required String text,
    required String sourceLanguage,
    required String targetLanguage,
    bool includeMedicalTerms = true,
  }) async {
    final resp = await _dio.post('/translations/realtime', data: {
      'text': text,
      'source_language': sourceLanguage,
      'target_language': targetLanguage,
      'include_medical_terms': includeMedicalTerms,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Health Record Extended (Sprint 12)
  // ==========================================================================

  /// Update a health record.
  Future<Map<String, dynamic>> updateHealthRecord(
    String recordId, {
    String? title,
    String? description,
    Map<String, String>? metadata,
  }) async {
    final resp = await _dio.put('/health-records/$recordId', data: {
      if (title != null) 'title': title,
      if (description != null) 'description': description,
      if (metadata != null) 'metadata': metadata,
    });
    return _asMap(resp);
  }

  /// Delete a health record.
  Future<void> deleteHealthRecord(String recordId) async {
    await _dio.delete('/health-records/$recordId');
  }

  /// Get health summary for a user.
  Future<Map<String, dynamic>> getHealthSummary(String userId) async {
    final resp = await _dio.get('/health-records/summary', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// Create a data sharing consent.
  Future<Map<String, dynamic>> createDataSharingConsent({
    required String userId,
    required String providerId,
    required List<String> dataTypes,
  }) async {
    final resp = await _dio.post('/health-records/consents', data: {
      'user_id': userId,
      'provider_id': providerId,
      'data_types': dataTypes,
    });
    return _asMap(resp);
  }

  /// Revoke a data sharing consent.
  Future<void> revokeDataSharingConsent(String consentId) async {
    await _dio.delete('/health-records/consents/$consentId');
  }

  /// List data sharing consents for a user.
  Future<Map<String, dynamic>> listDataSharingConsents(String userId) async {
    final resp =
        await _dio.get('/health-records/consents', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// Share a health record with a provider.
  Future<Map<String, dynamic>> shareWithProvider({
    required String recordId,
    required String providerId,
  }) async {
    final resp = await _dio.post('/health-records/share-provider', data: {
      'record_id': recordId,
      'provider_id': providerId,
    });
    return _asMap(resp);
  }

  /// Get data access log for a record.
  Future<Map<String, dynamic>> getDataAccessLog(String recordId) async {
    final resp = await _dio.get('/health-records/access-log', queryParameters: {
      'record_id': recordId,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Prescription Extended (Sprint 12)
  // ==========================================================================

  /// Create a prescription.
  Future<Map<String, dynamic>> createPrescription({
    required String consultationId,
    required List<Map<String, dynamic>> medications,
  }) async {
    final resp = await _dio.post('/prescriptions', data: {
      'consultation_id': consultationId,
      'medications': medications,
    });
    return _asMap(resp);
  }

  /// Get a prescription by ID.
  Future<Map<String, dynamic>> getPrescription(String prescriptionId) async {
    final resp = await _dio.get('/prescriptions/$prescriptionId');
    return _asMap(resp);
  }

  /// List prescriptions for a user.
  Future<Map<String, dynamic>> listPrescriptions({
    String? userId,
    String? status,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/prescriptions', queryParameters: {
      if (userId != null) 'user_id': userId,
      if (status != null) 'status': status,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Update prescription status.
  Future<Map<String, dynamic>> updatePrescriptionStatus(
    String prescriptionId,
    String status,
  ) async {
    final resp =
        await _dio.patch('/prescriptions/$prescriptionId/status', data: {
      'status': status,
    });
    return _asMap(resp);
  }

  /// Add a medication to a prescription.
  Future<Map<String, dynamic>> addMedication(
    String prescriptionId,
    Map<String, dynamic> medication,
  ) async {
    final resp = await _dio
        .post('/prescriptions/$prescriptionId/medications', data: medication);
    return _asMap(resp);
  }

  /// Remove a medication from a prescription.
  Future<void> removeMedication(
      String prescriptionId, String medicationId) async {
    await _dio
        .delete('/prescriptions/$prescriptionId/medications/$medicationId');
  }

  /// Check drug interactions between medications.
  Future<Map<String, dynamic>> checkDrugInteraction(
      List<String> medicationIds) async {
    final resp =
        await _dio.post('/prescriptions/check-drug-interaction', data: {
      'medication_ids': medicationIds,
    });
    return _asMap(resp);
  }

  /// Get medication reminders for a user.
  Future<Map<String, dynamic>> getMedicationReminders(
    String userId, {
    String? date,
  }) async {
    final resp = await _dio.get('/prescriptions/reminders', queryParameters: {
      'user_id': userId,
      if (date != null) 'date': date,
    });
    return _asMap(resp);
  }

  /// Update dispensary status of a prescription.
  Future<Map<String, dynamic>> updateDispensaryStatus(
    String prescriptionId,
    String status,
  ) async {
    final resp = await _dio
        .patch('/prescriptions/$prescriptionId/dispensary-status', data: {
      'status': status,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Community Extended – Comments & Challenges (Sprint 12)
  // ==========================================================================

  /// Create a comment on a post.
  Future<Map<String, dynamic>> createComment({
    required String postId,
    required String authorId,
    required String content,
  }) async {
    final resp = await _dio.post('/posts/$postId/comments', data: {
      'author_id': authorId,
      'content': content,
    });
    return _asMap(resp);
  }

  /// List comments on a post.
  Future<Map<String, dynamic>> listComments(
    String postId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio
        .get('/posts/$postId/comments', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Delete a post.
  Future<void> deletePost(String postId) async {
    await _dio.delete('/posts/$postId');
  }

  /// Create a health challenge.
  Future<Map<String, dynamic>> createChallenge({
    required String title,
    required String description,
    required String startDate,
    required String endDate,
  }) async {
    final resp = await _dio.post('/challenges', data: {
      'title': title,
      'description': description,
      'start_date': startDate,
      'end_date': endDate,
    });
    return _asMap(resp);
  }

  /// Get a challenge by ID.
  Future<Map<String, dynamic>> getChallenge(String challengeId) async {
    final resp = await _dio.get('/challenges/$challengeId');
    return _asMap(resp);
  }

  /// List all challenges.
  Future<Map<String, dynamic>> listChallenges({
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/challenges', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Family Extended (Sprint 12)
  // ==========================================================================

  /// Get a family group by ID.
  Future<Map<String, dynamic>> getFamilyGroup(String groupId) async {
    final resp = await _dio.get('/family/groups/$groupId');
    return _asMap(resp);
  }

  /// Invite a member to a family group.
  Future<Map<String, dynamic>> inviteFamilyMember({
    required String groupId,
    required String inviteePhone,
    String? role,
  }) async {
    final resp = await _dio.post('/family/groups/$groupId/invite', data: {
      'invitee_phone': inviteePhone,
      if (role != null) 'role': role,
    });
    return _asMap(resp);
  }

  /// Respond to a family invitation.
  Future<Map<String, dynamic>> respondToInvitation({
    required String invitationId,
    required bool accept,
  }) async {
    final resp = await _dio
        .post('/family/invitations/$invitationId/respond', data: {
      'accept': accept,
    });
    return _asMap(resp);
  }

  /// Remove a member from a family group.
  Future<void> removeFamilyMember(String groupId, String userId) async {
    await _dio.delete('/family/groups/$groupId/members/$userId');
  }

  /// Update a member's role in a family group.
  Future<Map<String, dynamic>> updateMemberRole({
    required String groupId,
    required String userId,
    required String role,
  }) async {
    final resp =
        await _dio.put('/family/groups/$groupId/members/$userId/role', data: {
      'role': role,
    });
    return _asMap(resp);
  }

  /// List members of a family group.
  Future<Map<String, dynamic>> listFamilyMembers(String groupId) async {
    final resp = await _dio.get('/family/groups/$groupId/members');
    return _asMap(resp);
  }

  /// Set sharing preferences for a family group.
  Future<void> setSharingPreferences(
    String groupId,
    Map<String, dynamic> preferences,
  ) async {
    await _dio.put('/family/groups/$groupId/sharing-prefs',
        data: preferences);
  }

  /// Validate sharing access for a user in a family group.
  Future<Map<String, dynamic>> validateSharingAccess({
    required String groupId,
    required String userId,
  }) async {
    final resp = await _dio
        .get('/family/groups/$groupId/sharing-access', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Video/WebRTC Extended (Sprint 12)
  // ==========================================================================

  /// Create a video room.
  Future<Map<String, dynamic>> createVideoRoom(
      String consultationId) async {
    final resp = await _dio.post('/video/rooms', data: {
      'consultation_id': consultationId,
    });
    return _asMap(resp);
  }

  /// Get a video room.
  Future<Map<String, dynamic>> getVideoRoom(String roomId) async {
    final resp = await _dio.get('/video/rooms/$roomId');
    return _asMap(resp);
  }

  /// End a video room.
  Future<Map<String, dynamic>> endVideoRoom(String roomId) async {
    final resp = await _dio.post('/video/rooms/$roomId/end');
    return _asMap(resp);
  }

  /// Send a WebRTC signal.
  Future<void> sendVideoSignal({
    required String roomId,
    required String signalType,
    required String payload,
  }) async {
    await _dio.post('/video/rooms/$roomId/signal', data: {
      'signal_type': signalType,
      'payload': payload,
    });
  }

  /// Get pending video signals for a user in a room.
  Future<Map<String, dynamic>> getVideoSignals({
    required String roomId,
    required String userId,
  }) async {
    final resp = await _dio.get('/video/rooms/$roomId/signals', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// List participants in a video room.
  Future<Map<String, dynamic>> listVideoParticipants(String roomId) async {
    final resp = await _dio.get('/video/rooms/$roomId/participants');
    return _asMap(resp);
  }

  /// Get video room statistics.
  Future<Map<String, dynamic>> getVideoRoomStats(String roomId) async {
    final resp = await _dio.get('/video/rooms/$roomId/stats');
    return _asMap(resp);
  }

  // ==========================================================================
  // Notification Extended (Sprint 12)
  // ==========================================================================

  /// Send a notification to a user.
  Future<Map<String, dynamic>> sendNotification({
    required String userId,
    required String title,
    required String body,
    String? type,
  }) async {
    final resp = await _dio.post('/notifications', data: {
      'user_id': userId,
      'title': title,
      'body': body,
      if (type != null) 'type': type,
    });
    return _asMap(resp);
  }

  /// Mark all notifications as read for a user.
  Future<void> markAllNotificationsAsRead(String userId) async {
    await _dio.post('/notifications/mark-all-read', data: {
      'user_id': userId,
    });
  }

  /// Update notification preferences.
  Future<void> updateNotificationPreferences(
    String userId,
    Map<String, dynamic> preferences,
  ) async {
    await _dio.put('/notifications/preferences', data: {
      'user_id': userId,
      ...preferences,
    });
  }

  /// Get notification preferences.
  Future<Map<String, dynamic>> getNotificationPreferences(
      String userId) async {
    final resp = await _dio
        .get('/notifications/preferences', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  /// Register a push notification token (FCM/APNs).
  Future<Map<String, dynamic>> registerPushToken({
    required String userId,
    required String token,
    String platform = 'fcm',
  }) async {
    final resp = await _dio.post('/notifications/push-token', data: {
      'user_id': userId,
      'token': token,
      'platform': platform,
    });
    return _asMap(resp);
  }

  /// Send a notification from a template.
  Future<Map<String, dynamic>> sendNotificationFromTemplate({
    required String templateId,
    required String userId,
    Map<String, String>? params,
  }) async {
    final resp = await _dio.post('/notifications/send-template', data: {
      'template_id': templateId,
      'user_id': userId,
      if (params != null) 'params': params,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Translation Extended (Sprint 12)
  // ==========================================================================

  /// Detect the language of a text.
  Future<Map<String, dynamic>> detectLanguage(String text) async {
    final resp = await _dio.post('/translations/detect-language', data: {
      'text': text,
    });
    return _asMap(resp);
  }

  /// List supported languages.
  Future<Map<String, dynamic>> listSupportedLanguages() async {
    final resp = await _dio.get('/translations/languages');
    return _asMap(resp);
  }

  /// Translate a batch of texts.
  Future<Map<String, dynamic>> translateBatch({
    required List<String> texts,
    required String targetLanguage,
    String? sourceLanguage,
  }) async {
    final resp = await _dio.post('/translations/batch', data: {
      'texts': texts,
      'target_language': targetLanguage,
      if (sourceLanguage != null) 'source_language': sourceLanguage,
    });
    return _asMap(resp);
  }

  /// Get translation history for a user.
  Future<Map<String, dynamic>> getTranslationHistory(
    String userId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp =
        await _dio.get('/translations/history', queryParameters: {
      'user_id': userId,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get translation usage statistics.
  Future<Map<String, dynamic>> getTranslationUsage(String userId) async {
    final resp =
        await _dio.get('/translations/usage', queryParameters: {
      'user_id': userId,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Telemedicine Extended (Sprint 12)
  // ==========================================================================

  /// Get a consultation by ID.
  Future<Map<String, dynamic>> getConsultation(
      String consultationId) async {
    final resp =
        await _dio.get('/telemedicine/consultations/$consultationId');
    return _asMap(resp);
  }

  /// List consultations.
  Future<Map<String, dynamic>> listConsultations({
    String? userId,
    String? status,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio
        .get('/telemedicine/consultations', queryParameters: {
      if (userId != null) 'user_id': userId,
      if (status != null) 'status': status,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Start a video session for a consultation.
  Future<Map<String, dynamic>> startVideoSession(
      String consultationId) async {
    final resp = await _dio
        .post('/telemedicine/consultations/$consultationId/start-video');
    return _asMap(resp);
  }

  /// End a video session for a consultation.
  Future<Map<String, dynamic>> endVideoSession(
      String consultationId) async {
    final resp = await _dio
        .post('/telemedicine/consultations/$consultationId/end-video');
    return _asMap(resp);
  }

  /// Rate a consultation.
  Future<void> rateConsultation({
    required String consultationId,
    required int rating,
    String? comment,
  }) async {
    await _dio
        .post('/telemedicine/consultations/$consultationId/rate', data: {
      'rating': rating,
      if (comment != null) 'comment': comment,
    });
  }

  // ==========================================================================
  // Subscription Extended (Sprint 12)
  // ==========================================================================

  /// Check if a user has access to a feature.
  Future<Map<String, dynamic>> checkFeatureAccess({
    required String userId,
    required String feature,
  }) async {
    final resp = await _dio
        .get('/subscriptions/$userId/feature-access', queryParameters: {
      'feature': feature,
    });
    return _asMap(resp);
  }

  /// Check if a user has access to a cartridge type.
  Future<Map<String, dynamic>> checkCartridgeAccess({
    required String userId,
    required String cartridgeType,
  }) async {
    final resp = await _dio
        .get('/subscriptions/$userId/cartridge-access', queryParameters: {
      'cartridge_type': cartridgeType,
    });
    return _asMap(resp);
  }

  /// List cartridges accessible to a user's subscription.
  Future<Map<String, dynamic>> listAccessibleCartridges(
      String userId) async {
    final resp =
        await _dio.get('/subscriptions/$userId/accessible-cartridges');
    return _asMap(resp);
  }

  // ==========================================================================
  // Coaching Extended (Sprint 12)
  // ==========================================================================

  /// List coaching messages.
  Future<Map<String, dynamic>> listCoachingMessages(
    String userId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/coaching/messages', queryParameters: {
      'user_id': userId,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get weekly health report.
  Future<Map<String, dynamic>> getWeeklyReport(String userId) async {
    final resp = await _dio.get('/coaching/weekly-report/$userId');
    return _asMap(resp);
  }

  // ==========================================================================
  // Admin Extended (Sprint 12)
  // ==========================================================================

  /// Create an admin.
  Future<Map<String, dynamic>> createAdmin({
    required String email,
    required String password,
    required String displayName,
    int? role,
    String? region,
    String? branch,
  }) async {
    final resp = await _dio.post('/admin/admins', data: {
      'email': email,
      'password': password,
      'display_name': displayName,
      if (role != null) 'role': role,
      if (region != null) 'region': region,
      if (branch != null) 'branch': branch,
    });
    return _asMap(resp);
  }

  /// List admins.
  Future<Map<String, dynamic>> listAdmins({
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/admin/admins', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Get admin by ID.
  Future<Map<String, dynamic>> getAdmin(String adminId) async {
    final resp = await _dio.get('/admin/admins/$adminId');
    return _asMap(resp);
  }

  /// Update admin role.
  Future<Map<String, dynamic>> updateAdminRole(
    String adminId,
    int newRole,
  ) async {
    final resp = await _dio.put('/admin/admins/$adminId/role', data: {
      'new_role': newRole,
    });
    return _asMap(resp);
  }

  /// Deactivate an admin.
  Future<Map<String, dynamic>> deactivateAdmin(String adminId) async {
    final resp = await _dio.post('/admin/admins/$adminId/deactivate');
    return _asMap(resp);
  }

  /// List admins by region.
  Future<Map<String, dynamic>> listAdminsByRegion({
    String? countryCode,
    String? regionCode,
  }) async {
    final resp =
        await _dio.get('/admin/admins/by-region', queryParameters: {
      if (countryCode != null) 'country_code': countryCode,
      if (regionCode != null) 'region_code': regionCode,
    });
    return _asMap(resp);
  }

  /// Get detailed audit log.
  Future<Map<String, dynamic>> getAuditLogDetails({
    String? adminId,
    String? action,
    int limit = 50,
    int offset = 0,
  }) async {
    final resp =
        await _dio.get('/admin/audit-log/details', queryParameters: {
      if (adminId != null) 'admin_id': adminId,
      if (action != null) 'action': action,
      'limit': limit,
      'offset': offset,
    });
    return _asMap(resp);
  }

  /// Set a system config value.
  Future<Map<String, dynamic>> setSystemConfig({
    required String key,
    required String value,
    String? description,
  }) async {
    final resp = await _dio.put('/admin/config', data: {
      'key': key,
      'value': value,
      if (description != null) 'description': description,
    });
    return _asMap(resp);
  }

  /// Get a system config value.
  Future<Map<String, dynamic>> getSystemConfig(String key) async {
    final resp = await _dio.get('/admin/config', queryParameters: {
      'key': key,
    });
    return _asMap(resp);
  }

  /// List system configs.
  Future<Map<String, dynamic>> listSystemConfigs({
    String? language,
    String? category,
    bool? includeSecrets,
  }) async {
    final resp = await _dio.get('/admin/configs', queryParameters: {
      if (language != null) 'language': language,
      if (category != null) 'category': category,
      if (includeSecrets != null) 'include_secrets': includeSecrets,
    });
    return _asMap(resp);
  }

  /// Get a config with metadata.
  Future<Map<String, dynamic>> getConfigWithMeta(
    String key, {
    String? language,
  }) async {
    final resp = await _dio.get('/admin/configs/$key', queryParameters: {
      if (language != null) 'language': language,
    });
    return _asMap(resp);
  }

  /// Validate a config value.
  Future<Map<String, dynamic>> validateConfigValue({
    required String key,
    required String value,
  }) async {
    final resp = await _dio.post('/admin/configs/validate', data: {
      'key': key,
      'value': value,
    });
    return _asMap(resp);
  }

  /// Bulk set config values.
  Future<Map<String, dynamic>> bulkSetConfigs({
    required List<Map<String, String>> configs,
    String? reason,
  }) async {
    final resp = await _dio.post('/admin/configs/bulk', data: {
      'configs': configs,
      if (reason != null) 'reason': reason,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Measurement Extended (Sprint 12)
  // ==========================================================================

  /// Export a single measurement in FHIR format.
  Future<Map<String, dynamic>> exportSingleMeasurement(
      String measurementId) async {
    final resp = await _dio.get('/measurements/$measurementId/export');
    return _asMap(resp);
  }

  /// Export measurements to FHIR Observations bundle.
  Future<Map<String, dynamic>> exportToFhirObservations({
    required String userId,
    List<String>? measurementIds,
  }) async {
    final resp =
        await _dio.post('/measurements/export/fhir-observations', data: {
      'user_id': userId,
      if (measurementIds != null) 'measurement_ids': measurementIds,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Device Extended (Sprint 12)
  // ==========================================================================

  /// Request an OTA firmware update for a device.
  Future<Map<String, dynamic>> requestOtaUpdate(
    String deviceId, {
    String? targetVersion,
  }) async {
    final resp = await _dio.post('/devices/$deviceId/ota', data: {
      if (targetVersion != null) 'target_version': targetVersion,
    });
    return _asMap(resp);
  }

  /// Update device status.
  Future<Map<String, dynamic>> updateDeviceStatus(
    String deviceId,
    String status,
  ) async {
    final resp = await _dio.put('/devices/$deviceId/status', data: {
      'status': status,
    });
    return _asMap(resp);
  }

  // ==========================================================================
  // Reservation Extended (Sprint 12)
  // ==========================================================================

  /// Cancel a reservation.
  Future<void> cancelReservation(String reservationId) async {
    await _dio.delete('/reservations/$reservationId');
  }

  // ==========================================================================
  // Helpers
  // ==========================================================================

  /// Safely cast Dio response data to Map<String, dynamic>.
  Map<String, dynamic> _asMap(Response<dynamic> resp) {
    if (resp.data is Map<String, dynamic>) {
      return resp.data as Map<String, dynamic>;
    }
    return <String, dynamic>{'data': resp.data};
  }
}
