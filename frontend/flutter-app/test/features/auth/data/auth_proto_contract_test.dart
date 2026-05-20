import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/generated/manpasik.pb.dart';

void main() {
  test('LoginResponse carries userId through the official generated proto', () {
    final encoded = (LoginResponse()
          ..accessToken = 'access-token'
          ..refreshToken = 'refresh-token'
          ..tokenType = 'Bearer'
          ..userId = 'user-123')
        .writeToBuffer();

    final decoded = LoginResponse.fromBuffer(encoded);

    expect(decoded.userId, 'user-123');
    expect(decoded.accessToken, 'access-token');
    expect(decoded.refreshToken, 'refresh-token');
  });
}
