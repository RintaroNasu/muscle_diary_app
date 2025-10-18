import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:frontend/controllers/record_detail_sheet_controller.dart';
import 'package:frontend/utils/data_picker.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import 'package:frontend/models/day_record.dart';
import 'package:frontend/models/exercise.dart';

import 'package:frontend/widgets/record/record_form.dart';
import 'package:frontend/controllers/common/record_form_controller.dart';

class RecordDetailSheet extends HookConsumerWidget {
  const RecordDetailSheet({
    super.key,
    required this.record,
    required this.exercises,
  });

  final DayRecord record;
  final List<Exercise> exercises;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final weightCtl = useTextEditingController(
      text: record.bodyWeight.toString(),
    );
    final dateCtl = useTextEditingController(
      text:
          '${record.trainedOn.year}-${record.trainedOn.month.toString().padLeft(2, '0')}-${record.trainedOn.day.toString().padLeft(2, '0')}',
    );
    final selected = useState<Exercise?>(
      exercises
          .where((e) => e.name == record.exerciseName)
          .cast<Exercise?>()
          .firstWhere((_) => true, orElse: () => null),
    );

    final keyStr = useMemoized(() => 'edit:${record.id}', [record.id]);

    final formState = ref.watch(recordFormControllerProvider(keyStr));
    final formCtl = ref.read(recordFormControllerProvider(keyStr).notifier);

    final detailState = ref.watch(recordDetailSheetControllerProvider(keyStr));
    final detailCtl = ref.read(
      recordDetailSheetControllerProvider(keyStr).notifier,
    );

    useEffect(() {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (formState.sets.isNotEmpty) return;

        formCtl.resetSets();
        final need = (record.sets.length - 1);
        for (int i = 0; i < need; i++) {
          formCtl.addSet();
        }

        final sets = ref.read(recordFormControllerProvider(keyStr)).sets;
        for (int i = 0; i < sets.length && i < record.sets.length; i++) {
          sets[i].weight.text = record.sets[i].exerciseWeight.toString();
          sets[i].reps.text = record.sets[i].reps.toString();
        }
      });

      return null;
    }, []);

    return UnFocus(
      child: Padding(
        padding: EdgeInsets.only(
          left: 16,
          right: 16,
          top: 16,
          bottom: MediaQuery.of(context).viewInsets.bottom + 16,
        ),
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text(
                '記録の編集',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),

              RecordForm(
                formKey: GlobalKey<FormState>(),
                weightController: weightCtl,
                dateController: dateCtl,
                selectedExercise: selected.value,
                exercises: exercises,
                onSelectExercise: (v) => selected.value = v,
                onPickDate: () async {
                  final ymd = await pickDateAsYmd(context);
                  if (ymd != null) dateCtl.text = ymd;
                },
                sets: formState.sets,
                onAddSet: formCtl.addSet,
                onRemoveSet: formCtl.removeSet,
                isSubmitting: false,
              ),

              const SizedBox(height: 24),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      style: OutlinedButton.styleFrom(
                        foregroundColor: Colors.red,
                        side: const BorderSide(color: Colors.red),
                      ),
                      onPressed: detailState.isDeleting
                          ? null
                          : () async {
                              await detailCtl.deleteRecord(
                                recordId: record.id,
                                onSuccess: () {
                                  if (!context.mounted) return;
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    const SnackBar(content: Text('記録を削除しました')),
                                  );
                                  Navigator.of(context).pop();
                                },
                              );
                              if (!context.mounted) return;
                              if (detailState.errorMessage != null) {
                                ScaffoldMessenger.of(context).showSnackBar(
                                  SnackBar(
                                    content: Text(detailState.errorMessage!),
                                    backgroundColor: Colors.red,
                                  ),
                                );
                              }
                            },
                      child: const Text('削除する'),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: ElevatedButton(
                      onPressed: detailState.isSubmitting
                          ? null
                          : () async {
                              if (selected.value == null) return;

                              await detailCtl.submitUpdate(
                                recordId: record.id,
                                bodyWeight:
                                    double.tryParse(weightCtl.text.trim()) ??
                                    0.0,
                                exerciseId: selected.value!.id,
                                trainedOn: dateCtl.text.trim(),
                                onSuccess: () {
                                  if (!context.mounted) return;
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    const SnackBar(content: Text('記録を更新しました')),
                                  );
                                  Navigator.of(context).pop();
                                },
                              );
                              if (!context.mounted) return;
                              if (detailState.errorMessage != null) {
                                ScaffoldMessenger.of(context).showSnackBar(
                                  SnackBar(
                                    content: Text(detailState.errorMessage!),
                                    backgroundColor: Colors.red,
                                  ),
                                );
                              }
                            },
                      child: const Text('保存する'),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
