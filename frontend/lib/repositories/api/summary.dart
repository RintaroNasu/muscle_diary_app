import 'dart:convert';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/models/summary.dart';

final _api = ApiClient();

Future<HomeSummary?> getHomeSummary() async {
  final res = await _api.get('/home/summary');

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
