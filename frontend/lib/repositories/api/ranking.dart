import 'dart:convert';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/provider.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class RankingApi {
  final ApiClient _api;
  RankingApi(this._api);

  Future<List<dynamic>> getMonthlyGymDays() async {
    final res = await _api.get('/ranking/monthly_gym_days');
    if (res.statusCode == 200) {
      final data = jsonDecode(res.body);
      if (data is List<dynamic>) {
        return data;
      }
    }

    if (res.statusCode == 401) {
      throw Exception('認証エラー: ログインし直してください');
    }

    throw Exception('サマリー取得に失敗しました: ${res.statusCode}');
  }
}

final rankingApiProvider = Provider<RankingApi>((ref) {
  final api = ref.watch(apiClientProvider);
  return RankingApi(api);
});
