import 'package:flutter/material.dart';
import 'package:frontend/pages/calendar/calendar_page.dart';
import 'package:frontend/pages/home/home_page.dart';
import 'package:frontend/pages/trend/trend_page.dart';
import 'package:go_router/go_router.dart';

String _titleFor(int index) {
  switch (index) {
    case 0:
      return 'ホーム';
    case 1:
      return 'カレンダー';
    case 2:
      return '記録の推移';
    default:
      return 'ホーム';
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

final goRouter = GoRouter(
  initialLocation: '/',
  routes: [
    ShellRoute(
      builder: (context, state, child) {
        final loc = state.uri.toString();
        final currentIndex = _indexFor(loc);

        return Scaffold(
          appBar: AppBar(title: Text(_titleFor(currentIndex))),
          body: child,
          bottomNavigationBar: BottomNavigationBar(
            currentIndex: currentIndex,
            onTap: (i) {
              switch (i) {
                case 0:
                  context.go('/');
                  break;
                case 1:
                  context.go('/calendar');
                  break;
                case 2:
                  context.go('/trends');
                  break;
              }
            },
            items: const [
              BottomNavigationBarItem(icon: Icon(Icons.home), label: 'ホーム'),
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
      routes: [
        GoRoute(
          path: '/',
          pageBuilder: (_, __) => const NoTransitionPage(child: HomePage()),
        ),
        GoRoute(
          path: '/calendar',
          pageBuilder: (_, __) => const NoTransitionPage(child: CalendarPage()),
        ),
        GoRoute(
          path: '/trends',
          pageBuilder: (_, __) => const NoTransitionPage(child: TrendPage()),
        ),
      ],
    ),
  ],
);
