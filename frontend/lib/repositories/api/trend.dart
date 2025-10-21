import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:frontend/models/workout_set_item.dart';
import 'package:frontend/repositories/api/auth.dart' show readStoredToken;

const _baseUrl = 'http://localhost:8080';

Future<List<WorkoutSetItem>> fetchWorkoutSetsByExercise(int exerciseId) async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのため取得できません（トークンなし）');
  }

  final res = await http.get(
    Uri.parse('$_baseUrl/training_records/exercises/$exerciseId'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
  );

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
