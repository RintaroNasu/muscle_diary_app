import 'package:frontend/models/timeline.dart';
import 'package:frontend/repositories/api/timeline.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

final timelineProvider = FutureProvider.autoDispose<List<TimelineItem>>((
  ref,
) async {
  final api = ref.watch(timelineApiProvider);
  return api.fetchTimeline();
});
