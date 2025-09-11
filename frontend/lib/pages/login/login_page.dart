import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/providers/auth_provider.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class LoginPage extends HookConsumerWidget {
  const LoginPage({super.key});

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
                  validator: (value) =>
                      value?.isEmpty == true ? '必須項目です' : null,
                  // onChanged: (value) => email.value = value,
                ),
                TextFormField(
                  controller: passwordController,
                  obscureText: true,
                  decoration: InputDecoration(labelText: 'パスワード'),
                  validator: (value) =>
                      value?.isEmpty == true ? '必須項目です' : null,
                  // onChanged: (value) => password.value = value,
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
