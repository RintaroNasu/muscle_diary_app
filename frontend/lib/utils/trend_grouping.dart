import 'package:frontend/models/workout_set_item.dart';

class ChartGroup {
  final DateTime date;
  final List<BarDatum> bars;
  const ChartGroup({required this.date, required this.bars});
}

class BarDatum {
  final int setNo;
  final int reps;
  final double weight;
  const BarDatum({
    required this.setNo,
    required this.reps,
    required this.weight,
  });
}

Map<DateTime, List<WorkoutSetItem>> groupSetsByDate(List<WorkoutSetItem> src) {
  final list = [...src]
    ..sort((a, b) {
      final c = a.trainedOn.compareTo(b.trainedOn);
      return c != 0 ? c : a.setNo.compareTo(b.setNo);
    });

  final map = <DateTime, List<WorkoutSetItem>>{};
  for (final it in list) {
    final t = it.trainedOn.toLocal();
    final key = DateTime(t.year, t.month, t.day);
    (map[key] ??= []).add(it);
  }
  for (final key in map.keys) {
    map[key]!.sort((a, b) => a.setNo.compareTo(b.setNo));
  }
  return map;
}

List<ChartGroup> toChartGroups(Map<DateTime, List<WorkoutSetItem>> grouped) {
  final dates = grouped.keys.toList()..sort();
  final out = <ChartGroup>[];
  for (final d in dates) {
    final sets = [...grouped[d]!]..sort((a, b) => a.setNo.compareTo(b.setNo));
    out.add(
      ChartGroup(
        date: d,
        bars: [
          for (final s in sets)
            BarDatum(setNo: s.setNo, reps: s.reps, weight: s.exerciseWeight),
        ],
      ),
    );
  }
  return out;
}

int maxReps(List<ChartGroup> groups) {
  if (groups.isEmpty) return 2;

  var m = 0;
  for (final g in groups) {
    for (final b in g.bars) {
      if (b.reps > m) m = b.reps;
    }
  }
  return m + 2;
}
