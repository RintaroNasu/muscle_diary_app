import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:intl/intl.dart';

import 'package:frontend/utils/data_picker.dart';

void main() {
  Widget host({required void Function(BuildContext) onTap}) {
    return MaterialApp(
      home: Scaffold(
        body: Builder(
          builder: (context) => Center(
            child: ElevatedButton(
              onPressed: () => onTap(context),
              child: const Text('OPEN'),
            ),
          ),
        ),
      ),
    );
  }

  testWidgets('【正常系】日付を選ぶと yyyy-MM-dd で返る', (tester) async {
    String? result;

    await tester.pumpWidget(
      host(onTap: (ctx) async => result = await pickDateAsYmd(ctx)),
    );

    await tester.tap(find.text('OPEN'));
    await tester.pumpAndSettle();

    final now = DateTime.now();
    final selected = DateTime(now.year, now.month, 15);
    await tester.tap(find.text('15'));
    await tester.pumpAndSettle();

    final dialogFinder = find.byType(DatePickerDialog);
    final loc = MaterialLocalizations.of(tester.element(dialogFinder));
    final okLabel = loc.okButtonLabel;

    await tester.tap(find.text(okLabel));
    await tester.pumpAndSettle();

    final expected = DateFormat('yyyy-MM-dd').format(selected);
    expect(result, expected);
  });

  testWidgets('【キャンセル】キャンセルすると null が返る', (tester) async {
    String? result;

    await tester.pumpWidget(
      host(onTap: (ctx) async => result = await pickDateAsYmd(ctx)),
    );

    await tester.tap(find.text('OPEN'));
    await tester.pumpAndSettle();

    final dialogFinder = find.byType(DatePickerDialog);
    final loc = MaterialLocalizations.of(tester.element(dialogFinder));
    final cancelLabel = loc.cancelButtonLabel;

    await tester.tap(find.text(cancelLabel));
    await tester.pumpAndSettle();

    expect(result, isNull);
  });
}
