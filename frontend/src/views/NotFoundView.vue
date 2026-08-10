<script setup>
import { useRouter } from 'vue-router'
import NButton from '@/components/ui/NButton.vue'

const router = useRouter()
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-bg p-6 text-center">
    <div class="ng-enter">
      <div class="text-[7rem] sm:text-[10rem] font-extrabold text-ink leading-none tracking-normal">
        4<span class="text-accent">0</span>4
      </div>
      <div class="w-16 h-1 bg-accent mx-auto mb-6" />
      <h1 class="text-xl font-bold text-ink">Страница не найдена</h1>
      <p class="text-muted mt-2 max-w-sm mx-auto">Возможно, она была перемещена или удалена. А может, её и не было никогда</p>
      <div class="flex items-center justify-center gap-2.5 mt-8">
        <NButton variant="secondary" icon="chevronLeft" @click="router.back()">Назад</NButton>
        <NButton icon="home" @click="router.push('/')">На главную</NButton>
      </div>
    </div>
  </div>
</template>
