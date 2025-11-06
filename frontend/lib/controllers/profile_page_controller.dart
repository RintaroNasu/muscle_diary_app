import 'package:frontend/repositories/api/profile.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class ProfileState {
  final double? height;
  final double? goalWeight;
  final bool isLoading;
  final String? successMessage;
  final String? error;

  const ProfileState({
    this.height,
    this.goalWeight,
    this.isLoading = false,
    this.successMessage,
    this.error,
  });

  ProfileState copyWith({
    double? height,
    double? goalWeight,
    bool? isLoading,
    String? successMessage,
    String? error,
  }) {
    return ProfileState(
      height: height ?? this.height,
      goalWeight: goalWeight ?? this.goalWeight,
      isLoading: isLoading ?? this.isLoading,
      successMessage: successMessage,
      error: error,
    );
  }
}

class ProfileNotifier extends StateNotifier<ProfileState> {
  final ProfileApi _profileApi;
  ProfileNotifier(this._profileApi) : super(const ProfileState());

  Future<void> loadProfile() async {
    try {
      state = state.copyWith(
        isLoading: true,
        error: null,
        successMessage: null,
      );
      final data = await _profileApi.getProfile();
      if (data != null) {
        state = state.copyWith(
          height: (data['height_cm'] as num?)?.toDouble(),
          goalWeight: (data['goal_weight_kg'] as num?)?.toDouble(),
          isLoading: false,
        );
      } else {
        state = state.copyWith(isLoading: false);
      }
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> updateProfile(double height, double goalWeight) async {
    try {
      state = state.copyWith(
        isLoading: true,
        error: null,
        successMessage: null,
      );
      await _profileApi.updateProfileApi(
        height: height,
        goalWeight: goalWeight,
      );
      state = state.copyWith(
        height: height,
        goalWeight: goalWeight,
        isLoading: false,
        successMessage: 'プロフィールを保存しました',
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}

final profileControllerProvider =
    StateNotifierProvider<ProfileNotifier, ProfileState>((ref) {
      final profileApi = ref.watch(profileApiProvider);
      final notifier = ProfileNotifier(profileApi);
      notifier.loadProfile();
      return notifier;
    });
