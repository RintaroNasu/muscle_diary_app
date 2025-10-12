import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:frontend/notifiers/auth_notifier.dart';

class SignupState {
  const SignupState({
    this.isLoading = false,
    this.errorMessage,
    this.successMessage,
  });

  final bool isLoading;
  final String? errorMessage;
  final String? successMessage;

  SignupState copyWith({
    bool? isLoading,
    String? errorMessage,
    String? successMessage,
  }) {
    return SignupState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      successMessage: successMessage,
    );
  }
}

final signupControllerProvider =
    StateNotifierProvider<SignupController, SignupState>(
      (ref) => SignupController(ref),
    );

class SignupController extends StateNotifier<SignupState> {
  SignupController(this.ref) : super(const SignupState());

  final Ref ref;

  bool isValidEmail(String email) {
    final trimmedEmail = email.trim();
    final regex = RegExp(r'^[^\s@]+@[^\s@]+\.[^\s@]+$');
    return regex.hasMatch(trimmedEmail);
  }

  Future<void> signup(String email, String password) async {
    state = state.copyWith(
      isLoading: true,
      errorMessage: null,
      successMessage: null,
    );

    await ref.read(authProvider.notifier).signup(email, password);
    final auth = ref.read(authProvider);
    if (auth.isLoggedIn) {
      state = state.copyWith(isLoading: false, successMessage: 'サインアップが完了しました');
    } else {
      state = state.copyWith(isLoading: false);
    }
  }
}
