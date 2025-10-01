import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/pages/workout_record/workout_record_page_controller.dart';

const _exercises = ['ベンチプレス', 'スクワット', 'デッドリフト'];
final _dateFmt = DateFormat('yyyy-MM-dd');

class WorkoutRecordPage extends HookConsumerWidget {
  const WorkoutRecordPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = useMemoized(() => GlobalKey<FormState>());
    final weightController = useTextEditingController();
    final dateController = useTextEditingController();
    final selectedExercise = useState<String?>(null);

    final state = ref.watch(workoutRecordControllerProvider);
    final controller = ref.read(workoutRecordControllerProvider.notifier);

    useListenable(weightController);
    useListenable(dateController);

    ref.listen(workoutRecordControllerProvider, (prev, next) {
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

    String? required(String? v) =>
        (v == null || v.trim().isEmpty) ? '必須項目です' : null;

    String? requiredDouble(String? v) {
      final t = v?.trim() ?? '';
      if (t.isEmpty) return '必須項目です';
      final n = double.tryParse(t);
      if (n == null) return '数値を入力してください';
      if (n <= 0) return '0より大きい値を入力してください';
      return null;
    }

    String? requiredInt(String? v) {
      final t = v?.trim() ?? '';
      if (t.isEmpty) return '必須項目です';
      final n = int.tryParse(t);
      if (n == null) return '整数を入力してください';
      if (n <= 0) return '0より大きい整数を入力してください';
      return null;
    }

    Future<void> pickDate() async {
      final now = DateTime.now();
      final picked = await showDatePicker(
        context: context,
        initialDate: now,
        firstDate: DateTime(now.year - 5),
        lastDate: DateTime(now.year + 5),
      );
      if (picked != null) {
        dateController.text = _dateFmt.format(picked);
      }
    }

    return Scaffold(
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(30),
          child: Form(
            key: formKey,
            child: Column(
              children: [
                TextFormField(
                  controller: weightController,
                  decoration: const InputDecoration(labelText: '体重(kg)'),
                  validator: requiredDouble,
                ),
                const SizedBox(height: 24),
                DropdownButtonFormField(
                  decoration: const InputDecoration(labelText: 'トレーニング名'),
                  initialValue: selectedExercise.value,
                  items: _exercises
                      .map((e) => DropdownMenuItem(value: e, child: Text(e)))
                      .toList(),
                  onChanged: (v) => selectedExercise.value = v,
                  validator: required,
                ),
                const SizedBox(height: 24),
                TextFormField(
                  controller: dateController,
                  readOnly: true,
                  decoration: const InputDecoration(
                    labelText: '実施日',
                    hintText: 'タップして日付を選択',
                    suffixIcon: Icon(Icons.calendar_today),
                  ),
                  onTap: pickDate,
                  validator: required,
                ),
                const SizedBox(height: 32),
                Row(
                  children: [
                    const Text(
                      'セット',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const Spacer(),
                    IconButton(
                      onPressed: state.isSubmitting ? null : controller.addSet,
                      icon: const Icon(Icons.add),
                      tooltip: 'セットを追加',
                    ),
                  ],
                ),
                for (int i = 0; i < state.sets.length; i++) ...[
                  Row(
                    children: [
                      Expanded(
                        flex: 2,
                        child: TextFormField(
                          controller: state.sets[i].weight,
                          decoration: const InputDecoration(
                            labelText: '重量(kg)',
                          ),
                          keyboardType: const TextInputType.numberWithOptions(
                            decimal: true,
                          ),
                          validator: requiredDouble,
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: TextFormField(
                          controller: state.sets[i].reps,
                          decoration: const InputDecoration(labelText: '回数'),
                          keyboardType: TextInputType.number,
                          validator: requiredInt,
                        ),
                      ),
                      IconButton(
                        onPressed: state.isSubmitting
                            ? null
                            : () => controller.removeSet(i),
                        icon: const Icon(Icons.delete_outline),
                        tooltip: 'このセットを削除',
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                ],
                const SizedBox(height: 24),
                ElevatedButton(
                  onPressed: state.isSubmitting
                      ? null
                      : () async {
                          final f = formKey.currentState;
                          if (f == null || !f.validate()) return;
                          if (selectedExercise.value == null) return;

                          await controller.submit(
                            bodyWeight: double.parse(
                              weightController.text.trim(),
                            ),
                            exerciseName: selectedExercise.value!,
                            trainedAtIso: '${dateController.text}T18:00:00Z',
                            onSuccess: () {
                              context.go('/');
                            },
                          );
                        },
                  child: state.isSubmitting
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
