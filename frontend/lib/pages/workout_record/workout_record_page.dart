import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/controllers/common/exercises_provider.dart';
import 'package:frontend/controllers/common/record_form_controller.dart';
import 'package:frontend/models/exercise.dart';
import 'package:frontend/utils/data_picker.dart';
import 'package:frontend/widgets/record/record_form.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/controllers/workout_record_page_controller.dart';

class WorkoutRecordPage extends HookConsumerWidget {
  const WorkoutRecordPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = useMemoized(() => GlobalKey<FormState>());
    final weightController = useTextEditingController();
    final dateController = useTextEditingController();

    final selectedExercise = useState<Exercise?>(null);

    final formState = ref.watch(recordFormControllerProvider('create'));
    final formCtl = ref.read(recordFormControllerProvider('create').notifier);
    final pageState = ref.watch(workoutRecordPageControllerProvider);
    final pageCtl = ref.read(workoutRecordPageControllerProvider.notifier);
    final exercises = ref.watch(exercisesProvider);

    useListenable(weightController);
    useListenable(dateController);

    ref.listen(exercisesProvider, (prev, next) {
      next.when(
        data: (_) {},
        loading: () {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (!context.mounted) return;
            ScaffoldMessenger.of(
              context,
            ).showSnackBar(const SnackBar(content: Text('種目を読み込み中です…')));
          });
        },
        error: (error, _) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (!context.mounted) return;
            ScaffoldMessenger.of(
              context,
            ).showSnackBar(SnackBar(content: Text('種目の取得に失敗しました: $error')));
          });
        },
      );
    });

    ref.listen(workoutRecordPageControllerProvider, (prev, next) {
      if (!context.mounted) return;
      if (prev?.successMessage != next.successMessage &&
          next.successMessage != null) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(next.successMessage!)));
      }
      if (prev?.errorMessage != next.errorMessage &&
          next.errorMessage != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(next.errorMessage!),
            backgroundColor: Colors.red,
          ),
        );
      }
    });

    return UnFocus(
      child: Scaffold(
        body: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(30),
            child: Column(
              children: [
                RecordForm(
                  formKey: formKey,
                  weightController: weightController,
                  dateController: dateController,
                  selectedExercise: selectedExercise.value,
                  exercises: exercises.maybeWhen(
                    data: (data) => data,
                    orElse: () => [],
                  ),
                  onSelectExercise: (v) => selectedExercise.value = v,
                  onPickDate: () async {
                    final ymd = await pickDateAsYmd(context);
                    if (ymd != null) dateController.text = ymd;
                  },
                  sets: formState.sets,
                  onAddSet: formCtl.addSet,
                  onRemoveSet: formCtl.removeSet,
                  isSubmitting: pageState.isSubmitting,
                ),
                ElevatedButton(
                  onPressed: pageState.isSubmitting
                      ? null
                      : () async {
                          final f = formKey.currentState;
                          if (f == null || !f.validate()) return;
                          if (selectedExercise.value == null) return;
                          print(dateController.text);
                          await pageCtl.submit(
                            bodyWeight: double.parse(
                              weightController.text.trim(),
                            ),
                            exerciseId: selectedExercise.value!.id,
                            trainedOn: dateController.text,
                            onSuccess: () {
                              context.go('/');
                            },
                          );
                        },
                  child: pageState.isSubmitting
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('記録する'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
