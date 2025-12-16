// lib/models/timeline.dart
class TimelineItem {
  final int recordId;
  final int userId;
  final String userEmail;
  final String exerciseName;
  final double? bodyWeight;
  final String trainedOn;
  final String? comment;
  final bool likedByMe;

  TimelineItem({
    required this.recordId,
    required this.userId,
    required this.userEmail,
    required this.exerciseName,
    this.bodyWeight,
    required this.trainedOn,
    this.comment,
    required this.likedByMe,
  });

  TimelineItem copyWith({bool? likedByMe}) {
    return TimelineItem(
      recordId: recordId,
      userId: userId,
      userEmail: userEmail,
      exerciseName: exerciseName,
      bodyWeight: bodyWeight,
      trainedOn: trainedOn,
      comment: comment,
      likedByMe: likedByMe ?? this.likedByMe,
    );
  }

  factory TimelineItem.fromJson(Map<String, dynamic> json) {
    return TimelineItem(
      recordId: json['record_id'] as int,
      userId: json['user_id'] as int,
      userEmail: json['user_email'] as String,
      exerciseName: json['exercise_name'] as String,
      bodyWeight: (json['body_weight'] as num?)?.toDouble(),
      trainedOn: json['trained_on'] as String,
      comment: json['comment'] as String?,
      likedByMe: json['liked_by_me'] as bool,
    );
  }
}
