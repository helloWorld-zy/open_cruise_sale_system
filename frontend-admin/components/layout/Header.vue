<template>
  <header class="bg-white border-b border-gray-200 sticky top-0 z-50">
    <div class="px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <div class="flex items-center">
          <NuxtLink to="/" class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-lg">C</span>
            </div>
            <span class="text-xl font-semibold text-gray-900">CruiseBooking</span>
          </NuxtLink>
        </div>

        <!-- Right side -->
        <div class="flex items-center space-x-4">
          <!-- User menu -->
          <UDropdown :items="userMenuItems">
            <UButton color="gray" variant="ghost">
              <template #leading>
                <UAvatar
                  :src="authStore.user?.avatarUrl"
                  :alt="authStore.user?.name"
                  size="sm"
                />
              </template>
              {{ authStore.user?.name || 'Admin' }}
            </UButton>
          </UDropdown>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const authStore = useAuthStore()

const userMenuItems = [
  [{
    label: 'Profile',
    icon: 'i-heroicons-user-circle',
    to: '/profile'
  }, {
    label: 'Settings',
    icon: 'i-heroicons-cog-6-tooth',
    to: '/settings'
  }],
  [{
    label: 'Sign out',
    icon: 'i-heroicons-arrow-right-on-rectangle',
    click: () => authStore.logout()
  }]
]
</script>
