import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/network/tenant_interceptor.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  group('TenantInterceptor — 활성 조직/사용자 저장', () {
    test('setActiveTenant 후 getActiveTenant 가 동일 값 반환', () async {
      await TenantInterceptor.setActiveTenant('hospital-A');
      final got = await TenantInterceptor.getActiveTenant();
      expect(got, 'hospital-A');
    });

    test('null 또는 빈 문자열은 저장값 제거', () async {
      await TenantInterceptor.setActiveTenant('hospital-A');
      await TenantInterceptor.setActiveTenant(null);
      expect(await TenantInterceptor.getActiveTenant(), isNull);

      await TenantInterceptor.setActiveTenant('hospital-B');
      await TenantInterceptor.setActiveTenant('');
      expect(await TenantInterceptor.getActiveTenant(), isNull);
    });

    test('setActiveUser/getActiveUser', () async {
      await TenantInterceptor.setActiveUser('u-001');
      expect(await TenantInterceptor.getActiveUser(), 'u-001');
      await TenantInterceptor.setActiveUser(null);
      expect(await TenantInterceptor.getActiveUser(), isNull);
    });

    test('clear 는 양쪽 모두 제거', () async {
      await TenantInterceptor.setActiveTenant('t-1');
      await TenantInterceptor.setActiveUser('u-1');
      await TenantInterceptor.clear();
      expect(await TenantInterceptor.getActiveTenant(), isNull);
      expect(await TenantInterceptor.getActiveUser(), isNull);
    });
  });

  group('TenantInterceptor — Dio 통합', () {
    test('활성 테넌트 설정 시 X-Tenant-ID 헤더 자동 주입', () async {
      await TenantInterceptor.setActiveTenant('clinic-7');
      await TenantInterceptor.setActiveUser('user-42');

      final dio = Dio();
      dio.interceptors.add(TenantInterceptor());

      Map<String, dynamic>? capturedHeaders;
      dio.interceptors.add(InterceptorsWrapper(
        onRequest: (options, handler) {
          capturedHeaders = Map<String, dynamic>.from(options.headers);
          // 실제 네트워크 호출 회피
          handler.reject(
            DioException(
              requestOptions: options,
              type: DioExceptionType.cancel,
            ),
          );
        },
      ));

      try {
        await dio.get('http://example.invalid/test');
      } catch (_) {
        // 인터셉터 cancel 예상됨
      }

      expect(capturedHeaders, isNotNull);
      expect(capturedHeaders!['X-Tenant-ID'], 'clinic-7');
      expect(capturedHeaders!['X-User-ID'], 'user-42');
    });

    test('테넌트 미설정 시 헤더 미부착 (기존 동작 유지)', () async {
      // SharedPreferences 가 비어 있음 (setUp 에서 초기화)
      final dio = Dio();
      dio.interceptors.add(TenantInterceptor());

      Map<String, dynamic>? captured;
      dio.interceptors.add(InterceptorsWrapper(
        onRequest: (options, handler) {
          captured = Map<String, dynamic>.from(options.headers);
          handler.reject(DioException(
            requestOptions: options,
            type: DioExceptionType.cancel,
          ));
        },
      ));

      try {
        await dio.get('http://example.invalid/test');
      } catch (_) {}

      expect(captured!.containsKey('X-Tenant-ID'), false);
      expect(captured!.containsKey('X-User-ID'), false);
    });

    test('테넌트만 설정되고 사용자는 미설정인 경우', () async {
      await TenantInterceptor.setActiveTenant('only-tenant');

      final dio = Dio();
      dio.interceptors.add(TenantInterceptor());
      Map<String, dynamic>? captured;
      dio.interceptors.add(InterceptorsWrapper(
        onRequest: (options, handler) {
          captured = Map<String, dynamic>.from(options.headers);
          handler.reject(DioException(
            requestOptions: options,
            type: DioExceptionType.cancel,
          ));
        },
      ));
      try {
        await dio.get('http://example.invalid/test');
      } catch (_) {}

      expect(captured!['X-Tenant-ID'], 'only-tenant');
      expect(captured!.containsKey('X-User-ID'), false);
    });

    test('테넌트 전환 후 새 헤더가 부착됨', () async {
      await TenantInterceptor.setActiveTenant('first');
      final dio = Dio();
      dio.interceptors.add(TenantInterceptor());

      String? lastTenant;
      dio.interceptors.add(InterceptorsWrapper(
        onRequest: (options, handler) {
          lastTenant = options.headers['X-Tenant-ID'] as String?;
          handler.reject(DioException(
            requestOptions: options,
            type: DioExceptionType.cancel,
          ));
        },
      ));

      try {
        await dio.get('http://example.invalid/r1');
      } catch (_) {}
      expect(lastTenant, 'first');

      await TenantInterceptor.setActiveTenant('second');
      try {
        await dio.get('http://example.invalid/r2');
      } catch (_) {}
      expect(lastTenant, 'second');
    });
  });
}
