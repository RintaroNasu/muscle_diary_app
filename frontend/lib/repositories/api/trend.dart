import 'dart:convert';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/models/workout_set_item.dart';

final _api = ApiClient();

Future<List<WorkoutSetItem>> fetchWorkoutSetsByExercise(int exerciseId) async {
  final res = await _api.get('/training_records/exercises/$exerciseId');

  if (res.statusCode == 200) {
    try {
      final data = jsonDecode(res.body);
      if (data is List) {
        return data.map((e) => WorkoutSetItem.fromJson(e)).toList();
      }
      return [];
    } catch (e) {
      throw Exception('レスポンスの解析に失敗しました: $e');
    }
  }

  if (res.statusCode == 401) {
    throw Exception('認証エラー: ログインし直してください');
  }
  if (res.statusCode == 400) {
    try {
      final data = jsonDecode(res.body);
      throw Exception(data['error'] ?? '種目IDが不正です');
    } catch (_) {
      throw Exception('種目IDが不正です: ${res.body}');
    }
  }

  throw Exception('取得に失敗しました: ${res.statusCode} ${res.body}');
}
