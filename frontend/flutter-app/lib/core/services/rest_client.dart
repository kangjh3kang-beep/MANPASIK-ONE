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
    double? latitude,
    double? longitude,
    double? radiusKm,
    int limit = 20,
    int offset = 0,
  }) async {
    final resp = await _dio.get('/facilities', queryParameters: {
      if (query != null) 'query': query,
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
