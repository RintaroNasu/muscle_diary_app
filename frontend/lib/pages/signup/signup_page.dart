import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/providers/auth_provider.dart';

class SignupPage extends HookConsumerWidget {
  const SignupPage({super.key});

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
      if (passwordController.text != confirmController.text) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('パスワードが一致しません')));
        return;
      }
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
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: passwordController,
                  obscureText: true,
                  decoration: const InputDecoration(labelText: 'パスワード'),
                  onChanged: (value) => passwordController.text = value,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: confirmController,
                  obscureText: true,
                  decoration: const InputDecoration(labelText: 'パスワード（確認）'),
                  onChanged: (value) => confirmController.text = value,
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
