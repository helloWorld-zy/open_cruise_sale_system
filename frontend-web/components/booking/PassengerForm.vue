<template>
  <div class="passenger-form">
    <h3 class="text-xl font-semibold mb-2">填写乘客信息</h3>
    <p class="text-gray-500 mb-6">请填写所有乘客的详细信息，确保与证件一致</p>

    <!-- Contact Information -->
    <div class="contact-section mb-8">
      <h4 class="font-medium text-lg mb-4">联系人信息</h4>
      <div class="grid grid-cols-3 gap-4">
        <div class="form-group">
          <label class="required">姓名</label>
          <input
            v-model="contact.name"
            type="text"
            placeholder="请输入联系人姓名"
            class="form-input"
          />
        </div>
        <div class="form-group">
          <label class="required">手机号码</label>
          <input
            v-model="contact.phone"
            type="tel"
            placeholder="请输入手机号码"
            class="form-input"
          />
        </div>
        <div class="form-group">
          <label class="required">电子邮箱</label>
          <input
            v-model="contact.email"
            type="email"
            placeholder="请输入电子邮箱"
            class="form-input"
          />
        </div>
      </div>
    </div>

    <!-- Passenger Forms -->
    <div class="passengers-section">
      <div
        v-for="(passenger, index) in passengers"
        :key="index"
        class="passenger-card"
      >
        <div class="passenger-header">
          <h4 class="font-medium">乘客 {{ index + 1 }}</h4>
          <span class="passenger-type">{{ getPassengerTypeLabel(passenger.type) }}</span>
        </div>

        <div class="form-grid">
          <div class="form-group">
            <label class="required">姓 (拼音/英文)</label>
            <input
              v-model="passenger.surname"
              type="text"
              placeholder="如: ZHANG"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="required">名 (拼音/英文)</label>
            <input
              v-model="passenger.givenName"
              type="text"
              placeholder="如: SAN"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="required">中文姓名</label>
            <input
              v-model="passenger.name"
              type="text"
              placeholder="如: 张三"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="required">性别</label>
            <select v-model="passenger.gender" class="form-input">
              <option value="">请选择</option>
              <option value="male">男</option>
              <option value="female">女</option>
            </select>
          </div>
          <div class="form-group">
            <label class="required">出生日期</label>
            <input
              v-model="passenger.birthDate"
              type="date"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>国籍</label>
            <input
              v-model="passenger.nationality"
              type="text"
              placeholder="如: 中国"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>护照号码</label>
            <input
              v-model="passenger.passportNumber"
              type="text"
              placeholder="请输入护照号码"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>护照有效期</label>
            <input
              v-model="passenger.passportExpiry"
              type="date"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>身份证号码</label>
            <input
              v-model="passenger.idNumber"
              type="text"
              placeholder="请输入身份证号码"
              class="form-input"
            />
          </div>
          <div class="form-group full-width">
            <label>特殊饮食要求</label>
            <input
              v-model="passenger.dietaryRequirements"
              type="text"
              placeholder="如有素食、过敏等要求请填写"
              class="form-input"
            />
          </div>
          <div class="form-group full-width">
            <label>医疗注意事项</label>
            <input
              v-model="passenger.medicalNotes"
              type="text"
              placeholder="如有特殊医疗需求请填写"
              class="form-input"
            />
          </div>
        </div>

        <div class="emergency-contact mt-4">
          <h5 class="text-sm font-medium text-gray-700 mb-2">紧急联系人</h5>
          <div class="grid grid-cols-2 gap-4">
            <div class="form-group">
              <label>姓名</label>
              <input
                v-model="passenger.emergencyContactName"
                type="text"
                placeholder="紧急联系人姓名"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>电话</label>
              <input
                v-model="passenger.emergencyContactPhone"
                type="tel"
                placeholder="紧急联系人电话"
                class="form-input"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="actions mt-8 flex justify-between">
      <button class="btn-prev" @click="$emit('prev')">
        上一步
      </button>
      <button
        class="btn-next"
        :disabled="!isValid"
        @click="onSubmit"
      >
        下一步
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const props = defineProps<{
  cabins: any[]
  initialData?: any[]
}>()

const emit = defineEmits<{
  submit: [data: { passengers: any[], contact: any }]
  prev: []
  next: []
}>()

const passengers = ref<any[]>([])
const contact = ref({
  name: '',
  phone: '',
  email: ''
})

onMounted(() => {
  if (props.initialData && props.initialData.length > 0) {
    passengers.value = [...props.initialData]
  } else {
    initializePassengers()
  }
})

const initializePassengers = () => {
  passengers.value = []
  
  props.cabins.forEach((cabin, cabinIndex) => {
    // Add adults
    for (let i = 0; i < cabin.adultCount; i++) {
      passengers.value.push(createPassenger('adult', cabinIndex))
    }
    // Add children
    for (let i = 0; i < cabin.childCount; i++) {
      passengers.value.push(createPassenger('child', cabinIndex))
    }
    // Add infants
    for (let i = 0; i < cabin.infantCount; i++) {
      passengers.value.push(createPassenger('infant', cabinIndex))
    }
  })
}

const createPassenger = (type: string, cabinIndex: number) => ({
  type,
  cabinIndex,
  surname: '',
  givenName: '',
  name: '',
  gender: '',
  birthDate: '',
  nationality: '中国',
  passportNumber: '',
  passportExpiry: '',
  idNumber: '',
  phone: '',
  email: '',
  emergencyContactName: '',
  emergencyContactPhone: '',
  dietaryRequirements: '',
  medicalNotes: ''
})

const getPassengerTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    adult: '成人',
    child: '儿童',
    infant: '婴儿'
  }
  return labels[type] || type
}

const isValid = computed(() => {
  // Validate contact info
  if (!contact.value.name || !contact.value.phone || !contact.value.email) {
    return false
  }

  // Validate passengers
  return passengers.value.every(p => 
    p.surname && 
    p.givenName && 
    p.name && 
    p.gender && 
    p.birthDate
  )
})

const onSubmit = () => {
  emit('submit', {
    passengers: passengers.value,
    contact: contact.value
  })
}
</script>

<style scoped>
.passenger-form {
  @apply w-full;
}

.contact-section {
  @apply p-4 bg-blue-50 rounded-lg;
}

.passengers-section {
  @apply space-y-6;
}

.passenger-card {
  @apply border border-gray-200 rounded-lg p-4;
}

.passenger-header {
  @apply flex justify-between items-center mb-4 pb-2 border-b;
}

.passenger-type {
  @apply px-2 py-1 bg-blue-100 text-blue-700 text-sm rounded;
}

.form-grid {
  @apply grid grid-cols-3 gap-4;
}

.form-group {
  @apply flex flex-col;
}

.form-group.full-width {
  @apply col-span-3;
}

.form-group label {
  @apply text-sm text-gray-600 mb-1;
}

.form-group label.required::after {
  content: ' *';
  @apply text-red-500;
}

.form-input {
  @apply px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent;
}

.emergency-contact {
  @apply p-3 bg-gray-50 rounded;
}

.btn-prev {
  @apply px-6 py-2 border border-gray-300 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors;
}

.btn-next {
  @apply px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors;
}
</style>
