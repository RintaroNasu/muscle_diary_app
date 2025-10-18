import 'package:flutter/material.dart';
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

class RecordFormState {
  const RecordFormState({this.sets = const []});
  final List<SetCtrls> sets;

  RecordFormState copyWith({List<SetCtrls>? sets}) =>
      RecordFormState(sets: sets ?? this.sets);
}

final recordFormControllerProvider = StateNotifierProvider.autoDispose
    .family<RecordFormController, RecordFormState, String>(
      (ref, key) => RecordFormController(),
    );

class RecordFormController extends StateNotifier<RecordFormState> {
  RecordFormController() : super(const RecordFormState());

  void addSet() {
    final next = [...state.sets, SetCtrls()];
    state = state.copyWith(sets: next);
  }

  void removeSet(int index) {
    final next = [...state.sets];
    final removed = next.removeAt(index);
    removed.dispose();
    state = state.copyWith(sets: next);
  }

  List<Map<String, dynamic>> buildSetsPayload() {
    return [
      for (int i = 0; i < state.sets.length; i++) state.sets[i].toPayload(i),
    ];
  }

  void resetSets() {
    for (final s in state.sets) {
      s.dispose();
    }
    state = state.copyWith(sets: [SetCtrls()]);
  }

  @override
  void dispose() {
    for (final s in state.sets) {
      s.dispose();
    }
    super.dispose();
  }
}
