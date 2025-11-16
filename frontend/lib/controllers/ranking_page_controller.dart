import 'package:frontend/models/ranking.dart';
import 'package:frontend/repositories/api/ranking.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

final gymDaysRankingProvider = FutureProvider.autoDispose<List<GymDaysRanking>>(
  (ref) async {
    final api = ref.watch(rankingApiProvider);
    final rows = await api.getMonthlyGymDays();
    return rows.map((json) => GymDaysRanking.fromJson(json)).toList();
  },
);
