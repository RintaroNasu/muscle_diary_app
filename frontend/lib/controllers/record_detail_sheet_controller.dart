import 'package:flutter/foundation.dart';
import 'package:frontend/controllers/common/record_form_controller.dart';
import 'package:frontend/repositories/api/workout_records.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class RecordDetailSheetState {
  const RecordDetailSheetState({
    this.isSubmitting = false,
    this.isDeleting = false,
    this.errorMessage,
    this.successMessage,
  });

  final bool isSubmitting;
  final bool isDeleting;
  final String? errorMessage;
  final String? successMessage;

  RecordDetailSheetState copyWith({
    bool? isSubmitting,
    bool? isDeleting,
    String? errorMessage,
    String? successMessage,
  }) {
    return RecordDetailSheetState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      isDeleting: isDeleting ?? this.isDeleting,
      errorMessage: errorMessage,
      successMessage: successMessage,
    );
  }
}

final recordDetailSheetControllerProvider =
    StateNotifierProvider.family<
      RecordDetailSheetController,
      RecordDetailSheetState,
      String
    >((ref, keyStr) => RecordDetailSheetController(ref, keyStr));

class RecordDetailSheetController
    extends StateNotifier<RecordDetailSheetState> {
  RecordDetailSheetController(this.ref, this.keyStr)
    : super(const RecordDetailSheetState());

  final Ref ref;
  final String keyStr;

  Future<void> submitUpdate({
    required int recordId,
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

      final setsPayload = ref
          .read(recordFormControllerProvider(keyStr).notifier)
          .buildSetsPayload();

      final body = {
        'body_weight': bodyWeight,
        'exercise_id': exerciseId,
        'trained_on': trainedOn,
        'sets': setsPayload,
      };

      final api = ref.read(workoutRecordsApiProvider);
      await api.updateWorkoutRecord(recordId, body);

      state = state.copyWith(isSubmitting: false, successMessage: '記録を更新しました');
      onSuccess?.call();
    } catch (e) {
      state = state.copyWith(isSubmitting: false, errorMessage: '更新に失敗: $e');
    }
  }

  Future<void> deleteRecord({
    required int recordId,
    VoidCallback? onSuccess,
  }) async {
    try {
      state = state.copyWith(
        isDeleting: true,
        errorMessage: null,
        successMessage: null,
      );
      final api = ref.read(workoutRecordsApiProvider);
      await api.deleteWorkoutRecord(recordId);

      state = state.copyWith(isDeleting: false, successMessage: '記録を削除しました');
      onSuccess?.call();
    } catch (e) {
      state = state.copyWith(isDeleting: false, errorMessage: '削除に失敗: $e');
    }
  }
}
