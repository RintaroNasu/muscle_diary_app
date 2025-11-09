import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:mocktail/mocktail.dart';

import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/api/profile.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late MockApiClient mockApi;
  late ProfileApi api;

  setUpAll(() {
    registerFallbackValue(<String, String>{});
  });

  setUp(() {
    mockApi = MockApiClient();
    api = ProfileApi(mockApi);
  });

  group('getProfile', () {
    test('【正常系】プロフィールを取得できること', () async {
      when(() => mockApi.get(any())).thenAnswer(
        (_) async => http.Response(jsonEncode({'height': 170.0}), 200),
      );

      final result = await api.getProfile();

      expect(result, isA<Map<String, dynamic>>());
      expect(result?['height'], 170.0);
      verify(() => mockApi.get('/profile')).called(1);
    });

    test('【準正常系】200でもMap以外の場合はnullを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response(jsonEncode(['invalid']), 200));

      final result = await api.getProfile();

      expect(result, isNull);
    });

    test('【異常系】200でも壊れたJSONの場合はエラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('{invalid json', 200));

      await expectLater(api.getProfile(), throwsA(isA<FormatException>()));
    });

    test('【異常系】401 の場合は認証エラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      await expectLater(
        api.getProfile(),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('認証エラー: ログインし直してください'),
          ),
        ),
      );
    });

    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('oops', 500));

      await expectLater(
        api.getProfile(),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('プロフィール取得に失敗しました: 500'),
          ),
        ),
      );
    });
  });

  group('updateProfileApi', () {
    test('【正常系】身長と目標体重を更新できること', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('', 200));

      await api.updateProfileApi(height: 170, goalWeight: 60);

      verify(
        () => mockApi.put(
          '/profile',
          body: {'height_cm': 170, 'goal_weight_kg': 60},
        ),
      ).called(1);
    });

    test('【異常系】400 でエラーメッセージが含まれる場合はその内容を返すこと', () async {
      when(() => mockApi.put(any(), body: any(named: 'body'))).thenAnswer(
        (_) async =>
            http.Response(jsonEncode({'error': 'invalid parameter'}), 400),
      );

      await expectLater(
        api.updateProfileApi(height: 170, goalWeight: 60),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('invalid parameter'),
          ),
        ),
      );
    });

    test('【異常系】400 でJSON以外の場合はボディを含むエラーを返すこと', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('bad request', 400));

      await expectLater(
        api.updateProfileApi(height: 170, goalWeight: 60),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('プロフィール更新エラー: bad request'),
          ),
        ),
      );
    });

    test('【異常系】401 の場合は認証エラーを返すこと', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      await expectLater(
        api.updateProfileApi(height: 170, goalWeight: 60),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('認証エラー: ログインし直してください'),
          ),
        ),
      );
    });

    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('oops', 500));

      await expectLater(
        api.updateProfileApi(height: 170, goalWeight: 60),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('プロフィール更新に失敗しました: 500'),
          ),
        ),
      );
    });
  });
}
