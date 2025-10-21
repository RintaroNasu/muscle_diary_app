import 'package:flutter/material.dart';
import 'package:frontend/controllers/common/exercises_provider.dart';
import 'package:frontend/controllers/trend_controller.dart';
import 'package:frontend/models/workout_set_item.dart';
import 'package:frontend/utils/format.dart';
import 'package:frontend/utils/trend_grouping.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';

List<BarChartGroupData> _buildBarGroups(List<ChartGroup> groups) {
  final barGroups = <BarChartGroupData>[];
  for (var i = 0; i < groups.length; i++) {
    final g = groups[i];
    final rods = <BarChartRodData>[];
    for (final b in g.bars) {
      rods.add(
        BarChartRodData(
          toY: b.reps.toDouble(),
          width: 14,
          borderRadius: BorderRadius.circular(4),
        ),
      );
    }
    barGroups.add(
      BarChartGroupData(
        x: i,
        barsSpace: 8,
        barRods: rods,
        groupVertically: false,
      ),
    );
  }
  return barGroups;
}

const double _pxPerGroup = 86;
const double _trailGap = 24;

class TrendPage extends HookConsumerWidget {
  const TrendPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final selectedExerciseId = ref.watch(selectedExerciseIdProvider);
    final exAsync = ref.watch(exercisesProvider);
    final setsAsync = selectedExerciseId == null
        ? const AsyncValue<List<WorkoutSetItem>>.data([])
        : ref.watch(exerciseSetsProvider(selectedExerciseId.toString()));

    return Scaffold(
      appBar: AppBar(title: Text('記録の推移')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              exAsync.when(
                loading: () => const LinearProgressIndicator(),
                error: (e, _) => Row(
                  children: [
                    Expanded(child: Text('種目の取得に失敗しました: $e')),
                    TextButton(
                      onPressed: () => ref.refresh(exercisesProvider),
                      child: const Text('再試行'),
                    ),
                  ],
                ),
                data: (exList) {
                  return Row(
                    children: [
                      const Text(
                        '種目',
                        style: TextStyle(fontWeight: FontWeight.bold),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: DropdownButtonFormField<int>(
                          value: selectedExerciseId,
                          items: exList
                              .map(
                                (e) => DropdownMenuItem<int>(
                                  value: e.id,
                                  child: Text(e.name),
                                ),
                              )
                              .toList(),
                          onChanged: (v) {
                            ref
                                    .read(selectedExerciseIdProvider.notifier)
                                    .state =
                                v;
                          },
                          decoration: const InputDecoration(
                            border: OutlineInputBorder(),
                            contentPadding: EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 8,
                            ),
                            hintText: '種目を選択',
                          ),
                        ),
                      ),
                    ],
                  );
                },
              ),
              const SizedBox(height: 16),
              Expanded(
                child: setsAsync.when(
                  loading: () =>
                      const Center(child: CircularProgressIndicator()),
                  error: (e, _) => Center(child: Text('取得に失敗しました: $e')),
                  data: (items) {
                    if (items.isEmpty) {
                      return Center(child: Text('記録はありません'));
                    }
                    final grouped = groupSetsByDate(items);
                    final groups = toChartGroups(grouped);
                    final maxY = maxReps(groups).toDouble();
                    final chartWidth = (groups.length * _pxPerGroup + _trailGap)
                        .clamp(360.0, double.infinity);
                    return DecoratedBox(
                      decoration: BoxDecoration(
                        border: Border.all(
                          color: Theme.of(context).dividerColor,
                        ),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: SingleChildScrollView(
                        scrollDirection: Axis.horizontal,
                        padding: const EdgeInsets.all(12),
                        child: SizedBox(
                          width: chartWidth,
                          height: 280,
                          child: BarChart(
                            BarChartData(
                              maxY: maxY,
                              barGroups: _buildBarGroups(groups),
                              gridData: FlGridData(
                                show: true,
                                horizontalInterval: 2,
                              ),
                              titlesData: FlTitlesData(
                                leftTitles: AxisTitles(
                                  sideTitles: SideTitles(
                                    showTitles: true,
                                    interval: 2,
                                    reservedSize: 30,
                                  ),
                                ),
                                rightTitles: const AxisTitles(
                                  sideTitles: SideTitles(showTitles: false),
                                ),
                                topTitles: const AxisTitles(
                                  sideTitles: SideTitles(showTitles: false),
                                ),
                                bottomTitles: AxisTitles(
                                  sideTitles: SideTitles(
                                    showTitles: true,
                                    reservedSize: 30,
                                    getTitlesWidget: (value, meta) {
                                      final i = value.toInt();
                                      if (i < 0 || i >= groups.length) {
                                        return const SizedBox.shrink();
                                      }
                                      final d = groups[i].date;
                                      return Padding(
                                        padding: const EdgeInsets.only(top: 6),
                                        child: Text(
                                          '${two(d.month)}/${two(d.day)}',
                                          style: const TextStyle(fontSize: 10),
                                        ),
                                      );
                                    },
                                  ),
                                ),
                              ),
                              barTouchData: BarTouchData(
                                enabled: true,
                                touchTooltipData: BarTouchTooltipData(
                                  tooltipBgColor: Colors.transparent,
                                  getTooltipItem:
                                      (group, groupIndex, rod, rodIndex) {
                                        final g = groups[groupIndex];
                                        final b = g.bars[rodIndex];
                                        final w = b.weight;
                                        final wStr = w % 1 == 0
                                            ? w.toStringAsFixed(0)
                                            : w.toStringAsFixed(1);
                                        return BarTooltipItem(
                                          '$wStr kg\n${b.reps} 回',
                                          const TextStyle(fontSize: 12),
                                        );
                                      },
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                    );
                  },
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
