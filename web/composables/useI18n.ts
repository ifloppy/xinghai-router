import { computed } from 'vue'

export type Locale = 'zh-CN' | 'en-US'

const messages = {
  'zh-CN': {
    home: '首页',
    rankings: '排行榜',
    quickStart: '快速开始',
    marketplace: '模型广场',
    console: '进入控制台',
    login: '登录',
    lightMode: '切换为浅色模式',
    darkMode: '切换为深色模式',
    language: '语言',
    chinese: '中文',
    english: 'English',
    titleLogin: '登录',
    titleMarketplace: '模型广场',
    titleRankings: '模型排行榜',
    overview: '概览',
    users: '用户管理',
    groups: '分组管理',
    keys: 'API 密钥管理',
    channels: '渠道管理',
    logs: '使用日志',
    account: 'API 密钥',
    profile: '个人资料',
    wallet: '钱包',
    usage: '使用日志',
    usageOverview: '用量概览',
    ledger: '余额流水',
    pricing: '模型定价',
    audit: '审计日志',
    siteSettings: '站点设置',
    general: '常规',
    management: '管理',
    personal: '个人',
    billing: '账户',
    refresh: '刷新',
    backHome: '返回首页',
    switchLanguage: '切换语言',
  },
  'en-US': {
    home: 'Home',
    rankings: 'Rankings',
    quickStart: 'Quick start',
    marketplace: 'Model catalog',
    console: 'Open console',
    login: 'Sign in',
    lightMode: 'Switch to light mode',
    darkMode: 'Switch to dark mode',
    language: 'Language',
    chinese: '中文',
    english: 'English',
    titleLogin: 'Sign in',
    titleMarketplace: 'Model catalog',
    titleRankings: 'Model rankings',
    overview: 'Overview',
    users: 'Users',
    groups: 'Groups',
    keys: 'API keys',
    channels: 'Channels',
    logs: 'Usage logs',
    account: 'API keys',
    profile: 'Profile',
    wallet: 'Wallet',
    usage: 'Usage logs',
    usageOverview: 'Usage overview',
    ledger: 'Balance ledger',
    pricing: 'Model pricing',
    audit: 'Audit logs',
    siteSettings: 'Site settings',
    general: 'General',
    management: 'Management',
    personal: 'Personal',
    billing: 'Account',
    refresh: 'Refresh',
    backHome: 'Back home',
    switchLanguage: 'Switch language',
  },
} as const

type MessageKey = keyof typeof messages['zh-CN']

export function useI18n() {
  const locale = useState<Locale>('xinghai-router-locale', () => 'zh-CN')
  const initialized = useState<boolean>('xinghai-router-locale-initialized', () => false)

  const t = (key: MessageKey) => messages[locale.value][key]
  const languageName = computed(() => locale.value === 'zh-CN' ? t('chinese') : t('english'))

  function setLocale(next: Locale) {
    locale.value = next
    if (import.meta.client) {
      localStorage.setItem('xinghai-router-locale', next)
      document.documentElement.lang = next
    }
  }

  function toggleLocale() {
    setLocale(locale.value === 'zh-CN' ? 'en-US' : 'zh-CN')
  }

  function initializeLocale() {
    if (!import.meta.client || initialized.value) return
    const saved = localStorage.getItem('xinghai-router-locale') as Locale | null
    const preferred = navigator.language.toLowerCase().startsWith('zh') ? 'zh-CN' : 'en-US'
    setLocale(saved === 'zh-CN' || saved === 'en-US' ? saved : preferred)
    initialized.value = true
  }

  return { locale, t, languageName, setLocale, toggleLocale, initializeLocale }
}
