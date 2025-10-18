import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

final _ymd = DateFormat('yyyy-MM-dd');

Future<String?> pickDateAsYmd(BuildContext context) async {
  final now = DateTime.now();
  final picked = await showDatePicker(
    context: context,
    initialDate: now,
    firstDate: DateTime(now.year - 5),
    lastDate: DateTime(now.year + 5),
  );
  return picked == null ? null : _ymd.format(picked);
}
