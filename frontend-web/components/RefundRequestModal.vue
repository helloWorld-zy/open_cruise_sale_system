<template>
  <Teleport to="body">
    <div class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">申请退款</h3>
          <button class="btn-close" @click="$emit('close')">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <!-- Order Summary -->
          <div class="order-summary">
            <h4>订单信息</h4>
            <p>订单号: {{ order?.order_number }}</p>
            <p class="order-amount">订单金额: ¥{{ formatPrice(order?.total_amount) }}</p>
          </div>

          <!-- Refund Amount -->
          <div class="form-group">
            <label class="required">退款金额</label>
            <div class="amount-input">
              <span class="currency">¥</span>
              <input
                v-model.number="refundAmount"
                type="number"
                :max="maxRefundableAmount"
                min="0.01"
                step="0.01"
                class="form-input"
              />
            </div>
            <p class="hint">最大可退款金额: ¥{{ formatPrice(maxRefundableAmount) }}</p>
            <div class="amount-presets">
              <button
                v-for="preset in amountPresets"
                :key="preset"
                class="preset-btn"
                :class="{ active: refundAmount === preset }"
                @click="refundAmount = preset"
              >
                ¥{{ formatPrice(preset) }}
              </button>
              <button
                class="preset-btn"
                :class="{ active: refundAmount === maxRefundableAmount }"
                @click="refundAmount = maxRefundableAmount"
              >
                全额
              </button>
            </div>
          </div>

          <!-- Refund Type -->
          <div class="form-group">
            <label class="required">退款类型</label>
            <div class="radio-group">
              <label class="radio-label">
                <input v-model="refundType" type="radio" value="full" />
                <span>全额退款</span>
              </label>
              <label class="radio-label">
                <input v-model="refundType" type="radio" value="partial" />
                <span>部分退款</span>
              </label>
            </div>
          </div>

          <!-- Cancellation Reason -->
          <div class="form-group">
            <label class="required">取消原因</label>
            <select v-model="cancellationReason" class="form-select">
              <option value="">请选择取消原因</option>
              <option value="customer_request">个人原因</option>
              <option value="voyage_cancelled">航次取消</option>
              <option value="cabin_upgrade">升级舱房</option>
              <option value="other">其他原因</option>
            </select>
          </div>

          <!-- Refund Reason -->
          <div class="form-group">
            <label class="required">退款说明</label>
            <textarea
              v-model="refundReason"
              rows="3"
              placeholder="请详细说明退款原因，以便我们更快处理..."
              class="form-textarea"
            ></textarea>
          </div>

          <!-- Refund Method -->
          <div class="form-group">
            <label>退款方式</label>
            <div class="radio-group">
              <label class="radio-label">
                <input v-model="refundMethod" type="radio" value="original" />
                <span>原路退回</span>
              </label>
              <label class="radio-label">
                <input v-model="refundMethod" type="radio" value="bank" />
                <span>银行卡</span>
              </label>
            </div>
          </div>

          <!-- Bank Info (if bank transfer selected) -->
          <div v-if="refundMethod === 'bank'" class="bank-info">
            <div class="form-group">
              <label class="required">开户银行</label>
              <input v-model="bankName" type="text" class="form-input" placeholder="如: 中国工商银行" />
            </div>
            <div class="form-group">
              <label class="required">银行卡号</label>
              <input v-model="bankAccount" type="text" class="form-input" placeholder="请输入银行卡号" />
            </div>
            <div class="form-group">
              <label class="required">开户姓名</label>
              <input v-model="accountHolder" type="text" class="form-input" placeholder="请输入开户姓名" />
            </div>
          </div>

          <!-- Warning -->
          <div class="warning-box">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
            <p>退款申请提交后，工作人员将在1-3个工作日内审核。审核通过后，款项将在3-7个工作日内退回至您的账户。</p>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-cancel" @click="$emit('close')">取消</button>
          <button
            class="btn-submit"
            :disabled="!isValid || submitting"
            @click="submit"
          >
            <span v-if="submitting" class="loading-spinner"></span>
            <span>{{ submitting ? '提交中...' : '提交申请' }}</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  order: any
}>()

const emit = defineEmits<{
  close: []
  submit: [data: any]
}>()

