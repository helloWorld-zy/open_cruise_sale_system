<template>
  <view class="login-page">
    <view class="logo">
      <image src="/static/logo.png" mode="aspectFit" />
      <text class="title">邮轮预订</text>
    </view>
    
    <view class="login-form">
      <button 
        class="btn-wechat" 
        open-type="getPhoneNumber" 
        @getphonenumber="onGetPhoneNumber"
      >
        <text class="icon">微信</text>
        <text>微信一键登录</text>
      </button>
      
      <view class="divider">
        <text>或</text>
      </view>
      
      <view class="phone-form">
        <input 
          v-model="phone" 
          type="number" 
          placeholder="请输入手机号码" 
          maxlength="11"
        />
        <view class="code-row">
          <input 
            v-model="code" 
            type="number" 
            placeholder="请输入验证码" 
            maxlength="6"
          />
          <button 
            class="btn-code" 
            :disabled="!canSend || sending"
            @click="sendCode"
          >
            {{ sendBtnText }}
          </button>
        </view>
        <button class="btn-login" :disabled="!canLogin" @click="phoneLogin">
          登录
        </button>
      </view>
    </view>
    
    <view class="agreement">
      <checkbox :checked="agreed" @click="agreed = !agreed" />
      <text>登录即表示同意</text>
      <text class="link" @click="showAgreement">《用户协议》</text>
      <text>和</text>
      <text class="link" @click="showPrivacy">《隐私政策》</text>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'

const phone = ref('')
const code = ref('')
const agreed = ref(false)
const sending = ref(false)
const countdown = ref(0)

const canSend = computed(() => /^1[3-9]\d{9}$/.test(phone.value) && countdown.value === 0)
const canLogin = computed(() => phone.value.length === 11 && code.value.length === 6 && agreed.value)
const sendBtnText = computed(() => countdown.value > 0 ? `${countdown.value}s` : '获取验证码')

// WeChat phone login
const onGetPhoneNumber = async (e) => {
  if (!agreed.value) {
    uni.showToast({ title: '请先同意用户协议', icon: 'none' })
    return
  }
  
  if (e.detail.errMsg.includes('fail')) {
    uni.showToast({ title: '授权失败', icon: 'none' })
    return
  }
  
  try {
    // Get login code
    const [loginRes] = await uni.login({ provider: 'weixin' })
    
    // Send to backend
    const res = await uni.request({
      url: `${getApp().globalData.apiBaseUrl}/auth/wechat/phone-login`,
      method: 'POST',
      data: {
        code: loginRes.code,
        encrypted_data: e.detail.encryptedData,
        iv: e.detail.iv
      }
    })
    
    if (res.data.data) {
      // Save tokens
      uni.setStorageSync('access_token', res.data.data.access_token)
      uni.setStorageSync('refresh_token', res.data.data.refresh_token)
      
      uni.showToast({ title: '登录成功' })
      uni.switchTab({ url: '/pages/index/index' })
    }
  } catch (err) {
    uni.showToast({ title: '登录失败', icon: 'none' })
  }
}

// Send SMS code
const sendCode = async () => {
  if (!canSend.value) return
  
  sending.value = true
  try {
    await uni.request({
      url: `${getApp().globalData.apiBaseUrl}/auth/sms/send`,
      method: 'POST',
      data: { phone: phone.value }
    })
    
    uni.showToast({ title: '验证码已发送' })
    
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch (err) {
    uni.showToast({ title: '发送失败', icon: 'none' })
  } finally {
    sending.value = false
  }
}

// Phone login
const phoneLogin = async () => {
  if (!agreed.value) {
    uni.showToast({ title: '请先同意用户协议', icon: 'none' })
    return
  }
  
  try {
    const res = await uni.request({
      url: `${getApp().globalData.apiBaseUrl}/auth/sms/login`,
      method: 'POST',
      data: { phone: phone.value, code: code.value }
    })
    
    if (res.data.data) {
      uni.setStorageSync('access_token', res.data.data.access_token)
      uni.setStorageSync('refresh_token', res.data.data.refresh_token)
      
      uni.showToast({ title: '登录成功' })
      uni.switchTab({ url: '/pages/index/index' })
    }
  } catch (err) {
    uni.showToast({ title: '登录失败', icon: 'none' })
  }
}

const showAgreement = () => {
  uni.navigateTo({ url: '/pages/agreement/index' })
}

const showPrivacy = () => {
  uni.navigateTo({ url: '/pages/privacy/index' })
}
</script>

<style>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx;
}

.logo {
  text-align: center;
  margin-bottom: 80rpx;
}

.logo image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 40rpx;
  background: white;
  padding: 20rpx;
}

.logo .title {
  display: block;
  color: white;
  font-size: 48rpx;
  font-weight: bold;
  margin-top: 30rpx;
}

.login-form {
  background: white;
  border-radius: 30rpx;
  padding: 60rpx 40rpx;
}

.btn-wechat {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20rpx;
  background: #07c160;
  color: white;
  height: 100rpx;
  border-radius: 50rpx;
  font-size: 32rpx;
}

.btn-wechat .icon {
  font-weight: bold;
}

.divider {
  text-align: center;
  margin: 40rpx 0;
  position: relative;
}

.divider::before,
.divider::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 30%;
  height: 2rpx;
  background: #e5e5e5;
}

.divider::before { left: 0; }
.divider::after { right: 0; }

.divider text {
  color: #999;
  font-size: 28rpx;
  padding: 0 20rpx;
}

.phone-form input {
  height: 90rpx;
  border-bottom: 2rpx solid #e5e5e5;
  font-size: 30rpx;
  margin-bottom: 30rpx;
}

.code-row {
  display: flex;
  gap: 20rpx;
  margin-bottom: 40rpx;
}

.code-row input {
  flex: 1;
  margin-bottom: 0;
}

.btn-code {
  width: 200rpx;
  height: 90rpx;
  line-height: 90rpx;
  background: #f0f0f0;
  color: #666;
  font-size: 26rpx;
  border-radius: 10rpx;
}

.btn-code[disabled] {
  color: #999;
}

.btn-login {
  background: #667eea;
  color: white;
  height: 100rpx;
  line-height: 100rpx;
  border-radius: 50rpx;
  font-size: 32rpx;
}

.btn-login[disabled] {
  opacity: 0.5;
}

.agreement {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 60rpx;
  color: white;
  font-size: 26rpx;
}

.agreement checkbox {
  margin-right: 10rpx;
}

.agreement .link {
  color: #ffd700;
  text-decoration: underline;
}
</style>
