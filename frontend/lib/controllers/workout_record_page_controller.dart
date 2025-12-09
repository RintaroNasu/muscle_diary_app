import 'package:frontend/repositories/api/workout_records.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:flutter/material.dart';
import 'package:frontend/controllers/common/record_form_controller.dart';

class WorkoutRecordPageState {
  const WorkoutRecordPageState({
    this.isSubmitting = false,
    this.errorMessage,
    this.successMessage,
  });

  final bool isSubmitting;
  final String? errorMessage;
  final String? successMessage;

  WorkoutRecordPageState copyWith({
    bool? isSubmitting,
    String? errorMessage,
    String? successMessage,
  }) {
    return WorkoutRecordPageState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      errorMessage: errorMessage,
      successMessage: successMessage,
    );
  }
}

final workoutRecordPageControllerProvider =
    StateNotifierProvider<WorkoutRecordPageController, WorkoutRecordPageState>(
      (ref) => WorkoutRecordPageController(ref),
    );

class WorkoutRecordPageController
    extends StateNotifier<WorkoutRecordPageState> {
  WorkoutRecordPageController(this.ref) : super(const WorkoutRecordPageState());

  final Ref ref;

  Future<void> submit({
    required double bodyWeight,
    required int exerciseId,
    required String trainedOn,
    required bool isPublic,
    required String comment,
    VoidCallback? onSuccess,
  }) async {
    try {
      state = state.copyWith(
        isSubmitting: true,
        errorMessage: null,
        successMessage: null,
      );

      final setsPayload = ref
          .read(recordFormControllerProvider('create').notifier)
          .buildSetsPayload();

      final body = {
        'body_weight': bodyWeight,
        'exercise_id': exerciseId,
        'sets': setsPayload,
        'trained_on': trainedOn,
        'is_public': isPublic,
        'comment': comment,
      };
      final api = ref.read(workoutRecordsApiProvider);
      await api.createWorkoutRecord(body);

      state = state.copyWith(isSubmitting: false, successMessage: '記録を保存しました');
      onSuccess?.call();
    } catch (e) {
      state = state.copyWith(isSubmitting: false, errorMessage: '保存に失敗: $e');
    }
  }
}
