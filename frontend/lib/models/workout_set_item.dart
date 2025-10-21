class WorkoutSetItem {
  final int recordId;
  final DateTime trainedOn;
  final int setNo;
  final int reps;
  final double exerciseWeight;
  final double bodyWeight;

  WorkoutSetItem({
    required this.recordId,
    required this.trainedOn,
    required this.setNo,
    required this.reps,
    required this.exerciseWeight,
    required this.bodyWeight,
  });

  factory WorkoutSetItem.fromJson(Map<String, dynamic> j) {
    final d = DateTime.parse(j['trained_on']);
    return WorkoutSetItem(
      recordId: j['record_id'] as int,
      trainedOn: DateTime(d.year, d.month, d.day),
      setNo: j['set'] as int,
      reps: j['reps'] as int,
      exerciseWeight: (j['exercise_weight'] as num).toDouble(),
      bodyWeight: (j['body_weight'] as num).toDouble(),
    );
  }
}
