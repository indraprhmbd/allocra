<script setup lang="ts">
export interface Column {
  key: string;
  label: string;
  class?: string;
  format?: (value: any) => string;
}

defineProps<{
  columns: Column[];
  items: any[];
  loading?: boolean;
}>();
</script>

<template>
  <div class="w-full">
    <!-- Desktop Table View -->
    <div class="hidden md:block overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-border bg-background/50">
            <th
              v-for="col in columns"
              :key="col.key"
              class="px-5 py-3 text-[10px] font-mono font-bold text-muted uppercase tracking-widest"
              :class="col.class"
            >
              {{ col.label }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-if="loading" v-for="i in 5" :key="i" class="animate-pulse">
            <td v-for="col in columns" :key="col.key" class="px-5 py-4">
              <div class="h-4 bg-border rounded-sm w-3/4"></div>
            </td>
          </tr>
          <tr v-else-if="items.length === 0">
            <td
              :colspan="columns.length"
              class="px-5 py-20 text-center text-muted font-mono text-xs"
            >
              NO_DATA_AVAILABLE
            </td>
          </tr>
          <tr
            v-else
            v-for="(item, idx) in items"
            :key="idx"
            class="hover:bg-accent/5 transition-colors group"
          >
            <td
              v-for="col in columns"
              :key="col.key"
              class="px-5 py-4 text-sm text-primary font-medium"
              :class="col.class"
            >
              <slot :name="col.key" :value="item[col.key]" :item="item">
                <span v-if="col.format">{{ col.format(item[col.key]) }}</span>
                <span v-else>{{ item[col.key] }}</span>
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Mobile Card View -->
    <div class="md:hidden space-y-4">
      <div
        v-if="loading"
        v-for="i in 3"
        :key="i"
        class="bg-surface border border-border rounded-sm p-4 animate-pulse space-y-3"
      >
        <div class="h-4 bg-border rounded-sm w-1/2"></div>
        <div class="h-3 bg-border rounded-sm w-3/4"></div>
        <div class="h-3 bg-border rounded-sm w-2/3"></div>
      </div>

      <div
        v-else-if="items.length === 0"
        class="text-center py-10 text-muted font-mono text-xs"
      >
        NO_DATA_AVAILABLE
      </div>

      <div
        v-else
        v-for="(item, idx) in items"
        :key="idx"
        class="bg-surface border border-border rounded-sm p-4 space-y-3"
      >
        <div
          v-for="col in columns"
          :key="col.key"
          class="flex justify-between items-start gap-4"
          :class="{ hidden: !col.label && !item[col.key] }"
        >
          <span
            class="text-[10px] font-mono font-bold text-muted uppercase tracking-widest shrink-0"
          >
            {{ col.label }}
          </span>
          <div class="text-sm text-primary font-medium text-right break-words">
            <slot :name="col.key" :value="item[col.key]" :item="item">
              <span v-if="col.format">{{ col.format(item[col.key]) }}</span>
              <span v-else>{{ item[col.key] }}</span>
            </slot>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
