<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Chart, registerables } from 'chart.js'
import { BookOpenIcon, ClockIcon, ImageIcon, PackageIcon, TrendingUpIcon, TriangleAlertIcon } from '@lucide/vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import PageHeader from '@/components/PageHeader.vue'
import StatCard from '@/components/StatCard.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { theme, accent } from '@/lib/theme'
import { isModerator } from '@/lib/auth'
import { CATEGORIES } from '@/lib/catalog'
import type { AnalyticsSummary, MaterialCategory } from '@/lib/types'

Chart.register(...registerables)

const summary = ref<AnalyticsSummary | null>(null)
const doughnutCanvas = ref<HTMLCanvasElement | null>(null)
const barCanvas = ref<HTMLCanvasElement | null>(null)
let doughnutChart: Chart<'doughnut'> | null = null
let barChart: Chart<'bar'> | null = null

const CATEGORY_COLORS: Record<MaterialCategory, string> = {
  photos: '--chart-1',
  videos: '--chart-2',
  library: '--chart-3',
  awards: '--chart-4',
  documents: '--chart-5',
}

const probe = document.createElement('canvas').getContext('2d', { willReadFrequently: true })

function cssColor(name: string, alpha = 1): string {
  const value = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  if (!probe || !value) return value
  probe.clearRect(0, 0, 1, 1)
  probe.fillStyle = value
  probe.fillRect(0, 0, 1, 1)
  const [r, g, b, a] = probe.getImageData(0, 0, 1, 1).data
  return `rgba(${r}, ${g}, ${b}, ${((a / 255) * alpha).toFixed(3)})`
}

function monthLabel(ym: string): string {
  const [y, m] = ym.split('-')
  const names = ['янв', 'фев', 'мар', 'апр', 'май', 'июн', 'июл', 'авг', 'сен', 'окт', 'ноя', 'дек']
  return `${names[Number(m) - 1]} ${y.slice(2)}`
}

function buildCharts() {
  if (!summary.value || !doughnutCanvas.value || !barCanvas.value) return
  const s = summary.value
  const muted = cssColor('--muted-foreground')
  const grid = cssColor('--border')
  const font = { family: getComputedStyle(document.body).fontFamily, size: 12 }
  const tooltip = {
    backgroundColor: cssColor('--popover'),
    titleColor: cssColor('--popover-foreground'),
    bodyColor: cssColor('--muted-foreground'),
    borderColor: cssColor('--border'),
    borderWidth: 1,
    padding: 10,
    cornerRadius: 8,
    boxPadding: 4,
    titleFont: { ...font, weight: 600 },
    bodyFont: font,
  }

  const byCat = CATEGORIES.map((c) => ({
    label: c.label,
    color: cssColor(CATEGORY_COLORS[c.id]),
    count: s.by_category.find((x) => x.category === c.id)?.count || 0,
  }))

  doughnutChart?.destroy()
  doughnutChart = new Chart(doughnutCanvas.value, {
    type: 'doughnut',
    data: {
      labels: byCat.map((c) => c.label),
      datasets: [
        {
          data: byCat.map((c) => c.count),
          backgroundColor: byCat.map((c) => c.color),
          borderColor: cssColor('--card'),
          borderWidth: 3,
          hoverOffset: 6,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: '65%',
      plugins: {
        legend: {
          position: 'bottom',
          labels: {
            color: muted,
            font,
            boxWidth: 8,
            boxHeight: 8,
            padding: 16,
            usePointStyle: true,
            pointStyle: 'circle',
          },
        },
        tooltip,
      },
    },
  })

  const monthly = s.monthly
  barChart?.destroy()
  barChart = new Chart(barCanvas.value, {
    type: 'bar',
    data: {
      labels: monthly.map((m) => monthLabel(m.month)),
      datasets: [
        {
          label: 'Добавлено материалов',
          data: monthly.map((m) => m.count),
          backgroundColor: cssColor('--primary'),
          hoverBackgroundColor: cssColor('--primary', 0.8),
          borderRadius: 6,
          maxBarThickness: 36,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false }, tooltip },
      scales: {
        x: { ticks: { color: muted, font }, grid: { display: false }, border: { display: false } },
        y: {
          beginAtZero: true,
          ticks: { color: muted, font, precision: 0 },
          grid: { color: grid },
          border: { display: false },
        },
      },
    },
  })
}

onMounted(async () => {
  try {
    summary.value = await api.get<AnalyticsSummary>('/analytics/summary')
    requestAnimationFrame(buildCharts)
  } catch (e) {
    notify.error('Не удалось загрузить аналитику', { text: errorMessage(e) })
  }
})

watch([theme, accent], () => requestAnimationFrame(buildCharts))

onBeforeUnmount(() => {
  doughnutChart?.destroy()
  barChart?.destroy()
})

function catCount(id: MaterialCategory): number {
  return summary.value?.by_category.find((x) => x.category === id)?.count || 0
}
</script>

<template>
  <div>
    <PageHeader title="Аналитика" description="Состав архива и динамика пополнения" />

    <template v-if="summary">
      <div class="mb-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard label="Всего единиц в базе" :value="summary.total" :icon="PackageIcon" />
        <StatCard
          label="За последний год"
          :value="summary.monthly.reduce((s, m) => s + m.count, 0)"
          :icon="TrendingUpIcon"
        />
        <StatCard v-if="isModerator" label="Взято книг" :value="summary.active_loans" :icon="ClockIcon" />
        <StatCard v-else label="Библиотека" :value="catCount('library')" :icon="BookOpenIcon" />
        <StatCard
          v-if="isModerator"
          label="Просрочено возвратов"
          :value="summary.overdue_loans"
          :icon="TriangleAlertIcon"
        />
        <StatCard v-else label="Фотографии" :value="catCount('photos')" :icon="ImageIcon" />
      </div>

      <div class="grid grid-cols-1 gap-6 lg:grid-cols-5">
        <Card class="lg:col-span-2">
          <CardHeader>
            <CardTitle>Соотношение по категориям</CardTitle>
            <CardDescription>Все материалы архива</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="h-72">
              <canvas ref="doughnutCanvas" />
            </div>
          </CardContent>
        </Card>
        <Card class="lg:col-span-3">
          <CardHeader>
            <CardTitle>Динамика добавления материалов</CardTitle>
            <CardDescription>За последние 12 месяцев</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="h-72">
              <canvas ref="barCanvas" />
            </div>
          </CardContent>
        </Card>
      </div>

      <div class="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-5">
        <Card v-for="c in CATEGORIES" :key="c.id" size="sm">
          <CardHeader>
            <CardDescription class="flex items-center gap-1.5">
              <span class="size-2 rounded-full" :style="{ background: `var(${CATEGORY_COLORS[c.id]})` }" />
              {{ c.label }}
            </CardDescription>
            <CardTitle class="text-2xl font-semibold tabular-nums group-data-[size=sm]/card:text-2xl">{{ catCount(c.id) }}</CardTitle>
          </CardHeader>
        </Card>
      </div>
    </template>

    <div v-else class="space-y-6">
      <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <Skeleton v-for="i in 4" :key="i" class="h-24 rounded-xl" />
      </div>
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-5">
        <Skeleton class="h-96 rounded-xl lg:col-span-2" />
        <Skeleton class="h-96 rounded-xl lg:col-span-3" />
      </div>
    </div>
  </div>
</template>
