import 'package:flutter/material.dart';
import 'package:frontend/pages/calendar/calendar_page.dart';
import 'package:frontend/pages/home/home_page.dart';
import 'package:frontend/pages/login/login_page.dart';
import 'package:frontend/pages/profile/profile_page.dart';
import 'package:frontend/pages/signup/signup_page.dart';
import 'package:frontend/pages/trend/trend_page.dart';
import 'package:frontend/pages/workout_record/workout_record_page.dart';
import 'package:frontend/notifiers/auth_notifier.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'refresh_listenable.dart';

enum Routes {
  home(path: '/', name: 'home'),
  login(path: '/login', name: 'login'),
  signup(path: '/signup', name: 'signup'),
  calendar(path: '/calendar', name: 'calendar'),
  trends(path: '/trends', name: 'trends'),
  profile(path: '/profile', name: 'profile'),
  record(path: '/record', name: 'record');

  const Routes({required this.path, required this.name});

  final String path;
  final String name;
}

String _titleFor(int index) {
  switch (index) {
    case 0:
      return '筋トレ日記';
    case 1:
      return 'カレンダー';
    case 2:
      return '記録の推移';
    default:
      return '筋トレ日記';
  }
}

int _indexFor(String loc) {
  switch (loc) {
    case '/':
      return 0;
    case '/calendar':
      return 1;
    case '/trends':
      return 2;
    default:
      return 0;
  }
}

final routerProvider = Provider<GoRouter>((ref) {
  final refresh = GoRouterRefreshStream(ref.read(authProvider.notifier).stream);

  return GoRouter(
    initialLocation: Routes.login.path,
    refreshListenable: refresh,
    redirect: (context, state) {
      final authState = ref.read(authProvider);
      final isLoggedIn = authState.isLoggedIn;
      final loc = state.matchedLocation;
      final isAuthRoute =
          (loc == Routes.login.path || loc == Routes.signup.path);

      if (!isLoggedIn && !isAuthRoute) return Routes.login.path;
      if (isLoggedIn && isAuthRoute) return Routes.home.path;

      return null;
    },
    routes: [
      ShellRoute(
        builder: (context, state, child) {
          final loc = state.uri.toString();
          final currentIndex = _indexFor(loc);
          final isAuthRoute =
              (loc == Routes.login.path || loc == Routes.signup.path);

          return Consumer(
            builder: (context, ref, _) {
              final authState = ref.watch(authProvider);
              return Scaffold(
                appBar: AppBar(
                  title: Text(
                    _titleFor(currentIndex),
                    style: const TextStyle(fontSize: 25),
                  ),
                  actions: isAuthRoute
                      ? null
                      : [
                          Padding(
                            padding: const EdgeInsets.only(right: 12),
                            child: IconButton(
                              iconSize: 40,
                              icon: const Icon(Icons.person),
                              onPressed: () {
                                context.go(Routes.profile.path);
                              },
                            ),
                          ),
                          if (authState.isLoggedIn)
                            Padding(
                              padding: const EdgeInsets.only(right: 12),
                              child: IconButton(
                                icon: const Icon(Icons.logout),
                                tooltip: 'ログアウト',
                                onPressed: () async {
                                  await ref
                                      .read(authProvider.notifier)
                                      .logout();
                                  if (context.mounted) {
                                    context.go(Routes.login.path);
                                  }
                                },
                              ),
                            ),
                        ],
                ),
                body: child,
                bottomNavigationBar: isAuthRoute
                    ? null
                    : BottomNavigationBar(
                        currentIndex: currentIndex,
                        onTap: (i) {
                          switch (i) {
                            case 0:
                              context.go(Routes.home.path);
                              break;
                            case 1:
                              context.go(Routes.calendar.path);
                              break;
                            case 2:
                              context.go(Routes.trends.path);
                              break;
                          }
                        },
                        items: const [
                          BottomNavigationBarItem(
                            icon: Icon(Icons.home),
                            label: 'ホーム',
                          ),
                          BottomNavigationBarItem(
                            icon: Icon(Icons.calendar_month),
                            label: 'カレンダー',
                          ),
                          BottomNavigationBarItem(
                            icon: Icon(Icons.show_chart),
                            label: '記録の推移',
                          ),
                        ],
                      ),
              );
            },
          );
        },
        routes: [
          GoRoute(
            path: Routes.home.path,
            pageBuilder: (_, __) => const NoTransitionPage(child: HomePage()),
          ),
          GoRoute(
            path: Routes.calendar.path,
            pageBuilder: (_, __) =>
                const NoTransitionPage(child: CalendarPage()),
          ),
          GoRoute(
            path: Routes.trends.path,
            pageBuilder: (_, __) => const NoTransitionPage(child: TrendPage()),
          ),
          GoRoute(
            path: Routes.record.path,
            pageBuilder: (_, __) =>
                const NoTransitionPage(child: WorkoutRecordPage()),
          ),
          GoRoute(
            path: Routes.login.path,
            pageBuilder: (_, __) => const NoTransitionPage(child: LoginPage()),
          ),
          GoRoute(
            path: Routes.signup.path,
            pageBuilder: (_, __) => const NoTransitionPage(child: SignupPage()),
          ),
          GoRoute(
            path: Routes.profile.path,
            pageBuilder: (_, __) =>
                const NoTransitionPage(child: ProfilePage()),
          ),
        ],
      ),
    ],
  );
});
