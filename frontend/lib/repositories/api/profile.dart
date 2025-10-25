import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:frontend/repositories/api/auth.dart' show readStoredToken;

const _baseUrl = 'http://localhost:8080';

Future<Map<String, dynamic>?> getProfile() async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのためプロフィールを取得できません（トークンなし）');
  }

  final res = await http.get(
    Uri.parse('$_baseUrl/profile'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
  );

  if (res.statusCode == 200) {
    final data = jsonDecode(res.body);
    if (data is Map<String, dynamic>) return data;
    return null;
  }

  if (res.statusCode == 401) {
    throw Exception('認証エラー: ログインし直してください');
  }

  throw Exception('プロフィール取得に失敗しました: ${res.statusCode}');
}

Future<void> updateProfileApi({
  required double height,
  required double goalWeight,
}) async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのためプロフィールを更新できません（トークンなし）');
  }

  final res = await http.put(
    Uri.parse('$_baseUrl/profile'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
    body: jsonEncode({'height_cm': height, 'goal_weight_kg': goalWeight}),
  );

  if (res.statusCode == 200 || res.statusCode == 204) return;

  if (res.statusCode == 400) {
    try {
      final data = jsonDecode(res.body);
      throw Exception(data['error'] ?? 'リクエストエラー');
    } catch (_) {
      throw Exception('プロフィール更新エラー: ${res.body}');
    }
  }

  if (res.statusCode == 401) {
    throw Exception('認証エラー: ログインし直してください');
  }

  throw Exception('プロフィール更新に失敗しました: ${res.statusCode}');
}
