// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appName => 'MANPASIK';

  @override
  String get appTagline => 'High-Precision Differential Measurement System';

  @override
  String get greeting => 'Hello,';

  @override
  String greetingWithName(String name) {
    return '$name';
  }

  @override
  String get login => 'Login';

  @override
  String get register => 'Sign Up';

  @override
  String get logout => 'Logout';

  @override
  String get logoutConfirm => 'Are you sure you want to logout?';

  @override
  String get email => 'Email';

  @override
  String get emailHint => 'example@manpasik.com';

  @override
  String get password => 'Password';

  @override
  String get passwordHint => 'At least 8 characters (letters + numbers)';

  @override
  String get passwordConfirm => 'Confirm Password';

  @override
  String get passwordConfirmHint => 'Re-enter your password';

  @override
  String get displayName => 'Name';

  @override
  String get displayNameHint => 'Enter your display name';

  @override
  String get noAccountYet => 'Don\'t have an account?';

  @override
  String get alreadyHaveAccount => 'Already have an account?';

  @override
  String get loginFailed =>
      'Login failed. Please check your email and password.';

  @override
  String get registerFailed => 'Registration failed. Please try again.';

  @override
  String get home => 'Home';

  @override
  String get measurement => 'Measure';

  @override
  String get devices => 'Devices';

  @override
  String get settings => 'Settings';

  @override
  String get newMeasurement => 'New Measurement';

  @override
  String get startMeasurement => 'Start Measurement';

  @override
  String get startMeasurementAction => 'Start Measuring';

  @override
  String get checkHealth => 'Check your\nhealth status';

  @override
  String get recentHistory => 'Recent History';

  @override
  String get viewAll => 'View All';

  @override
  String get preparingDevice => 'Prepare your device';

  @override
  String get preparingDeviceDesc =>
      'Insert the cartridge and\npress the measure button';

  @override
  String get connectingDevice => 'Connecting device...';

  @override
  String get connectingDeviceDesc => 'Attempting BLE connection';

  @override
  String get measuring => 'Measuring...';

  @override
  String get measuringDesc => 'Please wait';

  @override
  String get measurementComplete => 'Measurement Complete!';

  @override
  String get measurementFailed => 'Measurement Failed';

  @override
  String get viewResult => 'View Result';

  @override
  String get retryMeasurement => 'Retry';

  @override
  String get bloodSugar => 'Blood Sugar';

  @override
  String get diagnosis => 'Diagnosis';

  @override
  String get noDevicesRegistered => 'No devices registered';

  @override
  String get noDevicesDesc =>
      'Tap the + button at the top right\nto register a new device';

  @override
  String get searchDevices => 'Search Devices';

  @override
  String get addDevice => 'Add Device';

  @override
  String get connected => 'Connected';

  @override
  String get disconnected => 'Disconnected';

  @override
  String get deviceRegistrationComingSoon =>
      'Device registration will be available in the next update';

  @override
  String get profile => 'Profile';

  @override
  String get general => 'General';

  @override
  String get theme => 'Theme';

  @override
  String get themeSystem => 'System Default';

  @override
  String get themeLight => 'Light Mode';

  @override
  String get themeDark => 'Dark Mode';

  @override
  String get themeSelect => 'Select Theme';

  @override
  String get language => 'Language';

  @override
  String get languageSelect => 'Select Language';

  @override
  String get appInfo => 'App Info';

  @override
  String get version => 'Version';

  @override
  String get termsOfService => 'Terms of Service';

  @override
  String get privacyPolicy => 'Privacy Policy';

  @override
  String get account => 'Account';

  @override
  String get loginRequired => 'Login required';

  @override
  String get user => 'User';

  @override
  String get cancel => 'Cancel';

  @override
  String get resultNormal => 'Normal';

  @override
  String get resultWarning => 'Warning';

  @override
  String get resultDanger => 'Danger';

  @override
  String get validationEmailRequired => 'Please enter your email';

  @override
  String get validationEmailInvalid => 'Please enter a valid email address';

  @override
  String get validationPasswordRequired => 'Please enter your password';

  @override
  String get validationPasswordTooShort =>
      'Password must be at least 8 characters';

  @override
  String get validationPasswordNeedsLetter => 'Password must contain letters';

  @override
  String get validationPasswordNeedsNumber => 'Password must contain numbers';

  @override
  String get validationNameRequired => 'Please enter your name';

  @override
  String get validationNameLength => 'Name must be between 2 and 50 characters';

  @override
  String get validationPasswordMismatch => 'Passwords do not match';
}
