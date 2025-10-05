import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:frontend/repositories/api/auth.dart' show readStoredToken;

const _baseUrl = 'http://localhost:8080';

Future<void> createWorkoutRecord(Map<String, dynamic> body) async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのため記録を保存できません（トークンなし）');
  }

  final res = await http.post(
    Uri.parse('$_baseUrl/training_records'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
    body: jsonEncode(body),
  );

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
