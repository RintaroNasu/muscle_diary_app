import 'dart:convert';
import 'package:frontend/repositories/api.dart';
import 'package:frontend/repositories/provider.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class AuthApi {
  final ApiClient _api;
  AuthApi(this._api);
  Future<String?> loginApi(String email, String password) async {
    final response = await _api.post(
      '/login',
      body: {'email': email, 'password': password},
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['token'];
    } else if (response.statusCode == 401) {
      throw Exception('メールアドレスまたはパスワードが間違っています');
    } else {
      throw Exception('ログインに失敗しました: ${response.statusCode}');
    }
  }

  Future<String?> signupApi(String email, String password) async {
    final response = await _api.post(
      '/signup',
      body: {'email': email, 'password': password},
    );
    if (response.statusCode == 201) {
      final data = jsonDecode(response.body);
      return data['token'];
    } else if (response.statusCode == 401) {
      throw Exception('メールアドレスまたはパスワードが間違っています');
    } else {
      throw Exception('サインアップに失敗しました: ${response.statusCode}');
    }
  }
}

final authApiProvider = Provider<AuthApi>((ref) {
  final api = ref.watch(apiClientProvider);
  return AuthApi(api);
});
