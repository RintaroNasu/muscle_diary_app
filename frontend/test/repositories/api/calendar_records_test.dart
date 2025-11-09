import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:http/http.dart' as http;
import 'package:hooks_riverpod/hooks_riverpod.dart';

import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/api/calendar_records.dart';
import 'package:frontend/controllers/common/auth_controller.dart';
import 'package:frontend/controllers/common/auth_controller.dart' as auth;

class MockApiClient extends Mock implements ApiClient {}

class MockRef extends Mock implements Ref {}

class MockAuthController extends Mock implements AuthNotifier {}

void main() {
  late MockApiClient mockApi;
  late CalendarRecordsApi repo;

  setUpAll(() {
    registerFallbackValue(<String, String>{});
  });

  setUp(() {
    mockApi = MockApiClient();
    repo = CalendarRecordsApi(mockApi);
  });

  group('fetchDayRecords', () {
    test('【正常系】正しい日付で記録を取得できること', () async {
      final d = DateTime(2025, 10, 2);
      final expectedPath = '/training_records/date?date=2025-10-02';

      when(() => mockApi.get(any())).thenAnswer(
        (_) async => http.Response(
          jsonEncode([
            {'id': 1, 'setNo': 1, 'reps': 8, 'weight': 80},
            {'id': 2, 'setNo': 2, 'reps': 6, 'weight': 85},
          ]),
          200,
        ),
      );

      final ref = MockRef();
      final list = await repo.fetchDayRecords(ref, d);

      expect(list, isA<List<Map<String, dynamic>>>());
      expect(list.length, 2);

      verify(() => mockApi.get(expectedPath)).called(1);
    });

    test('【準正常系】200だが配列でない場合は空配列を返すこと', () async {
      final d = DateTime(2025, 10, 2);

      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response(jsonEncode({'ok': true}), 200));

      final ref = MockRef();
      final list = await repo.fetchDayRecords(ref, d);
      expect(list, isEmpty);
    });

    test('【異常系】200だが壊れたJSONの場合はエラーを返すこと', () async {
      final d = DateTime(2025, 10, 2);

      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('{invalid json', 200));

      final ref = MockRef();

      expect(
        () => repo.fetchDayRecords(ref, d),
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
      final d = DateTime(2025, 10, 2);
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('Unauthorized', 401));

      final ref = MockRef();
      final mockAuth = MockAuthController();

      when(() => ref.read(auth.authProvider.notifier)).thenReturn(mockAuth);
      when(() => mockAuth.logout()).thenAnswer((_) async {});

      await expectLater(
        repo.fetchDayRecords(ref, d),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('認証エラー: ログインし直してください'),
          ),
        ),
      );
      verify(() => mockAuth.logout()).called(1);
    });
    test('【異常系】500 の場合はステータスと本文を含む汎用メッセージを返すこと', () async {
      final d = DateTime(2025, 10, 2);
      when(
        () => mockApi.get(any()),
      ).thenAnswer((_) async => http.Response('oops', 500));

      final ref = MockRef();

      expect(
        () => repo.fetchDayRecords(ref, d),
        throwsA(
          isA<Exception>().having(
            (e) => e.toString(),
            'message',
            contains('記録の取得に失敗しました: 500'),
          ),
        ),
      );
    });
  });
}
