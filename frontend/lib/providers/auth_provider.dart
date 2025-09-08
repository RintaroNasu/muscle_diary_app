import 'package:frontend/repositories/api/auth.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class AuthState {
  const AuthState({this.isLoading = false, this.token, this.error});
  final bool isLoading;
  final String? token;
  final String? error;
  bool get isLoggedIn => token != null;
  AuthState copyWith({bool? isLoading, String? token, String? error}) =>
      AuthState(
        isLoading: isLoading ?? this.isLoading,
        token: token ?? this.token,
        error: error,
      );
}

class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier() : super(const AuthState()) {
    _restore();
  }
  Future<void> _restore() async {
    final t = await readStoredToken();
    if (t != null) state = state.copyWith(token: t);
  }

  Future<void> login(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final token = await loginApi(email, password);
      state = AuthState(isLoading: false, token: token);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> signup(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final token = await signupApi(email, password);
      state = AuthState(isLoading: false, token: token);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  // Future<void> logout() async {
  //   await logoutApi();
  //   state = const AuthState();
  // }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>(
  (ref) => AuthNotifier(),
);
