import 'package:frontend/models/exercise.dart';
import 'package:frontend/repositories/api/exercise_records.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

final exercisesProvider = FutureProvider<List<Exercise>>((ref) async {
  final list = await getExercises();
  return list.map((e) => Exercise.fromJson(e)).toList();
});
