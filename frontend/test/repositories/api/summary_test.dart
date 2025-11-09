import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:mocktail/mocktail.dart';

import 'package:frontend/models/summary.dart';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/api/summary.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late MockApiClient mockApi;
  late SummaryApi api;

  setUpAll(() {
    registerFallbackValue(<String, String>{});
  });

  setUp(() {
    mockApi = MockApiClient();
    api = SummaryApi(mockApi);
  });

  group('getHomeSummary', () {
    test('【正常系】サマリーを取得できること', () async {
      when(() => mockApi.get(any())).thenAnswer(
        (_) async => http.Response(
          jsonEncode({
            'total_training_days': 12,
            'latest_weight': 61.5,
            'goal_weight': 58.0,
            'height': 172.0,
          }),
          200,
        ),
      );

      final result = await api.getHomeSummary();

      expect(result, isA<HomeSummary>());
      expect(result?.trainingDays, 12);
      expect(result?.currentWeight, 61.5);
      expect(result?.goalWeight, 58.0);
      expect(result?.height, 172.0);
      verify(() => mockApi.get('/home/summary')).called(1);
    });

    test('【準正常系】200でもMap以外の場合はnullを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response(jsonEncode(['invalid']), 200));

      final result = await api.getHomeSummary();

      expect(result, isNull);
    });

    test('【異常系】200でも壊れたJSONの場合はエラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('{invalid json', 200));

      await expectLater(api.getHomeSummary(), throwsA(isA<FormatException>()));
    });

    test('【異常系】401 の場合は認証エラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      await expectLater(
        api.getHomeSummary(),
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
        api.getHomeSummary(),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('サマリー取得に失敗しました: 500'),
          ),
        ),
      );
    });
  });
}
