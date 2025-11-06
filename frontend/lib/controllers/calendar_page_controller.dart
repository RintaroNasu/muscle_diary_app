import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:frontend/models/day_record.dart';
import 'package:frontend/repositories/api/calendar_records.dart';

final selectedDateProvider = StateProvider<DateTime>((ref) {
  final now = DateTime.now();
  return DateTime(now.year, now.month, now.day);
});

final dayRecordsProvider = FutureProvider.autoDispose<List<DayRecord>>((
  ref,
) async {
  final selectedDate = ref.watch(selectedDateProvider);
  final api = ref.read(calendarRecordsApiProvider);
  final records = await api.fetchDayRecords(ref, selectedDate);
  return records.map((json) => DayRecord.fromJson(json)).toList();
});

final monthHasRecordDaysProvider = FutureProvider.autoDispose<Set<DateTime>>((
  ref,
) async {
  final selectedDate = ref.watch(selectedDateProvider);
  final year = selectedDate.year;
  final month = selectedDate.month;
  final api = ref.read(calendarRecordsApiProvider);
  return await api.fetchMonthRecordDays(year, month);
});
