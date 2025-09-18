import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:intl/intl.dart';

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

    final sets = useState<List<Map<String, TextEditingController>>>([
      {'weight': TextEditingController(), 'reps': TextEditingController()},
    ]);

    useEffect(() {
      return () {
        for (final set in sets.value) {
          set['weight']?.dispose();
          set['reps']?.dispose();
        }
      };
    }, []);

    useListenable(weightController);
    useListenable(dateController);

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

    void addSet() {
      sets.value = [
        ...sets.value,
        {'weight': TextEditingController(), 'reps': TextEditingController()},
      ];
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
                  value: selectedExercise.value,
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

                ...sets.value.asMap().entries.map((entry) {
                  final index = entry.key;
                  final set = entry.value;
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '${index + 1}セット目',
                        style: const TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: set['weight']!,
                              decoration: const InputDecoration(
                                labelText: '重量(kg)',
                              ),
                              validator: requiredDouble,
                              keyboardType:
                                  const TextInputType.numberWithOptions(
                                    decimal: true,
                                  ),
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: TextFormField(
                              controller: set['reps']!,
                              decoration: const InputDecoration(
                                labelText: '回数',
                              ),
                              validator: requiredInt,
                              keyboardType: TextInputType.number,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),
                    ],
                  );
                }),
                OutlinedButton.icon(
                  onPressed: addSet,
                  icon: const Icon(Icons.add),
                  label: const Text('セットを追加'),
                ),
                const SizedBox(height: 24),
                ElevatedButton(
                  child: const Text('記録する'),
                  onPressed: () {
                    final currentState = formKey.currentState;
                    if (currentState == null) return;
                    if (!currentState.validate()) return;

                    final requestBody = {
                      "body_weight":
                          double.tryParse(weightController.text) ?? 0.0,
                      "exercise_name": selectedExercise.value ?? "",
                      "sets": sets.value.asMap().entries.map((entry) {
                        final index = entry.key;
                        final set = entry.value;
                        return {
                          "set": index + 1,
                          "reps": int.tryParse(set['reps']!.text) ?? 0,
                          "exercise_weight":
                              double.tryParse(set['weight']!.text) ?? 0.0,
                        };
                      }).toList(),
                      "trained_at": "${dateController.text}T18:00:00Z",
                    };
                    print(requestBody);
                  },
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
