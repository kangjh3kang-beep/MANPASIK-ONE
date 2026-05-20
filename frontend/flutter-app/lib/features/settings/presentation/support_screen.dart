import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// SupportScreen — Sanggam Orbit 고객 지원
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar + TabBar(bottom) → body 내 헤더 + 독립 TabBar
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData params 2x 제거
// [Rule 4] theme.textTheme ~2x → 직접 TextStyle
// [Rule 4] theme.colorScheme ~2x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 3x → SanggamTheme.primary
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Card ~5x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 4] BottomSheet/TextField/FilledButton 다크 테마
// [Rule 2] borderRadius 12→16
// ───────────────────────────────────────────────────

/// 고객 지원 화면 (storyboard-support.md)
///
/// FAQ 아코디언 + 1:1 문의 + 전화 문의
class SupportScreen extends StatefulWidget {
  const SupportScreen({super.key});

  @override
  State<SupportScreen> createState() =>
      _SupportScreenState();
}

class _SupportScreenState
    extends State<SupportScreen>
    with SingleTickerProviderStateMixin {
  late final TabController _tabController;
  final _typeController =
      TextEditingController();
  final _titleController =
      TextEditingController();
  final _contentController =
      TextEditingController();
  final _imagePicker = ImagePicker();
  final List<XFile> _attachedImages = [];

  @override
  void initState() {
    super.initState();
    _tabController =
        TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _typeController.dispose();
    _titleController.dispose();
    _contentController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '고객 지원',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // TabBar
            TabBar(
              controller: _tabController,
              indicatorColor:
                  SanggamTheme.primary,
              labelColor:
                  SanggamTheme.primary,
              unselectedLabelColor:
                  SanggamTheme.onSurfaceDim,
              dividerColor:
                  SanggamTheme.surfaceVariant,
              tabs: const [
                Tab(text: '자주 묻는 질문'),
                Tab(text: '1:1 문의'),
              ],
            ),
            // 본문
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  _buildFaqTab(),
                  _buildInquiryTab(),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFaqTab() {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        // 검색바
        TextField(
          style: const TextStyle(
              color: Colors.white),
          decoration: InputDecoration(
            hintText: '궁금한 내용을 검색하세요',
            hintStyle: const TextStyle(
                color: SanggamTheme
                    .onSurfaceDim),
            prefixIcon: const Icon(
                Icons.search,
                color: SanggamTheme
                    .onSurfaceDim),
            filled: true,
            fillColor: SanggamTheme.surface,
            border: OutlineInputBorder(
              borderRadius:
                  BorderRadius.circular(16),
              borderSide: BorderSide.none,
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius:
                  BorderRadius.circular(16),
              borderSide: const BorderSide(
                  color: SanggamTheme
                      .surfaceVariant),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius:
                  BorderRadius.circular(16),
              borderSide: const BorderSide(
                  color:
                      SanggamTheme.primary),
            ),
          ),
        ),
        const SizedBox(height: 16),

        // FAQ 카테고리
        ..._faqCategories.entries
            .map((entry) {
          return Padding(
            padding: const EdgeInsets.only(
                bottom: 8),
            child: SanggamContainer(
              borderRadius: 16,
              padding: EdgeInsets.zero,
              child: Theme(
                data:
                    ThemeData.dark().copyWith(
                  dividerColor: SanggamTheme
                      .surfaceVariant,
                ),
                child: ExpansionTile(
                  leading: Icon(
                      entry.value.icon,
                      color: SanggamTheme
                          .primary),
                  title: Text(entry.key,
                      style:
                          const TextStyle(
                        color: Colors.white,
                        fontWeight:
                            FontWeight.w600,
                      )),
                  children: entry.value.items
                      .map((faq) {
                    return ExpansionTile(
                      tilePadding:
                          const EdgeInsets
                              .only(
                              left: 48,
                              right: 16),
                      title: Text(
                          faq.question,
                          style:
                              const TextStyle(
                            color:
                                Colors.white,
                            fontSize: 14,
                          )),
                      children: [
                        Padding(
                          padding:
                              const EdgeInsets
                                  .fromLTRB(
                                  48,
                                  0,
                                  16,
                                  16),
                          child: Text(
                            faq.answer,
                            style:
                                const TextStyle(
                              color: SanggamTheme
                                  .onSurfaceDim,
                              fontSize: 13,
                            ),
                          ),
                        ),
                      ],
                    );
                  }).toList(),
                ),
              ),
            ),
          );
        }),

        const SizedBox(height: 24),
        // 전화 문의
        SanggamContainer(
          borderRadius: 16,
          padding: EdgeInsets.zero,
          backgroundColor: SanggamTheme
              .primary
              .withValues(alpha: 0.1),
          child: ListTile(
            leading: const Icon(Icons.phone,
                color: SanggamTheme.primary),
            title: const Text('전화 문의',
                style: TextStyle(
                    color: Colors.white)),
            subtitle: const Text(
                '평일 09:00~18:00 (점심시간 12:00~13:00)',
                style: TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                )),
            trailing: const Text('1588-0000',
                style: TextStyle(
                  color: SanggamTheme.primary,
                  fontWeight: FontWeight.bold,
                )),
          ),
        ),
      ],
    );
  }

  Widget _buildInquiryTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.stretch,
        children: [
          DropdownButtonFormField<String>(
            dropdownColor:
                SanggamTheme.surface,
            style: const TextStyle(
                color: Colors.white),
            decoration: InputDecoration(
              labelText: '문의 유형',
              labelStyle: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim),
              filled: true,
              fillColor:
                  SanggamTheme.surface,
              border: OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: BorderSide.none,
              ),
              enabledBorder:
                  OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: const BorderSide(
                    color: SanggamTheme
                        .surfaceVariant),
              ),
              focusedBorder:
                  OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: const BorderSide(
                    color: SanggamTheme
                        .primary),
              ),
            ),
            items: const [
              DropdownMenuItem(
                  value: 'device',
                  child:
                      Text('기기/카트리지')),
              DropdownMenuItem(
                  value: 'subscription',
                  child: Text('구독/결제')),
              DropdownMenuItem(
                  value: 'account',
                  child: Text('계정/인증')),
              DropdownMenuItem(
                  value: 'measurement',
                  child: Text('측정/결과')),
              DropdownMenuItem(
                  value: 'other',
                  child: Text('기타')),
            ],
            onChanged: (v) =>
                _typeController.text =
                    v ?? '',
          ),
          const SizedBox(height: 16),
          TextFormField(
            controller: _titleController,
            style: const TextStyle(
                color: Colors.white),
            decoration: InputDecoration(
              labelText: '제목',
              labelStyle: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim),
              hintText: '문의 제목을 입력해주세요',
              hintStyle: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim),
              filled: true,
              fillColor:
                  SanggamTheme.surface,
              border: OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: BorderSide.none,
              ),
              enabledBorder:
                  OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: const BorderSide(
                    color: SanggamTheme
                        .surfaceVariant),
              ),
              focusedBorder:
                  OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: const BorderSide(
                    color: SanggamTheme
                        .primary),
              ),
            ),
          ),
          const SizedBox(height: 16),
          TextFormField(
            controller: _contentController,
            maxLines: 6,
            style: const TextStyle(
                color: Colors.white),
            decoration: InputDecoration(
              labelText: '문의 내용',
              labelStyle: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim),
              hintText:
                  '문의하실 내용을 상세히 작성해주세요.\n\n기기 관련 문의 시 기기 시리얼 번호를 함께 기재해주세요.',
              hintStyle: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim),
              alignLabelWithHint: true,
              filled: true,
              fillColor:
                  SanggamTheme.surface,
              border: OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: BorderSide.none,
              ),
              enabledBorder:
                  OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: const BorderSide(
                    color: SanggamTheme
                        .surfaceVariant),
              ),
              focusedBorder:
                  OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(16),
                borderSide: const BorderSide(
                    color: SanggamTheme
                        .primary),
              ),
            ),
          ),
          const SizedBox(height: 16),
          OutlinedButton.icon(
            onPressed: () async {
              final source =
                  await showModalBottomSheet<
                      ImageSource>(
                context: context,
                backgroundColor:
                    SanggamTheme.surface,
                shape:
                    const RoundedRectangleBorder(
                  borderRadius:
                      BorderRadius.vertical(
                          top: Radius.circular(
                              16)),
                ),
                builder: (ctx) => SafeArea(
                  child: Column(
                    mainAxisSize:
                        MainAxisSize.min,
                    children: [
                      ListTile(
                        leading: const Icon(
                            Icons.camera_alt,
                            color:
                                Colors.white),
                        title: const Text(
                            '카메라로 촬영',
                            style: TextStyle(
                                color: Colors
                                    .white)),
                        onTap: () =>
                            Navigator.pop(
                                ctx,
                                ImageSource
                                    .camera),
                      ),
                      const Divider(
                          height: 1,
                          color: SanggamTheme
                              .surfaceVariant),
                      ListTile(
                        leading: const Icon(
                            Icons
                                .photo_library,
                            color:
                                Colors.white),
                        title: const Text(
                            '갤러리에서 선택',
                            style: TextStyle(
                                color: Colors
                                    .white)),
                        onTap: () =>
                            Navigator.pop(
                                ctx,
                                ImageSource
                                    .gallery),
                      ),
                    ],
                  ),
                ),
              );
              if (source == null) return;
              final image =
                  await _imagePicker
                      .pickImage(
                          source: source,
                          maxWidth: 1024);
              if (image != null) {
                setState(() =>
                    _attachedImages
                        .add(image));
              }
            },
            style: OutlinedButton.styleFrom(
              foregroundColor:
                  SanggamTheme.onSurfaceDim,
              side: const BorderSide(
                  color: SanggamTheme
                      .surfaceVariant),
              shape: RoundedRectangleBorder(
                borderRadius:
                    BorderRadius.circular(16),
              ),
            ),
            icon: const Icon(
                Icons.attach_file),
            label: Text(
                _attachedImages.isEmpty
                    ? '사진 첨부 (선택)'
                    : '사진 ${_attachedImages.length}장 첨부됨'),
          ),
          const SizedBox(height: 24),
          FilledButton(
            onPressed: () {
              ScaffoldMessenger.of(context)
                  .showSnackBar(
                const SnackBar(
                    content: Text(
                        '문의가 접수되었습니다. 담당자가 확인 후 답변드리겠습니다.')),
              );
              _titleController.clear();
              _contentController.clear();
            },
            style: FilledButton.styleFrom(
              minimumSize:
                  const Size.fromHeight(48),
              backgroundColor:
                  SanggamTheme.primary,
              foregroundColor:
                  SanggamTheme.background,
              shape:
                  RoundedRectangleBorder(
                borderRadius:
                    BorderRadius.circular(16),
              ),
            ),
            child: const Text('문의 접수'),
          ),
        ],
      ),
    );
  }
}

