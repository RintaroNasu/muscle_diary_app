import 'package:flutter/material.dart';
import 'package:frontend/controllers/common/record_form_controller.dart';
import 'package:frontend/models/exercise.dart';

class RecordForm extends StatelessWidget {
  const RecordForm({
    super.key,
    required this.formKey,
    required this.weightController,
    required this.dateController,
    required this.selectedExercise,
    required this.exercises,
    required this.onSelectExercise,
    required this.onPickDate,
    required this.sets,
    required this.onAddSet,
    required this.onRemoveSet,
    this.isSubmitting = false,
  });

  final GlobalKey<FormState> formKey;

  final TextEditingController weightController;
  final TextEditingController dateController;

  final Exercise? selectedExercise;
  final List<Exercise> exercises;
  final ValueChanged<Exercise?> onSelectExercise;

  final VoidCallback onPickDate;

  final List<SetCtrls> sets;
  final VoidCallback onAddSet;
  final void Function(int index) onRemoveSet;

  final bool isSubmitting;

  String? _required(String? v) =>
      (v == null || v.trim().isEmpty) ? '必須項目です' : null;

  String? _requiredDouble(String? v) {
    final t = v?.trim() ?? '';
    if (t.isEmpty) return '必須項目です';
    final n = double.tryParse(t);
    if (n == null) return '数値を入力してください';
    if (n <= 0) return '0より大きい値を入力してください';
    return null;
  }

  String? _requiredInt(String? v) {
    final t = v?.trim() ?? '';
    if (t.isEmpty) return '必須項目です';
    final n = int.tryParse(t);
    if (n == null) return '整数を入力してください';
    if (n <= 0) return '0より大きい整数を入力してください';
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return Form(
      key: formKey,
      child: Column(
        children: [
          TextFormField(
            controller: weightController,
            decoration: const InputDecoration(labelText: '体重(kg)'),
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            validator: _requiredDouble,
          ),
          const SizedBox(height: 24),

          DropdownButtonFormField<Exercise>(
            value: selectedExercise,
            items: exercises
                .map((e) => DropdownMenuItem(value: e, child: Text(e.name)))
                .toList(),
            onChanged: onSelectExercise,
            decoration: const InputDecoration(labelText: '種目'),
            validator: (v) => v == null ? '必須項目です' : null,
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
            onTap: isSubmitting ? null : onPickDate,
            validator: _required,
          ),
          const SizedBox(height: 32),

          Row(
            children: [
              const Text(
                'セット',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              ),
              const Spacer(),
              IconButton(
                onPressed: isSubmitting ? null : onAddSet,
                icon: const Icon(Icons.add),
                tooltip: 'セットを追加',
              ),
            ],
          ),
          for (int i = 0; i < sets.length; i++) ...[
            Row(
              children: [
                Expanded(
                  flex: 2,
                  child: TextFormField(
                    controller: sets[i].weight,
                    decoration: const InputDecoration(labelText: '重量(kg)'),
                    keyboardType: const TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    validator: _requiredDouble,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: TextFormField(
                    controller: sets[i].reps,
                    decoration: const InputDecoration(labelText: '回数'),
                    keyboardType: TextInputType.number,
                    validator: _requiredInt,
                  ),
                ),
                IconButton(
                  onPressed: isSubmitting ? null : () => onRemoveSet(i),
                  icon: const Icon(Icons.delete_outline),
                  tooltip: 'このセットを削除',
                ),
              ],
            ),
            const SizedBox(height: 12),
          ],
        ],
      ),
    );
  }
}
