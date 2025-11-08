import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/utils/format.dart';

void main() {
  group('two()', () {
    final cases = <int, String>{
      0: '00',
      3: '03',
      9: '09',
      10: '10',
      12: '12',
      99: '99',
    };

    cases.forEach((input, want) {
      test('pads $input -> "$want"', () {
        expect(two(input), want);
      });
    });
  });

  group('kg()', () {
    final cases = <double, String>{
      60.0: '60kg',
      60.5: '60.5kg',
      60.25: '60.3kg',
      60.04: '60.0kg',
      0.0: '0kg',
    };

    cases.forEach((input, want) {
      test('formats $input -> "$want"', () {
        expect(kg(input), want);
      });
    });
  });

  group('ymd()', () {
    final cases = {
      DateTime(2025, 1, 9): '2025/01/09',
      DateTime(1999, 12, 31): '1999/12/31',
      DateTime(2024, 2, 29): '2024/02/29',
    };

    cases.forEach((input, want) {
      test('formats $input -> "$want"', () {
        expect(ymd(input), want);
      });
    });
  });
}
