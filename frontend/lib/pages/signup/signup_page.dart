import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/providers/auth_provider.dart';

class SignupPage extends HookConsumerWidget {
  const SignupPage({super.key});

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
    final confirmController = useTextEditingController();
    final authState = ref.watch(authProvider);

    useListenable(emailController);
    useListenable(passwordController);
    useListenable(confirmController);

    final isFormValid =
        emailController.text.trim().isNotEmpty &&
        passwordController.text.isNotEmpty &&
        confirmController.text.isNotEmpty;

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

    Future<void> onSignup() async {
      final currentState = formKey.currentState;
      if (currentState == null) return;
      if (!currentState.validate()) return;

      await ref
          .read(authProvider.notifier)
          .signup(emailController.text.trim(), passwordController.text);
    }

    return UnFocus(
      child: Scaffold(
        appBar: AppBar(title: const Text("新規登録")),
        body: Padding(
          padding: const EdgeInsets.all(16),
          child: Form(
            key: formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextFormField(
                  controller: emailController,
                  decoration: const InputDecoration(labelText: 'メールアドレス'),
                  onChanged: (value) => emailController.text = value,
                  validator: (value) {
                    final v = (value ?? '').trim();
                    if (v.isEmpty) return '必須項目です';
                    if (!_isValidEmail(v)) return 'メールアドレスの形式が正しくありません';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: passwordController,
                  obscureText: true,
                  decoration: const InputDecoration(labelText: 'パスワード'),
                  onChanged: (value) => passwordController.text = value,
                  validator: (value) {
                    final v = (value ?? '').trim();
                    if (v.isEmpty) return '必須項目です';
                    if (v.length < 6) return '6文字以上のパスワードを入力してください';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: confirmController,
                  obscureText: true,
                  decoration: const InputDecoration(labelText: 'パスワード（確認）'),
                  onChanged: (value) => confirmController.text = value,
                  validator: (value) {
                    final v = value ?? '';
                    if (v.isEmpty) return '必須項目です';
                    if (v != passwordController.text) return 'パスワードが一致しません';
                    return null;
                  },
                ),
                const SizedBox(height: 24),
                ElevatedButton(
                  onPressed: (authState.isLoading || !isFormValid)
                      ? null
                      : onSignup,
                  child: authState.isLoading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('新規登録'),
                ),
                const SizedBox(height: 12),
                TextButton(
                  onPressed: () => context.go('/login'),
                  child: const Text('ログインはこちら'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
