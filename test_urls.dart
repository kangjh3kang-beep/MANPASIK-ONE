import 'dart:io';

void main() async {
  final urls = [
    'http://127.0.0.1:8080/models/Holo_M1.glb',
    'http://127.0.0.1:8080/assets/web/models/Holo_M1.glb',
    'http://127.0.0.1:8080/human_body.glb',
    'http://127.0.0.1:8080/assets/web/human_body.glb',
  ];

  for (final url in urls) {
    print('Testing $url');
    try {
      final request = await HttpClient().getUrl(Uri.parse(url));
      final response = await request.close();
      print('Status: ${response.statusCode}');
      
      // Read a few bytes
      final firstBytes = await response.take(1).toList();
      print('Got data: ${firstBytes.isNotEmpty}');
    } catch (e) {
      print('Error: $e');
    }
  }
}
