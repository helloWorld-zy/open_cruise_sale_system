<template>
  <div class="passengers-page">
    <div class="container">
      <div class="page-header">
        <button class="btn-back" @click="$router.back()">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
          返回
        </button>
        <h1 class="text-2xl font-bold">常用乘客</h1>
        <button class="btn-add" @click="showForm = true">添加</button>
      </div>

      <!-- Empty State -->
      <div v-if="passengers.length === 0" class="empty-state">
        <p class="text-gray-500">暂无常用乘客</p>
        <button class="btn-primary mt-4" @click="showForm = true">添加乘客</button>
      </div>

      <!-- Passenger List -->
      <div v-else class="passenger-list">
        <div v-for="p in passengers" :key="p.id" class="passenger-card">
          <div class="passenger-info">
            <div class="name">{{ p.name }}</div>
            <div class="detail">{{ p.surname }} {{ p.given_name }} | {{ p.gender === 'male' ? '男' : '女' }}</div>
            <div class="doc">{{ p.id_number || p.passport_number }}</div>
            <span v-if="p.is_default" class="badge">默认</span>
          </div>
          <div class="actions">
            <button @click="editPassenger(p)">编辑</button>
            <button class="text-red-500" @click="deletePassenger(p.id)">删除</button>
          </div>
        </div>
      </div>

      <!-- Add/Edit Form Modal -->
      <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
        <div class="modal-content">
          <h3 class="text-xl font-bold mb-4">{{ editing ? '编辑乘客' : '添加乘客' }}</h3>
          <form @submit.prevent="savePassenger">
            <div class="form-group">
              <label>中文姓名</label>
              <input v-model="form.name" required />
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>姓(拼音)</label>
                <input v-model="form.surname" required />
              </div>
              <div class="form-group">
                <label>名(拼音)</label>
                <input v-model="form.given_name" />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>性别</label>
                <select v-model="form.gender" required>
                  <option value="male">男</option>
                  <option value="female">女</option>
                </select>
              </div>
              <div class="form-group">
                <label>出生日期</label>
                <input v-model="form.birth_date" type="date" required />
              </div>
            </div>
            <div class="form-group">
              <label>证件号码</label>
              <input v-model="form.id_number" placeholder="身份证号或护照号" />
            </div>
            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="showForm = false">取消</button>
              <button type="submit" class="btn-save">保存</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const passengers = ref<any[]>([])
const showForm = ref(false)
const editing = ref(false)
const form = ref({
  id: '',
  name: '',
  surname: '',
  given_name: '',
  gender: 'male',
  birth_date: '',
  id_number: '',
  is_default: false
})

onMounted(loadPassengers)

async function loadPassengers() {
  const res = await $fetch('/api/user/passengers')
  passengers.value = res.data || []
}

function editPassenger(p: any) {
  form.value = { ...p }
  editing.value = true
  showForm.value = true
}

async function savePassenger() {
  try {
    if (editing.value) {
      await $fetch(`/api/user/passengers/${form.value.id}`, {
        method: 'PUT',
        body: form.value
      })
    } else {
      await $fetch('/api/user/passengers', {
        method: 'POST',
        body: form.value
      })
    }
    showForm.value = false
    resetForm()
    loadPassengers()
  } catch (err: any) {
    alert(err.message || '保存失败')
  }
}

async function deletePassenger(id: string) {
  if (!confirm('确定要删除此乘客吗？')) return
  
  try {
    await $fetch(`/api/user/passengers/${id}`, { method: 'DELETE' })
    loadPassengers()
  } catch (err: any) {
    alert(err.message || '删除失败')
  }
}

function resetForm() {
  form.value = {
    id: '',
    name: '',
    surname: '',
    given_name: '',
    gender: 'male',
    birth_date: '',
    id_number: '',
    is_default: false
  }
  editing.value = false
}
</script>

<style scoped>
.passengers-page {
  @apply min-h-screen bg-gray-50 py-8;
}

.container {
  @apply max-w-2xl mx-auto px-4;
}

.page-header {
  @apply flex items-center gap-4 mb-6;
}

.btn-back {
  @apply flex items-center gap-2 text-gray-600;
}

.btn-add {
  @apply ml-auto px-4 py-2 bg-blue-500 text-white rounded-lg;
}

.empty-state {
  @apply text-center py-12;
}

.btn-primary {
  @apply px-6 py-2 bg-blue-500 text-white rounded-lg;
}

.passenger-list {
  @apply space-y-4;
}

.passenger-card {
  @apply bg-white rounded-xl p-4 shadow-sm flex justify-between items-center;
}

.passenger-info .name {
  @apply font-bold;
}

.passenger-info .detail,
.passenger-info .doc {
  @apply text-sm text-gray-500;
}

.badge {
  @apply inline-block px-2 py-1 bg-blue-100 text-blue-600 text-xs rounded-full mt-1;
}

.actions {
  @apply flex gap-2;
}

.actions button {
  @apply px-3 py-1 text-sm text-blue-500;
}

.modal-overlay {
  @apply fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4;
}

.modal-content {
  @apply bg-white rounded-xl p-6 w-full max-w-md;
}

.form-group {
  @apply mb-4;
}

.form-group label {
  @apply block text-sm font-medium mb-1;
}

.form-group input,
.form-group select {
  @apply w-full px-3 py-2 border border-gray-300 rounded-lg;
}

.form-row {
  @apply grid grid-cols-2 gap-4;
}

.form-actions {
  @apply flex gap-4 mt-6;
}

.btn-cancel {
  @apply flex-1 py-2 border border-gray-300 rounded-lg;
}

.btn-save {
  @apply flex-1 py-2 bg-blue-500 text-white rounded-lg;
}
</style>
