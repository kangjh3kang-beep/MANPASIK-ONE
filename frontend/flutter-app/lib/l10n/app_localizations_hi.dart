// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Hindi (`hi`).
class AppLocalizationsHi extends AppLocalizations {
  AppLocalizationsHi([String locale = 'hi']) : super(locale);

  @override
  String get appName => 'मानपासिक';

  @override
  String get appTagline => 'स्वस्थ जीवन के लिए AI हेल्थकेयर';

  @override
  String get greeting => 'नमस्ते,';

  @override
  String greetingWithName(String name) {
    return '$name';
  }

  @override
  String get login => 'लॉग इन';

  @override
  String get register => 'साइन अप';

  @override
  String get logout => 'लॉग आउट';

  @override
  String get logoutConfirm => 'क्या आप वाकई लॉग आउट करना चाहते हैं?';

  @override
  String get email => 'ईमेल';

  @override
  String get emailHint => 'example@manpasik.com';

  @override
  String get password => 'पासवर्ड';

  @override
  String get passwordHint => 'कम से कम 8 अक्षर (अक्षर + अंक)';

  @override
  String get passwordConfirm => 'पासवर्ड की पुष्टि';

  @override
  String get passwordConfirmHint => 'अपना पासवर्ड दोबारा दर्ज करें';

  @override
  String get displayName => 'नाम';

  @override
  String get displayNameHint => 'अपना प्रदर्शन नाम दर्ज करें';

  @override
  String get noAccountYet => 'खाता नहीं है?';

  @override
  String get alreadyHaveAccount => 'पहले से खाता है?';

  @override
  String get loginFailed => 'लॉगिन विफल। कृपया अपना ईमेल और पासवर्ड जांचें।';

  @override
  String get registerFailed => 'पंजीकरण विफल। कृपया पुनः प्रयास करें।';

  @override
  String get home => 'होम';

  @override
  String get measurement => 'माप';

  @override
  String get devices => 'डिवाइस';

  @override
  String get settings => 'सेटिंग्स';

  @override
  String get newMeasurement => 'नया माप';

  @override
  String get startMeasurement => 'माप शुरू करें';

  @override
  String get startMeasurementAction => 'मापना शुरू करें';

  @override
  String get checkHealth => 'अपनी स्वास्थ्य\nस्थिति जांचें';

  @override
  String get recentHistory => 'हाल का इतिहास';

  @override
  String get viewAll => 'सभी देखें';

  @override
  String get preparingDevice => 'अपना डिवाइस तैयार करें';

  @override
  String get preparingDeviceDesc => 'कार्ट्रिज डालें और\nमाप बटन दबाएं';

  @override
  String get connectingDevice => 'डिवाइस कनेक्ट हो रहा है...';

  @override
  String get connectingDeviceDesc => 'BLE कनेक्शन का प्रयास किया जा रहा है';

  @override
  String get measuring => 'माप जारी...';

  @override
  String get measuringDesc => 'कृपया प्रतीक्षा करें';

  @override
  String get measurementComplete => 'माप पूर्ण!';

  @override
  String get measurementFailed => 'माप विफल';

  @override
  String get viewResult => 'परिणाम देखें';

  @override
  String get retryMeasurement => 'पुनः प्रयास';

  @override
  String get bloodSugar => 'रक्त शर्करा';

  @override
  String get diagnosis => 'निदान';

  @override
  String get noDevicesRegistered => 'कोई डिवाइस पंजीकृत नहीं';

  @override
  String get noDevicesDesc =>
      'नया डिवाइस पंजीकृत करने के लिए\nऊपर दाईं ओर + बटन दबाएं';

  @override
  String get searchDevices => 'डिवाइस खोजें';

  @override
  String get addDevice => 'डिवाइस जोड़ें';

  @override
  String get connected => 'कनेक्टेड';

  @override
  String get disconnected => 'डिस्कनेक्टेड';

  @override
  String get deviceRegistrationComingSoon =>
      'डिवाइस पंजीकरण सुविधा अगले अपडेट में उपलब्ध होगी';

  @override
  String get profile => 'प्रोफ़ाइल';

  @override
  String get general => 'सामान्य';

  @override
  String get theme => 'थीम';

  @override
  String get themeSystem => 'सिस्टम डिफ़ॉल्ट';

  @override
  String get themeLight => 'लाइट मोड';

  @override
  String get themeDark => 'डार्क मोड';

  @override
  String get themeSelect => 'थीम चुनें';

  @override
  String get language => 'भाषा';

  @override
  String get languageSelect => 'भाषा चुनें';

  @override
  String get appInfo => 'ऐप जानकारी';

  @override
  String get version => 'संस्करण';

  @override
  String get termsOfService => 'सेवा की शर्तें';

  @override
  String get privacyPolicy => 'गोपनीयता नीति';

  @override
  String get account => 'खाता';

  @override
  String get loginRequired => 'लॉगिन आवश्यक';

  @override
  String get user => 'उपयोगकर्ता';

  @override
  String get cancel => 'रद्द करें';

  @override
  String get resultNormal => 'सामान्य';

  @override
  String get resultWarning => 'सावधानी';

  @override
  String get resultDanger => 'खतरा';

  @override
  String get validationEmailRequired => 'कृपया अपना ईमेल दर्ज करें';

  @override
  String get validationEmailInvalid => 'कृपया एक मान्य ईमेल पता दर्ज करें';

  @override
  String get validationPasswordRequired => 'कृपया अपना पासवर्ड दर्ज करें';

  @override
  String get validationPasswordTooShort =>
      'पासवर्ड कम से कम 8 अक्षर का होना चाहिए';

  @override
  String get validationPasswordNeedsLetter => 'पासवर्ड में अक्षर होने चाहिए';

  @override
  String get validationPasswordNeedsNumber => 'पासवर्ड में अंक होने चाहिए';

  @override
  String get validationNameRequired => 'कृपया अपना नाम दर्ज करें';

  @override
  String get validationNameLength => 'नाम 2 से 50 अक्षर के बीच होना चाहिए';

  @override
  String get validationPasswordMismatch => 'पासवर्ड मेल नहीं खाते';
}
