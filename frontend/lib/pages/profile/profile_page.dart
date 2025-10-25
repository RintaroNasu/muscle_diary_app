import 'package:flutter/material.dart';
import 'package:frontend/controllers/profile_page_controller.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:flutter_hooks/flutter_hooks.dart';

class ProfilePage extends HookConsumerWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final profile = ref.watch(profileControllerProvider);
    final profileCtl = ref.read(profileControllerProvider.notifier);

    final heightCtl = useTextEditingController(
      text: profile.height?.toString() ?? '',
    );
    final goalCtl = useTextEditingController(
      text: profile.goalWeight?.toString() ?? '',
    );

    ref.listen(profileControllerProvider, (prev, next) {
      if (!context.mounted) return;
      if (prev?.successMessage != next.successMessage &&
          next.successMessage != null) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(next.successMessage!)));
      }
      if (prev?.error != next.error && next.error != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(next.error!), backgroundColor: Colors.red),
        );
      }
    });

    useEffect(() {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final hasFocus = FocusScope.of(context).hasPrimaryFocus == false;

        if (!heightCtl.selection.isValid ||
            !goalCtl.selection.isValid ||
            hasFocus) {
          if (profile.height != null) {
            heightCtl.text = profile.height!.toString();
          }
          if (profile.goalWeight != null) {
            goalCtl.text = profile.goalWeight!.toString();
          }
        }
      });

      return null;
    }, [profile.height, profile.goalWeight]);

    return UnFocus(
      child: Scaffold(
        appBar: AppBar(title: const Text('マイページ')),
        body: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (profile.isLoading) const LinearProgressIndicator(),
              const SizedBox(height: 8),

              const Text('身長 (cm)'),
              TextField(
                controller: heightCtl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(hintText: '例: 170'),
              ),
              const SizedBox(height: 16),

              const Text('目標体重 (kg)'),
              TextField(
                controller: goalCtl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(hintText: '例: 65'),
              ),
              const SizedBox(height: 24),

              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: profile.isLoading
                      ? null
                      : () async {
                          final h = double.tryParse(heightCtl.text);
                          final g = double.tryParse(goalCtl.text);
                          if (h == null || g == null) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('数値を正しく入力してください')),
                            );
                            return;
                          }
                          await profileCtl.updateProfile(h, g);
                        },
                  child: const Text('保存'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
