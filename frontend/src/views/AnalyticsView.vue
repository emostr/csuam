<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Chart, registerables } from 'chart.js'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NCard from '@/components/ui/NCard.vue'
import NTile from '@/components/ui/NTile.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { theme, accent } from '@/lib/theme'
import { isModerator } from '@/lib/auth'
import { CATEGORIES, categoryMeta } from '@/lib/catalog'
import type { AnalyticsSummary, MaterialCategory } from '@/lib/types'

Chart.register(...registerables)

const summary = ref<AnalyticsSummary | null>(null)
const doughnutCanvas = ref<HTMLCanvasElement | null>(null)
const barCanvas = ref<HTMLCanvasElement | null>(null)
let doughnutChart: Chart<'doughnut'> | null = null
let barChart: Chart<'bar'> | null = null

function cssVar(name: string): string {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

function monthLabel(ym: string): string {
  const [y, m] = ym.split('-')
  const names = ['янв', 'фев', 'мар', 'апр', 'май', 'июн', 'июл', 'авг', 'сен', 'окт', 'ноя', 'дек']
  return `${names[Number(m) - 1]} ${y.slice(2)}`
}

function buildCharts() {
  if (!summary.value || !doughnutCanvas.value || !barCanvas.value) return
  const s = summary.value
  const muted = cssVar('--ng-muted')
  const line = cssVar('--ng-line')
  const accentColor = cssVar('--ng-accent')

  const catColors: Record<MaterialCategory, string> = {
    photos: cssVar('--ng-info'),
    videos: cssVar('--ng-warning'),
    library: accentColor,
    awards: cssVar('--ng-success'),
    documents: cssVar('--ng-danger'),
  }

  const byCat = CATEGORIES.map((c) => ({
    label: c.label,
    id: c.id,
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
          backgroundColor: byCat.map((c) => catColors[c.id]),
          borderColor: cssVar('--ng-surface'),
          borderWidth: 2,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: '62%',
      plugins: {
        legend: {
          position: 'bottom',
          labels: { color: muted, boxWidth: 12, boxHeight: 12, padding: 14 },
        },
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
          backgroundColor: accentColor,
          maxBarThickness: 36,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: muted }, grid: { display: false } },
        y: {
          beginAtZero: true,
          ticks: { color: muted, precision: 0 },
          grid: { color: line },
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
    <NPageHeader title="Аналитика" />

    <template v-if="summary">
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <NTile label="Всего единиц в базе" :value="summary.total" icon="box" solid />
        <NTile label="За последний год" :value="summary.monthly.reduce((s, m) => s + m.count, 0)" icon="trendUp" />
        <NTile v-if="isModerator" label="Сейчас на прокате" :value="summary.active_loans" icon="clock" />
        <NTile v-else label="Библиотека" :value="catCount('library')" icon="book" />
        <NTile
          v-if="isModerator"
          label="Просрочено возвратов"
          :value="summary.overdue_loans"
          icon="alert"
        />
        <NTile v-else label="Фотографии" :value="catCount('photos')" icon="image" />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
        <NCard title="Соотношение по категориям" class="lg:col-span-2">
          <div class="h-72">
            <canvas ref="doughnutCanvas" />
          </div>
        </NCard>
        <NCard title="Динамика добавления материалов" subtitle="За последние 12 месяцев" class="lg:col-span-3">
          <div class="h-72">
            <canvas ref="barCanvas" />
          </div>
        </NCard>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-5 gap-4 mt-6">
        <div
          v-for="c in CATEGORIES"
          :key="c.id"
          class="bg-surface border border-line p-4 text-center"
        >
          <div class="text-2xl font-extrabold text-ink tabular-nums">{{ catCount(c.id) }}</div>
          <div class="text-xs text-muted mt-1">{{ categoryMeta(c.id).label }}</div>
        </div>
      </div>
    </template>
    <div v-else class="text-center text-muted py-16">Загрузка…</div>
  </div>
</template>
