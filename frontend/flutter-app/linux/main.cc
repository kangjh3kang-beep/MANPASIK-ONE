#include <cstring>

#include <webview_cef/webview_cef_plugin.h>
#include "my_application.h"

namespace {
bool IsCefSubProcess(int argc, char** argv) {
  for (int i = 1; i < argc; ++i) {
    if (argv[i] == nullptr) {
      continue;
    }
    if (std::strncmp(argv[i], "--type=", 7) == 0 ||
        std::strcmp(argv[i], "--type") == 0) {
      return true;
    }
  }
  return false;
}
}  // namespace

int main(int argc, char** argv) {
  initCEFProcesses(argc, argv);
  if (IsCefSubProcess(argc, argv)) {
    // CEF helper subprocess should not start the Flutter UI window.
    return 0;
  }

  g_autoptr(MyApplication) app = my_application_new();
  return g_application_run(G_APPLICATION(app), argc, argv);
}
