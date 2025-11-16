import 'package:flutter/material.dart';

class RankingRow extends StatelessWidget {
  const RankingRow({
    super.key,
    required this.rank,
    required this.email,
    required this.days,
  });

  final int rank;
  final String email;
  final int days;

  @override
  Widget build(BuildContext context) {
    IconData? icon;
    Color? iconColor;

    switch (rank) {
      case 1:
        icon = Icons.emoji_events;
        iconColor = Colors.amber;
        break;
      case 2:
        icon = Icons.emoji_events;
        iconColor = Colors.grey;
        break;
      case 3:
        icon = Icons.emoji_events;
        iconColor = Colors.brown;
        break;
      default:
        icon = null;
        iconColor = null;
    }

    return Card(
      elevation: 1,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
        child: Row(
          children: [
            SizedBox(
              width: 40,
              child: icon != null
                  ? Icon(icon, color: iconColor)
                  : Text(
                      '$rank位',
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    email,
                    style: const TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                    ),
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'ジム日数：$days日',
                    style: const TextStyle(fontSize: 12, color: Colors.black54),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
