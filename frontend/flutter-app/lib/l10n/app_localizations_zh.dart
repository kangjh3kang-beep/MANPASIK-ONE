// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Chinese (`zh`).
class AppLocalizationsZh extends AppLocalizations {
  AppLocalizationsZh([String locale = 'zh']) : super(locale);

  @override
  String get appName => '万波息';

  @override
  String get appTagline => '为健康生活而生的AI医疗保健';

  @override
  String get greeting => '您好，';

  @override
  String greetingWithName(String name) {
    return '$name';
  }

  @override
  String get login => '登录';

  @override
  String get register => '注册';

  @override
  String get logout => '退出登录';

  @override
  String get logoutConfirm => '确定要退出登录吗？';

  @override
  String get email => '邮箱';

  @override
  String get emailHint => 'example@manpasik.com';

  @override
  String get password => '密码';

  @override
  String get passwordHint => '至少8位（字母+数字）';

  @override
  String get passwordConfirm => '确认密码';

  @override
  String get passwordConfirmHint => '请再次输入密码';

  @override
  String get displayName => '姓名';

  @override
  String get displayNameHint => '请输入您的显示名称';

  @override
  String get noAccountYet => '还没有账号？';

  @override
  String get alreadyHaveAccount => '已有账号？';

  @override
  String get loginFailed => '登录失败，请检查邮箱和密码。';

  @override
  String get registerFailed => '注册失败，请重试。';

  @override
  String get home => '首页';

  @override
  String get measurement => '测量';

  @override
  String get devices => '设备';

  @override
  String get settings => '设置';

  @override
  String get newMeasurement => '新测量';

  @override
  String get startMeasurement => '开始测量';

  @override
  String get startMeasurementAction => '开始测量';

  @override
  String get checkHealth => '检查您的\n健康状况';

  @override
  String get recentHistory => '最近记录';

  @override
  String get viewAll => '查看全部';

  @override
  String get preparingDevice => '请准备设备';

  @override
  String get preparingDeviceDesc => '安装试剂盒并\n按下测量按钮';

  @override
  String get connectingDevice => '正在连接设备...';

  @override
  String get connectingDeviceDesc => '正在尝试BLE连接';

  @override
  String get measuring => '测量中...';

  @override
  String get measuringDesc => '请稍候';

  @override
  String get measurementComplete => '测量完成！';

  @override
  String get measurementFailed => '测量失败';

  @override
  String get viewResult => '查看结果';

  @override
  String get retryMeasurement => '重新测量';

  @override
  String get bloodSugar => '血糖';

  @override
  String get diagnosis => '诊断';

  @override
  String get noDevicesRegistered => '没有已注册的设备';

  @override
  String get noDevicesDesc => '点击右上角的+按钮\n注册新设备';

  @override
  String get searchDevices => '搜索设备';

  @override
  String get addDevice => '添加设备';

  @override
  String get connected => '已连接';

  @override
  String get disconnected => '未连接';

  @override
  String get deviceRegistrationComingSoon => '设备注册功能将在下次更新中提供';

  @override
  String get profile => '个人资料';

  @override
  String get general => '通用';

  @override
  String get theme => '主题';

  @override
  String get themeSystem => '跟随系统';

  @override
  String get themeLight => '浅色模式';

  @override
  String get themeDark => '深色模式';

  @override
  String get themeSelect => '选择主题';

  @override
  String get language => '语言';

  @override
  String get languageSelect => '选择语言';

  @override
  String get appInfo => '应用信息';

  @override
  String get version => '版本';

  @override
  String get termsOfService => '服务条款';

  @override
  String get privacyPolicy => '隐私政策';

  @override
  String get account => '账户';

  @override
  String get loginRequired => '需要登录';

  @override
  String get user => '用户';

  @override
  String get cancel => '取消';

  @override
  String get resultNormal => '正常';

  @override
  String get resultWarning => '注意';

  @override
  String get resultDanger => '危险';

  @override
  String get validationEmailRequired => '请输入邮箱地址';

  @override
  String get validationEmailInvalid => '请输入有效的邮箱地址';

  @override
  String get validationPasswordRequired => '请输入密码';

  @override
  String get validationPasswordTooShort => '密码至少需要8个字符';

  @override
  String get validationPasswordNeedsLetter => '密码需要包含字母';

  @override
  String get validationPasswordNeedsNumber => '密码需要包含数字';

  @override
  String get validationNameRequired => '请输入姓名';

  @override
  String get validationNameLength => '姓名长度需在2到50个字符之间';

  @override
  String get validationPasswordMismatch => '两次输入的密码不一致';
}
