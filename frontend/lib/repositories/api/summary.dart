import 'dart:convert';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/models/summary.dart';
import 'package:frontend/repositories/provider.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class SummaryApi {
  final ApiClient _api;
  SummaryApi(this._api);

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
}

final summaryApiProvider = Provider<SummaryApi>((ref) {
  final api = ref.watch(apiClientProvider);
  return SummaryApi(api);
});
