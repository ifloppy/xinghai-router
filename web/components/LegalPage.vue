<script setup lang="ts">
import { Bot, ChevronLeft, Moon, Sun } from 'lucide-vue-next'

const props = defineProps<{ kind: 'terms' | 'privacy' }>()
const { locale, toggleLocale, initializeLocale } = useI18n()
const theme = ref<'light' | 'dark'>('light')

const isTerms = computed(() => props.kind === 'terms')
const title = computed(() => locale.value === 'zh-CN' ? (isTerms.value ? '用户协议' : '隐私协议') : (isTerms.value ? 'Terms of Service' : 'Privacy Policy'))
const updatedAt = '2026-07-16'

function setTheme(next: 'light' | 'dark') {
  theme.value = next
  document.documentElement.dataset.theme = next
  localStorage.setItem('xinghai-router-theme', next)
}

onMounted(() => {
  initializeLocale()
  const saved = localStorage.getItem('xinghai-router-theme')
  setTheme(saved === 'dark' || saved === 'light' ? saved : window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
})
</script>

<template>
  <main class="legal-shell">
    <nav class="legal-nav">
      <a class="landing-logo" href="/"><span class="brand-mark small"><Bot :size="19" /></span><span>Xinghai Router</span></a>
      <div class="legal-nav-actions">
        <button class="theme-toggle" :aria-label="theme === 'dark' ? '切换为浅色模式' : '切换为深色模式'" @click="setTheme(theme === 'dark' ? 'light' : 'dark')"><Sun v-if="theme === 'dark'" :size="16" /><Moon v-else :size="16" /></button>
        <button class="language-toggle" aria-label="切换语言" @click="toggleLocale">{{ locale === 'zh-CN' ? 'EN' : '中' }}</button>
        <a class="button ghost" href="/"><ChevronLeft :size="15" />{{ locale === 'zh-CN' ? '返回首页' : 'Back home' }}</a>
      </div>
    </nav>

    <article class="legal-document">
      <header class="legal-header">
        <p class="eyebrow">XINGHAI ROUTER / {{ isTerms ? 'TERMS' : 'PRIVACY' }}</p>
        <h1>{{ title }}</h1>
        <p>{{ locale === 'zh-CN' ? `更新日期：${updatedAt}` : `Last updated: ${updatedAt}` }}</p>
      </header>

      <template v-if="locale === 'zh-CN' && isTerms">
        <p class="legal-lead">欢迎使用 Xinghai Router。请在注册、登录或使用服务前仔细阅读本协议。完成注册或继续使用服务，即表示您已阅读并同意本协议。</p>
        <section><h2>1. 服务说明</h2><p>Xinghai Router 提供模型 API 网关、请求路由、用量统计、账户管理及相关控制台功能。具体可用模型、渠道、价格和功能以页面实际展示为准。</p></section>
        <section><h2>2. 账户与安全</h2><p>您应提供真实、准确且完整的注册信息，并妥善保管账户密码和 API 密钥。因您保管不善、主动披露或设备被他人使用造成的损失，由您自行承担。发现未经授权的使用时，请立即联系我们。</p></section>
        <section><h2>3. 合法使用</h2><p>您不得利用服务从事违法活动，不得绕过访问控制、探测或攻击系统，不得干扰服务运行，不得滥用 API、转售服务或提交侵犯他人权益的内容。您应同时遵守适用法律及上游模型服务商的使用规则。</p></section>
        <section><h2>4. 费用与额度</h2><p>模型调用可能按照控制台展示的价格和倍率计费。请求产生的费用、余额预扣和退款以系统记录为准。余额不足、账户停用或渠道不可用时，相关请求可能失败。</p></section>
        <section><h2>5. 内容与知识产权</h2><p>您对提交的输入和生成结果承担相应责任，并应确保拥有必要的权利。Xinghai Router 及其界面、软件和品牌的相关权利归服务提供方或合法权利人所有。</p></section>
        <section><h2>6. 服务变更与终止</h2><p>我们可能基于维护、安全、合规或产品调整需要变更、暂停或终止部分服务，并尽量提前通知。您违反本协议时，我们可以限制或终止访问。</p></section>
        <section><h2>7. 免责声明与责任限制</h2><p>服务按现状提供。我们不保证服务持续不中断、完全无错误或满足特定用途。除法律另有规定外，对于间接损失、数据丢失或因上游服务造成的损失，我们不承担责任。</p></section>
        <section><h2>8. 协议更新与联系我们</h2><p>我们可能更新本协议，更新后的版本会在本页面发布。若您继续使用服务，即视为接受更新内容。运营主体、联系邮箱及适用法律信息请由部署方在正式上线前补充并替换本段内容。</p></section>
      </template>

      <template v-else-if="locale === 'zh-CN'">
        <p class="legal-lead">本隐私协议说明 Xinghai Router 如何收集、使用、存储和保护您在使用服务时提供的信息。</p>
        <section><h2>1. 我们收集的信息</h2><p>注册时我们会收集姓名、邮箱和密码哈希。使用服务时，我们可能记录账户、API 密钥标识、请求时间、模型、渠道、状态码、耗时、Token 用量和费用等运行数据。我们不会以明文保存您的账户密码或完整 API 密钥。</p></section>
        <section><h2>2. 信息使用目的</h2><p>我们使用这些信息提供登录认证、请求转发、计量结算、故障排查、安全防护、审计和服务改进。除非获得您的同意或法律要求，我们不会将个人信息用于无关目的。</p></section>
        <section><h2>3. 上游处理</h2><p>为完成模型请求，必要的请求内容会按您选择的渠道转发给相应的上游模型服务商。请在使用前阅读相关上游服务商的条款和隐私政策，并避免提交不必要的敏感个人信息。</p></section>
        <section><h2>4. 存储与安全</h2><p>信息会在实现服务所需的期限内保存。我们通过密码哈希、密钥脱敏、权限控制和审计日志等措施保护信息，但互联网传输和存储不存在绝对安全。</p></section>
        <section><h2>5. 共享与披露</h2><p>我们仅在提供服务所必需、获得授权、履行法律义务或保护服务及用户安全时共享或披露信息。我们不会出售您的个人信息。</p></section>
        <section><h2>6. 您的权利</h2><p>在适用法律允许的范围内，您可以请求访问、更正或删除账户信息，也可以申请注销账户。注销可能不影响法律要求保留的记录。运营主体和联系邮箱请由部署方在正式上线前补充。</p></section>
        <section><h2>7. 协议更新</h2><p>我们可能根据服务或法律变化更新本协议，并在本页面标注更新日期。继续使用服务表示您接受更新后的协议。</p></section>
      </template>

      <template v-else-if="isTerms">
        <p class="legal-lead">Welcome to Xinghai Router. By registering, signing in, or using the service, you confirm that you have read and accepted these Terms of Service.</p>
        <section><h2>1. Service</h2><p>Xinghai Router provides an API gateway, model routing, usage tracking, account management, and related console features. Available models, channels, prices, and features are shown in the product.</p></section>
        <section><h2>2. Accounts and security</h2><p>Keep your registration information accurate and protect your password and API keys. You are responsible for activity caused by sharing or failing to secure your credentials.</p></section>
        <section><h2>3. Acceptable use</h2><p>Do not use the service unlawfully, bypass access controls, attack or disrupt systems, abuse the API, resell the service, or infringe others' rights. You must follow applicable law and upstream provider rules.</p></section>
        <section><h2>4. Charges</h2><p>Requests may be charged according to the prices and multipliers displayed in the console. System records determine charges, reservations, and refunds.</p></section>
        <section><h2>5. Content and intellectual property</h2><p>You are responsible for your inputs and outputs and must have the necessary rights. Xinghai Router and its software, interface, and brand remain the property of their respective rights holders.</p></section>
        <section><h2>6. Changes and termination</h2><p>We may change, suspend, or terminate features for maintenance, security, compliance, or product reasons. Access may be restricted for violations of these terms.</p></section>
        <section><h2>7. Disclaimer and liability</h2><p>The service is provided as-is. To the extent permitted by law, we disclaim warranties and liability for indirect loss, data loss, or upstream service failures.</p></section>
        <section><h2>8. Updates and contact</h2><p>Updates will be posted on this page. The deploying operator must add its legal entity, contact email, and governing-law details before production launch.</p></section>
      </template>

      <template v-else>
        <p class="legal-lead">This Privacy Policy explains how Xinghai Router collects, uses, stores, and protects information when you use the service.</p>
        <section><h2>1. Information we collect</h2><p>We collect your name, email, and a password hash during registration. We may record account identifiers, API key identifiers, request time, model, channel, status, latency, token usage, and charges. We do not store your password or complete API keys in plain text.</p></section>
        <section><h2>2. How we use information</h2><p>We use information for authentication, request routing, metering, billing, troubleshooting, security, auditing, and service improvement.</p></section>
        <section><h2>3. Upstream processing</h2><p>To fulfill a model request, necessary request content is sent to the upstream provider selected for that request. Review the provider's policies and avoid submitting unnecessary sensitive personal information.</p></section>
        <section><h2>4. Storage and security</h2><p>Information is retained as needed to operate the service. We use password hashing, key redaction, access controls, and audit logs, but no internet system is completely secure.</p></section>
        <section><h2>5. Sharing</h2><p>We share information only as needed to provide the service, with authorization, to comply with law, or to protect users and the service. We do not sell personal information.</p></section>
        <section><h2>6. Your rights</h2><p>Where required by law, you may request access, correction, deletion, or account closure. The deploying operator must add its legal entity and contact email before production launch.</p></section>
        <section><h2>7. Updates</h2><p>We may update this policy as the service or law changes. The updated date will be shown on this page.</p></section>
      </template>
    </article>
  </main>
</template>
