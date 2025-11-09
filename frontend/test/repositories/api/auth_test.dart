// test/repositories/auth_api_test.dart
import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;

import 'package:frontend/repositories/provider.dart';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/api/auth.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late MockApiClient mockApi;
  late ProviderContainer container;

  setUpAll(() {
    registerFallbackValue(<String, String>{});
  });

  setUp(() {
    mockApi = MockApiClient();
    container = ProviderContainer(
      overrides: [apiClientProvider.overrideWithValue(mockApi)],
    );
  });

  tearDown(() => container.dispose());

  group('AuthApi.loginApi', () {
    test('【正常系】正しい認証情報でログインできること', () async {
      when(
        () => mockApi.post(
          '/login',
          headers: any(named: 'headers'),
          body: any(named: 'body'),
        ),
      ).thenAnswer(
        (_) async => http.Response(jsonEncode({'token': 'abc'}), 200),
      );

      final auth = container.read(authApiProvider);
      final token = await auth.loginApi('a@example.com', 'pw');

      expect(token, 'abc');
      verify(
        () => mockApi.post(
          '/login',
          headers: any(named: 'headers'),
          body: {'email': 'a@example.com', 'password': 'pw'},
        ),
      ).called(1);
    });

    test('【異常系】存在しないユーザーでログインした場合はエラーを返すこと', () async {
      when(
        () => mockApi.post(
          any(),
          headers: any(named: 'headers'),
          body: any(named: 'body'),
        ),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      final auth = container.read(authApiProvider);

      expect(
        () => auth.loginApi('a@example.com', 'wrong'),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('メールアドレスまたはパスワードが間違っています'),
          ),
        ),
      );
    });

    test('【異常系】サービス側の内部エラーが発生した場合はエラーを返すこと', () async {
      when(
        () => mockApi.post(
          any(),
          headers: any(named: 'headers'),
          body: any(named: 'body'),
        ),
      ).thenAnswer((_) async => http.Response('oops', 500));

      final auth = container.read(authApiProvider);

      expect(
        () => auth.loginApi('a@example.com', 'pw'),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('ログインに失敗しました: 500'),
          ),
        ),
      );
    });
  });

  group('AuthApi.signupApi', () {
    test('【正常系】正しい認証情報でサインアップできること', () async {
      when(
        () => mockApi.post(
          any(),
          headers: any(named: 'headers'),
          body: any(named: 'body'),
        ),
      ).thenAnswer(
        (_) async => http.Response(jsonEncode({'token': 'xyz'}), 201),
      );

      final auth = container.read(authApiProvider);
      final token = await auth.signupApi('b@example.com', 'pw');
      expect(token, 'xyz');

      verify(
        () => mockApi.post(
          '/signup',
          headers: any(named: 'headers'),
          body: {'email': 'b@example.com', 'password': 'pw'},
        ),
      ).called(1);
    });
    test('【異常系】メールアドレスまたはパスワードが間違っている場合はエラーを返すこと', () async {
      when(
        () => mockApi.post(
          any(),
          headers: any(named: 'headers'),
          body: any(named: 'body'),
        ),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      final auth = container.read(authApiProvider);

      expect(
        () => auth.signupApi('b@example.com', 'wrong'),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('メールアドレスまたはパスワードが間違っています'),
          ),
        ),
      );

      verify(
        () => mockApi.post(
          '/signup',
          headers: any(named: 'headers'),
          body: {'email': 'b@example.com', 'password': 'wrong'},
        ),
      ).called(1);
    });
    test('【異常系】サービス側の内部エラーが発生した場合はエラーを返すこと', () async {
      when(
        () => mockApi.post(
          any(),
          headers: any(named: 'headers'),
          body: any(named: 'body'),
        ),
      ).thenAnswer((_) async => http.Response('oops', 500));

      final auth = container.read(authApiProvider);

      expect(
        () => auth.signupApi('b@example.com', 'pw'),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('サインアップに失敗しました: 500'),
          ),
        ),
      );
    });
  });
}
