import 'package:flutter/material.dart';
import 'package:frontend/controllers/common/exercises_provider.dart';
import 'package:frontend/pages/calendar/%20record_detail_sheet.dart';
import 'package:frontend/widgets/unfocus.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import '../../controllers/calendar_page_controller.dart';
import 'package:table_calendar/table_calendar.dart';

String _kg(double v) {
  return '${v % 1 == 0 ? v.toStringAsFixed(0) : v.toStringAsFixed(1)}kg';
}

String _ymd(DateTime d) => '${d.year}/${_two(d.month)}/${_two(d.day)}';
String _two(int n) => n.toString().padLeft(2, '0');

class CalendarPage extends HookConsumerWidget {
  const CalendarPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final selectedDate = ref.watch(selectedDateProvider);
    final dayRecords = ref.watch(dayRecordsProvider);
    final monthMarks = ref.watch(monthHasRecordDaysProvider);
    final exercises = ref.watch(exercisesProvider);

    return UnFocus(
      child: Scaffold(
        appBar: AppBar(title: Text('カレンダー')),
        body: Column(
          children: [
            SizedBox(
              height: 400,
              child: Padding(
                padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
                child: Material(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(16),
                  elevation: 1,
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: TableCalendar(
                      firstDay: DateTime.utc(2020, 1, 1),
                      lastDay: DateTime.utc(2035, 12, 31),
                      focusedDay: selectedDate,
                      selectedDayPredicate: (d) => isSameDay(d, selectedDate),
                      onDaySelected: (selected, focused) {
                        ref
                            .read(selectedDateProvider.notifier)
                            .state = DateTime(
                          selected.year,
                          selected.month,
                          selected.day,
                        );
                      },
                      eventLoader: (day) {
                        final marks = monthMarks.maybeWhen(
                          data: (s) => s,
                          orElse: () => <DateTime>{},
                        );
                        final has = marks.any((d) => isSameDay(d, day));
                        return has ? const ['has'] : const [];
                      },
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 12),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  Text(
                    _ymd(selectedDate),
                    style: TextStyle(fontWeight: FontWeight.w700, fontSize: 16),
                  ),
                  const SizedBox(width: 8),
                  const Text(
                    'の記録',
                    style: TextStyle(fontWeight: FontWeight.w600),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 8),
            Expanded(
              child: dayRecords.when(
                loading: () => const Center(child: CircularProgressIndicator()),
                error: (e, _) => Center(child: Text('読み込みに失敗しました: $e')),
                data: (items) {
                  if (items.isEmpty) {
                    return const Center(child: Text('この日には記録がありません'));
                  }
                  return ListView.separated(
                    padding: const EdgeInsets.fromLTRB(16, 4, 16, 16),
                    itemCount: items.length,
                    separatorBuilder: (_, __) => const SizedBox(height: 10),
                    itemBuilder: (_, i) {
                      final r = items[i];
                      final sets = r.sets
                          .map(
                            (s) =>
                                '${s.setNo}set: ${_kg(s.exerciseWeight)} × ${s.reps}',
                          )
                          .join(' / ');

                      return Material(
                        color: Colors.white,
                        elevation: 1,
                        borderRadius: BorderRadius.circular(12),
                        child: ListTile(
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                          title: Text(
                            r.exerciseName,
                            style: const TextStyle(fontWeight: FontWeight.w600),
                          ),
                          subtitle: Text('BW ${_kg(r.bodyWeight)} • $sets'),
                          onTap: () {
                            exercises.when(
                              data: (data) {
                                showModalBottomSheet(
                                  context: context,
                                  isScrollControlled: true,
                                  backgroundColor: Colors.white,
                                  shape: const RoundedRectangleBorder(
                                    borderRadius: BorderRadius.vertical(
                                      top: Radius.circular(20),
                                    ),
                                  ),
                                  builder: (context) => RecordDetailSheet(
                                    record: r,
                                    exercises: data,
                                  ),
                                );
                              },
                              loading: () {
                                if (context.mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    const SnackBar(
                                      content: Text('種目を読み込み中です…'),
                                    ),
                                  );
                                }
                              },
                              error: (error, _) {
                                if (context.mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                      content: Text('種目の取得に失敗しました: $error'),
                                    ),
                                  );
                                }
                              },
                            );
                          },
                        ),
                      );
                    },
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
