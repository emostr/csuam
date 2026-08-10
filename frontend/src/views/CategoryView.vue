<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NButton from '@/components/ui/NButton.vue'
import NInput from '@/components/ui/NInput.vue'
import NDateInput from '@/components/ui/NDateInput.vue'
import MaterialCard from '@/components/MaterialCard.vue'
import { api } from '@/lib/api'
import { notify } from '@/lib/notify'
import { isHeadTeacher } from '@/lib/auth'
import { categoryMeta } from '@/lib/catalog'

const route = useRoute()
const router = useRouter()

const category = computed(() => route.meta.category)
const meta = computed(() => categoryMeta(category.value))
// подзаголовки убраны
const subtitles = {
  photos: '',
  videos: '',
  library: '',
  awards: '',
  documents: '',
}

const canAdd = computed(() => category.value !== 'documents' || isHeadTeacher.value)

const materials = ref([])
const loading = ref(true)
const q = ref('')
const from = ref('')
const to = ref('')

async function load() {
  loading.value = true
  try {
    const params = new URLSearchParams({ category: category.value })
    if (q.value) params.set('q', q.value)
    if (from.value) params.set('from', from.value)
    if (to.value) params.set('to', to.value)
    materials.value = await api.get('/materials?' + params.toString())
  } catch (e) {
    notify.error('Не удалось загрузить материалы', { text: e.message })
  } finally {
    loading.value = false
  }
}

watch([category, from, to], load)
onMounted(load)
</script>

<template>
  <div>
    <NPageHeader :title="meta.label" :subtitle="subtitles[category] || ''">
      <template #actions>
        <NButton v-if="canAdd" icon="plus" @click="router.push({ path: '/add', query: { category } })">
          Добавить
        </NButton>
      </template>
    </NPageHeader>

    <div class="bg-surface border border-line p-4 mb-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 items-end">
      <form class="lg:col-span-2" @submit.prevent="load">
        <NInput v-model="q" label="Поиск" placeholder="Ключевые слова…" icon="search" @keyup.enter="load" />
      </form>
      <NDateInput v-model="from" label="Происхождение с" />
      <NDateInput v-model="to" label="Происхождение по" />
    </div>

    <div v-if="loading" class="text-center text-muted py-16">Загрузка…</div>
    <div v-else-if="!materials.length" class="bg-surface border border-line py-16 text-center text-muted">
      {{ category === 'documents' && !isHeadTeacher
        ? 'Нет доступных документов. Доступ к ним выдаёт завуч.'
        : 'В этой категории пока нет материалов.' }}
    </div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <MaterialCard v-for="m in materials" :key="m.id" :material="m" />
    </div>
  </div>
</template>
