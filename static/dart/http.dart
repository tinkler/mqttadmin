import 'package:dio/dio.dart';
// config file
// baseUrl
// basePicUrl
import './http_config.dart';

class D {
  static D? _instance;

  late Dio dio;

  static D get instance {
    _instance ??= D._init();
    _instance!._iniDio();
    return _instance!;
  }

  D._init();

  _iniDio() {
    _instance!.dio = Dio(
      BaseOptions(baseUrl: baseUrl),
    );
    _instance!.dio.interceptors.add(_ErrorInterceptor());
  }
}

extension _ErrorTypeExtension on DioErrorType {}

class _ErrorInterceptor extends Interceptor {
  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    if (response.data['code'] != 0) {
      handler.reject(_ModelError.fromDioError(
          DioError(
            requestOptions: response.requestOptions,
          ),
          response.data['message']));
      return;
    }
    super.onResponse(response, handler);
  }

  @override
  void onError(DioError err, ErrorInterceptorHandler handler) {
    if (err.response?.statusCode == 404) {
      handler.resolve(Response(
        requestOptions: err.requestOptions,
        data: {
          'code': 1,
          'message': '远程服务器未响应',
        },
      ));
      return;
      // Navigator.pushNamed(context, SiginPage.routeName);
    }
    super.onError(err, handler);
  }
}

class _ModelError extends DioError {
  _ModelError.fromDioError(DioError dioError, String message)
      : super(
            requestOptions: dioError.requestOptions,
            response: dioError.response,
            error: dioError.error,
            type: dioError.type,
            message: message);

  @override
  String toString() {
    if (message != null) {
      return message!;
    }
    return super.toString();
  }
}
