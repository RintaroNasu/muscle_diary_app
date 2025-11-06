import 'dart:convert';
import 'package:frontend/repositories/api.dart';

final _api = ApiClient();

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
    body: ({'email': email, 'password': password}),
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
