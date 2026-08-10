<script setup>
defineProps({
  columns: { type: Array, default: () => [] },
  rows: { type: Array, default: () => [] },
  hover: { type: Boolean, default: true },
})
</script>

<template>
  <div class="w-full overflow-x-auto border border-line bg-surface">
    <table class="w-full text-sm border-collapse">
      <thead>
        <tr class="border-b border-line">
          <th
            v-for="col in columns"
            :key="col.key"
            class="ng-label text-muted text-left px-4 py-3 whitespace-nowrap"
            :class="col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : ''"
            :style="col.width ? { width: col.width } : {}"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(row, i) in rows"
          :key="row.id ?? i"
          class="border-b border-line last:border-0 transition-colors"
          :class="hover ? 'hover:bg-surface-2' : ''"
        >
          <td
            v-for="col in columns"
            :key="col.key"
            class="px-4 py-3 text-ink align-middle"
            :class="col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : ''"
          >
            <slot :name="col.key" :row="row" :value="row[col.key]">
              {{ row[col.key] }}
            </slot>
          </td>
        </tr>
        <tr v-if="!rows.length">
          <td :colspan="columns.length" class="px-4 py-10 text-center text-muted text-sm">
            Нет данных
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
