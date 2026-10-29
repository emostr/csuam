<script setup lang="ts">
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import PoweredBy from '@/components/PoweredBy.vue'
</script>

<template>
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset class="min-w-0">
      <AppHeader />

      <div class="mx-auto w-full max-w-screen-2xl flex-1 p-4 md:p-6 lg:p-8">
        <RouterView v-slot="{ Component }">
          <Transition name="view-fade" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </div>

      <footer
        class="flex flex-col items-center justify-between gap-3 border-t px-4 py-4 text-xs text-muted-foreground sm:flex-row md:px-6"
      >
        <p>ЦСУАМ «Архивли»</p>
        <PoweredBy />
        <span>2026</span>
      </footer>
    </SidebarInset>
  </SidebarProvider>
</template>
