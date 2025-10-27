import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:frontend/repositories/api/summary.dart';

class HomeSummaryState {
  final bool isLoading;
  final String? error;
  final int trainingDays;
  final double? currentWeight;
  final double? goalWeight;
  final double? height;
  final double? bmi;
  final String? bmiLabel;
  final double? diffKg;

  const HomeSummaryState({
    required this.isLoading,
    this.error,
    this.trainingDays = 0,
    this.currentWeight,
    this.goalWeight,
    this.height,
    this.bmi,
    this.bmiLabel,
    this.diffKg,
  });

  HomeSummaryState copyWith({
    bool? isLoading,
    String? error,
    int? trainingDays,
    double? currentWeight,
    double? goalWeight,
    double? height,
    double? bmi,
    String? bmiLabel,
    double? diffKg,
  }) {
    return HomeSummaryState(
      isLoading: isLoading ?? this.isLoading,
      error: error,
      trainingDays: trainingDays ?? this.trainingDays,
      currentWeight: currentWeight ?? this.currentWeight,
      goalWeight: goalWeight ?? this.goalWeight,
      height: height ?? this.height,
      bmi: bmi ?? this.bmi,
      bmiLabel: bmiLabel ?? this.bmiLabel,
      diffKg: diffKg ?? this.diffKg,
    );
  }
}

class HomeSummaryNotifier extends StateNotifier<HomeSummaryState> {
  HomeSummaryNotifier() : super(const HomeSummaryState(isLoading: false));

  Future<void> fetchSummary() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final s = await getHomeSummary();
      if (s == null) {
        state = state.copyWith(isLoading: false, error: 'データがありません');
        return;
      }

      final bmi = _calcBmi(s.height, s.currentWeight);
      final label = _bmiCategory(bmi);
      final diff = (s.goalWeight != null && s.currentWeight != null)
          ? (s.goalWeight! - s.currentWeight!)
          : null;

      state = state.copyWith(
        isLoading: false,
        trainingDays: s.trainingDays,
        currentWeight: s.currentWeight,
        goalWeight: s.goalWeight,
        height: s.height,
        bmi: bmi,
        bmiLabel: label,
        diffKg: diff == null ? null : double.parse(diff.toStringAsFixed(1)),
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  double? _calcBmi(double? heightCm, double? weightKg) {
    if (heightCm == null ||
        heightCm <= 0 ||
        weightKg == null ||
        weightKg <= 0) {
      return null;
    }
    final h = heightCm / 100;
    final bmi = weightKg / (h * h);
    return double.parse(bmi.toStringAsFixed(1));
  }

  String? _bmiCategory(double? bmi) {
    if (bmi == null) return null;
    if (bmi < 18.5) return 'やせ';
    if (bmi < 25.0) return '普通';
    if (bmi < 30.0) return '肥満(1度)';
    if (bmi < 35.0) return '肥満(2度)';
    if (bmi < 40.0) return '肥満(3度)';
    return '肥満(4度)';
  }
}

final homeSummaryProvider =
    StateNotifierProvider.autoDispose<HomeSummaryNotifier, HomeSummaryState>((ref) {
      final n = HomeSummaryNotifier();
      n.fetchSummary();
      return n;
    });
