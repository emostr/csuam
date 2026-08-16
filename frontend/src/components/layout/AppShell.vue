<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './Sidebar.vue'
import Topbar from './Topbar.vue'
import PoweredBy from '@/components/PoweredBy.vue'

const sidebarOpen = ref(false)
</script>

<template>
  <div class="min-h-screen flex bg-bg">
    <Transition name="modal">
      <div
        v-if="sidebarOpen"
        class="fixed inset-0 z-30 bg-black/60 lg:hidden"
        @click="sidebarOpen = false"
      />
    </Transition>

    <Sidebar :open="sidebarOpen" @close="sidebarOpen = false" />

    <div class="flex-1 min-w-0 flex flex-col">
      <Topbar @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="flex-1 p-4 sm:p-6 lg:p-8">
        <RouterView v-slot="{ Component }">
          <Transition name="view-fade" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>

      <footer class="border-t border-line px-4 sm:px-8 py-4 flex flex-col sm:flex-row items-center justify-between gap-3">
        <p class="text-xs text-faint">ЦСУАМ "Архивли"</p>
        <PoweredBy />
        <div class="flex items-center gap-4 text-xs text-faint">
          <span>2026</span>
        </div>
      </footer>
    </div>
  </div>
</template>
