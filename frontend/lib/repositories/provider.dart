import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:frontend/repositories/api.dart';

final apiClientProvider = Provider<ApiClient>((ref) => ApiClient());
