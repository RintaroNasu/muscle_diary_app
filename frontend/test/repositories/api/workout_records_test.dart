import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:mocktail/mocktail.dart';

import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/api/workout_records.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late MockApiClient mockApi;
  late WorkoutRecordsApi api;

  setUpAll(() {
    registerFallbackValue(<String, String>{});
  });

  setUp(() {
    mockApi = MockApiClient();
    api = WorkoutRecordsApi(mockApi);
  });

  group('createWorkoutRecord', () {
    final body = {
      'trained_on': '2025-10-03',
      'exercise_id': 1,
      'sets': [
        {'set_no': 1, 'reps': 8, 'weight': 80},
      ],
    };

    test('【正常系】記録を作成できること', () async {
      when(
        () => mockApi.post(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('', 201));

      await api.createWorkoutRecord(body);

      verify(() => mockApi.post('/training_records', body: body)).called(1);
    });

    test('【異常系】401 の場合は認証エラーを返すこと', () async {
      when(
        () => mockApi.post(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      await expectLater(
        api.createWorkoutRecord(body),
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
      when(() => mockApi.post(any(), body: any(named: 'body'))).thenAnswer(
        (_) async =>
            http.Response(jsonEncode({'message': 'invalid sets'}), 400),
      );

      await expectLater(
        api.createWorkoutRecord(body),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('invalid sets'),
          ),
        ),
      );
    });

    test('【異常系】400 でJSON以外の場合はボディを含むエラーを返すこと', () async {
      when(
        () => mockApi.post(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('bad request', 400));

      await expectLater(
        api.createWorkoutRecord(body),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('入力エラー: bad request'),
          ),
        ),
      );
    });

    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      when(
        () => mockApi.post(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('oops', 500));

      await expectLater(
        api.createWorkoutRecord(body),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('記録の作成に失敗しました: 500'),
          ),
        ),
      );
    });
  });

  group('updateWorkoutRecord', () {
    final body = {
      'trained_on': '2025-10-03',
      'exercise_id': 1,
      'sets': [
        {'set_no': 1, 'reps': 10, 'weight': 82.5},
      ],
    };

    test('【正常系】記録を更新できること', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('', 200));

      await api.updateWorkoutRecord(99, body);

      verify(() => mockApi.put('/training_records/99', body: body)).called(1);
    });

    test('【異常系】401 の場合は認証エラーを返すこと', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      await expectLater(
        api.updateWorkoutRecord(99, body),
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
      when(() => mockApi.put(any(), body: any(named: 'body'))).thenAnswer(
        (_) async => http.Response(jsonEncode({'message': 'invalid set'}), 400),
      );

      await expectLater(
        api.updateWorkoutRecord(99, body),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('invalid set'),
          ),
        ),
      );
    });

    test('【異常系】400 でJSON以外の場合はボディを含むエラーを返すこと', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('bad request', 400));

      await expectLater(
        api.updateWorkoutRecord(99, body),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('入力エラー: bad request'),
          ),
        ),
      );
    });

    test('【異常系】404 の場合は見つからないエラーを返すこと', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('Not Found', 404));

      await expectLater(
        api.updateWorkoutRecord(99, body),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('記録が見つかりません'),
          ),
        ),
      );
    });

    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      when(
        () => mockApi.put(any(), body: any(named: 'body')),
      ).thenAnswer((_) async => http.Response('oops', 500));

      await expectLater(
        api.updateWorkoutRecord(99, body),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('記録の更新に失敗しました: 500'),
          ),
        ),
      );
    });
  });

  group('deleteWorkoutRecord', () {
    test('【正常系】記録を削除できること', () async {
      when(
        () => mockApi.delete(any()),
      ).thenAnswer((_) async => http.Response('', 204));

      await api.deleteWorkoutRecord(123);

      verify(() => mockApi.delete('/training_records/123')).called(1);
    });

    test('【異常系】401 の場合は認証エラーを返すこと', () async {
      when(
        () => mockApi.delete(any()),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      await expectLater(
        api.deleteWorkoutRecord(123),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('認証エラー: ログインし直してください'),
          ),
        ),
      );
    });

    test('【異常系】404 の場合は見つからないエラーを返すこと', () async {
      when(
        () => mockApi.delete(any()),
      ).thenAnswer((_) async => http.Response('Not Found', 404));

      await expectLater(
        api.deleteWorkoutRecord(123),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('記録が見つかりません'),
          ),
        ),
      );
    });

    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      when(
        () => mockApi.delete(any()),
      ).thenAnswer((_) async => http.Response('oops', 500));

      await expectLater(
        api.deleteWorkoutRecord(123),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('記録の削除に失敗しました: 500'),
          ),
        ),
      );
    });
  });
}
