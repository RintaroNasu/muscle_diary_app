import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/utils/trend_grouping.dart';
import 'package:frontend/models/workout_set_item.dart';

void main() {
  WorkoutSetItem make({
    required int recordId,
    required DateTime on,
    required int setNo,
    required int reps,
    required double weight,
    double bodyWeight = 70.0,
  }) {
    return WorkoutSetItem(
      recordId: recordId,
      trainedOn: on,
      setNo: setNo,
      reps: reps,
      exerciseWeight: weight,
      bodyWeight: bodyWeight,
    );
  }

  group('groupSetsByDate', () {
    test('同日のセットが日付キーでまとまり、日付→setNoの順で並ぶ', () {
      final d1 = DateTime(2025, 10, 2, 8);
      final d2 = DateTime(2025, 10, 2, 20);
      final d3 = DateTime(2025, 10, 5, 9);

      final src = [
        make(recordId: 3, on: d3, setNo: 2, reps: 6, weight: 105),
        make(recordId: 1, on: d1, setNo: 2, reps: 8, weight: 95),
        make(recordId: 2, on: d2, setNo: 1, reps: 10, weight: 90),
        make(recordId: 4, on: d3, setNo: 1, reps: 7, weight: 100),
      ];

      final grouped = groupSetsByDate(src);

      expect(grouped.length, 2);
      final key1 = DateTime(2025, 10, 2);
      final key2 = DateTime(2025, 10, 5);
      expect(grouped.containsKey(key1), isTrue);
      expect(grouped.containsKey(key2), isTrue);

      final day2 = grouped[key1]!;
      expect(day2.length, 2);
      expect(day2[0].setNo, 1);
      expect(day2[1].setNo, 2);

      final day5 = grouped[key2]!;
      expect(day5.length, 2);
      expect(day5[0].setNo, 1);
      expect(day5[1].setNo, 2);
    });
  });

  group('toChartGroups', () {
    test('日付昇順で ChartGroup が並び、bars は setNo 昇順', () {
      final keyA = DateTime(2025, 10, 2);
      final keyB = DateTime(2025, 10, 5);

      final grouped = <DateTime, List<WorkoutSetItem>>{
        keyB: [
          make(recordId: 3, on: keyB, setNo: 2, reps: 6, weight: 105),
          make(recordId: 4, on: keyB, setNo: 1, reps: 7, weight: 100),
        ],
        keyA: [
          make(recordId: 2, on: keyA, setNo: 1, reps: 10, weight: 90),
          make(recordId: 1, on: keyA, setNo: 2, reps: 8, weight: 95),
        ],
      };

      final groups = toChartGroups(grouped);
      expect(groups.length, 2);

      expect(groups[0].date, keyA);
      expect(groups[1].date, keyB);

      expect(groups[0].bars.map((b) => b.setNo).toList(), [1, 2]);
      expect(groups[1].bars.map((b) => b.setNo).toList(), [1, 2]);

      expect(groups[0].bars[0].reps, 10);
      expect(groups[0].bars[0].weight, 90);
      expect(groups[1].bars[1].reps, 6);
      expect(groups[1].bars[1].weight, 105);
    });
  });

  group('maxReps', () {
    test('最大repsに+2した値を返す', () {
      final key = DateTime(2025, 10, 2);
      final groups = toChartGroups({
        key: [
          make(recordId: 1, on: key, setNo: 1, reps: 8, weight: 80),
          make(recordId: 2, on: key, setNo: 2, reps: 12, weight: 85),
        ],
      });

      expect(maxReps(groups), 14);
    });

    test('空配列の場合は 2 を返す（0 + 2）', () {
      expect(maxReps(<ChartGroup>[]), 2);
    });
  });
}
