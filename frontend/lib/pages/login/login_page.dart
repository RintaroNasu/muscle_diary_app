import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/providers/auth_provider.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class LoginPage extends HookConsumerWidget {
  const LoginPage({super.key});

  bool _isValidEmail(String v) {
    final email = v.trim();
    final regex = RegExp(r'^[^\s@]+@[^\s@]+\.[^\s@]+$');
    return regex.hasMatch(email);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = useMemoized(() => GlobalKey<FormState>());
    final emailController = useTextEditingController();
    final passwordController = useTextEditingController();
    final authState = ref.watch(authProvider);

    useListenable(emailController);
    useListenable(passwordController);

    final isFormFilled =
        emailController.text.trim().isNotEmpty &&
        passwordController.text.isNotEmpty;

    ref.listen<bool>(authProvider.select((s) => s.isLoggedIn), (
      prev,
      loggedIn,
    ) {
      if (loggedIn) {
        context.go('/');
      }
    });
    ref.listen<String?>(authProvider.select((s) => s.error), (prev, err) {
      if (err == null || err == prev) return;
      if (!context.mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text(err), backgroundColor: Colors.red));
    });

    Future<void> onLogin() async {
      final currentState = formKey.currentState;
      if (currentState == null) return;
      if (!currentState.validate()) return;

      await ref
          .read(authProvider.notifier)
          .login(emailController.text.trim(), passwordController.text);
    }

    return UnFocus(
      child: Scaffold(
        body: Padding(
          padding: const EdgeInsets.all(16),
          child: Form(
            key: formKey,

            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextFormField(
                  controller: emailController,
                  decoration: InputDecoration(labelText: 'メールアドレス'),
                  validator: (value) {
                    final v = (value ?? '').trim();
                    if (v.isEmpty) return '必須項目です';
                    if (!_isValidEmail(v)) return '有効なメールアドレスを入力してください';
                    return null;
                  },
                ),
                TextFormField(
                  controller: passwordController,
                  obscureText: true,
                  decoration: InputDecoration(labelText: 'パスワード'),
                  validator: (value) {
                    final v = (value ?? '').trim();
                    if (v.isEmpty) return '必須項目です';
                    if (v.length < 6) return '6文字以上のパスワードを入力してください';
                    return null;
                  },
                ),
                const SizedBox(height: 24),
                ElevatedButton(
                  onPressed: (!isFormFilled || authState.isLoading)
                      ? null
                      : onLogin,
                  child: Text(authState.isLoading ? 'ログイン中...' : 'ログイン'),
                ),
                const SizedBox(height: 12),
                TextButton(
                  onPressed: () => context.go('/signup'),
                  child: const Text('新規登録はこちら'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
