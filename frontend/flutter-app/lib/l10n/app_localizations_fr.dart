// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for French (`fr`).
class AppLocalizationsFr extends AppLocalizations {
  AppLocalizationsFr([String locale = 'fr']) : super(locale);

  @override
  String get appName => 'ManPaSik';

  @override
  String get appTagline => 'Soins de santé IA pour une vie saine';

  @override
  String get greeting => 'Bonjour,';

  @override
  String greetingWithName(String name) {
    return '$name';
  }

  @override
  String get login => 'Connexion';

  @override
  String get register => 'Inscription';

  @override
  String get logout => 'Déconnexion';

  @override
  String get logoutConfirm => 'Êtes-vous sûr de vouloir vous déconnecter ?';

  @override
  String get email => 'E-mail';

  @override
  String get emailHint => 'example@manpasik.com';

  @override
  String get password => 'Mot de passe';

  @override
  String get passwordHint => 'Au moins 8 caractères (lettres + chiffres)';

  @override
  String get passwordConfirm => 'Confirmer le mot de passe';

  @override
  String get passwordConfirmHint => 'Saisissez à nouveau votre mot de passe';

  @override
  String get displayName => 'Nom';

  @override
  String get displayNameHint => 'Entrez votre nom d\'affichage';

  @override
  String get noAccountYet => 'Pas encore de compte ?';

  @override
  String get alreadyHaveAccount => 'Vous avez déjà un compte ?';

  @override
  String get loginFailed =>
      'Échec de la connexion. Veuillez vérifier votre e-mail et mot de passe.';

  @override
  String get registerFailed => 'Échec de l\'inscription. Veuillez réessayer.';

  @override
  String get home => 'Accueil';

  @override
  String get measurement => 'Mesure';

  @override
  String get devices => 'Appareils';

  @override
  String get settings => 'Paramètres';

  @override
  String get newMeasurement => 'Nouvelle mesure';

  @override
  String get startMeasurement => 'Démarrer la mesure';

  @override
  String get startMeasurementAction => 'Commencer la mesure';

  @override
  String get checkHealth => 'Vérifiez votre\nétat de santé';

  @override
  String get recentHistory => 'Historique récent';

  @override
  String get viewAll => 'Tout voir';

  @override
  String get preparingDevice => 'Préparez votre appareil';

  @override
  String get preparingDeviceDesc =>
      'Insérez la cartouche et\nappuyez sur le bouton de mesure';

  @override
  String get connectingDevice => 'Connexion à l\'appareil...';

  @override
  String get connectingDeviceDesc => 'Tentative de connexion BLE';

  @override
  String get measuring => 'Mesure en cours...';

  @override
  String get measuringDesc => 'Veuillez patienter';

  @override
  String get measurementComplete => 'Mesure terminée !';

  @override
  String get measurementFailed => 'Échec de la mesure';

  @override
  String get viewResult => 'Voir le résultat';

  @override
  String get retryMeasurement => 'Réessayer';

  @override
  String get bloodSugar => 'Glycémie';

  @override
  String get diagnosis => 'Diagnostic';

  @override
  String get noDevicesRegistered => 'Aucun appareil enregistré';

  @override
  String get noDevicesDesc =>
      'Appuyez sur le bouton + en haut à droite\npour enregistrer un nouvel appareil';

  @override
  String get searchDevices => 'Rechercher des appareils';

  @override
  String get addDevice => 'Ajouter un appareil';

  @override
  String get connected => 'Connecté';

  @override
  String get disconnected => 'Déconnecté';

  @override
  String get deviceRegistrationComingSoon =>
      'L\'enregistrement d\'appareils sera disponible dans la prochaine mise à jour';

  @override
  String get profile => 'Profil';

  @override
  String get general => 'Général';

  @override
  String get theme => 'Thème';

  @override
  String get themeSystem => 'Paramètre système';

  @override
  String get themeLight => 'Mode clair';

  @override
  String get themeDark => 'Mode sombre';

  @override
  String get themeSelect => 'Sélectionner le thème';

  @override
  String get language => 'Langue';

  @override
  String get languageSelect => 'Sélectionner la langue';

  @override
  String get appInfo => 'Infos sur l\'application';

  @override
  String get version => 'Version';

  @override
  String get termsOfService => 'Conditions d\'utilisation';

  @override
  String get privacyPolicy => 'Politique de confidentialité';

  @override
  String get account => 'Compte';

  @override
  String get loginRequired => 'Connexion requise';

  @override
  String get user => 'Utilisateur';

  @override
  String get cancel => 'Annuler';

  @override
  String get resultNormal => 'Normal';

  @override
  String get resultWarning => 'Attention';

  @override
  String get resultDanger => 'Danger';

  @override
  String get validationEmailRequired => 'Veuillez saisir votre adresse e-mail';

  @override
  String get validationEmailInvalid =>
      'Veuillez saisir une adresse e-mail valide';

  @override
  String get validationPasswordRequired => 'Veuillez saisir votre mot de passe';

  @override
  String get validationPasswordTooShort =>
      'Le mot de passe doit contenir au moins 8 caractères';

  @override
  String get validationPasswordNeedsLetter =>
      'Le mot de passe doit contenir des lettres';

  @override
  String get validationPasswordNeedsNumber =>
      'Le mot de passe doit contenir des chiffres';

  @override
  String get validationNameRequired => 'Veuillez saisir votre nom';

  @override
  String get validationNameLength =>
      'Le nom doit contenir entre 2 et 50 caractères';

  @override
  String get validationPasswordMismatch =>
      'Les mots de passe ne correspondent pas';
}
