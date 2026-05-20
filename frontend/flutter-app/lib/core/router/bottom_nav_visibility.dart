/// Global bottom dock visibility policy by route prefix.
///
/// Keep this list as the single source of truth to avoid duplicated route
/// checks across shell scaffolds.
const Set<String> kHideGlobalDockRoutePrefixes = {
  '/home',
  '/measure',
  '/medical/video-call',
};

/// Global top frame visibility policy by route prefix.
///
/// Routes that render their own top chrome should hide the global frame to
/// avoid double-header layout and size mismatch.
const Set<String> kHideGlobalTopFrameRoutePrefixes = {
  '/home',
  '/measure',
  '/medical/video-call',
};

bool _matchesAnyPrefix(String matchedLocation, Set<String> prefixes) {
  for (final prefix in prefixes) {
    if (matchedLocation == prefix || matchedLocation.startsWith('$prefix/')) {
      return true;
    }
  }
  return false;
}

bool shouldHideGlobalDock(String matchedLocation) {
  return _matchesAnyPrefix(matchedLocation, kHideGlobalDockRoutePrefixes);
}

bool shouldHideGlobalTopFrame(String matchedLocation) {
  return _matchesAnyPrefix(matchedLocation, kHideGlobalTopFrameRoutePrefixes);
}
