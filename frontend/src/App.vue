<script setup lang="ts">
import { RouterView } from 'vue-router'
import { Toaster } from '@/components/ui/sonner'
import NotifyDialogs from '@/components/NotifyDialogs.vue'
import { theme } from '@/lib/theme'
</script>

<template>
  <RouterView />
  <NotifyDialogs />
  <Toaster :theme="theme" close-button />
</template>
