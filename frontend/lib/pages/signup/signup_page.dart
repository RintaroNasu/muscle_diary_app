import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/controllers/signup_page_controller.dart';
import 'package:frontend/controllers/common/auth_controller.dart';

class SignupPage extends HookConsumerWidget {
  const SignupPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = useMemoized(() => GlobalKey<FormState>());
    final emailController = useTextEditingController();
    final passwordController = useTextEditingController();
    final confirmController = useTextEditingController();
    final state = ref.watch(signupControllerProvider);
    final controller = ref.read(signupControllerProvider.notifier);

    useListenable(emailController);
    useListenable(passwordController);
    useListenable(confirmController);

    final isFormValid =
        emailController.text.trim().isNotEmpty &&
        passwordController.text.isNotEmpty &&
        confirmController.text.isNotEmpty;

    ref.listen(authProvider, (prev, next) {
      if (!context.mounted) return;
      if (prev?.isLoggedIn != next.isLoggedIn && next.isLoggedIn) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('サインアップが完了しました')));
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
                  validator: (value) {
                    final v = (value ?? '').trim();
                    if (v.isEmpty) return '必須項目です';
                    if (!controller.isValidEmail(v)) {
                      return 'メールアドレスの形式が正しくありません';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: passwordController,
                  obscureText: true,
                  decoration: const InputDecoration(labelText: 'パスワード'),
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
                  validator: (value) {
                    final v = value ?? '';
                    if (v.isEmpty) return '必須項目です';
                    if (v != passwordController.text) return 'パスワードが一致しません';
                    return null;
                  },
                ),
                const SizedBox(height: 24),
                ElevatedButton(
                  onPressed: (state.isLoading || !isFormValid)
                      ? null
                      : () async {
                          final currentState = formKey.currentState;
                          if (currentState == null) return;
                          if (!currentState.validate()) return;

                          await controller.signup(
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
