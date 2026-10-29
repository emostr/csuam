<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ChevronLeftIcon, HouseIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Empty, EmptyContent, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'

const router = useRouter()
</script>

<template>
  <div class="flex min-h-svh items-center justify-center bg-background p-6">
    <Empty>
      <EmptyHeader>
        <p class="text-7xl font-semibold tracking-tighter text-primary tabular-nums sm:text-8xl">404</p>
        <EmptyTitle class="text-lg">Страница не найдена</EmptyTitle>
        <EmptyDescription>
          Возможно, она была перемещена или удалена. А может, её и не было никогда
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent>
        <div class="flex gap-2">
          <Button variant="outline" @click="router.back()">
            <ChevronLeftIcon />
            Назад
          </Button>
          <Button @click="router.push('/')">
            <HouseIcon />
            На главную
          </Button>
        </div>
      </EmptyContent>
    </Empty>
  </div>
</template>
