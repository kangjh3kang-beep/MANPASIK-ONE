// Validators 유틸리티 테스트
// - 이메일, 비밀번호, 표시이름 검증 로직 단위 테스트

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/utils/validators.dart';

void main() {
  group('Validators.validateEmail 테스트', () {
    // null 입력
    test('null 입력 시 에러 메시지를 반환한다', () {
      expect(Validators.validateEmail(null), isNotNull);
    });

    // 빈 문자열
    test('빈 문자열 입력 시 에러 메시지를 반환한다', () {
      final result = Validators.validateEmail('');
      expect(result, '이메일을 입력해주세요');
    });

    // 잘못된 형식 (@ 없음)
    test('@ 기호가 없는 이메일은 유효하지 않다', () {
      final result = Validators.validateEmail('invalid-email');
      expect(result, '올바른 이메일 형식이 아닙니다');
    });

    // 잘못된 형식 (도메인 없음)
    test('도메인이 없는 이메일은 유효하지 않다', () {
      final result = Validators.validateEmail('user@');
      expect(result, isNotNull);
    });

    // 정상 이메일
    test('올바른 이메일 형식은 null을 반환한다', () {
      expect(Validators.validateEmail('user@manpasik.com'), isNull);
      expect(Validators.validateEmail('test.user@example.co.kr'), isNull);
    });
  });

  group('Validators.validatePassword 테스트', () {
    // null 입력
    test('null 입력 시 에러 메시지를 반환한다', () {
      expect(Validators.validatePassword(null), isNotNull);
    });

    // 빈 문자열
    test('빈 문자열 입력 시 에러 메시지를 반환한다', () {
      final result = Validators.validatePassword('');
      expect(result, '비밀번호를 입력해주세요');
    });

    // 8자 미만
    test('8자 미만 비밀번호는 유효하지 않다', () {
      final result = Validators.validatePassword('Ab1');
      expect(result, '비밀번호는 8자 이상이어야 합니다');
    });

    // 영문 미포함
    test('영문자가 없는 비밀번호는 유효하지 않다', () {
      final result = Validators.validatePassword('12345678');
      expect(result, '영문자를 포함해야 합니다');
    });

    // 숫자 미포함
    test('숫자가 없는 비밀번호는 유효하지 않다', () {
      final result = Validators.validatePassword('abcdefgh');
      expect(result, '숫자를 포함해야 합니다');
    });

    // 유효한 비밀번호
    test('영문+숫자 8자 이상 비밀번호는 유효하다', () {
      expect(Validators.validatePassword('Password1'), isNull);
      expect(Validators.validatePassword('abc12345'), isNull);
      expect(Validators.validatePassword('MyP@ssw0rd'), isNull);
    });
  });

  group('Validators.validateDisplayName 테스트', () {
    // null 입력
    test('null 입력 시 에러 메시지를 반환한다', () {
      expect(Validators.validateDisplayName(null), isNotNull);
    });

    // 빈 문자열
    test('빈 문자열 입력 시 에러 메시지를 반환한다', () {
      final result = Validators.validateDisplayName('');
      expect(result, '이름을 입력해주세요');
    });

    // 1자 (너무 짧음)
    test('1자 이름은 유효하지 않다', () {
      final result = Validators.validateDisplayName('A');
      expect(result, '이름은 2~50자 사이여야 합니다');
    });

    // 51자 (너무 김)
    test('51자 이상 이름은 유효하지 않다', () {
      final longName = 'A' * 51;
      final result = Validators.validateDisplayName(longName);
      expect(result, '이름은 2~50자 사이여야 합니다');
    });

    // 정상 이름
    test('2~50자 사이 이름은 유효하다', () {
      expect(Validators.validateDisplayName('홍길동'), isNull);
      expect(Validators.validateDisplayName('AB'), isNull);
      expect(Validators.validateDisplayName('A' * 50), isNull);
    });
  });
}
