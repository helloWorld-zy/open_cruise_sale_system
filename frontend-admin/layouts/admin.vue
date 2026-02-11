<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Sidebar -->
    <Sidebar />
    
    <!-- Main Content -->
    <div class="lg:pl-64">
      <!-- Header -->
      <Header />
      
      <!-- Page Content -->
      <main class="py-6 px-4 sm:px-6 lg:px-8">
        <slot />
      </main>
      
      <!-- Footer -->
      <Footer />
    </div>
  </div>
</template>

<script setup lang="ts">
// Admin layout with sidebar navigation
// Automatically wraps all admin pages
</script>
