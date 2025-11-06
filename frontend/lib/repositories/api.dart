import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;
import 'package:frontend/config/env.dart';

class ApiClient {
  final http.Client _client;
  final String baseUrl;
  final _storage = const FlutterSecureStorage();

  ApiClient({http.Client? client, String? baseUrl})
    : _client = client ?? http.Client(),
      baseUrl = baseUrl ?? Env.apiBaseUrl;

  Future<Map<String, String>> _headers() async {
    final token = await _storage.read(key: 'token');

    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  Future<http.Response> get(String path, {Map<String, String>? headers}) async {
    final headers = await _headers();
    return _client.get(Uri.parse('$baseUrl$path'), headers: headers);
  }

  Future<http.Response> post(
    String path, {
    Map<String, String>? headers,
    Object? body,
  }) async {
    final headers = await _headers();
    return _client.post(
      Uri.parse('$baseUrl$path'),
      headers: headers,
      body: jsonEncode(body),
    );
  }

  Future<http.Response> put(
    String path, {
    Map<String, String>? headers,
    Object? body,
  }) async {
    final headers = await _headers();
    return _client.put(
      Uri.parse('$baseUrl$path'),
      headers: headers,
      body: jsonEncode(body),
    );
  }

  Future<http.Response> delete(
    String path, {
    Map<String, String>? headers,
  }) async {
    final headers = await _headers();
    return _client.delete(Uri.parse('$baseUrl$path'), headers: headers);
  }
}
