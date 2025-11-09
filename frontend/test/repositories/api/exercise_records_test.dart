import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:mocktail/mocktail.dart';

import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/api/exercise_records.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late MockApiClient mockApi;
  late ExerciseRecordsApi repo;

  setUpAll(() {
    registerFallbackValue(<String, String>{});
  });

  setUp(() {
    mockApi = MockApiClient();
    repo = ExerciseRecordsApi(mockApi);
  });

  group('getExercises', () {
    test('【正常系】一覧を取得できること', () async {
      when(() => mockApi.get(any())).thenAnswer(
        (_) async => http.Response(
          jsonEncode([
            {'id': 1, 'name': 'bench press'},
            {'id': 2, 'name': 'squat'},
          ]),
          200,
        ),
      );

      final list = await repo.getExercises();

      expect(list, isA<List<Map<String, dynamic>>>());
      expect(list.length, 2);
      expect(list.first['name'], 'bench press');
      verify(() => mockApi.get('/exercises')).called(1);
    });

    test('【準正常系】200だが配列でない場合は空配列を返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response(jsonEncode({'ok': true}), 200));

      final list = await repo.getExercises();
      expect(list, isEmpty);
    });

    test('【異常系】200だが壊れたJSONの場合はエラーを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('{invalid json', 200));

      expect(
        () => repo.getExercises(),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('レスポンスの解析に失敗しました'),
          ),
        ),
      );
    });

    test('【異常系】401 の場合は認証エラーを投げること', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      expect(
        () => repo.getExercises(),
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
            http.Response(jsonEncode({'error': 'invalid parameter'}), 400),
      );

      expect(
        () => repo.getExercises(),
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
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('bad request', 400));

      expect(
        () => repo.getExercises(),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('リクエストエラー: bad request'),
          ),
        ),
      );
    });

    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('oops', 500));

      expect(
        () => repo.getExercises(),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('エクササイズ一覧の取得に失敗しました: 500'),
          ),
        ),
      );
    });
  });
}
