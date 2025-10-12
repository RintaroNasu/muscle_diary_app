class WorkoutSetSummary {
  WorkoutSetSummary({
    required this.setNo,
    required this.reps,
    required this.exerciseWeight,
  });

  final int setNo;
  final int reps;
  final double exerciseWeight;

  factory WorkoutSetSummary.fromJson(Map<String, dynamic> json) {
    return WorkoutSetSummary(
      setNo: json['set'] as int,
      reps: json['reps'] as int,
      exerciseWeight: (json['exercise_weight'] as num).toDouble(),
    );
  }
}

class DayRecord {
  DayRecord({
    required this.id,
    required this.exerciseName,
    required this.bodyWeight,
    required this.trainedOn,
    required this.sets,
  });

  final int id;
  final String exerciseName;
  final double bodyWeight;
  final DateTime trainedOn;
  final List<WorkoutSetSummary> sets;

  factory DayRecord.fromJson(Map<String, dynamic> json) {
    return DayRecord(
      id: json['id'] as int,
      exerciseName: json['exercise_name'] as String,
      bodyWeight: (json['body_weight'] as num).toDouble(),
      trainedOn: DateTime.parse(json['trained_on'] as String),
      sets: (json['sets'] as List)
          .map((s) => WorkoutSetSummary.fromJson(s as Map<String, dynamic>))
          .toList(),
    );
  }
}
