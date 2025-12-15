import 'package:frontend/models/timeline.dart';
import 'package:frontend/repositories/api/timeline.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class TimelineState {
  final AsyncValue<List<TimelineItem>> timelineAsync;
  final String? errorMessage;

  const TimelineState({required this.timelineAsync, this.errorMessage});

  factory TimelineState.initial() =>
      const TimelineState(timelineAsync: AsyncLoading(), errorMessage: null);

  TimelineState copyWith({
    AsyncValue<List<TimelineItem>>? timelineAsync,
    String? errorMessage,
    bool clearError = false,
  }) {
    return TimelineState(
      timelineAsync: timelineAsync ?? this.timelineAsync,
      errorMessage: clearError ? null : (errorMessage ?? this.errorMessage),
    );
  }
}

final timelineControllerProvider =
    StateNotifierProvider.autoDispose<TimelineController, TimelineState>((ref) {
      final api = ref.watch(timelineApiProvider);
      return TimelineController(api)..fetch();
    });

class TimelineController extends StateNotifier<TimelineState> {
  final TimelineApi _api;

  TimelineController(this._api) : super(TimelineState.initial());

  Future<void> fetch() async {
    state = state.copyWith(timelineAsync: const AsyncLoading());
    try {
      final items = await _api.fetchTimeline();
      state = state.copyWith(timelineAsync: AsyncData(items));
    } catch (e) {
      state = state.copyWith(
        timelineAsync: AsyncError(e, StackTrace.current),
        errorMessage: e.toString(),
      );
    }
  }

  void clearError() {
    state = state.copyWith(clearError: true);
  }

  Future<void> toggleLike(int recordId) async {
    final current = state.timelineAsync;
    if (!current.hasValue) return;

    final items = current.value!;
    final index = items.indexWhere((e) => e.recordId == recordId);
    if (index == -1) return;

    final before = items[index];
    final optimistic = before.copyWith(likedByMe: !before.likedByMe);

    // 楽観更新
    final next = [...items];
    next[index] = optimistic;
    state = state.copyWith(timelineAsync: AsyncData(next));

    try {
      if (!before.likedByMe) {
        await _api.like(recordId);
      } else {
        await _api.unlike(recordId);
      }
    } catch (e) {
      final rollback = [...next];
      rollback[index] = before;

      state = state.copyWith(
        timelineAsync: AsyncData(rollback),
        errorMessage: e.toString(),
      );
    }
  }
}
