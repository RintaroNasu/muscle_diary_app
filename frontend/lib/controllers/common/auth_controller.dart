import 'package:frontend/repositories/api/auth.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:jwt_decoder/jwt_decoder.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthState {
  const AuthState({this.isLoading = false, this.token, this.error});
  final bool isLoading;
  final String? token;
  final String? error;
  bool get isLoggedIn {
    if (token == null) return false;
    try {
      return !JwtDecoder.isExpired(token!);
    } catch (e) {
      return false;
    }
  }

  AuthState copyWith({bool? isLoading, String? token, String? error}) =>
      AuthState(
        isLoading: isLoading ?? this.isLoading,
        token: token ?? this.token,
        error: error,
      );
}

class AuthNotifier extends StateNotifier<AuthState> {
  static const _storage = FlutterSecureStorage();

  AuthNotifier() : super(const AuthState()) {
    _restore();
  }
  Future<void> _restore() async {
    final t = await readStoredToken();
    if (t != null) {
      // トークンの有効期限をチェック
      if (JwtDecoder.isExpired(t)) {
        // 期限切れの場合はトークンを削除
        await _storage.delete(key: 'token');
        state = const AuthState();
      } else {
        state = state.copyWith(token: t);
      }
    }
  }

  Future<void> login(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final token = await loginApi(email, password);
      if (token != null) {
        await _storage.write(key: 'token', value: token);
        state = AuthState(isLoading: false, token: token);
      } else {
        state = state.copyWith(isLoading: false, error: 'ログインに失敗しました');
      }
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> signup(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final token = await signupApi(email, password);
      if (token != null) {
        await _storage.write(key: 'token', value: token);
        state = AuthState(isLoading: false, token: token);
      } else {
        state = state.copyWith(isLoading: false, error: 'サインアップに失敗しました');
      }
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> logout() async {
    await _storage.delete(key: 'token');
    state = const AuthState();
  }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>(
  (ref) => AuthNotifier(),
);
