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

Future<void> updateWorkoutRecord(
  int recordId,
  Map<String, dynamic> body,
) async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのため記録を更新できません（トークンなし）');
  }

  final res = await http.put(
    Uri.parse('$_baseUrl/training_records/$recordId'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
    body: jsonEncode(body),
  );

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
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのため記録を削除できません（トークンなし）');
  }

  final res = await http.delete(
    Uri.parse('$_baseUrl/training_records/$recordId'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
  );

  if (res.statusCode == 200 || res.statusCode == 204) return;

  if (res.statusCode == 401) {
    throw Exception('認証エラー: ログインし直してください');
  }
  if (res.statusCode == 404) {
    throw Exception('記録が見つかりません');
  }

  throw Exception('記録の削除に失敗しました: ${res.statusCode} ${res.body}');
}
