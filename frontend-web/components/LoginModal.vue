<template>
  <Teleport to="body">
    <div class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-content">
        <div class="modal-header">
          <h2 class="text-xl font-bold">{{ showSms ? '手机验证码登录' : '登录' }}</h2>
          <button class="btn-close" @click="$emit('close')">×</button>
        </div>

        <div class="modal-body">
          <!-- Login Options -->
          <div v-if="!showSms" class="login-options">
            <button class="btn-wechat" @click="wechatLogin">
              <svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
                <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.045c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.027-.407-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z"/>
              </svg>
              微信一键登录
            </button>
            <div class="divider"><span>或</span></div>
            <button class="btn-sms" @click="showSms = true">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/>
              </svg>
              手机验证码登录
            </button>
          </div>

          <!-- SMS Login -->
          <div v-else class="sms-form">
            <div class="form-group">
              <input v-model="phone" type="tel" placeholder="请输入手机号码" maxlength="11" />
            </div>
            <div class="form-group code-group">
              <input v-model="code" type="text" placeholder="请输入验证码" maxlength="6" />
              <button 
                class="btn-send" 
                :disabled="!canSend || sending"
                @click="sendCode"
              >
                {{ sendBtnText }}
              </button>
            </div>
            <button class="btn-login" :disabled="!canLogin" @click="smsLogin">
              登录
            </button>
            <button class="btn-back" @click="showSms = false">
              其他登录方式
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const emit = defineEmits<{
  close: []
  login: [user: any]
}>()

const showSms = ref(false)
const phone = ref('')
const code = ref('')
const sending = ref(false)
const countdown = ref(0)

const canSend = computed(() => /^1[3-9]\d{9}$/.test(phone.value) && countdown.value === 0)
const canLogin = computed(() => phone.value.length === 11 && code.value.length === 6)
const sendBtnText = computed(() => countdown.value > 0 ? `${countdown.value}s后重试` : '获取验证码')

const wechatLogin = () => {
  // Trigger WeChat login
  alert('请使用微信小程序登录')
}

const sendCode = async () => {
  if (!canSend.value) return
  
  sending.value = true
  try {
    await $fetch('/api/auth/sms/send', {
      method: 'POST',
      body: { phone: phone.value }
    })
    
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch (err: any) {
    alert(err.message || '发送失败')
  } finally {
    sending.value = false
  }
}

const smsLogin = async () => {
  try {
    const res = await $fetch('/api/auth/sms/login', {
      method: 'POST',
      body: { phone: phone.value, code: code.value }
    })
    
    emit('login', res.data)
    emit('close')
  } catch (err: any) {
    alert(err.message || '登录失败')
  }
}
</script>

<style scoped>
.modal-overlay {
  @apply fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4;
}

.modal-content {
  @apply bg-white rounded-2xl shadow-xl max-w-md w-full;
}

.modal-header {
  @apply flex justify-between items-center p-6 border-b;
}

.btn-close {
  @apply text-2xl text-gray-400 hover:text-gray-600;
}

.modal-body {
  @apply p-6;
}

.login-options {
  @apply space-y-4;
}

.btn-wechat,
.btn-sms {
  @apply w-full flex items-center justify-center gap-3 py-3 px-4 rounded-xl font-medium transition-colors;
}

.btn-wechat {
  @apply bg-green-500 text-white hover:bg-green-600;
}

.btn-sms {
  @apply bg-blue-500 text-white hover:bg-blue-600;
}

.divider {
  @apply relative text-center;
}

.divider::before {
  content: '';
  @apply absolute top-1/2 left-0 right-0 h-px bg-gray-200;
}

.divider span {
  @apply relative bg-white px-4 text-gray-400 text-sm;
}

.sms-form {
  @apply space-y-4;
}

.form-group input {
  @apply w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.code-group {
  @apply flex gap-2;
}

.code-group input {
  @apply flex-1;
}

.btn-send {
  @apply px-4 py-2 bg-gray-100 text-gray-700 rounded-xl text-sm font-medium whitespace-nowrap;
}

.btn-send:disabled {
  @apply opacity-50 cursor-not-allowed;
}

.btn-login {
  @apply w-full py-3 bg-blue-500 text-white rounded-xl font-medium hover:bg-blue-600 disabled:opacity-50;
}

.btn-back {
  @apply w-full py-2 text-gray-500 text-sm hover:text-gray-700;
}
</style>
