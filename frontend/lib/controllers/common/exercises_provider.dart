import 'package:frontend/models/exercise.dart';
import 'package:frontend/repositories/api/exercise_records.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

final exercisesProvider = FutureProvider<List<Exercise>>((ref) async {
  final api = ref.read(exerciseRecordsApiProvider);
  final list = await api.getExercises();
  return list.map((e) => Exercise.fromJson(e)).toList();
});
