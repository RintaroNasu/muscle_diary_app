// lib/repositories/api/timeline.dart
import 'dart:convert';

import 'package:frontend/models/timeline.dart';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/provider.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class TimelineApi {
  final ApiClient _api;
  TimelineApi(this._api);

  Future<List<TimelineItem>> fetchTimeline() async {
    final res = await _api.get('/timeline');

    if (res.statusCode == 200) {
      final data = jsonDecode(res.body);
      if (data is List) {
        return data
            .whereType<Map<String, dynamic>>()
            .map(TimelineItem.fromJson)
            .toList();
      }
      return [];
    }

    if (res.statusCode == 401) {
      throw Exception('認証エラー: ログインし直してください');
    }

    throw Exception('タイムライン取得に失敗しました: ${res.statusCode}');
  }

  Future<void> like(int recordId) async {
    final res = await _api.post('/timeline/$recordId/like');
    if (res.statusCode == 200) return;
    if (res.statusCode == 401) throw Exception('認証エラー: ログインし直してください');
    throw Exception('いいねに失敗しました: ${res.statusCode}');
  }

  Future<void> unlike(int recordId) async {
    final res = await _api.delete('/timeline/$recordId/like');
    if (res.statusCode == 200) return;
    if (res.statusCode == 401) throw Exception('認証エラー: ログインし直してください');
    throw Exception('いいね解除に失敗しました: ${res.statusCode}');
  }
}

final timelineApiProvider = Provider<TimelineApi>((ref) {
  final api = ref.watch(apiClientProvider);
  return TimelineApi(api);
});
