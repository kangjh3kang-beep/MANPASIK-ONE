/// 입력 검증 유틸리티
///
/// OWASP 기반 입력 검증. 모든 사용자 입력에 적용.
class Validators {
  Validators._();

  /// 이메일 검증
  static String? validateEmail(String? value) {
    if (value == null || value.isEmpty) {
      return '이메일을 입력해주세요';
    }
    final emailRegex = RegExp(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');
    if (!emailRegex.hasMatch(value)) {
      return '올바른 이메일 형식이 아닙니다';
    }
    return null;
  }

  /// 비밀번호 검증 (최소 8자, 영문+숫자+특수문자)
  static String? validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return '비밀번호를 입력해주세요';
    }
    if (value.length < 8) {
      return '비밀번호는 8자 이상이어야 합니다';
    }
    if (!RegExp(r'[A-Za-z]').hasMatch(value)) {
      return '영문자를 포함해야 합니다';
    }
    if (!RegExp(r'[0-9]').hasMatch(value)) {
      return '숫자를 포함해야 합니다';
    }
    return null;
  }

  /// 표시 이름 검증
  static String? validateDisplayName(String? value) {
    if (value == null || value.isEmpty) {
      return '이름을 입력해주세요';
    }
    if (value.length < 2 || value.length > 50) {
      return '이름은 2~50자 사이여야 합니다';
    }
    return null;
  }
}
