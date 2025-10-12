import 'package:flutter/material.dart';
import 'package:frontend/repositories/api/workout_records.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class SetCtrls {
  final TextEditingController weight = TextEditingController();
  final TextEditingController reps = TextEditingController();

  Map<String, dynamic> toPayload(int index) => {
    'set': index + 1,
    'reps': int.tryParse(reps.text.trim()) ?? 0,
    'exercise_weight': double.tryParse(weight.text.trim()) ?? 0.0,
  };

  void dispose() {
    weight.dispose();
    reps.dispose();
  }
}

class WorkoutRecordState {
  const WorkoutRecordState({
    this.isSubmitting = false,
    this.errorMessage,
    this.successMessage,
    this.sets = const [],
  });

  final bool isSubmitting;
  final String? errorMessage;
  final String? successMessage;
  final List<SetCtrls> sets;

  WorkoutRecordState copyWith({
    bool? isSubmitting,
    String? errorMessage,
    String? successMessage,
    List<SetCtrls>? sets,
  }) {
    return WorkoutRecordState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      errorMessage: errorMessage,
      successMessage: successMessage,
      sets: sets ?? this.sets,
    );
  }
}

final workoutRecordControllerProvider =
    StateNotifierProvider<WorkoutRecordController, WorkoutRecordState>(
      (ref) => WorkoutRecordController(),
    );

class WorkoutRecordController extends StateNotifier<WorkoutRecordState> {
  WorkoutRecordController() : super(const WorkoutRecordState());

  void addSet() {
    final next = [...state.sets, SetCtrls()];
    state = state.copyWith(
      sets: next,
      errorMessage: null,
      successMessage: null,
    );
  }

  void removeSet(int index) {
    if (state.sets.length <= 1) return;
    final next = [...state.sets];
    final removed = next.removeAt(index);
    removed.dispose();
    state = state.copyWith(
      sets: next,
      errorMessage: null,
      successMessage: null,
    );
  }

  Future<void> submit({
    required double bodyWeight,
    required int exerciseId,
    required String trainedOn,
    VoidCallback? onSuccess,
  }) async {
    try {
      state = state.copyWith(
        isSubmitting: true,
        errorMessage: null,
        successMessage: null,
      );

      final setsPayload = [
        for (int i = 0; i < state.sets.length; i++) state.sets[i].toPayload(i),
      ];

      final body = {
        'body_weight': bodyWeight,
        'exercise_id': exerciseId,
        'sets': setsPayload,
        'trained_on': trainedOn,
      };
      await createWorkoutRecord(body);

      state = state.copyWith(isSubmitting: false, successMessage: '記録を保存しました');
      onSuccess?.call();
    } catch (e) {
      state = state.copyWith(isSubmitting: false, errorMessage: '保存に失敗: $e');
    }
  }

  void resetSets() {
    for (final s in state.sets) {
      s.dispose();
    }
    state = state.copyWith(
      sets: [SetCtrls()],
      errorMessage: null,
      successMessage: null,
    );
  }

  @override
  void dispose() {
    for (final s in state.sets) {
      s.dispose();
    }
    super.dispose();
  }
}
