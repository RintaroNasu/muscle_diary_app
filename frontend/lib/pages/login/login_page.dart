import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/providers/auth_provider.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class LoginPage extends HookConsumerWidget {
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final emailController = useTextEditingController();
    final passwordController = useTextEditingController();
    final isLoading = useState(false);
    final authState = ref.watch(authProvider);

    ref.listen(authProvider, (previous, next) {
      if (previous?.token == null && next.token != null) {
        context.go('/');
      }
    });

    Future<void> onLogin() async {
      await ref
          .read(authProvider.notifier)
          .login(emailController.text.trim(), passwordController.text);
    }

    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextFormField(
              controller: emailController,
              decoration: InputDecoration(labelText: 'メールアドレス'),
              validator: (value) => value?.isEmpty == true ? '必須項目です' : null,
            ),
            TextFormField(
              controller: passwordController,
              obscureText: true,
              decoration: InputDecoration(labelText: 'パスワード'),
              validator: (value) => value?.isEmpty == true ? '必須項目です' : null,
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: authState.isLoading ? null : onLogin,
              child: Text(isLoading.value ? 'ログイン中...' : 'ログイン'),
            ),
            const SizedBox(height: 12),
            TextButton(
              onPressed: () => context.go('/signup'),
              child: const Text('新規登録はこちら'),
            ),
          ],
        ),
      ),
    );
  }
}
