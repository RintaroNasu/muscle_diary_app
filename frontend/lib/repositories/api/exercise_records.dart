import 'dart:convert';
import 'package:frontend/repositories/api.dart';

final _api = ApiClient();

Future<List<Map<String, dynamic>>> getExercises() async {
  final res = await _api.get('/exercises');

  if (res.statusCode == 200) {
    try {
      final data = jsonDecode(res.body);
      if (data is List) {
        return List<Map<String, dynamic>>.from(data);
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
      throw Exception(data['error'] ?? 'リクエストエラー');
    } catch (_) {
      throw Exception('リクエストエラー: ${res.body}');
    }
  }

  throw Exception('エクササイズ一覧の取得に失敗しました: ${res.statusCode} ${res.body}');
}
