import 'package:flutter/material.dart';
import 'package:frontend/controllers/ranking_page_controller.dart';
import 'package:frontend/widgets/ranking_row.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class RankingPage extends HookConsumerWidget {
  const RankingPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(gymDaysRankingProvider);

    return Scaffold(
      body: SafeArea(
        child: state.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (err, _) => Center(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Text(
                '今月のジム日数ランキングの取得に失敗しました。\n${err.toString()}',
                textAlign: TextAlign.center,
              ),
            ),
          ),
          data: (rows) {
            if (rows.isEmpty) {
              return const Center(child: Text('今月のジム日数ランキングはまだありません。'));
            }

            final top5 = rows.take(5).toList();

            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Padding(
                  padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
                  child: Text(
                    '今月のジム日数ランキング',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                ),
                const SizedBox(height: 8),
                Expanded(
                  child: ListView.separated(
                    padding: const EdgeInsets.all(16),
                    itemCount: top5.length,
                    separatorBuilder: (_, __) => const SizedBox(height: 8),
                    itemBuilder: (context, index) {
                      final item = top5[index];
                      final rank = index + 1;
                      return RankingRow(
                        rank: rank,
                        email: item.email.isNotEmpty ? item.email : '名無しのトレーニー',
                        days: item.totalTrainingDays,
                      );
                    },
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }
}
