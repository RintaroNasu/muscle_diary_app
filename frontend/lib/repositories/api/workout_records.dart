import 'dart:convert';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/provider.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class WorkoutRecordsApi {
  final ApiClient _api;
  WorkoutRecordsApi(this._api);

  Future<void> createWorkoutRecord(Map<String, dynamic> body) async {
    final res = await _api.post('/training_records', body: body);

    if (res.statusCode == 200 || res.statusCode == 201) return;

    if (res.statusCode == 401) {
      throw Exception('認証エラー: ログインし直してください');
    }
    if (res.statusCode == 400) {
      try {
        final data = jsonDecode(res.body);
        throw Exception(data['message'] ?? '入力エラー');
      } catch (_) {
        throw Exception('入力エラー: ${res.body}');
      }
    }

    throw Exception('記録の作成に失敗しました: ${res.statusCode} ${res.body}');
  }

  Future<void> updateWorkoutRecord(
    int recordId,
    Map<String, dynamic> body,
  ) async {
    final res = await _api.put('/training_records/$recordId', body: body);

    if (res.statusCode == 200 || res.statusCode == 204) return;

    if (res.statusCode == 401) {
      throw Exception('認証エラー: ログインし直してください');
    }
    if (res.statusCode == 400) {
      try {
        final data = jsonDecode(res.body);
        throw Exception(data['message'] ?? '入力エラー');
      } catch (_) {
        throw Exception('入力エラー: ${res.body}');
      }
    }
    if (res.statusCode == 404) {
      throw Exception('記録が見つかりません');
    }

    throw Exception('記録の更新に失敗しました: ${res.statusCode} ${res.body}');
  }

  Future<void> deleteWorkoutRecord(int recordId) async {
    final res = await _api.delete('/training_records/$recordId');

    if (res.statusCode == 200 || res.statusCode == 204) return;

    if (res.statusCode == 401) {
      throw Exception('認証エラー: ログインし直してください');
    }
    if (res.statusCode == 404) {
      throw Exception('記録が見つかりません');
    }

    throw Exception('記録の削除に失敗しました: ${res.statusCode} ${res.body}');
  }
}

final workoutRecordsApiProvider = Provider<WorkoutRecordsApi>((ref) {
  final api = ref.watch(apiClientProvider);
  return WorkoutRecordsApi(api);
});
