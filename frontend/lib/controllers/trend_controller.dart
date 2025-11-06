import 'package:frontend/models/workout_set_item.dart';
import 'package:frontend/repositories/api/trend.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

final selectedExerciseIdProvider = StateProvider<int?>((ref) => null);

final exerciseSetsProvider = FutureProvider.autoDispose
    .family<List<WorkoutSetItem>, String>((ref, exerciseId) async {
      final id = int.parse(exerciseId);
      final api = ref.read(trendApiProvider);
      final items = await api.fetchWorkoutSetsByExercise(id);

      items.sort((a, b) {
        final c = a.trainedOn.compareTo(b.trainedOn);
        if (c != 0) return c;
        return a.setNo.compareTo(b.setNo);
      });
      return items;
    });
