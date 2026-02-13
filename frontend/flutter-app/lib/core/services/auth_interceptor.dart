/// JWT 토큰을 gRPC 요청 메타데이터에 자동 첨부하는 클라이언트 인터셉터
import 'package:grpc/grpc.dart';

/// 액세스 토큰을 Authorization 헤더에 붙여 보냄.
/// [tokenProvider]가 null이거나 빈 문자열이면 메타데이터에 토큰을 넣지 않음.
class AuthInterceptor extends ClientInterceptor {
  AuthInterceptor(this.tokenProvider);

  /// 현재 액세스 토큰을 반환 (예: Riverpod에서 읽은 값)
  final String? Function() tokenProvider;

  @override
  ResponseFuture<R> interceptUnary<Q, R>(
    ClientMethod<Q, R> method,
    Q request,
    CallOptions options,
    ClientUnaryInvoker<Q, R> invoker,
  ) {
    final token = tokenProvider();
    final nextOptions = token != null && token.isNotEmpty
        ? options.mergedWith(
            CallOptions(metadata: {'authorization': 'Bearer $token'}),
          )
        : options;
    return invoker(method, request, nextOptions);
  }

  @override
  ResponseStream<R> interceptStreaming<Q, R>(
    ClientMethod<Q, R> method,
    Stream<Q> request,
    CallOptions options,
    ClientStreamingInvoker<Q, R> invoker,
  ) {
    final token = tokenProvider();
    final nextOptions = token != null && token.isNotEmpty
        ? options.mergedWith(
            CallOptions(metadata: {'authorization': 'Bearer $token'}),
          )
        : options;
    return invoker(method, request, nextOptions);
  }
}