// ── FAQ 데이터 ──

class _FaqItem {
  final String question;
  final String answer;
  const _FaqItem(this.question, this.answer);
}

class _FaqCategory {
  final IconData icon;
  final List<_FaqItem> items;
  const _FaqCategory(this.icon, this.items);
}

final _faqCategories =
    <String, _FaqCategory>{
  '기기 관련': _FaqCategory(Icons.devices, [
    _FaqItem(
        '리더기를 처음 연결하는 방법은?',
        'ManPaSik 앱 > 설정 > 기기 관리에서 "기기 추가" 버튼을 눌러주세요. '
            'BLE(블루투스 저전력)를 통해 자동으로 주변 리더기를 검색합니다.'),
    _FaqItem(
        '카트리지가 인식되지 않아요',
        'NFC 안테나 위치를 확인하고, 카트리지를 리더기 NFC 부분에 가까이 대주세요. '
            '카트리지의 유효기간이 지나지 않았는지도 확인해주세요.'),
    _FaqItem(
        '리더기 펌웨어 업데이트는 어떻게 하나요?',
        '기기 관리 화면에서 "펌웨어 업데이트" 알림이 표시되면, '
            'Wi-Fi 연결 상태에서 업데이트 버튼을 눌러주세요. 약 5분 정도 소요됩니다.'),
  ]),
  '측정 관련': _FaqCategory(Icons.science, [
    _FaqItem(
        '측정 결과가 "분석 중"으로 표시돼요',
        'AI 분석에 약간의 시간이 소요될 수 있습니다. '
            '네트워크 연결 상태를 확인해주세요. 오프라인에서도 로컬 AI로 기본 분석이 가능합니다.'),
    _FaqItem(
        '정확한 측정을 위한 조건은?',
        '실온(20~25°C)에서, 샘플을 넣은 후 30초 이상 안정화 시간을 두고 '
            '측정하시면 최적의 결과를 얻을 수 있습니다.'),
  ]),
  '구독/결제': _FaqCategory(Icons.payment, [
    _FaqItem(
        '구독 플랜 차이점은?',
        'Free: 기본 측정 1회/일, Basic: 무제한 측정+AI코칭, '
            'Pro: 가족 공유+원격진료, Clinical: 의료기관급 분석 지원'),
    _FaqItem(
        '구독 해지 방법은?',
        '설정 > 구독 관리에서 "구독 해지" 버튼을 눌러주세요. '
            '남은 기간까지 서비스 이용이 가능합니다.'),
    _FaqItem('환불 정책은?',
        '결제 후 7일 이내 미사용 시 전액 환불이 가능합니다. 고객 지원으로 문의해주세요.'),
  ]),
  '계정/보안': _FaqCategory(Icons.security, [
    _FaqItem(
        '비밀번호를 잊어버렸어요',
        '로그인 화면에서 "비밀번호를 잊으셨나요?"를 눌러 '
            '이메일 인증을 통해 재설정할 수 있습니다.'),
    _FaqItem(
        '데이터는 안전하게 보관되나요?',
        '모든 건강 데이터는 AES-256 암호화로 저장되며, TLS 1.3 프로토콜로 전송됩니다. '
            'ISO 13485 및 IEC 62304 규정을 준수합니다.'),
  ]),
};
