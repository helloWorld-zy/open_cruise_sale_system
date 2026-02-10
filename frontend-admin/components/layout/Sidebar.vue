<template>
  <aside class="w-64 bg-gray-900 min-h-screen flex flex-col">
    <!-- Logo area -->
    <div class="h-16 flex items-center px-6 border-b border-gray-800">
      <span class="text-white text-lg font-semibold">管理后台</span>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-4 py-6 space-y-1">
      <NuxtLink
        v-for="item in navigation"
        :key="item.to"
        :to="item.to"
        :class="[
          'flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
          $route.path === item.to
            ? 'bg-primary-600 text-white'
            : 'text-gray-300 hover:bg-gray-800 hover:text-white'
        ]"
      >
        <UIcon :name="item.icon" class="w-5 h-5" />
        <span>{{ item.label }}</span>
      </NuxtLink>
    </nav>

    <!-- Footer -->
    <div class="px-6 py-4 border-t border-gray-800">
      <p class="text-gray-400 text-xs">
        © 2026 CruiseBooking
      </p>
    </div>
  </aside>
</template>

<script setup lang="ts">
const navigation = [
  { label: 'Dashboard', to: '/', icon: 'i-heroicons-home' },
  { label: 'Cruises', to: '/cruises', icon: 'i-heroicons-ship' },
  { label: 'Cabin Types', to: '/cabin-types', icon: 'i-heroicons-square-2-stack' },
  { label: 'Facilities', to: '/facilities', icon: 'i-heroicons-building-storefront' },
  { label: 'Routes', to: '/routes', icon: 'i-heroicons-map' },
  { label: 'Voyages', to: '/voyages', icon: 'i-heroicons-calendar' },
  { label: 'Cabins', to: '/cabins', icon: 'i-heroicons-key' },
  { label: 'Orders', to: '/orders', icon: 'i-heroicons-shopping-cart' },
  { label: 'Staff', to: '/staffs', icon: 'i-heroicons-users' },
  { label: 'Settings', to: '/settings', icon: 'i-heroicons-cog-6-tooth' },
]
</script>
