// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Japanese (`ja`).
class AppLocalizationsJa extends AppLocalizations {
  AppLocalizationsJa([String locale = 'ja']) : super(locale);

  @override
  String get appName => 'マンパシク';

  @override
  String get appTagline => '健康な日常のためのAIヘルスケア';

  @override
  String get greeting => 'こんにちは、';

  @override
  String greetingWithName(String name) {
    return '$nameさん';
  }

  @override
  String get login => 'ログイン';

  @override
  String get register => '新規登録';

  @override
  String get logout => 'ログアウト';

  @override
  String get logoutConfirm => '本当にログアウトしますか？';

  @override
  String get email => 'メール';

  @override
  String get emailHint => 'example@manpasik.com';

  @override
  String get password => 'パスワード';

  @override
  String get passwordHint => '8文字以上（英数字）';

  @override
  String get passwordConfirm => 'パスワード確認';

  @override
  String get passwordConfirmHint => 'パスワードを再入力してください';

  @override
  String get displayName => '名前';

  @override
  String get displayNameHint => '表示名を入力してください';

  @override
  String get noAccountYet => 'アカウントをお持ちでないですか？';

  @override
  String get alreadyHaveAccount => 'すでにアカウントをお持ちですか？';

  @override
  String get loginFailed => 'ログインに失敗しました。メールとパスワードをご確認ください。';

  @override
  String get registerFailed => '登録に失敗しました。もう一度お試しください。';

  @override
  String get home => 'ホーム';

  @override
  String get measurement => '測定';

  @override
  String get devices => 'デバイス';

  @override
  String get settings => '設定';

  @override
  String get newMeasurement => '新しい測定';

  @override
  String get startMeasurement => '測定開始';

  @override
  String get startMeasurementAction => '測定を始める';

  @override
  String get checkHealth => '健康状態を\nチェックしましょう';

  @override
  String get recentHistory => '最近の記録';

  @override
  String get viewAll => 'すべて表示';

  @override
  String get preparingDevice => 'デバイスを準備してください';

  @override
  String get preparingDeviceDesc => 'カートリッジを装着し\n測定ボタンを押してください';

  @override
  String get connectingDevice => 'デバイス接続中...';

  @override
  String get connectingDeviceDesc => 'BLE接続を試みています';

  @override
  String get measuring => '測定中...';

  @override
  String get measuringDesc => '少々お待ちください';

  @override
  String get measurementComplete => '測定完了！';

  @override
  String get measurementFailed => '測定失敗';

  @override
  String get viewResult => '結果を見る';

  @override
  String get retryMeasurement => '再測定';

  @override
  String get bloodSugar => '血糖';

  @override
  String get diagnosis => '判定';

  @override
  String get noDevicesRegistered => '登録されたデバイスがありません';

  @override
  String get noDevicesDesc => '右上の+ボタンをタップして\n新しいデバイスを登録してください';

  @override
  String get searchDevices => 'デバイス検索';

  @override
  String get addDevice => 'デバイス追加';

  @override
  String get connected => '接続済み';

  @override
  String get disconnected => '未接続';

  @override
  String get deviceRegistrationComingSoon => 'デバイス登録機能は次のアップデートで利用可能になります';

  @override
  String get profile => 'プロフィール';

  @override
  String get general => '一般';

  @override
  String get theme => 'テーマ';

  @override
  String get themeSystem => 'システム設定';

  @override
  String get themeLight => 'ライトモード';

  @override
  String get themeDark => 'ダークモード';

  @override
  String get themeSelect => 'テーマ選択';

  @override
  String get language => '言語';

  @override
  String get languageSelect => '言語選択';

  @override
  String get appInfo => 'アプリ情報';

  @override
  String get version => 'バージョン';

  @override
  String get termsOfService => '利用規約';

  @override
  String get privacyPolicy => 'プライバシーポリシー';

  @override
  String get account => 'アカウント';

  @override
  String get loginRequired => 'ログインが必要です';

  @override
  String get user => 'ユーザー';

  @override
  String get cancel => 'キャンセル';

  @override
  String get resultNormal => '正常';

  @override
  String get resultWarning => '注意';

  @override
  String get resultDanger => '危険';

  @override
  String get validationEmailRequired => 'メールアドレスを入力してください';

  @override
  String get validationEmailInvalid => '正しいメール形式を入力してください';

  @override
  String get validationPasswordRequired => 'パスワードを入力してください';

  @override
  String get validationPasswordTooShort => 'パスワードは8文字以上必要です';

  @override
  String get validationPasswordNeedsLetter => '英字を含めてください';

  @override
  String get validationPasswordNeedsNumber => '数字を含めてください';

  @override
  String get validationNameRequired => '名前を入力してください';

  @override
  String get validationNameLength => '名前は2〜50文字の間にしてください';

  @override
  String get validationPasswordMismatch => 'パスワードが一致しません';
}
