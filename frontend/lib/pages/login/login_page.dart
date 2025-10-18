import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:frontend/controllers/login_page_controller.dart';
import 'package:frontend/controllers/common/auth_controller.dart';

class LoginPage extends HookConsumerWidget {
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = useMemoized(() => GlobalKey<FormState>());
    final emailController = useTextEditingController();
    final passwordController = useTextEditingController();
    final state = ref.watch(loginControllerProvider);
    final controller = ref.read(loginControllerProvider.notifier);

    useListenable(emailController);
    useListenable(passwordController);

    final isFormFilled =
        emailController.text.trim().isNotEmpty &&
        passwordController.text.isNotEmpty;

    ref.listen(authProvider, (prev, next) {
      if (!context.mounted) return;
      if (prev?.isLoggedIn != next.isLoggedIn && next.isLoggedIn) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('ログインが完了しました')));
        context.go('/');
      }
      if (prev?.error != next.error && next.error != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(next.error!), backgroundColor: Colors.red),
        );
      }
    });

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
                    if (!controller.isValidEmail(v)) {
                      return '有効なメールアドレスを入力してください';
                    }
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
                  onPressed: (!isFormFilled || state.isLoading)
                      ? null
                      : () async {
                          final currentState = formKey.currentState;
                          if (currentState == null) return;
                          if (!currentState.validate()) return;

                          await controller.login(
                            emailController.text.trim(),
                            passwordController.text,
                          );
                        },
                  child: state.isLoading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('ログイン'),
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
