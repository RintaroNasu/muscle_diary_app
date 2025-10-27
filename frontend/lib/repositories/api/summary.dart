import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:frontend/repositories/api/auth.dart' show readStoredToken;
import 'package:frontend/models/summary.dart';

const _baseUrl = 'http://localhost:8080';

Future<HomeSummary?> getHomeSummary() async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのためサマリーを取得できません（トークンなし）');
  }

  final res = await http.get(
    Uri.parse('$_baseUrl/home/summary'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
  );

  if (res.statusCode == 200) {
    final data = jsonDecode(res.body);
    if (data is Map<String, dynamic>) {
      return HomeSummary.fromJson(data);
    }
    return null;
  }

  if (res.statusCode == 401) {
    throw Exception('認証エラー: ログインし直してください');
  }

  throw Exception('サマリー取得に失敗しました: ${res.statusCode}');
}
