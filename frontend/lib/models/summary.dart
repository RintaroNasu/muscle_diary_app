class HomeSummary {
  final int trainingDays;
  final double? currentWeight;
  final double? goalWeight;
  final double? height;

  HomeSummary({
    required this.trainingDays,
    this.currentWeight,
    this.goalWeight,
    this.height,
  });

  factory HomeSummary.fromJson(Map<String, dynamic> json) {
    return HomeSummary(
      trainingDays: json['total_training_days'] as int,
      currentWeight: (json['latest_weight'] as num?)?.toDouble(),
      goalWeight: (json['goal_weight'] as num?)?.toDouble(),
      height: (json['height'] as num?)?.toDouble(),
    );
  }
}
