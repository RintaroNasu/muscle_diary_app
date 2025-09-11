import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

const _baseUrl = 'http://localhost:8080';
const _tokenKey = 'token';
final _storage = const FlutterSecureStorage();

Future<String?> readStoredToken() => _storage.read(key: _tokenKey);

Future<String?> loginApi(String email, String password) async {
  final response = await http.post(
    Uri.parse('$_baseUrl/login'),
    headers: {'Content-Type': 'application/json'},
    body: jsonEncode({'email': email, 'password': password}),
  );

  if (response.statusCode == 200) {
    return response.body;
  } else if (response.statusCode == 401) {
    throw Exception('メールアドレスまたはパスワードが間違っています');
  } else {
    throw Exception('ログインに失敗しました: ${response.statusCode}');
  }
}

Future<String?> signupApi(String email, String password) async {
  final response = await http.post(
    Uri.parse('$_baseUrl/signup'),
    headers: {'Content-Type': 'application/json'},
    body: jsonEncode({'email': email, 'password': password}),
  );
  if (response.statusCode == 201) {
    return response.body;
  } else if (response.statusCode == 401) {
    throw Exception('メールアドレスまたはパスワードが間違っています');
  } else {
    throw Exception('サインアップに失敗しました: ${response.statusCode}');
  }
}
