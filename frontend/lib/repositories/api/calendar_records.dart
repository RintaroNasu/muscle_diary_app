import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:frontend/controllers/common/auth_controller.dart';
import 'package:frontend/repositories/api/auth.dart' show readStoredToken;
import 'package:hooks_riverpod/hooks_riverpod.dart';

const _baseUrl = 'http://localhost:8080';

Future<List<Map<String, dynamic>>> fetchDayRecords(
  Ref ref,
  DateTime date,
) async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのため記録を取得できません（トークンなし）');
  }

  final dateStr =
      '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';
  final res = await http.get(
    Uri.parse('$_baseUrl/training_records/date?date=$dateStr'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
  );

  if (res.statusCode == 200) {
    try {
      final data = jsonDecode(res.body);
      if (data is List) {
        return List<Map<String, dynamic>>.from(data);
      }
      return [];
    } catch (e) {
      throw Exception('レスポンスの解析に失敗しました: $e');
    }
  }

  if (res.statusCode == 401) {
    await ref.read(authProvider.notifier).logout();
    throw Exception('認証エラー: ログインし直してください');
  }
  if (res.statusCode == 400) {
    try {
      final data = jsonDecode(res.body);
      throw Exception(data['error'] ?? 'リクエストエラー');
    } catch (_) {
      throw Exception('リクエストエラー: ${res.body}');
    }
  }

  throw Exception('記録の取得に失敗しました: ${res.statusCode} ${res.body}');
}

Future<Set<DateTime>> fetchMonthRecordDays(int year, int month) async {
  final token = await readStoredToken();
  if (token == null || token.isEmpty) {
    throw Exception('未ログインのため記録を取得できません（トークンなし）');
  }

  final res = await http.get(
    Uri.parse(
      '$_baseUrl/training_records/monthly_days?year=$year&month=$month',
    ),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
  );

  if (res.statusCode == 200) {
    try {
      final data = jsonDecode(res.body);
      if (data is List) {
        return data.map<DateTime>((dateStr) {
          final parts = dateStr.split('-');
          return DateTime(
            int.parse(parts[0]),
            int.parse(parts[1]),
            int.parse(parts[2]),
          );
        }).toSet();
      }
      return {};
    } catch (e) {
      throw Exception('レスポンスの解析に失敗しました: $e');
    }
  }

  if (res.statusCode == 401) {
    throw Exception('認証エラー: ログインし直してください');
  }

  throw Exception('記録日の取得に失敗しました: ${res.statusCode} ${res.body}');
}
