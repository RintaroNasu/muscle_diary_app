class GymDaysRanking {
  final int userId;
  final String email;
  final int totalTrainingDays;

  const GymDaysRanking({
    required this.userId,
    required this.email,
    required this.totalTrainingDays,
  });

  factory GymDaysRanking.fromJson(Map<String, dynamic> json) {
    return GymDaysRanking(
      userId: json['user_id'] as int,
      email: json['email'] as String,
      totalTrainingDays: json['total_training_days'] as int,
    );
  }
}
