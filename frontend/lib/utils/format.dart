String two(int n) => n.toString().padLeft(2, '0');

String kg(double v) =>
    '${v % 1 == 0 ? v.toStringAsFixed(0) : v.toStringAsFixed(1)}kg';
