<template>
  <div class="booking-wizard">
    <div class="wizard-header">
      <h2 class="text-2xl font-bold text-gray-800 mb-2">预订邮轮</h2>
      <div class="progress-steps">
        <div
          v-for="(step, index) in steps"
          :key="index"
          class="step"
          :class="{
            'active': currentStep === index,
            'completed': currentStep > index
          }"
        >
          <div class="step-number">{{ index + 1 }}</div>
          <div class="step-label">{{ step.label }}</div>
        </div>
      </div>
    </div>

    <div class="wizard-content">
      <!-- Step 1: Select Voyage -->
      <SelectVoyage
        v-if="currentStep === 0"
        :cruise-id="cruiseId"
        :initial-selection="bookingData.voyage"
        @select="onVoyageSelect"
        @next="nextStep"
      />

      <!-- Step 2: Select Cabin -->
      <SelectCabin
        v-if="currentStep === 1"
        :voyage="bookingData.voyage"
        :initial-selection="bookingData.cabins"
        @select="onCabinSelect"
        @prev="prevStep"
        @next="nextStep"
      />

      <!-- Step 3: Passenger Form -->
      <PassengerForm
        v-if="currentStep === 2"
        :cabins="bookingData.cabins"
        :initial-data="bookingData.passengers"
        @submit="onPassengerSubmit"
        @prev="prevStep"
        @next="nextStep"
      />

      <!-- Step 4: Order Confirmation -->
      <OrderConfirm
        v-if="currentStep === 3"
        :booking-data="bookingData"
        @submit="onOrderSubmit"
        @prev="prevStep"
      />
    </div>

    <!-- Error Message -->
    <div v-if="error" class="error-message mt-4 p-4 bg-red-50 text-red-600 rounded-lg">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import SelectVoyage from './SelectVoyage.vue'
import SelectCabin from './SelectCabin.vue'
import PassengerForm from './PassengerForm.vue'
import OrderConfirm from './OrderConfirm.vue'

const props = defineProps<{
  cruiseId: string
}>()

const router = useRouter()
const currentStep = ref(0)
const error = ref('')

const steps = [
  { label: '选择航次' },
  { label: '选择舱房' },
  { label: '填写乘客' },
  { label: '确认订单' }
]

const bookingData = reactive({
  voyage: null as any,
  cabins: [] as any[],
  passengers: [] as any[],
  contact: {
    name: '',
    phone: '',
    email: ''
  }
})

const nextStep = () => {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const onVoyageSelect = (voyage: any) => {
  bookingData.voyage = voyage
}

const onCabinSelect = (cabins: any[]) => {
  bookingData.cabins = cabins
}

const onPassengerSubmit = (data: { passengers: any[], contact: any }) => {
  bookingData.passengers = data.passengers
  bookingData.contact = data.contact
  nextStep()
}

const onOrderSubmit = async () => {
  try {
    error.value = ''
    
    // Build order request
    const orderRequest = {
      cruise_id: props.cruiseId,
      voyage_id: bookingData.voyage.id,
      items: bookingData.cabins.map(cabin => ({
        cabin_id: cabin.id,
        cabin_type_id: cabin.cabin_type_id,
        adult_count: cabin.adultCount,
        child_count: cabin.childCount || 0,
        infant_count: cabin.infantCount || 0
      })),
      passengers: bookingData.passengers.map(p => ({
        name: p.name,
        surname: p.surname,
        given_name: p.givenName,
        gender: p.gender,
        birth_date: p.birthDate,
        nationality: p.nationality,
        passport_number: p.passportNumber,
        passport_expiry: p.passportExpiry,
        id_number: p.idNumber,
        phone: p.phone,
        email: p.email,
        passenger_type: p.passengerType,
        emergency_contact_name: p.emergencyContactName,
        emergency_contact_phone: p.emergencyContactPhone,
        dietary_requirements: p.dietaryRequirements,
        medical_notes: p.medicalNotes
      })),
      contact_name: bookingData.contact.name,
      contact_phone: bookingData.contact.phone,
      contact_email: bookingData.contact.email
    }

    // Create order
    const response = await $fetch('/api/orders', {
      method: 'POST',
      body: orderRequest
    })

    // Navigate to payment page
    router.push(`/payment/${response.data.id}`)
  } catch (err: any) {
    error.value = err.message || '创建订单失败，请重试'
  }
}
</script>

<style scoped>
.booking-wizard {
  @apply max-w-4xl mx-auto p-6;
}

.progress-steps {
  @apply flex justify-between mb-8;
}

.step {
  @apply flex flex-col items-center flex-1 relative;
}

.step-number {
  @apply w-10 h-10 rounded-full bg-gray-200 text-gray-600 flex items-center justify-center font-bold transition-all duration-300;
}

.step.active .step-number {
  @apply bg-blue-500 text-white;
}

.step.completed .step-number {
  @apply bg-green-500 text-white;
}

.step-label {
  @apply mt-2 text-sm text-gray-600;
}

.step.active .step-label {
  @apply text-blue-500 font-medium;
}

.wizard-content {
  @apply bg-white rounded-lg shadow-lg p-6;
}

.error-message {
  @apply border border-red-200;
}
</style>
