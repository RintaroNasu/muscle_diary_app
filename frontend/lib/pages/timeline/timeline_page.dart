import 'package:flutter/material.dart';
import 'package:frontend/controllers/timeline_page_controller.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class TimelinePage extends HookConsumerWidget {
  const TimelinePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.listen(timelineControllerProvider, (prev, next) {
      final msg = next.errorMessage;
      final prevMsg = prev?.errorMessage;

      if (!context.mounted) return;
      if (msg != null && msg != prevMsg) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(msg), backgroundColor: Colors.red),
        );
        ref.read(timelineControllerProvider.notifier).clearError();
      }
    });

    final state = ref.watch(timelineControllerProvider);
    final timelineAsync = state.timelineAsync;
    final ctrl = ref.read(timelineControllerProvider.notifier);

    return Scaffold(
      appBar: AppBar(title: const Text('タイムライン')),
      body: timelineAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) => Center(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Text(err.toString(), textAlign: TextAlign.center),
          ),
        ),
        data: (items) {
          if (items.isEmpty) {
            return const Center(child: Text('まだタイムラインに投稿はありません'));
          }

          return ListView.separated(
            padding: const EdgeInsets.all(10),
            itemCount: items.length,
            separatorBuilder: (_, __) => const SizedBox(height: 12),
            itemBuilder: (context, index) {
              final item = items[index];
              return Card(
                child: ListTile(
                  title: Text(
                    item.exerciseName,
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const SizedBox(height: 4),
                      Text(
                        '${item.userEmail} ・ ${item.trainedOn}',
                        style: const TextStyle(fontSize: 12),
                      ),
                      if (item.comment != null && item.comment!.isNotEmpty)
                        Padding(
                          padding: const EdgeInsets.only(top: 4),
                          child: Text(item.comment!),
                        ),
                      if (item.bodyWeight != null)
                        Padding(
                          padding: const EdgeInsets.only(top: 4),
                          child: Text(
                            '体重: ${item.bodyWeight!.toStringAsFixed(1)} kg',
                            style: const TextStyle(fontSize: 12),
                          ),
                        ),
                    ],
                  ),
                  trailing: IconButton(
                    icon: Icon(
                      item.likedByMe ? Icons.favorite : Icons.favorite_border,
                    ),
                    onPressed: () => ctrl.toggleLike(item.recordId),
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
