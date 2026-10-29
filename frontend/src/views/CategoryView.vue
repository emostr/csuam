<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PlusIcon, SearchIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { InputGroup, InputGroupAddon, InputGroupInput } from '@/components/ui/input-group'
import PageHeader from '@/components/PageHeader.vue'
import DatePicker from '@/components/DatePicker.vue'
import MaterialGrid from '@/components/MaterialGrid.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { isHeadTeacher } from '@/lib/auth'
import { categoryMeta } from '@/lib/catalog'
import type { Material, MaterialCategory } from '@/lib/types'

const route = useRoute()
const router = useRouter()

const category = computed<MaterialCategory>(() => route.meta.category ?? 'photos')
const meta = computed(() => categoryMeta(category.value))

const canAdd = computed(() => category.value !== 'documents' || isHeadTeacher.value)

const materials = ref<Material[]>([])
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
    materials.value = await api.get<Material[]>('/materials?' + params.toString())
  } catch (e) {
    notify.error('Не удалось загрузить материалы', { text: errorMessage(e) })
  } finally {
    loading.value = false
  }
}

watch([category, from, to], load)
onMounted(load)
</script>

<template>
  <div>
    <PageHeader :title="meta.label">
      <template #actions>
        <Button v-if="canAdd" @click="router.push({ path: '/add', query: { category } })">
          <PlusIcon />
          Добавить
        </Button>
      </template>
    </PageHeader>

    <div class="mb-6 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-[minmax(0,2fr)_repeat(2,minmax(0,1fr))]">
      <form class="sm:col-span-2 lg:col-span-1" role="search" @submit.prevent="load">
        <InputGroup>
          <InputGroupAddon>
            <SearchIcon />
          </InputGroupAddon>
          <InputGroupInput v-model="q" placeholder="Поиск по названию и описанию…" aria-label="Поиск" />
        </InputGroup>
      </form>
      <DatePicker v-model="from" placeholder="Происхождение с" />
      <DatePicker v-model="to" placeholder="Происхождение по" />
    </div>

    <MaterialGrid
      :materials="materials"
      :loading="loading"
      :empty-title="category === 'documents' && !isHeadTeacher ? 'Нет доступных документов' : 'Пока пусто'"
      :empty-text="
        category === 'documents' && !isHeadTeacher
          ? 'Доступ к документам выдаёт завуч.'
          : 'В этой категории пока нет материалов.'
      "
    />
  </div>
</template>
