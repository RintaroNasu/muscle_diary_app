class Exercise {
  final int id;
  final String name;

  const Exercise({required this.id, required this.name});

  factory Exercise.fromJson(Map<String, dynamic> json) {
    return Exercise(id: json['id'] as int, name: json['name'] as String);
  }

  Map<String, dynamic> toJson() => {'id': id, 'name': name};
}
