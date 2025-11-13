import 'package:flutter/material.dart';
import 'package:frontend/widgets/status_card.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:frontend/controllers/home_page_controller.dart';

class HomePage extends HookConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(homeSummaryProvider);
    final notifier = ref.read(homeSummaryProvider.notifier);

    ref.listen(homeSummaryProvider, (prev, next) {
      if (!context.mounted) return;
      if (prev?.error != next.error && next.error != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(next.error!), backgroundColor: Colors.red),
        );
      }
    });

    String formatToOneDecimal(double? v) {
      if (v == null) return '--';
      return v.toStringAsFixed(1);
    }

    String formatDiffText(double? v) {
      if (v == null) return '--';
      return '${v >= 0 ? '+' : ''}${formatToOneDecimal(v)} kg';
    }

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (state.isLoading) const LinearProgressIndicator(),
            if (state.error != null)
              Padding(
                padding: const EdgeInsets.only(top: 8, bottom: 8),
                child: Row(
                  children: [
                    const Icon(Icons.error_outline, color: Colors.red),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        state.error!,
                        style: const TextStyle(color: Colors.red),
                      ),
                    ),
                    TextButton(
                      onPressed: state.isLoading ? null : notifier.fetchSummary,
                      child: const Text('再試行'),
                    ),
                  ],
                ),
              ),

            SizedBox(
              width: double.infinity,
              height: 100,
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color.fromARGB(255, 21, 148, 253),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                onPressed: () => context.go('/record'),
                child: const Text(
                  '今日の記録を始める',
                  style: TextStyle(fontSize: 30, fontWeight: FontWeight.w900),
                ),
              ),
            ),
            const SizedBox(height: 24),

            StatusCard(
              title: '累計ジム日数',
              child: Text(
                '${state.trainingDays} 日',
                style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),

            StatusCard(
              title: '目標体重 / 現在体重（差分）',
              child: Text(
                '${formatToOneDecimal(state.goalWeight)} kg /'
                '${formatToOneDecimal(state.currentWeight)} kg'
                '（${formatDiffText(state.diffKg)}）',
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),

            StatusCard(
              title: 'BMI',
              child: Text(
                state.bmi == null
                    ? '—（身長 or 体重が未設定）'
                    : '${formatToOneDecimal(state.bmi)}（${state.bmiLabel ?? '—'}）',
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
