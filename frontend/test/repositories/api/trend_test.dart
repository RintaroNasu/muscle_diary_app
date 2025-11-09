import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:mocktail/mocktail.dart';

import 'package:frontend/models/workout_set_item.dart';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/api/trend.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late MockApiClient mockApi;
  late TrendApi api;

  setUpAll(() {
    registerFallbackValue(<String, String>{});
  });

  setUp(() {
    mockApi = MockApiClient();
    api = TrendApi(mockApi);
  });

  group('fetchWorkoutSetsByExercise', () {
    test('【正常系】指定種目のセット一覧を取得できること', () async {
      when(() => mockApi.get(any())).thenAnswer(
        (_) async => http.Response(
          jsonEncode([
            {
              'record_id': 1,
              'trained_on': '2025-10-03',
              'set': 1,
              'reps': 8,
              'exercise_weight': 80,
              'body_weight': 60.5,
            },
            {
              'record_id': 1,
              'trained_on': '2025-10-03',
              'set': 2,
              'reps': 6,
              'exercise_weight': 85,
              'body_weight': 60.5,
            },
          ]),
          200,
        ),
      );

      final list = await api.fetchWorkoutSetsByExercise(42);

      expect(list, isA<List<WorkoutSetItem>>());
      expect(list.length, 2);
      expect(list.first.recordId, 1);
      expect(list.first.setNo, 1);
      expect(list.first.exerciseWeight, 80);
      verify(() => mockApi.get('/training_records/exercises/42')).called(1);
    });

    test('【準正常系】200でも配列でない場合は空配列を返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response(jsonEncode({'ok': true}), 200));

      final list = await api.fetchWorkoutSetsByExercise(1);

      expect(list, isEmpty);
    });

    test('【異常系】200でも壊れたJSONの場合はエラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('{invalid json', 200));

      await expectLater(
        api.fetchWorkoutSetsByExercise(1),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('レスポンスの解析に失敗しました'),
          ),
        ),
      );
    });

    test('【異常系】401 の場合は認証エラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      await expectLater(
        api.fetchWorkoutSetsByExercise(1),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('認証エラー: ログインし直してください'),
          ),
        ),
      );
    });

    test('【異常系】400 でエラーメッセージが含まれる場合はその内容を返すこと', () async {
      when(() => mockApi.get(any())).thenAnswer(
        (_) async =>
            http.Response(jsonEncode({'error': 'invalid exercise id'}), 400),
      );

      await expectLater(
        api.fetchWorkoutSetsByExercise(1),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('invalid exercise id'),
          ),
        ),
      );
    });

    test('【異常系】400 でJSON以外の場合はボディを含むエラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('bad request', 400));

      await expectLater(
        api.fetchWorkoutSetsByExercise(1),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('種目IDが不正です: bad request'),
          ),
        ),
      );
    });

    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('oops', 500));

      await expectLater(
        api.fetchWorkoutSetsByExercise(1),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('取得に失敗しました: 500'),
          ),
        ),
      );
    });
  });
}
