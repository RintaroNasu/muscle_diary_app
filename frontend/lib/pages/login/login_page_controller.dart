import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:frontend/notifiers/auth_notifier.dart';

class LoginState {
  const LoginState({
    this.isLoading = false,
    this.errorMessage,
    this.successMessage,
  });

  final bool isLoading;
  final String? errorMessage;
  final String? successMessage;

  LoginState copyWith({
    bool? isLoading,
    String? errorMessage,
    String? successMessage,
  }) {
    return LoginState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      successMessage: successMessage,
    );
  }
}

final loginControllerProvider =
    StateNotifierProvider<LoginController, LoginState>(
      (ref) => LoginController(ref),
    );

class LoginController extends StateNotifier<LoginState> {
  LoginController(this.ref) : super(const LoginState());

  final Ref ref;

  bool isValidEmail(String email) {
    final trimmedEmail = email.trim();
    final regex = RegExp(r'^[^\s@]+@[^\s@]+\.[^\s@]+$');
    return regex.hasMatch(trimmedEmail);
  }

  Future<void> login(String email, String password) async {
    state = state.copyWith(
      isLoading: true,
      errorMessage: null,
      successMessage: null,
    );
    await ref.read(authProvider.notifier).login(email, password);
    final auth = ref.read(authProvider);

    if (auth.isLoggedIn) {
      state = state.copyWith(isLoading: false, successMessage: 'ログインが完了しました');
    } else {
      state = state.copyWith(isLoading: false);
    }
  }
}
