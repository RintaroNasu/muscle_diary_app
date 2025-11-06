import 'dart:convert';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/provider.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class ProfileApi {
  final ApiClient _api;
  ProfileApi(this._api);

  Future<Map<String, dynamic>?> getProfile() async {
    final res = await _api.get('/profile');

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
    final res = await _api.put(
      '/profile',
      body: {'height_cm': height, 'goal_weight_kg': goalWeight},
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
}

final profileApiProvider = Provider<ProfileApi>((ref) {
  final api = ref.watch(apiClientProvider);
  return ProfileApi(api);
});