const refundAmount = ref(0)
const refundType = ref('partial')
const cancellationReason = ref('')
const refundReason = ref('')
const refundMethod = ref('original')
const bankName = ref('')
const bankAccount = ref('')
const accountHolder = ref('')
const submitting = ref(false)

const maxRefundableAmount = computed(() => {
  return props.order?.total_amount || 0
})

const amountPresets = computed(() => {
  const max = maxRefundableAmount.value
  return [
    Math.round(max * 0.25 * 100) / 100,
    Math.round(max * 0.5 * 100) / 100,
    Math.round(max * 0.75 * 100) / 100
  ]
})

const isValid = computed(() => {
  if (refundAmount.value <= 0 || refundAmount.value > maxRefundableAmount.value) return false
  if (!cancellationReason.value) return false
  if (!refundReason.value.trim()) return false
  if (refundMethod.value === 'bank') {
    if (!bankName.value.trim()) return false
    if (!bankAccount.value.trim()) return false
    if (!accountHolder.value.trim()) return false
  }
  return true
})

const submit = () => {
  if (!isValid.value) return

  submitting.value = true

  const data = {
    order_id: props.order.id,
    refund_amount: refundAmount.value,
    refund_reason: refundReason.value,
    refund_type: refundType.value,
    cancellation_reason: cancellationReason.value,
    refund_method: refundMethod.value,
    bank_name: refundMethod.value === 'bank' ? bankName.value : undefined,
    bank_account: refundMethod.value === 'bank' ? bankAccount.value : undefined,
    account_holder: refundMethod.value === 'bank' ? accountHolder.value : undefined
  }

  emit('submit', data)
  submitting.value = false
}

const formatPrice = (price: number) => {
  return price?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}
</script>

<style scoped>
.modal-overlay {
  @apply fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4;
}

.modal-content {
  @apply bg-white rounded-xl shadow-xl max-w-lg w-full max-h-[90vh] overflow-y-auto;
}

.modal-header {
  @apply flex justify-between items-center p-4 border-b;
}

.modal-title {
  @apply text-xl font-bold text-gray-800;
}

.btn-close {
  @apply text-gray-400 hover:text-gray-600;
}

.modal-body {
  @apply p-4 space-y-4;
}

.order-summary {
  @apply p-3 bg-gray-50 rounded-lg;
}

.order-summary h4 {
  @apply font-medium text-gray-800 mb-2;
}

.order-summary p {
  @apply text-sm text-gray-600 mb-1;
}

.order-amount {
  @apply font-bold text-orange-500 text-lg;
}

.form-group {
  @apply space-y-2;
}

.form-group label {
  @apply block text-sm font-medium text-gray-700;
}

.form-group label.required::after {
  content: ' *';
  @apply text-red-500;
}

.form-input,
.form-select,
.form-textarea {
  @apply w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent;
}

.amount-input {
  @apply relative flex items-center;
}

.amount-input .currency {
  @apply absolute left-3 text-gray-500;
}

.amount-input input {
  @apply pl-8;
}

.hint {
  @apply text-xs text-gray-500;
}

.amount-presets {
  @apply flex gap-2 flex-wrap;
}

.preset-btn {
  @apply px-3 py-1 bg-gray-100 text-sm text-gray-700 rounded-lg hover:bg-gray-200 transition-colors;
}

.preset-btn.active {
  @apply bg-blue-500 text-white;
}

.radio-group {
  @apply flex gap-4;
}

.radio-label {
  @apply flex items-center gap-2 cursor-pointer;
}

.radio-label input {
  @apply w-4 h-4 text-blue-500;
}

.bank-info {
  @apply p-3 bg-gray-50 rounded-lg space-y-3;
}

.warning-box {
  @apply flex gap-3 p-3 bg-yellow-50 rounded-lg text-sm text-yellow-700;
}

.warning-box svg {
  @apply flex-shrink-0 text-yellow-500;
}

.modal-footer {
  @apply flex gap-3 p-4 border-t;
}

.btn-cancel {
  @apply flex-1 py-2 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50;
}

.btn-submit {
  @apply flex-1 py-2 bg-blue-500 text-white font-medium rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2;
}

.loading-spinner {
  @apply w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin;
}
</style>
